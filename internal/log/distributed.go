package log

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	api "github.com/justagabriel/proglog/api/v1"
	"google.golang.org/protobuf/proto"
)

type DistributedLog struct {
	config Config
	log    *Log
	raft   *raft.Raft
}

func NewDistributedLog(dataDir string, config Config) (*DistributedLog, error) {
	l := &DistributedLog{config: config}

	err := l.setupLog(dataDir)
	if err != nil {
		return nil, err
	}

	err = l.setupRaft(dataDir)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *DistributedLog) setupLog(dataDir string) error {
	logDir := filepath.Join(dataDir, "log")
	err := os.MkdirAll(logDir, 0775)
	if err != nil {
		return err
	}

	l.log, err = NewLog(logDir, l.config)
	return err
}

func (l *DistributedLog) setupRaft(dataDir string) error {
	fsm := &fsm{log: l.log}

	logDir := filepath.Join(dataDir, "raft", "log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	logConfig := l.config
	logConfig.Segment.InitialOffset = 1
	logStore, err := newLogStore(logDir, logConfig)
	if err != nil {
		return err
	}

	stableStorePath := filepath.Join(dataDir, "raft", "stable")
	stableStore, err := raftboltdb.NewBoltStore(stableStorePath)
	if err != nil {
		return err
	}

	retain := 1
	snapshotFilePath := filepath.Join(dataDir, "raft")
	snapshotStore, err := raft.NewFileSnapshotStore(snapshotFilePath, retain, os.Stderr)
	if err != nil {
		return err
	}

	maxPool := 5
	timeout := 10 * time.Second
	transport := raft.NewNetworkTransport(
		l.config.Raft.StreamLayer,
		maxPool,
		timeout,
		os.Stderr,
	)

	config := raft.DefaultConfig()
	config.LocalID = l.config.Raft.LocalID
	if l.config.Raft.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = l.config.Raft.LeaderLeaseTimeout
	}
	if l.config.Raft.CommitTimeout != 0 {
		config.CommitTimeout = l.config.Raft.CommitTimeout
	}

	l.raft, err = raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return err
	}

	hasState, err := raft.HasExistingState(logStore, stableStore, snapshotStore)
	if err != nil {
		return err
	}

	if l.config.Raft.Bootstrap && !hasState {
		config := raft.Configuration{
			Servers: []raft.Server{{
				ID:      config.LocalID,
				Address: raft.ServerAddress(l.config.Raft.BindAddr),
			}},
		}
		err = l.raft.BootstrapCluster(config).Error()
	}

	return err
}

func (l *DistributedLog) Append(record *api.Record) (uint64, error) {
	res, err := l.apply(AppendRequestType, &api.CreateRecordRequest{Record: record})
	if err != nil {
		return 0, err
	}
	return res.(*api.CreateRecordResponse).Offset, nil
}

func (l *DistributedLog) apply(reqType RequestType, req proto.Message) (interface{}, error) {
	var buf bytes.Buffer
	_, err := buf.Write([]byte{byte(reqType)})
	if err != nil {
		return nil, err
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}

	timemout := 10 * time.Second
	future := l.raft.Apply(buf.Bytes(), timemout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res := future.Response()
	err, isErr := res.(error)
	if isErr {
		return nil, err
	}

	return res, nil
}

func (l *DistributedLog) Read(offset uint64) (*api.Record, error) {
	return l.log.Read(offset)
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	log *Log
}

func (l *DistributedLog) Join(id, addr string) error {
	configFuture := l.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}

	serverID := raft.ServerID(id)
	serverAddr := raft.ServerAddress(addr)
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == serverID || srv.Address == serverAddr {
			if srv.ID == serverID && srv.Address == serverAddr {
				// server has already joined
				return nil
			}
			//remove existing server
			removeFuture := l.raft.RemoveServer(serverID, 0, 0)
			if err := removeFuture.Error(); err != nil {
				return err
			}
		}
	}
	addFuture := l.raft.AddVoter(serverID, serverAddr, 0, 0)
	if err := addFuture.Error(); err != nil {
		return err
	}
	return nil
}

func (l *DistributedLog) Leave(id string) error {
	removeFuture := l.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (l *DistributedLog) WaitForLeader(timeout time.Duration) error {
	timeoutc := time.After(timeout)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeoutc:
			return fmt.Errorf("timed out")
		case <-ticker.C:
			if leader := l.raft.Leader(); leader != "" {
				return nil
			}
		}
	}
}

// Close disconnects from the Raft cluster and shut's down the replication service.
func (l *DistributedLog) Close() error {
	f := l.raft.Shutdown()
	if err := f.Error(); err != nil {
		return err
	}
	return l.log.Close()
}

// GetServers returns information about all servers of the app.
func (l *DistributedLog) GetServers() ([]*api.Server, error) {
	future := l.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}
	var servers []*api.Server
	for _, srv := range future.Configuration().Servers {
		servers = append(servers, &api.Server{
			Id:       string(srv.ID),
			RpcAddr:  string(srv.Address),
			IsLeader: l.raft.Leader() == srv.Address,
		})
	}
	return servers, nil
}

type RequestType uint8

const (
	AppendRequestType RequestType = 0
)

// Apply implements raft.FSM.
func (l *fsm) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqTyep := RequestType(buf[0])
	switch reqTyep {
	case AppendRequestType:
		return l.applyAppend(buf[1:])
	}
	return nil
}

func (l *fsm) applyAppend(b []byte) interface{} {
	var req api.CreateRecordRequest
	err := proto.Unmarshal(b, &req)
	if err != nil {
		return err
	}
	offset, err := l.log.Append(req.Record)
	if err != nil {
		return err
	}
	return &api.CreateRecordResponse{
		Offset: offset,
	}
}

// Snapshot implements raft.FSM.
func (m *fsm) Snapshot() (raft.FSMSnapshot, error) {
	r := m.log.Reader()
	return &snapshot{reader: r}, nil
}

var _ raft.FSMSnapshot = (*snapshot)(nil)

type snapshot struct {
	reader io.Reader
}

// Persist implements raft.FSMSnapshot.
func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	_, err := io.Copy(sink, s.reader)
	if err != nil {
		return err
	}
	return sink.Close()
}

// Release implements raft.FSMSnapshot.
func (*snapshot) Release() {}

func (f *fsm) Restore(r io.ReadCloser) error {
	b := make([]byte, lenWidth)
	var buf bytes.Buffer
	for i := 0; ; i++ {
		_, err := io.ReadFull(r, b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		size := int64(enc.Uint64(b))
		_, err = io.CopyN(&buf, r, size)
		if err != nil {
			return err
		}

		record := &api.Record{}
		err = proto.Unmarshal(buf.Bytes(), record)
		if err != nil {
			return err
		}

		if i == 0 {
			f.log.Config.Segment.InitialOffset = record.Offset
			err = f.log.Reset()
			if err != nil {
				return err
			}
		}

		_, err = f.log.Append(record)
		if err != nil {
			return err
		}
		buf.Reset()
	}

	return nil
}

var _ raft.LogStore = (*logStore)(nil)

type logStore struct {
	*Log
}

func newLogStore(dir string, c Config) (*logStore, error) {
	log, err := NewLog(dir, c)
	if err != nil {
		return nil, err
	}
	return &logStore{log}, nil
}

func (l *logStore) FirstIndex() (uint64, error) {
	return l.LowestOffset()
}

// LastIndex implements raft.LogStore.
func (l *logStore) LastIndex() (uint64, error) {
	return l.HighestOffset()
}

func (l *logStore) GetLog(index uint64, out *raft.Log) error {
	in, err := l.Read(index)
	if err != nil {
		return err
	}

	out.Data = in.Value
	out.Index = in.Offset
	out.Type = raft.LogType(in.Type)
	out.Term = in.Term
	return nil
}

func (l *logStore) StoreLog(record *raft.Log) error {
	return l.StoreLogs([]*raft.Log{record})
}

func (l *logStore) StoreLogs(records []*raft.Log) error {
	for _, record := range records {
		apiRec := &api.Record{
			Value: record.Data,
			Term:  record.Term,
			Type:  uint32(record.Type),
		}
		_, err := l.Append(apiRec)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteRange implements raft.LogStore.
func (l *logStore) DeleteRange(min uint64, max uint64) error {
	panic("unimplemented")
}

var _ raft.StreamLayer = new(StreamLayer)

type StreamLayer struct {
	listener        net.Listener
	serverTLSConfig *tls.Config
	peerTLSConfig   *tls.Config
}

func NewStreamLayer(
	listener net.Listener,
	serverTLSConfig *tls.Config,
	peerTLSConfig *tls.Config,
) *StreamLayer {
	return &StreamLayer{
		listener:        listener,
		serverTLSConfig: serverTLSConfig,
		peerTLSConfig:   peerTLSConfig,
	}
}

const RaftRPC = 1

func (s *StreamLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", string(addr))
	if err != nil {
		return nil, err
	}

	// identify to mux this is a raft rpc
	_, err = conn.Write([]byte{byte(RaftRPC)})
	if err != nil {
		return nil, err
	}

	if s.peerTLSConfig != nil {
		conn = tls.Client(conn, s.peerTLSConfig)
	}

	return conn, err
}

func (s *StreamLayer) Accept() (net.Conn, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	b := make([]byte, 1)
	_, err = conn.Read(b)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal([]byte{byte(RaftRPC)}, b) {
		return nil, fmt.Errorf("not a raft rpc")
	}

	if s.serverTLSConfig != nil {
		return tls.Server(conn, s.serverTLSConfig), nil
	}

	return conn, nil
}

func (s *StreamLayer) Close() error {
	return s.listener.Close()
}

func (s *StreamLayer) Addr() net.Addr {
	return s.listener.Addr()
}

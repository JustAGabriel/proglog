package log

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal"
	"github.com/stretchr/testify/require"
)

func TestMultipleNodes(t *testing.T) {
	// arrange
	var logs []*DistributedLog
	nodeCount := 3

	for i := 0; i < nodeCount; i++ {
		dataDir := internal.GetTempDir(t, "distributed-log-test")
		defer func(dir string) {
			_ = os.RemoveAll(dir)
		}(dataDir)

		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", internal.FreePort(t)))
		require.NoError(t, err)

		config := Config{}
		config.Raft.StreamLayer = NewStreamLayer(ln, nil, nil)
		config.Raft.LocalID = raft.ServerID(fmt.Sprintf("%d", i))
		config.Raft.HeartbeatTimeout = 50 * time.Millisecond
		config.Raft.ElectionTimeout = 50 * time.Millisecond
		config.Raft.LeaderLeaseTimeout = 50 * time.Millisecond
		config.Raft.CommitTimeout = 5 * time.Millisecond
		config.Raft.BindAddr = ln.Addr().String()

		if i == 0 {
			config.Raft.Bootstrap = true
		}

		dlog, err := NewDistributedLog(dataDir, config)
		require.NoError(t, err)

		if i != 0 {
			err = logs[0].Join(
				fmt.Sprintf("%d", i),
				ln.Addr().String(),
			)
			require.NoError(t, err)
		} else {
			err = dlog.WaitForLeader(3 * time.Second)
			require.NoError(t, err)
		}

		logs = append(logs, dlog)
	}

	records := []*api.Record{
		{Value: []byte("first")},
		{Value: []byte("second")},
	}

	// assert  that logs are replicated to followers
	for _, r := range records {
		off, err := logs[0].Append(r)
		require.NoError(t, err)
		require.Eventually(
			t,
			func() bool {
				for i := 0; i < nodeCount; i++ {
					got, err := logs[i].Read(off)
					if err != nil {
						return false
					}

					r.Offset = off
					if !reflect.DeepEqual(got.Value, r.Value) {
						return false
					}
				}
				return true
			},
			500*time.Millisecond,
			50*time.Millisecond,
		)
	}

	servers, err := logs[0].GetServers()
	require.NoError(t, err)
	require.Equal(t, 3, len(servers))
	require.True(t, servers[0].IsLeader)
	require.False(t, servers[1].IsLeader)
	require.False(t, servers[2].IsLeader)

	err = logs[0].Leave("1")
	require.NoError(t, err)

	servers, err = logs[0].GetServers()
	require.NoError(t, err)
	require.Equal(t, 2, len(servers))
	require.True(t, servers[0].IsLeader)
	require.False(t, servers[1].IsLeader)

	time.Sleep(50 * time.Millisecond)

	off, err := logs[0].Append(&api.Record{
		Value: []byte("third"),
	})
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	record, err := logs[1].Read(off)
	require.IsType(t, api.ErrOffsetOutOfRange{}, err)
	require.Nil(t, record)

	record, err = logs[2].Read(off)
	require.NoError(t, err)
	require.Equal(t, []byte("third"), record.Value)
	require.Equal(t, off, record.Offset)
}

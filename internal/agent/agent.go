package agent

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/justagabriel/proglog/internal/auth"
	"github.com/justagabriel/proglog/internal/discovery"
	"github.com/justagabriel/proglog/internal/log"
	"github.com/justagabriel/proglog/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/justagabriel/proglog/api/v1"
)

// Agent encapsulates all components of a Node.
type Agent struct {
	Config

	log        *log.Log
	server     *grpc.Server
	membership *discovery.Membership
	replicator *log.Replicator

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

// Config aggregates all configurations needed to create a new Agent.
type Config struct {
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	DataDir         string
	BindAddr        string
	RPCPort         int
	NodeName        string
	StartJoinAddr   []string
	ACLModelFile    string
	ACLPolicyFile   string
}

// RPCAddr returns the URI of the Agent client.
func (c *Config) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

// New creates a new Agent.
func New(config Config) (*Agent, error) {
	a := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
	}

	// order is crucial
	setup := []func() error{
		a.setupLogger,
		a.setupLog,
		a.setupServer,
		a.setupMembership,
	}

	for _, fn := range setup {
		err := fn()
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *Agent) setupLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

func (a *Agent) setupLog() error {
	var err error
	a.log, err = log.NewLog(a.DataDir, log.Config{})
	return err
}

func (a *Agent) setupServer() error {
	authorizer, err := auth.New(a.Config.ACLModelFile, a.Config.ACLPolicyFile)
	if err != nil {
		return err
	}

	serverConfig := &server.Config{
		CommitLog:  a.log,
		Authorizer: authorizer,
	}

	var opts []grpc.ServerOption
	if a.Config.ServerTLSConfig != nil {
		creds := credentials.NewTLS(a.Config.ServerTLSConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	a.server, err = server.NewGRPCServer(serverConfig, opts...)
	if err != nil {
		return err
	}

	rpcAddr, err := a.RPCAddr()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}

	go func() {
		err := a.server.Serve(listener)
		if err != nil {
			_ = a.Shutdown()
		}
	}()

	return err
}

func (a *Agent) setupMembership() error {
	rpcAddr, err := a.Config.RPCAddr()
	if err != nil {
		return err
	}

	var opts []grpc.DialOption
	if a.PeerTLSConfig != nil {
		creds := credentials.NewTLS(a.PeerTLSConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.Dial(rpcAddr, opts...)
	if err != nil {
		return err
	}

	client := api.NewLogClient(conn)
	a.replicator = &log.Replicator{
		DialOptions: opts,
		LocalServer: client,
	}

	discoveryConfig := discovery.Config{
		NodeName: a.NodeName,
		BindAddr: a.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: a.StartJoinAddr,
	}
	a.membership, err = discovery.New(a.replicator, discoveryConfig)
	return err
}

// Shutdown free's up all resources hold by this Agent.
func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}
	a.shutdown = true
	close(a.shutdowns)

	shutdownFuncs := []func() error{
		a.membership.Leave,
		a.replicator.Close,
		func() error {
			a.server.GracefulStop()
			return nil
		},
		a.log.Close,
	}

	for _, fn := range shutdownFuncs {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

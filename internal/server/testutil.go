package server

import (
	"net"
	"testing"
	"time"

	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal"
	"github.com/justagabriel/proglog/internal/auth"
	"github.com/justagabriel/proglog/internal/config"
	"github.com/justagabriel/proglog/internal/log"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/examples/exporter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// LogServerTestSetup contains all entities necessary to create a test for a LogServer.
type LogServerTestSetup struct {
	// AuthorizedClient is an authenicated grpc client, able to communicate with the created server.
	AuthorizedClient api.LogClient

	// UnauthorizedClient is an authenicated grpc client, unable to communicate with the created server.
	UnauthorizedClient api.LogClient

	// Config represents internal LogServer entities.
	Config *Config

	// Teardown will release all resources bound to the test server instance.
	Teardown func()

	// LogServerAddr contains the address under which the LogServer is reachable.
	LogServerAddr string
}

// SetupTest creates a new 'Log' server and returns an authorized and an unauthorized client of it
// additionally the config of the server itself, and a teardown function is returned.
func SetupTest(t *testing.T, fn func(*Config), isDebugMode *bool) LogServerTestSetup {
	t.Helper()

	setup := LogServerTestSetup{}

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	serverAddr := listener.Addr().String()
	setup.LogServerAddr = serverAddr

	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: serverAddr,
		Server:        true,
	})
	require.NoError(t, err)
	serverCreds := credentials.NewTLS(serverTLSConfig)

	dir := internal.GetTempDir(t, "server-test")
	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)

	authorizer, err := auth.New(config.ACLModelFile, config.ACLPolicyFile)
	require.NoError(t, err)

	var telemetryExporter *exporter.LogExporter
	if *isDebugMode {
		metricsLogFile := internal.GetTempFile(t, "", "metrics-*.log")
		t.Logf("metrics log file: %q", metricsLogFile.Name())

		tracesLogFile := internal.GetTempFile(t, "", "traces-*.log")
		t.Logf("traces log file: %q", metricsLogFile.Name())

		telemetryExporter, err = exporter.NewLogExporter(exporter.Options{
			ReportingInterval: time.Second,
			MetricsLogFile:    metricsLogFile.Name(),
			TracesLogFile:     tracesLogFile.Name(),
		})
		require.NoError(t, err)
		err = telemetryExporter.Start()
		require.NoError(t, err)
	}

	setup.Config = &Config{
		CommitLog:  clog,
		Authorizer: authorizer,
	}
	if fn != nil {
		fn(setup.Config)
	}
	server, err := NewGRPCServer(setup.Config, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go func() {
		server.Serve(listener)
	}()

	newClient := func(crtPath, keyPath string) (*grpc.ClientConn, api.LogClient, []grpc.DialOption) {
		clientTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
			CertFile: crtPath,
			KeyFile:  keyPath,
			CAFile:   config.CAFile,
		})
		require.NoError(t, err)

		clientCreds := credentials.NewTLS(clientTLSConfig)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(clientCreds)}
		cc, err := grpc.Dial(listener.Addr().String(), opts...)
		require.NoError(t, err)

		client := api.NewLogClient(cc)
		return cc, client, opts
	}

	rootConn, rootClient, _ := newClient(config.RootClientCertFile, config.RootClientKeyFile)
	nobodyConn, nobodyClient, _ := newClient(config.NobodyClientCertFile, config.NobodyClientKeyFile)

	setup.AuthorizedClient = rootClient
	setup.UnauthorizedClient = nobodyClient
	setup.Teardown = func() {
		server.Stop()
		rootConn.Close()
		nobodyConn.Close()
		listener.Close()
		clog.Remove()

		if telemetryExporter != nil {
			time.Sleep(1500 * time.Millisecond)
			telemetryExporter.Stop()
			telemetryExporter.Close()
		}
	}

	return setup
}

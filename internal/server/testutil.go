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

// SetupTest creates a new 'Log' server and returns an authorized and an unauthorized client of it
// additionally the config of the server itself, and a teardown function is returned.
func SetupTest(t *testing.T, fn func(*Config), isDebugMode *bool) (rootClient api.LogClient, nobodyClient api.LogClient, cfg *Config, teardown func()) {
	t.Helper()

	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: l.Addr().String(),
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

	cfg = &Config{
		CommitLog:  clog,
		Authorizer: authorizer,
	}
	if fn != nil {
		fn(cfg)
	}
	server, err := NewGRPCServer(cfg, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go func() {
		server.Serve(l)
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
		cc, err := grpc.Dial(l.Addr().String(), opts...)
		require.NoError(t, err)

		client := api.NewLogClient(cc)
		return cc, client, opts
	}

	rootConn, rootClient, _ := newClient(config.RootClientCertFile, config.RootClientKeyFile)
	nobodyConn, nobodyClient, _ := newClient(config.NobodyClientCertFile, config.NobodyClientKeyFile)

	return rootClient, nobodyClient, cfg, func() {
		server.Stop()
		rootConn.Close()
		nobodyConn.Close()
		l.Close()
		clog.Remove()

		if telemetryExporter != nil {
			time.Sleep(1500 * time.Millisecond)
			telemetryExporter.Stop()
			telemetryExporter.Close()
		}
	}
}

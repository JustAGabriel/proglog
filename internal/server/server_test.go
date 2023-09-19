package server

import (
	"context"
	"net"
	"testing"

	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal"
	"github.com/justagabriel/proglog/internal/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	scenarios := map[string]func(t *testing.T, client api.LogClient, config *Config){
		"create/get a message from/to the log succeeds": testCreateGet,
		"consume past log boundary fails":               testGetPastBoundary,
		// "create/get a stream succeeds":                  testCreateGetStream,
	}

	for title, scenario := range scenarios {
		t.Run(title, func(t *testing.T) {
			client, config, teardown := setupTest(t, nil)
			defer teardown()
			scenario(t, client, config)
		})
	}
}

func setupTest(t *testing.T, fn func(*Config)) (client api.LogClient, cfg *Config, teardown func()) {
	t.Helper()
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	dir := internal.GetTempDir("server-test")
	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)

	cfg = &Config{
		CommitLog: clog,
	}
	if fn != nil {
		fn(cfg)
	}
	server, err := NewGRPCServer(cfg)
	require.NoError(t, err)

	go func() {
		server.Serve(l)
	}()

	clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	client = api.NewLogClient(cc)
	return client, cfg, func() {
		server.Stop()
		cc.Close()
		l.Close()
		clog.Remove()
	}
}

func testCreateGet(t *testing.T, client api.LogClient, config *Config) {
	ctx := context.Background()
	want := &api.Record{
		Value: []byte("hello world"),
	}

	createResp, err := client.Create(
		ctx,
		&api.CreateRecordRequest{
			Record: want,
		},
	)
	require.NoError(t, err)

	getResp, err := client.Get(
		ctx,
		&api.GetRecordRequest{
			Offset: createResp.Offset,
		},
	)

	require.NoError(t, err)
	require.Equal(t, want.Value, getResp.Record.Value)
	require.Equal(t, want.Offset, getResp.Record.Offset)
}

func testGetPastBoundary(t *testing.T, client api.LogClient, config *Config) {
	ctx := context.Background()

	createReq := &api.CreateRecordRequest{
		Record: &api.Record{
			Value: []byte("hello world!"),
		},
	}
	createResp, err := client.Create(ctx, createReq)
	require.NoError(t, err)

	getReq := &api.GetRecordRequest{
		Offset: createResp.Offset + 1,
	}
	getResp, err := client.Get(ctx, getReq)
	if getResp != nil {
		t.Fatal("expected no response, since invalid request!")
	}

	got := status.Code(err)
	want := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	if got != want {
		t.Fatalf("got err: %v, want: %v", got, want)
	}
}

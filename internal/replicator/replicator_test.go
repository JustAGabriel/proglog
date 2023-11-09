package replicator

import (
	"context"
	"sync"
	"testing"

	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestReplicator(t *testing.T) {
	// arrange
	records := []api.Record{
		{
			Value:  []byte("hello world 1!"),
			Offset: 0,
		},
		{
			Value:  []byte("hello world 2!"),
			Offset: 1,
		},
	}

	server1Client, _, _, teardownServer1 := server.SetupTest(t, nil, nil)
	ctx := context.Background()
	createStream, err := server1Client.CreateStream(ctx)
	require.NoError(t, err)

	for offset, record := range records {
		createReq := &api.CreateRecordRequest{
			Record: &record,
		}
		err = createStream.Send(createReq)
		require.NoError(t, err)

		createResp, err := createStream.Recv()
		require.NoError(t, err)
		if createResp.Offset != uint64(offset) {
			t.Fatalf("got offset: %d, want: %d", createResp.Offset, offset)
		}
	}

	server2Client, _, server2Config, teardownServer2 := server.SetupTest(t, nil, nil)
	getStream, err := server1Client.GetStream(ctx)
	require.NoError(t, err)

	replicator := Replicator{
		DialOptions: []grpc.DialOption{},
		LocalServer: nil,
		logger:      &zap.Logger{},
		mtx:         sync.Mutex{},
		servers:     map[string]chan struct{}{},
		closed:      false,
		close:       make(chan struct{}),
	}
	// setup second log server
	// create replicator
	// connect replicator to second server
	// configure replicator to replicate from frist server

	// act

	//assert
	// get record created in first server by requesting from second server
	teardownServer1()
	teardownServer2()
}

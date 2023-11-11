package replicator

import (
	"context"
	"sync"
	"testing"
	"time"

	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal/config"
	"github.com/justagabriel/proglog/internal/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	falsePtr := new(bool)
	*falsePtr = false

	server1Setup := server.SetupTest(t, nil, falsePtr)

	server1Client := server1Setup.AuthorizedClient
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

	server2Setup := server.SetupTest(t, nil, falsePtr)
	server2Client := server2Setup.AuthorizedClient

	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: "127.0.0.1",
		Server:        false,
	})

	creds := credentials.NewTLS(peerTLSConfig)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	replicator := Replicator{
		DialOptions: []grpc.DialOption{grpc.WithTransportCredentials(creds)},
		LocalServer: server2Client,
		logger:      logger,
		mtx:         sync.Mutex{},
		servers:     map[string]chan struct{}{},
		closed:      false,
		close:       make(chan struct{}),
	}

	replicatedServerName := "server1"

	// act
	replicator.Join(replicatedServerName, server1Setup.LogServerAddr)

	// //assert
	time.Sleep(3 * time.Second)
	// replicator.Leave(replicatedServerName)

	// getStream, err := server2Client.GetStream(ctx)
	// require.NoError(t, err)
	// for offset, record := range records {
	// 	getReq := &api.GetRecordRequest{
	// 		Offset: uint64(offset),
	// 	}
	// 	err = getStream.Send(getReq)
	// 	require.NoError(t, err)

	// 	getResp, err := getStream.Recv()
	// 	require.NoError(t, err)
	// 	if getResp.Record.Offset != uint64(offset) {
	// 		t.Fatalf("got offset: %d, want: %d", getResp.Record.Offset, offset)
	// 	}

	// 	givenRecordValue := string(getResp.Record.Value)
	// 	expectedRecordValue := string(record.Value)
	// 	require.EqualValues(t, expectedRecordValue, givenRecordValue)
	// }

	server1Setup.Teardown()
	server1Setup.Teardown()
}

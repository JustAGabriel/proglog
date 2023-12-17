package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"testing"
	"time"

	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal"
	"github.com/justagabriel/proglog/internal/config"
	"github.com/justagabriel/proglog/internal/loadbalance"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestAgent(t *testing.T) {
	host := "localhost" // todo: check why "127.0.0.1" doesn't work
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: host,
		Server:        true,
	})
	require.NoError(t, err)

	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: host,
		Server:        false,
	})
	require.NoError(t, err)

	var agents []*Agent
	for i := 0; i < 3; i++ {
		bindPort := internal.FreePort(t)
		bindAddr := fmt.Sprintf("%s:%d", host, bindPort)
		rpcPort := internal.FreePort(t)

		dataDir := internal.GetTempDir(t, "agent-test-log")

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(startJoinAddrs, agents[0].Config.BindAddr)
		}

		isLeader := i == 0
		agent, err := New(Config{
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			DataDir:         dataDir,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			NodeName:        fmt.Sprintf("%d", i),
			StartJoinAddr:   startJoinAddrs,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
			Bootstrap:       isLeader,
		})
		require.NoError(t, err)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir))
		}
	}()

	time.Sleep(3 * time.Second)

	leaderClient := client(t, agents[0], peerTLSConfig)
	createReq := api.CreateRecordRequest{
		Record: &api.Record{
			Value: []byte("foo"),
		},
	}
	createResp, err := leaderClient.Create(context.Background(), &createReq)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	getReq := api.GetRecordRequest{
		Offset: createResp.Offset,
	}
	getResp, err := leaderClient.Get(context.Background(), &getReq)
	require.NoError(t, err)
	require.Equal(t, getResp.Record.Value, createReq.Record.Value)

	time.Sleep(3 * time.Second)

	followerClient := client(t, agents[1], peerTLSConfig)
	getResp2, err := followerClient.Get(context.Background(), &getReq)
	require.NoError(t, err)
	require.Equal(t, getResp2.Record.Value, createReq.Record.Value)

	getReqOutOfBounds := api.GetRecordRequest{
		Offset: createResp.GetOffset() + 1,
	}
	getRespOutOfBounds, err := leaderClient.Get(context.Background(), &getReqOutOfBounds)
	require.Nil(t, getRespOutOfBounds)
	require.Error(t, err)
	got := status.Code(err)
	want := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	require.Equal(t, got, want)
}

func client(t *testing.T, agent *Agent, tlsConfig *tls.Config) api.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}

	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", loadbalance.Name, rpcAddr), opts...)
	require.NoError(t, err)

	return api.NewLogClient(conn)
}

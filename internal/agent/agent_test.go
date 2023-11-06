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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestAgent(t *testing.T) {
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: "127.0.0.1",
		Server:        true,
	})
	require.NoError(t, err)

	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: "127.0.0.1",
		Server:        false,
	})
	require.NoError(t, err)

	var agents []*Agent
	for i := 0; i < 3; i++ {
		bindPort := internal.FreePort(t)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", bindPort)
		rpcPort := internal.FreePort(t)

		dataDir := internal.GetTempDir(t, "agent-test-log")

		var startJoinAddr []string
		if i != 0 {
			startJoinAddr = append(startJoinAddr, agents[0].Config.BindAddr)
		}

		agent, err := New(Config{
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			DataDir:         dataDir,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			NodeName:        fmt.Sprintf("%d", i),
			StartJoinAddr:   startJoinAddr,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
		})
		require.NoError(t, err)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.DataDir))
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
}

func client(t *testing.T, agent *Agent, tlsConfig *tls.Config) api.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}

	rpcAddr, err := agent.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(rpcAddr, opts...)
	require.NoError(t, err)

	return api.NewLogClient(conn)
}

package loadbalance

import (
	"context"
	"fmt"
	"sync"
	"time"

	api "github.com/justagabriel/proglog/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

const (
	Name string = "proglog"
)

type Resolver struct {
	mu            sync.Mutex
	clientConn    resolver.ClientConn
	resolverConn  *grpc.ClientConn
	serviceConfig *serviceconfig.ParseResult
	logger        *zap.Logger
}

// Close implements resolver.Resolver.
func (r *Resolver) Close() {
	if err := r.resolverConn.Close(); err != nil {
		r.logger.Error("failed to close conn", zap.Error(err))
	}
}

// ResolveNow implements resolver.Resolver.
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()
	client := api.NewLogClient(r.resolverConn)
	ctx := context.Background()

	retries := 3
	for i := 0; i < retries; i++ {
		resp, err := client.GetServers(ctx, &api.GetServersRequest{})
		if err != nil {
			r.logger.Error("failed to resolve server", zap.Error(err))
			return
		}

		var foundLeader bool
		var addrs []resolver.Address
		for _, server := range resp.Servers {
			if server.IsLeader {
				foundLeader = true
			}

			addr := resolver.Address{
				Addr:       server.RpcAddr,
				Attributes: attributes.New("is_leader", server.IsLeader),
			}
			addrs = append(addrs, addr)
		}

		if !foundLeader {
			isLastRetry := (i == (retries - 1))
			if isLastRetry {
				r.logger.Error("no leader found, no more retries")
			} else {
				r.logger.Warn("no leader found")
			}
			time.Sleep(time.Duration(i*300) * time.Millisecond)
			continue
		}

		resolverState := resolver.State{
			Addresses:     addrs,
			ServiceConfig: r.serviceConfig,
		}
		r.clientConn.UpdateState(resolverState)
		break
	}
}

// Scheme implements resolver.Builder.
func (*Resolver) Scheme() string {
	return Name
}

// Build implements resolver.Builder.
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.logger = zap.L().Named("resolver")
	r.clientConn = cc
	var dialOpts []grpc.DialOption
	if opts.DialCreds != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(opts.DialCreds))
	}
	configStr := fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, Name)
	r.serviceConfig = r.clientConn.ParseServiceConfig(configStr)
	var err error
	r.resolverConn, err = grpc.Dial(target.Endpoint(), dialOpts...)
	if err != nil {
		return nil, err
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

var _ resolver.Builder = (*Resolver)(nil)

func init() {
	resolver.Register(&Resolver{})
}

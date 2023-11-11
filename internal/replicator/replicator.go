package replicator

import (
	"context"
	"sync"

	api "github.com/justagabriel/proglog/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Replicator struct {
	DialOptions []grpc.DialOption
	LocalServer api.LogClient
	logger      *zap.Logger
	mtx         sync.Mutex
	servers     map[string]chan struct{}
	closed      bool
	close       chan struct{}
}

func (r *Replicator) Join(name, addr string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.init()

	if r.closed {
		return nil
	}

	if _, ok := r.servers[name]; ok {
		// already replicating so skip
		return nil
	}

	r.servers[name] = make(chan struct{})
	go r.replicate(addr, r.servers[name])
	return nil
}

func (r *Replicator) replicate(addr string, leave chan struct{}) {
	clientConn, err := grpc.Dial(addr, r.DialOptions...)
	if err != nil {
		r.logError(err, "failed to dial", addr)
		return
	}
	defer clientConn.Close()

	client := api.NewLogClient(clientConn)
	ctx := context.Background()

	// stream, err := client.GetStream(ctx)
	// stream.Send(&api.GetRecordRequest{
	// 	Offset: 0,
	// })

	// if err != nil {
	// 	r.logError(err, "failed to consume", addr)
	// 	return
	// }

	records := make(chan *api.Record)
	go func() {
		getReq := api.GetRecordRequest{Offset: 0}

		for {
			resp, err := client.Get(ctx, &getReq)
			if err != nil {
				r.logError(err, "failed to receive", addr)
				return
			}
			r.logger.Sugar().Debugf("received record: %+v", resp.Record)
			records <- resp.Record
			nextOffset := getReq.Offset + 1
			getReq = api.GetRecordRequest{Offset: nextOffset}
		}
	}()

	for {
		select {
		case <-r.close:
			return
		case <-leave:
			return
		case record := <-records:
			_, err = r.LocalServer.Create(ctx, &api.CreateRecordRequest{Record: record})
			if err != nil {
				r.logError(err, "failed to create", addr)
				return
			}
		}
	}
}

func (r *Replicator) Leave(name string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.init()
	if _, ok := r.servers[name]; !ok {
		return nil
	}

	close(r.servers[name])
	delete(r.servers, name)
	return nil
}

func (r *Replicator) init() {
	if r.logger == nil {
		r.logger = zap.L().Named("replicator")
	}

	if r.servers == nil {
		r.servers = make(map[string]chan struct{})
	}

	if r.close == nil {
		r.close = make(chan struct{})
	}
}

func (r *Replicator) Close() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.init()

	if r.closed {
		return nil
	}

	r.closed = true
	close(r.close)
	return nil
}

func (r *Replicator) logError(err error, msg, addr string) {
	r.logger.Error(
		msg,
		zap.String("addr", addr),
		zap.Error(err),
	)
}

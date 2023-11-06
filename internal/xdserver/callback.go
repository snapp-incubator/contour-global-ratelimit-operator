package xdserver

import (
	"context"
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

// Callbacks is a struct that handles xDS callbacks.
type Callbacks struct {
	Signal         chan struct{}
	Debug          bool
	Fetches        int
	Requests       int
	DeltaRequests  int
	DeltaResponses int
	mu             sync.Mutex
}

var _ server.Callbacks = &Callbacks{}

// Report logs callback statistics.
func (cb *Callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	logger.Infof("Server callbacks: fetches=%d requests=%d\n", cb.Fetches, cb.Requests)
}

// OnStreamOpen handles the opening of a stream.
func (cb *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	if cb.Debug {
		logger.Infof("Stream %d opened for %s\n", id, typ)
	}
	return nil
}

// OnStreamClosed handles the closing of a stream.
func (cb *Callbacks) OnStreamClosed(id int64, node *core.Node) {
	if cb.Debug {
		logger.Infof("Stream %d of node %s closed\n", id, node.Id)
	}
}

// OnDeltaStreamOpen handles the opening of a delta stream.
func (cb *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	if cb.Debug {
		logger.Infof("Delta stream %d opened for %s\n", id, typ)
	}
	return nil
}

// OnDeltaStreamClosed handles the closing of a delta stream.
func (cb *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	if cb.Debug {
		logger.Infof("Delta stream %d of node %s closed\n", id, node.Id)
	}
}

// OnStreamRequest handles a stream request.
func (cb *Callbacks) OnStreamRequest(int64, *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Requests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

// OnStreamResponse handles a stream response.
func (cb *Callbacks) OnStreamResponse(context.Context, int64, *discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {
}

// OnStreamDeltaResponse handles a delta stream response.
func (cb *Callbacks) OnStreamDeltaResponse(int64, *discovery.DeltaDiscoveryRequest, *discovery.DeltaDiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.DeltaResponses++
}

// OnStreamDeltaRequest handles a delta stream request.
func (cb *Callbacks) OnStreamDeltaRequest(int64, *discovery.DeltaDiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.DeltaRequests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

// OnFetchRequest handles a fetch request.
func (cb *Callbacks) OnFetchRequest(context.Context, *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Fetches++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

// OnFetchResponse handles a fetch response.
func (cb *Callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}

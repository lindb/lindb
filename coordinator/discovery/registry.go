package discovery

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pathutil"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
)

// Registry represents server node register
type Registry interface {
	// Register registers node info, add it to active node list for discovery
	Register(node models.Node) error
	// Deregister deregister node info, remove it from active list
	Deregister(node models.Node) error
	// Close closes registry, releases resources
	Close() error
}

// registry implements registry interface for server node register with prefix
type registry struct {
	prefix string
	ttl    int64
	repo   state.Repository

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewRegistry returns a new registry with prefix and ttl
func NewRegistry(repo state.Repository, prefix string, ttl int64) Registry {
	ctx, cancel := context.WithCancel(context.Background())
	return &registry{
		prefix: prefix,
		ttl:    ttl,
		repo:   repo,
		ctx:    ctx,
		cancel: cancel,
		log:    logger.GetLogger("coordinator/registry"),
	}
}

// Register registers node info, add it to active node list for discovery
func (r *registry) Register(node models.Node) error {
	// register node info
	path := pathutil.GetNodePath(r.prefix, node.Indicator())
	// register node if fail retry it
	go r.register(path, node)
	return nil
}

// Deregister deregisters node info, remove it from active list
func (r *registry) Deregister(node models.Node) error {
	return r.repo.Delete(r.ctx, pathutil.GetNodePath(r.prefix, node.Indicator()))
}

// Close closes registry, releases resources
func (r *registry) Close() error {
	r.cancel()
	return nil
}

// register registers node info, if fail do retry
func (r *registry) register(path string, node models.Node) {
	for {
		// if ctx happen err, exit register loop
		if r.ctx.Err() != nil {
			return
		}
		nodeBytes, _ := json.Marshal(&models.ActiveNode{OnlineTime: timeutil.Now(), Node: node})

		closed, err := r.repo.Heartbeat(r.ctx, path, nodeBytes, r.ttl)
		if err != nil {
			r.log.Error("register node error", logger.Error(err))
			time.Sleep(500 * time.Millisecond)
			continue
		}

		r.log.Info("register node successfully", logger.String("path", path))

		select {
		case <-r.ctx.Done():
			r.log.Warn("context is canceled, exit register loop")
			return
		case <-closed:
			r.log.Warn("the heartbeat channel is closed, retry register")
		}
	}
}

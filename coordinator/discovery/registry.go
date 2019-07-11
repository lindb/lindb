package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"

	"go.uber.org/zap"
)

// Registry represents server node register
type Registry interface {
	// Register registers node info, add it to acitve node list for discovery
	Register(node models.Node) error
	// Deregister deregisters node info, remove it from active list
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

	log *zap.Logger
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
		log:    logger.GetLogger(),
	}
}

// Register registers node info, add it to acitve node list for discovery
func (r *registry) Register(node models.Node) error {
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		r.log.Error("convert node to byte error when register node info", zap.Error(err))
		return err
	}
	// register node info
	path := r.nodePath(node)
	// register node if fail retry it
	go r.register(path, nodeBytes)
	return nil
}

// Deregister deregisters node info, remove it from active list
func (r *registry) Deregister(node models.Node) error {
	return r.repo.Delete(r.ctx, r.nodePath(node))
}

// Close closes registry, releases resources
func (r *registry) Close() error {
	r.cancel()
	return nil
}

// register registers node unfo, if fail do retry
func (r *registry) register(path string, node []byte) {
	for {
		// if ctx happen err, exit register loop
		if r.ctx.Err() != nil {
			return
		}
		closed, err := r.repo.Heartbeat(r.ctx, path, node, r.ttl)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		select {
		case <-r.ctx.Done():
			r.log.Warn("context is canceled, exit register loop")
			return
		case <-closed:
			r.log.Warn("the heartbeat channel is closed, retry register")
		}
	}
}

// nodePath retruns node register path
func (r *registry) nodePath(node models.Node) string {
	return fmt.Sprintf("%s/%s", r.prefix, node.String())
}

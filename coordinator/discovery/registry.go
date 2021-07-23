// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package discovery

import (
	"context"
	"io"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./registry.go -destination=./registry_mock.go -package=discovery

// Registry represents server node register.
type Registry interface {
	io.Closer
	// Register registers node info, add it to active node list for discovery.
	Register(node models.Node) error
	// Deregister deregister node info, remove it from active list.
	Deregister(node models.Node) error
}

// registry implements registry interface for server node register with prefix.
type registry struct {
	ttl  time.Duration
	repo state.Repository

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewRegistry returns a new registry with prefix and ttl.
func NewRegistry(repo state.Repository, ttl time.Duration) Registry {
	ctx, cancel := context.WithCancel(context.Background())
	return &registry{
		ttl:    ttl,
		repo:   repo,
		ctx:    ctx,
		cancel: cancel,
		log:    logger.GetLogger("coordinator", "Registry"),
	}
}

// Register registers node info, add it to active node list for discovery.
func (r *registry) Register(node models.Node) error {
	// register node info
	path := constants.GetLiveNodePath(node.Indicator())
	r.log.Info("starting register live node", logger.String("path", path))
	// register node if fail retry it
	go r.register(path, node)
	return nil
}

// Deregister deregisters node info, remove it from active list.
func (r *registry) Deregister(node models.Node) error {
	return r.repo.Delete(r.ctx, constants.GetLiveNodePath(node.Indicator()))
}

// Close closes registry, releases resources.
func (r *registry) Close() error {
	r.cancel()
	return nil
}

// register registers node info, if fail do retry.
func (r *registry) register(path string, node models.Node) {
	for {
		// if ctx happen err, exit register loop
		if r.ctx.Err() != nil {
			return
		}
		nodeBytes := encoding.JSONMarshal(node)

		closed, err := r.repo.Heartbeat(r.ctx, path, nodeBytes, int64(r.ttl.Seconds()))
		if err != nil {
			r.log.Error("register node error", logger.Error(err))
			time.Sleep(500 * time.Millisecond)
			continue
		}

		r.log.Info("register node successfully", logger.String("path", path), logger.Any("node", node))

		select {
		case <-r.ctx.Done():
			r.log.Warn("context is canceled, exit register loop")
			return
		case <-closed:
			r.log.Warn("the heartbeat channel is closed, retry register")
		}
	}
}

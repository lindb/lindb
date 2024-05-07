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
	"runtime/pprof"
	"time"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./registry.go -destination=./registry_mock.go -package=discovery

// Registry represents server node register.
type Registry interface {
	io.Closer
	// Register registers node info, add it to active node list for discovery.
	Register() error
	// Deregister deregister node info, remove it from active list.
	Deregister() error
	// IsSuccess returns if registry successfully.
	IsSuccess() bool
}

// registry implements registry interface for server node register with prefix.
type registry struct {
	node            models.Node
	repo            state.Repository
	ctx             context.Context
	log             logger.Logger
	cancel          context.CancelFunc
	path            string
	ttl             time.Duration
	registrySuccess atomic.Bool
}

// NewRegistry returns a new registry with path, node and ttl.
func NewRegistry(repo state.Repository, path string, node models.Node, ttl time.Duration) Registry {
	ctx, cancel := context.WithCancel(context.Background())
	return &registry{
		path:            path,
		node:            node,
		ttl:             ttl,
		repo:            repo,
		ctx:             ctx,
		cancel:          cancel,
		registrySuccess: *atomic.NewBool(false),
		log:             logger.GetLogger("Coordinator", "Registry"),
	}
}

// Register registers node info, add it to active node list for discovery.
func (r *registry) Register() error {
	// register node info
	r.log.Info("starting register node", logger.String("path", r.path))
	// register node if fail retry it
	go func() {
		registerLabels := pprof.Labels("path", r.path,
			"timestamp", timeutil.FormatTimestamp(timeutil.Now(), timeutil.DataTimeFormat2))
		pprof.Do(r.ctx, registerLabels, func(_ context.Context) {
			r.register()
		})
	}()
	return nil
}

// Deregister deregisters node info, remove it from active list.
func (r *registry) Deregister() error {
	return r.repo.Delete(r.ctx, r.path)
}

// IsSuccess returns if registry successfully.
func (r *registry) IsSuccess() bool {
	return r.registrySuccess.Load()
}

// Close closes registry, releases resources.
func (r *registry) Close() error {
	r.cancel()
	return nil
}

// register registers node info, if fail do retry.
func (r *registry) register() {
	for {
		// if ctx happen err, exit register loop
		if r.ctx.Err() != nil {
			return
		}
		r.registrySuccess.Store(false)
		r.node.Online() // reset online timestamp
		nodeBytes := encoding.JSONMarshal(r.node)

		closed, err := r.repo.Heartbeat(r.ctx, r.path, nodeBytes, int64(r.ttl.Seconds()))
		if err != nil {
			r.log.Error("register node error", logger.String("path", r.path), logger.Error(err))
			time.Sleep(500 * time.Millisecond)
			continue
		}

		r.registrySuccess.Store(true)
		r.log.Info("register node successfully", logger.String("path", r.path), logger.String("node", string(nodeBytes)))

		select {
		case <-r.ctx.Done():
			r.log.Warn("context is canceled, exit register loop", logger.String("path", r.path))
			return
		case <-closed:
			r.log.Warn("the heartbeat channel is closed, retry register", logger.String("path", r.path))
		}
	}
}

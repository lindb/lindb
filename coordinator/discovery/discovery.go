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
	"fmt"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./discovery.go -destination=./discovery_mock.go -package=discovery

// Factory represents a discovery create factory.
type Factory interface {
	// GetRepo returns the repo of discovery used.
	GetRepo() state.Repository
	// CreateDiscovery creates a discovery who will watch the changes with the given prefix.
	CreateDiscovery(prefix string, listener Listener) Discovery
}

// factory implements factory interface using state repo.
type factory struct {
	repo state.Repository
}

// NewFactory creates a factory.
func NewFactory(repo state.Repository) Factory {
	return &factory{repo: repo}
}

// GetRepo returns the repo of discovery used.
func (f *factory) GetRepo() state.Repository {
	return f.repo
}

// CreateDiscovery creates a discovery who will watch the changes with the given prefix.
func (f *factory) CreateDiscovery(prefix string, listener Listener) Discovery {
	ctx, cancel := context.WithCancel(context.Background())
	r := &discovery{
		prefix:   prefix,
		repo:     f.repo,
		ctx:      ctx,
		cancel:   cancel,
		listener: listener,
		logger:   logger.GetLogger("coordinator", "Discovery"),
	}

	r.logger.Info("create new discovery", logger.String("watch", prefix))
	return r
}

// Discovery represents discovery resources, through watch resource's prefix.
type Discovery interface {
	// Discovery starts discovery resources change, includes create/delete/clean,
	// if init = true, need list all resources before watch path.
	Discovery(init bool) error
	// Close stops watch, trigger all resource cleanup callback.
	Close()
}

// discovery implements discovery interface.
type discovery struct {
	prefix   string
	repo     state.Repository
	listener Listener

	ctx    context.Context
	cancel context.CancelFunc

	logger *logger.Logger
}

// Discovery starts discovery resources change, includes create/delete/clean.
func (d *discovery) Discovery(init bool) error {
	if d.prefix == "" {
		return fmt.Errorf("watch prefix is empth for discovery resource")
	}

	if init {
		kvs, err := d.repo.List(d.ctx, d.prefix)
		if err != nil {
			return fmt.Errorf("list resource error:%s", err)
		}

		// init exist resource.
		for _, kv := range kvs {
			d.listener.OnCreate(kv.Key, kv.Value)
		}
	}

	watchEventCh := d.repo.WatchPrefix(d.ctx, d.prefix, false)
	go func() {
		d.handlerResourceChange(watchEventCh)
		d.logger.Warn("exit discovery loop", logger.String("prefix", d.prefix))
	}()
	return nil
}

// Close cleans all resources
func (d *discovery) Close() {
	d.cancel()
}

// handlerResourceChange handles the changes of event for resources.
func (d *discovery) handlerResourceChange(eventCh state.WatchEventChan) {
	for event := range eventCh {
		if event.Err != nil {
			continue
		}
		switch event.Type {
		case state.EventTypeDelete:
			for _, kv := range event.KeyValues {
				d.listener.OnDelete(kv.Key)
			}
		case state.EventTypeModify:
			for _, kv := range event.KeyValues {
				d.listener.OnCreate(kv.Key, kv.Value)
			}
		}
	}
}

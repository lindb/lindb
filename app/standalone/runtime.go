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

package standalone

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/lindb/common/pkg/logger"
	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap/zapcore"

	"github.com/lindb/lindb/app/broker"
	"github.com/lindb/lindb/app/storage"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/state"
)

// for testing
var (
	startEtcdFn = embed.StartEtcd
)

var log = logger.GetLogger("Standalone", "Runtime")

// runtime represents the runtime dependency of standalone mode
type runtime struct {
	version   string
	embedEtcd bool

	state       server.State
	repoFactory state.RepositoryFactory
	cfg         *config.Standalone
	etcd        *embed.Etcd
	broker      server.Service
	storage     server.Service

	initializer bootstrap.ClusterInitializer
	delayInit   time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	pusher monitoring.NativePusher
}

// NewStandaloneRuntime creates the runtime
func NewStandaloneRuntime(version string, cfg *config.Standalone, embedEtcd bool) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		version:     version,
		embedEtcd:   embedEtcd,
		state:       server.New,
		delayInit:   5 * time.Second,
		repoFactory: state.NewRepositoryFactory("standalone"),
		broker: broker.NewBrokerRuntime(version,
			&config.Broker{
				Query:       cfg.Query,
				Coordinator: cfg.Coordinator,
				BrokerBase:  cfg.BrokerBase,
				Monitor:     cfg.Monitor,
				Logging:     cfg.Logging,
				Prometheus:  cfg.Prometheus,
			}, true),
		storage: storage.NewStorageRuntime(version,
			1, // default: 1
			&config.Storage{
				Query:       cfg.Query,
				Coordinator: cfg.Coordinator,
				StorageBase: cfg.StorageBase,
				Monitor:     cfg.Monitor,
				Logging:     cfg.Logging,
			}),
		cfg:         cfg,
		initializer: bootstrap.NewClusterInitializer(fmt.Sprintf("http://localhost:%d", cfg.BrokerBase.HTTP.Port)),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Config returns the configure of standalone.
func (r *runtime) Config() any {
	return r.cfg
}

// Name returns the cluster mode
func (r *runtime) Name() string {
	return "standalone"
}

// Run runs the cluster as standalone mode
func (r *runtime) Run() error {
	config.StandaloneMode = true

	if r.embedEtcd {
		if err := r.startETCD(); err != nil {
			log.Error("failed to start ETCD", logger.Error(err))
			r.state = server.Failed
			return err
		}
	}

	// cleanup state for previous embed etcd server state
	if err := r.cleanupState(); err != nil {
		r.state = server.Failed
		return err
	}
	if err := r.runServer(); err != nil {
		return err
	}

	r.state = server.Running

	time.AfterFunc(r.delayInit, func() {
		if err := r.initializer.InitInternalDatabase("create database _internal engine metric"); err != nil {
			log.Error("init _internal database with error", logger.Error(err))
		} else {
			log.Info("initialized _internal database successfully")
		}
	})

	return nil
}

func (r *runtime) runServer() error {
	if err := r.storage.Run(); err != nil {
		r.state = server.Failed
		return err
	}
	if err := r.broker.Run(); err != nil {
		r.state = server.Failed
		return err
	}
	return nil
}

// State returns the state of cluster
func (r *runtime) State() server.State {
	return r.state
}

// Stop stops the cluster
func (r *runtime) Stop() {
	defer r.cancel()
	if r.broker != nil {
		r.broker.Stop()
	}
	if r.storage != nil {
		r.storage.Stop()
	}
	if r.etcd != nil {
		r.etcd.Close()
		log.Info("stopped etcd server")
	}
	if r.pusher != nil {
		r.pusher.Stop()
		log.Info("stopped native linmetric pusher successfully")
	}
	r.state = server.Terminated
}

// startETCD starts embed etcd server
func (r *runtime) startETCD() error {
	cfg := embed.NewConfig()
	lcurl, _ := url.Parse(r.cfg.ETCD.URL)
	cfg.ListenClientUrls = []url.URL{*lcurl}
	cfg.Dir = r.cfg.ETCD.Dir
	// always set etcd runtime to error level
	cfg.LogLevel = zapcore.ErrorLevel.String()

	e, err := startEtcdFn(cfg)
	if err != nil {
		return err
	}
	r.etcd = e
	select {
	case <-e.Server.ReadyNotify():
		log.Info("etcd server is ready")
	case <-time.After(time.Minute):
		e.Server.Stop() // trigger a shutdown
		log.Error("etcd server took too long to start")
	case err := <-e.Err():
		log.Error("etcd server error", logger.Error(err))
	}
	return nil
}

// cleanupState cleans the state of previous standalone process.
// 1. master node in etcd, because etcd will trigger master node expire event
// 2. stateful node in etcd
func (r *runtime) cleanupState() error {
	repo, err := r.repoFactory.CreateNormalRepo(&r.cfg.Coordinator)
	if err != nil {
		return fmt.Errorf("start broker state repo error:%s", err)
	}
	defer func() {
		err = repo.Close()
		if err != nil {
			log.Error("close broker state repo when do cleanup", logger.Error(err))
		}
	}()
	err = repo.Delete(context.TODO(), constants.MasterPath)
	if err != nil {
		return fmt.Errorf("delete old master error")
	}

	kvs, err := repo.List(context.TODO(), constants.StorageLiveNodesPath)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := repo.Delete(context.TODO(), kv.Key); err != nil {
			return fmt.Errorf("delete stateful node info error")
		}
	}
	return nil
}

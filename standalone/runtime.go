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

	"go.etcd.io/etcd/embed"
	"go.uber.org/zap/zapcore"

	"github.com/lindb/lindb/broker"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/storage"
)

var log = logger.GetLogger("standalone", "Runtime")

const storageClusterName = "standalone"

// runtime represents the runtime dependency of standalone mode
type runtime struct {
	version     string
	state       server.State
	repoFactory state.RepositoryFactory
	cfg         *config.Standalone
	etcd        *embed.Etcd
	broker      server.Service
	storage     server.Service

	initializer *bootstrap.ClusterInitializer
	delayInit   time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStandaloneRuntime creates the runtime
func NewStandaloneRuntime(version string, cfg *config.Standalone) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		version:     version,
		state:       server.New,
		delayInit:   5 * time.Second,
		repoFactory: state.NewRepositoryFactory("standalone"),
		broker: broker.NewBrokerRuntime(version,
			&config.Broker{
				BrokerBase: cfg.BrokerBase,
				Monitor:    cfg.Monitor,
			}),
		storage: storage.NewStorageRuntime(version,
			&config.Storage{
				StorageBase: cfg.StorageBase,
				Monitor:     cfg.Monitor,
			}),
		cfg:         cfg,
		initializer: bootstrap.NewClusterInitializer(fmt.Sprintf("http://localhost:%d", cfg.BrokerBase.HTTP.Port)),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Name returns the cluster mode
func (r *runtime) Name() string {
	return "standalone"
}

// Run runs the cluster as standalone mode
func (r *runtime) Run() error {
	config.StandaloneMode = true

	if err := r.startETCD(); err != nil {
		log.Error("failed to start ETCD", logger.Error(err))
		r.state = server.Failed
		return err
	}

	// cleanup state for previous embed etcd server state
	if err := r.cleanupState(); err != nil {
		return err
	}
	if err := r.runServer(); err != nil {
		return err
	}
	r.state = server.Running

	time.AfterFunc(r.delayInit, func() {
		log.Info("initializing standalone internal database")
		if err := r.initializer.InitStorageCluster(config.StorageCluster{
			Name:   "standalone",
			Config: r.cfg.StorageBase.Coordinator}); err != nil {
			log.Error("initialized standalone storage cluster with error", logger.Error(err))
		} else {
			log.Info("initialized standalone storage cluster successfully")
		}

		if err := r.initializer.InitInternalDatabase(models.Database{
			Name:          "_internal",
			Cluster:       storageClusterName,
			NumOfShard:    1,
			ReplicaFactor: 1,
			Option: option.DatabaseOption{
				Interval: "10s",
			},
		}); err != nil {
			log.Error("init _internal database with error", logger.Error(err))
		} else {
			log.Info("initialized _internal database successfully")
		}
	})

	return nil
}

func (r *runtime) runServer() error {
	// start storage server
	if err := r.storage.Run(); err != nil {
		r.state = server.Failed
		return err
	}
	// start broker server
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
	r.state = server.Terminated
}

// startETCD starts embed etcd server
func (r *runtime) startETCD() error {
	cfg := embed.NewConfig()
	lcurl, _ := url.Parse(r.cfg.ETCD.URL)
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.Dir = r.cfg.ETCD.Dir
	// always set etcd runtime to error level
	cfg.LogLevel = zapcore.ErrorLevel.String()

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return err
	}
	r.etcd = e
	select {
	case <-e.Server.ReadyNotify():
		log.Info("etcd server is ready")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Error("etcd server took too long to start")
	case err := <-e.Err():
		log.Error("etcd server error", logger.Error(err))
	}
	return nil
}

// cleanupState cleans the state of previous standalone process.
// 1. master node in etcd, because etcd will trigger master node expire event
func (r *runtime) cleanupState() error {
	repo, err := r.repoFactory.CreateRepo(r.cfg.BrokerBase.Coordinator)
	if err != nil {
		return fmt.Errorf("start broker state repo error:%s", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Error("close state repo when do cleanup", logger.Error(err))
		}
	}()
	if err := repo.Delete(context.TODO(), constants.MasterPath); err != nil {
		return fmt.Errorf("delete old master error")
	}
	return nil
}

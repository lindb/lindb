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

	"github.com/lindb/lindb/app/broker"
	"github.com/lindb/lindb/app/storage"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/series/tag"
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
	pusher monitoring.NativePusher
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
				Query:       cfg.Query,
				Coordinator: cfg.Coordinator,
				BrokerBase:  cfg.BrokerBase,
				Monitor:     config.Monitor{}, // empty to disable broker monitor
			}),
		storage: storage.NewStorageRuntime(version,
			&config.Storage{
				Query:       cfg.Query,
				Coordinator: cfg.Coordinator,
				StorageBase: cfg.StorageBase,
				Monitor:     config.Monitor{}, // empty to disable storage monitor
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

	// start a standalone pusher
	r.nativePusher()
	r.state = server.Running

	time.AfterFunc(r.delayInit, func() {
		log.Info("initializing standalone internal database")
		if err := r.initializer.InitStorageCluster(config.StorageCluster{
			Name:   "standalone",
			Config: r.cfg.Coordinator}); err != nil {
			log.Error("initialized standalone storage cluster with error", logger.Error(err))
		} else {
			log.Info("initialized standalone storage cluster successfully")
		}

		if err := r.initializer.InitInternalDatabase(models.Database{
			Name:          "_internal",
			Storage:       storageClusterName,
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
// 2. stateful node in etcd
func (r *runtime) cleanupState() error {
	brokerRepo, err := r.repoFactory.CreateBrokerRepo(r.cfg.Coordinator)
	if err != nil {
		return fmt.Errorf("start broker state repo error:%s", err)
	}
	defer func() {
		if err := brokerRepo.Close(); err != nil {
			log.Error("close broker state repo when do cleanup", logger.Error(err))
		}
	}()
	if err := brokerRepo.Delete(context.TODO(), constants.MasterPath); err != nil {
		return fmt.Errorf("delete old master error")
	}

	storageRepo, err := r.repoFactory.CreateStorageRepo(r.cfg.Coordinator)
	if err != nil {
		return fmt.Errorf("start storage state repo error:%s", err)
	}
	defer func() {
		if err := storageRepo.Close(); err != nil {
			log.Error("close storage state repo when do cleanup", logger.Error(err))
		}
	}()
	kvs, err := storageRepo.List(context.TODO(), constants.LiveNodesPath)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := storageRepo.Delete(context.TODO(), kv.Key); err != nil {
			return fmt.Errorf("delete stateful node info error")
		}
	}
	return nil
}

func (r *runtime) nativePusher() {
	log.Info("disable pusher of both broker and storage")
	monitorEnabled := r.cfg.Monitor.ReportInterval > 0
	if !monitorEnabled {
		log.Info("pusher won't start because report-interval is 0")
		return
	}
	log.Info("pusher is running",
		logger.String("interval", r.cfg.Monitor.ReportInterval.String()))

	ip, _ := hostutil.GetHostIP()

	r.pusher = monitoring.NewNativeProtoPusher(
		r.ctx,
		r.cfg.Monitor.URL,
		r.cfg.Monitor.ReportInterval.Duration(),
		r.cfg.Monitor.PushTimeout.Duration(),
		tag.KeyValues{
			{Key: "node", Value: ip},
			{Key: "role", Value: "standalone"},
		},
	)
	go r.pusher.Start()
}

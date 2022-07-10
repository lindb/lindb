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

package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/tsdb"
)

var cfg = config.Storage{
	Coordinator: config.RepoState{
		Namespace: "/test/2222",
	},
	StorageBase: config.StorageBase{
		WAL:       config.WAL{RemoveTaskInterval: ltoml.Duration(time.Minute)},
		Indicator: 1,
		GRPC: config.GRPC{
			Port: 7777,
		},
		HTTP: config.HTTP{
			Port: 8888,
		},
		TSDB: config.TSDB{Dir: "/tmp/test/data"},
	}, Monitor: *config.NewDefaultMonitor(),
}

func TestStorageRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	cluster := mock.StartEtcdCluster(t, "http://localhost:8100")
	defer func() {
		newDatabaseLifecycleFn = NewDatabaseLifecycle
		cluster.Terminate(t)
		ctrl.Finish()
	}()

	dbLifecycle := NewMockDatabaseLifecycle(ctrl)
	dbLifecycle.EXPECT().Startup()
	dbLifecycle.EXPECT().Shutdown()
	newDatabaseLifecycleFn = func(ctx context.Context, repo state.Repository,
		walMgr replica.WriteAheadLogManager, engine tsdb.Engine) DatabaseLifecycle {
		return dbLifecycle
	}

	// test normal storage run
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.Coordinator.Timeout = ltoml.Duration(time.Second * 10)
	cfg.StorageBase.GRPC.Port = 9997
	config.SetGlobalStorageConfig(&cfg.StorageBase)
	storage := NewStorageRuntime("test-version", &cfg)
	err := storage.Run()
	assert.NoError(t, err)
	assert.Equal(t, server.Running, storage.State())
	// wait register success
	time.Sleep(500 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := constants.GetLiveNodePath(strconv.Itoa(int(runtime.node.ID)))
	nodeBytes, err := runtime.repo.Get(context.TODO(), nodePath)
	assert.NoError(t, err)

	nodeInfo := models.StatefulNode{}
	_ = encoding.JSONUnmarshal(nodeBytes, &nodeInfo)

	assert.Equal(t, *runtime.node, nodeInfo)
	assert.Equal(t, "storage", storage.Name())

	storage.Stop()
	assert.Equal(t, server.Terminated, storage.State())
	time.Sleep(500 * time.Millisecond)
}

func TestStorageRun_GetHost_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	cluster := mock.StartEtcdCluster(t, "http://localhost:8101")
	defer cluster.Terminate(t)

	defer func() {
		getHostIP = hostutil.GetHostIP
		newDatabaseLifecycleFn = NewDatabaseLifecycle
		hostName = os.Hostname
		ctrl.Finish()
	}()
	dbLifecycle := NewMockDatabaseLifecycle(ctrl)
	dbLifecycle.EXPECT().Startup()
	dbLifecycle.EXPECT().Shutdown()
	newDatabaseLifecycleFn = func(ctx context.Context, repo state.Repository,
		walMgr replica.WriteAheadLogManager, engine tsdb.Engine) DatabaseLifecycle {
		return dbLifecycle
	}
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.StorageBase.GRPC.Port = 8889
	cfg.StorageBase.Indicator = 2
	config.SetGlobalStorageConfig(&cfg.StorageBase)
	storage := NewStorageRuntime("test-version", &cfg)
	getHostIP = func() (string, error) {
		return "test-ip", fmt.Errorf("err")
	}
	err := storage.Run()
	assert.Error(t, err)

	getHostIP = func() (string, error) {
		return "ip", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	cfg.StorageBase.GRPC.Port = 8887
	cfg.StorageBase.Indicator = 3

	storage = NewStorageRuntime("test-version", &cfg)
	err = storage.Run()
	assert.NoError(t, err)
	// wait grpc server start and register success
	time.Sleep(500 * time.Millisecond)
	storage.Stop()
	assert.NoError(t, err)
}

func TestStorageRun_Err(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8102")
	defer cluster.Terminate(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg.StorageBase.GRPC.Port = 8889
	cfg.StorageBase.Indicator = 0
	config.SetGlobalStorageConfig(&cfg.StorageBase)
	storage := NewStorageRuntime("test-version", &cfg)
	err := storage.Run()
	assert.Error(t, err)

	cfg.StorageBase.GRPC.Port = 8886
	cfg.StorageBase.Indicator = 4
	storage = NewStorageRuntime("test-version", &cfg)
	s := storage.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	s.repoFactory = repoFactory
	repoFactory.EXPECT().CreateStorageRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = s.Run()
	assert.Error(t, err)
	// wait grpc server start and register success
	time.Sleep(500 * time.Millisecond)

	repo := state.NewMockRepository(ctrl)
	s.repo = repo

	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	repo.EXPECT().Close().Return(fmt.Errorf("err"))

	s.Stop()
	assert.Error(t, err)
}

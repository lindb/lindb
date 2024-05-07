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
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/ltoml"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	storagepkg "github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

var cfg = config.Storage{
	Coordinator: config.RepoState{
		Namespace: "/test/2222",
	},
	StorageBase: config.StorageBase{
		WAL: config.WAL{RemoveTaskInterval: ltoml.Duration(time.Minute)},
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
		walMgr replica.WriteAheadLogManager, engine tsdb.Engine,
	) DatabaseLifecycle {
		return dbLifecycle
	}

	// test normal storage run
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.Coordinator.Timeout = ltoml.Duration(time.Second * 10)
	cfg.StorageBase.GRPC.Port = 9997
	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "1")
	config.SetGlobalStorageConfig(&cfg.StorageBase)
	storage := NewStorageRuntime("test-version", 1, &cfg)
	err := storage.Run()
	assert.NoError(t, err)
	assert.NotNil(t, storage.Config())
	assert.Equal(t, server.Running, storage.State())
	// wait register success
	time.Sleep(500 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := constants.GetStorageLiveNodePath(strconv.Itoa(int(runtime.node.ID)))
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
		walMgr replica.WriteAheadLogManager, engine tsdb.Engine,
	) DatabaseLifecycle {
		return dbLifecycle
	}
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.StorageBase.GRPC.Port = 8889
	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "2")
	storage := NewStorageRuntime("test-version", 2, &cfg)
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

	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "3")
	storage = NewStorageRuntime("test-version", 3, &cfg)
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
	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "0")
	config.SetGlobalStorageConfig(&cfg.StorageBase)
	storage := NewStorageRuntime("test-version", 0, &cfg)
	err := storage.Run()
	assert.Error(t, err)

	cfg.StorageBase.GRPC.Port = 8886
	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "4")
	storage = NewStorageRuntime("test-version", 4, &cfg)
	s := storage.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	s.repoFactory = repoFactory
	repoFactory.EXPECT().CreateNormalRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
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

	// create engine failure
	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "6")
	storage = NewStorageRuntime("test-version", 6, &cfg)
	defer func() {
		newEngineFn = tsdb.NewEngine
		newWriteAheadLogManagerFn = replica.NewWriteAheadLogManager
	}()
	newEngineFn = func() (tsdb.Engine, error) {
		return nil, fmt.Errorf("err")
	}
	err = storage.Run()
	assert.Error(t, err)

	// wal recovery failure
	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	newWriteAheadLogManagerFn = func(_ context.Context, _ config.WAL,
		_ models.NodeID, _ tsdb.Engine, _ rpc.ClientStreamFactory,
		_ storagepkg.StateManager,
	) replica.WriteAheadLogManager {
		return walMgr
	}
	walMgr.EXPECT().Recovery().Return(fmt.Errorf("err"))
	cfg.StorageBase.TSDB.Dir = filepath.Join(t.TempDir(), "7")
	storage = NewStorageRuntime("test-version", 7, &cfg)
	newEngineFn = func() (tsdb.Engine, error) {
		return nil, nil
	}
	err = storage.Run()
	assert.Error(t, err)
}

func TestStorage_MyID(t *testing.T) {
	defer func() {
		existFn = fileutil.Exist
		atoiFn = strconv.Atoi
		mkDirIfNotExistFn = fileutil.MkDirIfNotExist
		readFileFn = os.ReadFile
		writeFileFn = os.WriteFile
	}()
	r := &runtime{myID: 99}
	err0 := fmt.Errorf("err")
	_, err1 := strconv.Atoi("abc")
	testCases := []struct {
		err     error
		prepare func()
		desc    string
		id      int
	}{
		{
			desc: "mk parent path failure",
			prepare: func() {
				mkDirIfNotExistFn = func(_ string) error {
					return err0
				}
			},
			err: err0,
			id:  0,
		},
		{
			desc: "write data failure",
			prepare: func() {
				mkDirIfNotExistFn = func(_ string) error {
					return nil
				}
				existFn = func(_ string) bool {
					return false
				}
				writeFileFn = func(_ string, _ []byte, _ os.FileMode) error {
					return err0
				}
			},
			err: err0,
			id:  0,
		},
		{
			desc: "write data successfully",
			prepare: func() {
				mkDirIfNotExistFn = func(_ string) error {
					return nil
				}
				existFn = func(_ string) bool {
					return false
				}
				writeFileFn = func(_ string, _ []byte, _ os.FileMode) error {
					return nil
				}
			},
			id: 99,
		},
		{
			desc: "read file failure",
			prepare: func() {
				mkDirIfNotExistFn = func(_ string) error {
					return nil
				}
				existFn = func(_ string) bool {
					return true
				}
				readFileFn = func(_ string) ([]byte, error) {
					return nil, err0
				}
			},
			err: err0,
			id:  0,
		},
		{
			desc: "read wrong value",
			prepare: func() {
				mkDirIfNotExistFn = func(_ string) error {
					return nil
				}
				existFn = func(_ string) bool {
					return true
				}
				readFileFn = func(_ string) ([]byte, error) {
					return []byte("abc"), nil
				}
			},
			err: err1,
			id:  0,
		},
		{
			desc: "read value successfully",
			prepare: func() {
				mkDirIfNotExistFn = func(_ string) error {
					return nil
				}
				existFn = func(_ string) bool {
					return true
				}
				readFileFn = func(_ string) ([]byte, error) {
					return []byte("100"), nil
				}
			},
			id: 100,
		},
	}
	for _, tC := range testCases {
		// reset test func
		existFn = fileutil.Exist
		atoiFn = strconv.Atoi
		mkDirIfNotExistFn = fileutil.MkDirIfNotExist
		readFileFn = os.ReadFile
		writeFileFn = os.WriteFile

		t.Run(tC.desc, func(t *testing.T) {
			tC.prepare()
			id, err := r.initMyID()
			assert.Equal(t, tC.id, id)
			assert.Equal(t, tC.err, err)
		})
	}
}

func TestStorage_Run_With_Wrong_MyID(t *testing.T) {
	defer func() {
		mkDirIfNotExistFn = fileutil.MkDirIfNotExist
	}()
	mkDirIfNotExistFn = func(_ string) error {
		return fmt.Errorf("err")
	}
	r := &runtime{}
	err := r.Run()
	assert.Error(t, err)
	assert.Equal(t, server.Failed, r.State())
}

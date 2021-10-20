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
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

func newDefaultStandaloneConfig(t *testing.T) config.Standalone {
	saCfg := config.Standalone{
		Query:       *config.NewDefaultQuery(),
		Coordinator: *config.NewDefaultCoordinator(),
		StorageBase: *config.NewDefaultStorageBase(),
		BrokerBase:  *config.NewDefaultBrokerBase(),
		Logging:     *config.NewDefaultLogging(),
		ETCD:        *config.NewDefaultETCD(),
		Monitor:     *config.NewDefaultMonitor(),
	}
	saCfg.StorageBase.TSDB.Dir = t.TempDir()
	saCfg.StorageBase.GRPC.Port = 3901
	saCfg.StorageBase.Indicator = 1
	saCfg.StorageBase.HTTPPort = 3902
	saCfg.StorageBase.WAL.RemoveTaskInterval = ltoml.Duration(time.Minute)
	return saCfg
}

func TestRuntime_Run(t *testing.T) {
	cfg := newDefaultStandaloneConfig(t)
	cfg.StorageBase.GRPC.Port = 3901
	standalone := NewStandaloneRuntime("test-version", &cfg)
	s := standalone.(*runtime)
	s.delayInit = 100 * time.Millisecond

	err := standalone.Run()
	assert.Equal(t, server.Running, standalone.State())
	assert.NoError(t, err)

	standalone.Stop()
	assert.Equal(t, server.Terminated, standalone.State())
	assert.Equal(t, "standalone", standalone.Name())
	time.Sleep(500 * time.Millisecond)
}

func TestRuntime_RunWithoutPusher(t *testing.T) {
	cfg := newDefaultStandaloneConfig(t)

	cfg.StorageBase.GRPC.Port = 3901
	cfg.Monitor.ReportInterval = ltoml.Duration(0)
	standalone := NewStandaloneRuntime("test-version", &cfg)
	s := standalone.(*runtime)
	s.delayInit = 100 * time.Millisecond

	_ = standalone.Run()
	standalone.Stop()
}

func TestRuntime_Run_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := newDefaultStandaloneConfig(t)
	cfg.StorageBase.GRPC.Port = 3902
	standalone := NewStandaloneRuntime("test-version", &cfg)
	s := standalone.(*runtime)
	storage := server.NewMockService(ctrl)
	s.storage = storage
	storage.EXPECT().Run().Return(fmt.Errorf("err"))
	err := standalone.Run()
	assert.Error(t, err)

	standalone = NewStandaloneRuntime("test-version", &cfg)
	// restart etcd err
	err = standalone.Run()
	assert.Error(t, err)
	storage.EXPECT().Stop().Return().AnyTimes()
	s.Stop()
}

func TestRuntime_runServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := newDefaultStandaloneConfig(t)
	cfg.StorageBase.GRPC.Port = 3903
	standalone := NewStandaloneRuntime("test-version", &cfg)
	s := standalone.(*runtime)
	storage := server.NewMockService(ctrl)
	s.storage = storage
	broker := server.NewMockService(ctrl)
	s.broker = broker
	storage.EXPECT().Run().Return(nil).AnyTimes()
	broker.EXPECT().Run().Return(fmt.Errorf("err"))
	err := s.runServer()
	assert.Error(t, err)
	storage.EXPECT().Stop().Return()
	broker.EXPECT().Stop().Return()
	s.Stop()
}

func TestRuntime_cleanupState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := newDefaultStandaloneConfig(t)
	cfg.StorageBase.GRPC.Port = 3904
	standalone := NewStandaloneRuntime("test-version", &cfg)
	s := standalone.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	s.repoFactory = repoFactory
	repoFactory.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := standalone.Run()
	assert.Error(t, err)
	s.Stop()

	repo := state.NewMockRepository(ctrl)
	repoFactory.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil).AnyTimes()
	repoFactory.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil).AnyTimes()
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	repo.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	err = s.cleanupState()
	assert.Error(t, err)
}

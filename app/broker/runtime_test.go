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

package broker

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	brokerpkg "github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

var cfg = config.Broker{
	Monitor: config.Monitor{
		ReportInterval: ltoml.Duration(10 * time.Second),
	},
	Coordinator: config.RepoState{
		Namespace: "/test/",
	},
	BrokerBase: config.BrokerBase{
		HTTP: config.HTTP{
			Port: 9999,
		},
		GRPC: config.GRPC{
			Port: 2881,
		},
	}}

func TestBrokerRun(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8200")
	defer cluster.Terminate(t)

	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.Coordinator.Timeout = ltoml.Duration(time.Second * 10)
	cfg.BrokerBase.HTTP.Port = 9876

	broker := NewBrokerRuntime("test-version", &cfg, true)
	err := broker.Run()
	assert.NoError(t, err)
	// wait run finish
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, server.Running, broker.State())
	assert.Equal(t, "broker", broker.Name())

	broker.Stop()
	assert.Equal(t, server.Terminated, broker.State())
}

func TestBrokerRun_GetHost_Err(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8201")
	defer cluster.Terminate(t)

	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
	}()
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.BrokerBase.HTTP.Port = 9875

	broker := NewBrokerRuntime("test-version", &cfg, false)
	getHostIP = func() (string, error) {
		return "ip1", fmt.Errorf("err")
	}
	err := broker.Run()
	assert.Error(t, err)

	getHostIP = func() (string, error) {
		return "ip2", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	cfg.BrokerBase.HTTP.Port = 9874
	broker = NewBrokerRuntime("test-version", &cfg, false)
	err = broker.Run()
	assert.NoError(t, err)

	broker.Stop()
	assert.Equal(t, server.Terminated, broker.State())
}

func TestBroker_Run_Err(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8202")
	defer cluster.Terminate(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		newStateMachineFactory = brokerpkg.NewStateMachineFactory
		newRegistry = discovery.NewRegistry
		if err := recover(); err != nil {
			assert.NotNil(t, err)
		}
	}()
	smFct := discovery.NewMockStateMachineFactory(ctrl)
	smFct.EXPECT().Stop().AnyTimes()

	newStateMachineFactory = func(ctx context.Context,
		discoveryFactory discovery.Factory,
		stateMgr brokerpkg.StateManager,
	) discovery.StateMachineFactory {
		return smFct
	}
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.BrokerBase.HTTP.Port = 9873

	broker := NewBrokerRuntime("test-version", &cfg, false)
	b := broker.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	b.repoFactory = repoFactory
	repoFactory.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := broker.Run()
	assert.Error(t, err)
	broker.Stop()

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	repoFactory.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil).AnyTimes()

	broker = NewBrokerRuntime("test-version", &cfg, false)
	b = broker.(*runtime)
	b.repoFactory = repoFactory
	smFct.EXPECT().Start().Return(fmt.Errorf("err"))
	err = broker.Run()
	// wait run finish
	time.Sleep(500 * time.Millisecond)
	assert.Error(t, err)
	broker.Stop()

	broker = NewBrokerRuntime("test-version", &cfg, false)
	b = broker.(*runtime)
	b.repoFactory = repoFactory
	smFct.EXPECT().Start().Return(nil)
	registry := discovery.NewMockRegistry(ctrl)
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	newRegistry = func(repo state.Repository, ttl time.Duration) discovery.Registry {
		return registry
	}
	registry.EXPECT().Register(gomock.Any()).Return(fmt.Errorf("err"))
	err = broker.Run()
	assert.Error(t, err)
	// wait run finish
	time.Sleep(500 * time.Millisecond)
	broker.Stop()
}

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
	"gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	brokerpkg "github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

type testBrokerRuntimeSuite struct {
	mock.RepoTestSuite
	t *testing.T
}

func TestBrokerRuntime(t *testing.T) {
	check.Suite(&testBrokerRuntimeSuite{t: t})
	check.TestingT(t)
}

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

func (ts *testBrokerRuntimeSuite) TestBrokerRun(c *check.C) {
	cfg.Coordinator.Endpoints = ts.Cluster.Endpoints
	cfg.Coordinator.Timeout = ltoml.Duration(time.Second * 10)

	broker := NewBrokerRuntime("test-version", &cfg)
	err := broker.Run()
	if err != nil {
		c.Fatal(err)
	}
	// wait run finish
	time.Sleep(500 * time.Millisecond)

	c.Assert(server.Running, check.Equals, broker.State())
	c.Assert("broker", check.Equals, broker.Name())

	broker.Stop()
	c.Assert(server.Terminated, check.Equals, broker.State())
}

func (ts *testBrokerRuntimeSuite) TestBrokerRun_GetHost_Err(c *check.C) {
	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
	}()
	broker := NewBrokerRuntime("test-version", &cfg)
	getHostIP = func() (string, error) {
		return "ip1", fmt.Errorf("err")
	}
	err := broker.Run()
	c.Assert(err, check.NotNil)

	getHostIP = func() (string, error) {
		return "ip2", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	err = broker.Run()
	assert.NoError(ts.t, err)

	broker.Stop()
	c.Assert(server.Terminated, check.Equals, broker.State())
}

func (ts *testBrokerRuntimeSuite) TestBroker_Run_Err(_ *check.C) {
	ctrl := gomock.NewController(ts.t)
	defer ctrl.Finish()

	defer func() {
		newStateMachineFactory = brokerpkg.NewStateMachineFactory
		newRegistry = discovery.NewRegistry
		if err := recover(); err != nil {
			assert.NotNil(ts.t, err)
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

	broker := NewBrokerRuntime("test-version", &cfg)
	b := broker.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	b.repoFactory = repoFactory
	repoFactory.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := broker.Run()
	assert.Error(ts.t, err)
	broker.Stop()

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	repoFactory.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil).AnyTimes()

	broker = NewBrokerRuntime("test-version", &cfg)
	b = broker.(*runtime)
	b.repoFactory = repoFactory
	smFct.EXPECT().Start().Return(fmt.Errorf("err"))
	err = broker.Run()
	// wait run finish
	time.Sleep(500 * time.Millisecond)
	assert.Error(ts.t, err)
	broker.Stop()

	broker = NewBrokerRuntime("test-version", &cfg)
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
	assert.Error(ts.t, err)
	// wait run finish
	time.Sleep(500 * time.Millisecond)
	broker.Stop()
}

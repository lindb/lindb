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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/server"
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
	BrokerBase: config.BrokerBase{
		HTTP: config.HTTP{
			Port: 9999,
		},
		Coordinator: config.RepoState{
			Namespace: "/test/broker",
		},
		GRPC: config.GRPC{
			Port: 2881,
			TTL:  1,
		},

		ReplicationChannel: config.ReplicationChannel{
			Dir:                "/tmp/broker/replication",
			DataSizeLimit:      128,
			RemoveTaskInterval: ltoml.Duration(time.Minute),
			CheckFlushInterval: ltoml.Duration(time.Second),
			FlushInterval:      ltoml.Duration(time.Second * 5),
			BufferSize:         128,
		},
	}}

func (ts *testBrokerRuntimeSuite) TestBrokerRun(c *check.C) {
	cfg.BrokerBase.Coordinator.Endpoints = ts.Cluster.Endpoints
	cfg.BrokerBase.Coordinator.Timeout = ltoml.Duration(time.Second * 10)

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
	if err != nil {
		c.Fatal(err)
	}
}

func (ts *testBrokerRuntimeSuite) TestBroker_Run_Err(c *check.C) {
	ctrl := gomock.NewController(ts.t)
	defer ctrl.Finish()

	defer func() {
		if err := recover(); err != nil {
			assert.NotNil(ts.t, err)
		}
	}()

	broker := NewBrokerRuntime("test-version", &cfg)
	b := broker.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	b.repoFactory = repoFactory
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := broker.Run()
	assert.Error(ts.t, err)
	broker.Stop()

	repo := state.NewMockRepository(ctrl)
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(repo, nil).AnyTimes()
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = broker.Run()
	assert.Error(ts.t, err)
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	registry := discovery.NewMockRegistry(ctrl)
	b.registry = registry
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	broker.Stop()
}

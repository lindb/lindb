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
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
)

type testStorageRuntimeSuite struct {
	mock.RepoTestSuite
	t *testing.T
}

func TestStorageRuntime(t *testing.T) {
	check.Suite(&testStorageRuntimeSuite{t: t})
	check.TestingT(t)
}

var cfg = config.Storage{
	StorageBase: config.StorageBase{
		GRPC: config.GRPC{
			Port: 9999,
			TTL:  1,
		},
		TSDB: config.TSDB{Dir: "/tmp/storage/data"},
		Coordinator: config.RepoState{
			Namespace: "/test/storage",
		},
	}, Monitor: *config.NewDefaultMonitor(),
}

func (ts *testStorageRuntimeSuite) TestStorageRun(c *check.C) {
	fmt.Println("run TestStorageRun...")
	// test normal storage run
	cfg.StorageBase.Coordinator.Endpoints = ts.Cluster.Endpoints
	cfg.StorageBase.Coordinator.Timeout = ltoml.Duration(time.Second * 10)
	cfg.StorageBase.GRPC.Port = 9999
	storage := NewStorageRuntime("test-version", &cfg)
	err := storage.Run()
	assert.NoError(ts.t, err)
	c.Assert(server.Running, check.Equals, storage.State())
	// wait register success
	time.Sleep(500 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := constants.GetNodePath(constants.ActiveNodesPath, runtime.node.Indicator())
	nodeBytes, err := runtime.repo.Get(context.TODO(), nodePath)
	assert.NoError(ts.t, err)

	nodeInfo := models.ActiveNode{}
	_ = json.Unmarshal(nodeBytes, &nodeInfo)

	c.Assert(runtime.node, check.Equals, nodeInfo.Node)
	c.Assert("storage", check.Equals, storage.Name())

	storage.Stop()
	c.Assert(server.Terminated, check.Equals, storage.State())
	time.Sleep(500 * time.Millisecond)
}

func (ts *testStorageRuntimeSuite) TestStorageRun_GetHost_Err(c *check.C) {
	fmt.Println("run TestStorageRun_GetHost_Err...")
	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
	}()
	cfg.StorageBase.GRPC.Port = 8889
	storage := NewStorageRuntime("test-version", &cfg)
	getHostIP = func() (string, error) {
		return "test-ip", fmt.Errorf("err")
	}
	err := storage.Run()
	assert.Error(ts.t, err)

	getHostIP = func() (string, error) {
		return "ip", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	cfg.StorageBase.GRPC.Port = 8887
	cfg.StorageBase.Coordinator.Endpoints = ts.Cluster.Endpoints
	storage = NewStorageRuntime("test-version", &cfg)
	err = storage.Run()
	assert.NoError(ts.t, err)
	// wait grpc server start and register success
	time.Sleep(500 * time.Millisecond)
	storage.Stop()
	assert.NoError(ts.t, err)
}

func (ts *testStorageRuntimeSuite) TestStorageRun_Err(c *check.C) {
	fmt.Println("run TestStorageRun_Err...")
	ctrl := gomock.NewController(ts.t)
	defer ctrl.Finish()

	cfg.StorageBase.GRPC.Port = 8886
	storage := NewStorageRuntime("test-version", &cfg)
	s := storage.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	s.repoFactory = repoFactory
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := s.Run()
	assert.Error(ts.t, err)
	// wait grpc server start and register success
	time.Sleep(500 * time.Millisecond)

	registry := discovery.NewMockRegistry(ctrl)
	s.registry = registry
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	repo := state.NewMockRepository(ctrl)
	s.repo = repo
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	s.Stop()
	assert.Error(ts.t, err)
}

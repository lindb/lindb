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

package root

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
)

var cfg = config.Root{
	Coordinator: config.RepoState{
		Namespace: "/test/2222",
	},
	Monitor: *config.NewDefaultMonitor(),
}

func TestRootRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	cluster := mock.StartEtcdCluster(t, "http://localhost:8100")
	defer func() {
		cluster.Terminate(t)
		ctrl.Finish()
	}()

	// test normal storage run
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.Coordinator.Timeout = ltoml.Duration(time.Second * 10)
	cfg.GRPC.Port = 3997
	cfg.HTTP.Port = 3990
	root := NewRootRuntime("test-version", &cfg)
	err := root.Run()
	assert.NoError(t, err)
	assert.Equal(t, server.Running, root.State())

	assert.Equal(t, "root", root.Name())

	root.Stop()
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, server.Terminated, root.State())
}

func TestRootRun_GetHost_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	cluster := mock.StartEtcdCluster(t, "http://localhost:8101")
	defer cluster.Terminate(t)

	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
		ctrl.Finish()
	}()
	cfg.Coordinator.Endpoints = cluster.Endpoints
	cfg.GRPC.Port = 8889
	cfg.HTTP.Port = 3991
	root := NewRootRuntime("test-version", &cfg)
	getHostIP = func() (string, error) {
		return "test-ip", fmt.Errorf("err")
	}
	err := root.Run()
	assert.Error(t, err)

	getHostIP = func() (string, error) {
		return "ip", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	cfg.GRPC.Port = 8887

	root = NewRootRuntime("test-version", &cfg)
	err = root.Run()
	assert.NoError(t, err)
	// wait grpc server start and register success
	time.Sleep(500 * time.Millisecond)
	root.Stop()
	assert.NoError(t, err)
}

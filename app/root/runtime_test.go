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
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/hostutil"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

var cfg = config.Root{
	Coordinator: config.RepoState{
		Namespace: "/test/2222",
	},
	Monitor: *config.NewDefaultMonitor(),
}

func TestRootRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newRepositoryFactory = state.NewRepositoryFactory
		newStateMachineFactory = root.NewStateMachineFactory
		newRegistry = discovery.NewRegistry
		ctrl.Finish()
	}()
	registry := discovery.NewMockRegistry(ctrl)
	newRegistry = func(_ state.Repository, _ string, _ time.Duration) discovery.Registry {
		return registry
	}
	registry.EXPECT().Register(gomock.Any()).Return(nil)
	registry.EXPECT().IsSuccess().Return(true)
	registry.EXPECT().Deregister(gomock.Any()).Return(fmt.Errorf("err"))
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	repoFct := state.NewMockRepositoryFactory(ctrl)
	newRepositoryFactory = func(_ string) state.RepositoryFactory {
		return repoFct
	}
	repo := state.NewMockRepository(ctrl)
	repoFct.EXPECT().CreateRootRepo(gomock.Any()).Return(repo, nil)

	stateMachineFct := discovery.NewMockStateMachineFactory(ctrl)
	stateMachineFct.EXPECT().Start().Return(nil)
	stateMachineFct.EXPECT().Stop()
	newStateMachineFactory = func(_ context.Context, _ discovery.Factory, _ root.StateManager) discovery.StateMachineFactory {
		return stateMachineFct
	}

	cfg.Coordinator.Timeout = ltoml.Duration(time.Second * 10)
	cfg.HTTP.Port = 3990
	r := NewRootRuntime("test-version", &cfg)
	err := r.Run()
	assert.NotNil(t, r.Config())
	assert.NoError(t, err)
	assert.Equal(t, server.Running, r.State())

	assert.Equal(t, "root", r.Name())

	r.Stop()
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, server.Terminated, r.State())
}

func TestRootRun_Err(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
		newRepositoryFactory = state.NewRepositoryFactory
		newStateMachineFactory = root.NewStateMachineFactory
		newRegistry = discovery.NewRegistry
		ctrl.Finish()
	}()
	registry := discovery.NewMockRegistry(ctrl)
	registry.EXPECT().IsSuccess().Return(true)
	newRegistry = func(_ state.Repository, _ string, _ time.Duration) discovery.Registry {
		return registry
	}
	cfg.HTTP.Port = 3991
	t.Run("get host ip fail", func(t *testing.T) {
		r := NewRootRuntime("test-version", &cfg)
		getHostIP = func() (string, error) {
			return "test-ip", fmt.Errorf("err")
		}
		err := r.Run()
		assert.Error(t, err)
	})
	getHostIP = func() (string, error) {
		return "ip", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}

	repoFct := state.NewMockRepositoryFactory(ctrl)
	newRepositoryFactory = func(owner string) state.RepositoryFactory {
		return repoFct
	}
	stateMachineFct := discovery.NewMockStateMachineFactory(ctrl)
	newStateMachineFactory = func(_ context.Context, _ discovery.Factory, _ root.StateManager) discovery.StateMachineFactory {
		return stateMachineFct
	}
	t.Run("start repo fail", func(t *testing.T) {
		repoFct.EXPECT().CreateRootRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
		r := NewRootRuntime("test-version", &cfg)
		err := r.Run()
		assert.Error(t, err)
	})
	t.Run("register node fail", func(t *testing.T) {
		repoFct.EXPECT().CreateRootRepo(gomock.Any()).Return(nil, nil)
		registry.EXPECT().Register(gomock.Any()).Return(fmt.Errorf("err"))
		r := NewRootRuntime("test-version", &cfg)
		err := r.Run()
		assert.Error(t, err)
	})
	t.Run("start state machine fail", func(t *testing.T) {
		repoFct.EXPECT().CreateRootRepo(gomock.Any()).Return(nil, nil)
		registry.EXPECT().Register(gomock.Any()).Return(nil)
		stateMachineFct.EXPECT().Start().Return(fmt.Errorf("err"))
		r := NewRootRuntime("test-version", &cfg)
		err := r.Run()
		assert.Error(t, err)
		registry.EXPECT().Deregister(gomock.Any()).Return(nil)
		registry.EXPECT().Close().Return(nil)
		r.Stop()
	})
}

func TestHttpServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("no http server", func(_ *testing.T) {
		r := &runtime{
			config: &config.Root{
				HTTP: config.HTTP{},
			},
			logger: logger.GetLogger("Test", "Root"),
		}
		r.startHTTPServer()
	})

	t.Run("http server panic", func(_ *testing.T) {
		s := httppkg.NewMockServer(ctrl)
		s.EXPECT().Run().Return(fmt.Errorf("err"))
		r := &runtime{
			config: &config.Root{
				HTTP: config.HTTP{
					Port: 8000,
				},
			},
			httpServer: s,
			deps:       &deps{},
			logger:     logger.GetLogger("Test", "Root"),
		}
		assert.Panics(t, func() {
			r.runHTTPServer()
		})
	})

	t.Run("http server stop fail", func(_ *testing.T) {
		defer func() {
			newHTTPServer = httppkg.NewServer
		}()
		ctx, cancel := context.WithCancel(context.Background())
		r := &runtime{
			ctx:    ctx,
			cancel: cancel,
			config: &config.Root{
				HTTP: config.HTTP{
					Port: 8000,
				},
			},
			deps:   &deps{},
			logger: logger.GetLogger("Test", "Root"),
		}
		s := httppkg.NewMockServer(ctrl)
		s.EXPECT().Run().Return(nil)
		s.EXPECT().Close(gomock.Any()).Return(fmt.Errorf("err"))
		s.EXPECT().GetAPIRouter().Return(gin.New().Group("api"))
		newHTTPServer = func(_ config.HTTP, _ bool, _ *linmetric.Registry) httppkg.Server {
			return s
		}
		r.startHTTPServer()
		r.Stop()
		time.Sleep(100 * time.Millisecond)
	})
}

func TestRuntime_MustRegisterNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		maxRetries = 20
		retryInterval = time.Second
	}()

	maxRetries = 2
	retryInterval = time.Millisecond

	register := discovery.NewMockRegistry(ctrl)
	ctx, cancel := context.WithCancel(context.TODO())
	r := &runtime{
		ctx:      ctx,
		registry: register,
	}
	register.EXPECT().Register(gomock.Any()).Return(fmt.Errorf("err"))
	err := r.MustRegisterStatelessNode()
	assert.Error(t, err)

	register.EXPECT().Register(gomock.Any()).Return(nil)
	register.EXPECT().IsSuccess().Return(true)
	err = r.MustRegisterStatelessNode()
	assert.NoError(t, err)

	register.EXPECT().Register(gomock.Any()).Return(nil)
	register.EXPECT().IsSuccess().Return(false).MaxTimes(2)
	err = r.MustRegisterStatelessNode()
	assert.Error(t, err)

	cancel()
	register.EXPECT().Register(gomock.Any()).Return(nil)
	err = r.MustRegisterStatelessNode()
	assert.NoError(t, err)
}

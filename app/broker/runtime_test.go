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

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	brokerpkg "github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
)

var cfg = config.Broker{
	Monitor: config.Monitor{
		ReportInterval: 0,
	},
	Coordinator: config.RepoState{
		Namespace: "/test/",
	},
	BrokerBase: config.BrokerBase{
		HTTP: config.HTTP{
			Port: 9998,
		},
		GRPC: config.GRPC{
			Port: 2881,
		},
	},
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestBrokerRuntime_New(t *testing.T) {
	defer func() {
		newRepositoryFactory = state.NewRepositoryFactory
	}()
	newRepositoryFactory = func(owner string) state.RepositoryFactory {
		return nil
	}
	r := NewBrokerRuntime("version", &cfg, false)
	assert.NotNil(t, r)
	assert.Equal(t, "broker", r.Name())
	assert.NotNil(t, r.Config())
}

func TestBrokerRuntime_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoFct := state.NewMockRepositoryFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get host ip failure",
			prepare: func() {
				getHostIP = func() (string, error) {
					return "", fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "get host name/create state repo failure",
			prepare: func() {
				hostName = func() (name string, err error) {
					return "", fmt.Errorf("err")
				}
				repoFct.EXPECT().CreateNormalRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "registry alive node failure",
			prepare: func() {
				repoFct.EXPECT().CreateNormalRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register().Return(fmt.Errorf("err"))
				newRegistry = func(repo state.Repository, path string, node models.Node, ttl time.Duration) discovery.Registry {
					return registry
				}
			},
			wantErr: true,
		},
		{
			name: "start master controller failure",
			prepare: func() {
				repoFct.EXPECT().CreateNormalRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register().Return(nil)
				newRegistry = func(repo state.Repository, path string, node models.Node, ttl time.Duration) discovery.Registry {
					return registry
				}
				mc := coordinator.NewMockMasterController(ctrl)
				newMasterController = func(cfg *coordinator.MasterCfg) coordinator.MasterController {
					return mc
				}
				mc.EXPECT().WatchMasterElected(gomock.Any())
				mc.EXPECT().Start().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "broker state machine start failure, after master election",
			prepare: func() {
				repoFct.EXPECT().CreateNormalRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register().Return(nil)
				newRegistry = func(repo state.Repository, path string, node models.Node, ttl time.Duration) discovery.Registry {
					return registry
				}
				mc := coordinator.NewMockMasterController(ctrl)
				newMasterController = func(cfg *coordinator.MasterCfg) coordinator.MasterController {
					return mc
				}
				mc.EXPECT().WatchMasterElected(gomock.Any()).DoAndReturn(func(fn func(_ *models.Master)) {
					fn(&models.Master{})
				})
				mc.EXPECT().Start()
				smFct := discovery.NewMockStateMachineFactory(ctrl)
				smFct.EXPECT().Start().Return(fmt.Errorf("err"))
				newStateMachineFactory = func(ctx context.Context, discoveryFactory discovery.Factory,
					stateMgr brokerpkg.StateManager,
				) discovery.StateMachineFactory {
					return smFct
				}
			},
			wantErr: true,
		},
		{
			name: "broker successfully",
			prepare: func() {
				repoFct.EXPECT().CreateNormalRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register().Return(nil)
				newRegistry = func(repo state.Repository, path string, node models.Node, ttl time.Duration) discovery.Registry {
					return registry
				}
				mc := coordinator.NewMockMasterController(ctrl)
				newMasterController = func(cfg *coordinator.MasterCfg) coordinator.MasterController {
					return mc
				}
				mc.EXPECT().WatchMasterElected(gomock.Any()).DoAndReturn(func(fn func(_ *models.Master)) {
					fn(&models.Master{})
				})
				mc.EXPECT().Start()
				smFct := discovery.NewMockStateMachineFactory(ctrl)
				smFct.EXPECT().Start().Return(nil)
				newStateMachineFactory = func(ctx context.Context, discoveryFactory discovery.Factory,
					stateMgr brokerpkg.StateManager,
				) discovery.StateMachineFactory {
					return smFct
				}
				httpSrv := httppkg.NewMockServer(ctrl)
				httpSrv.EXPECT().GetAPIRouter().Return(gin.New().Group("/api"))
				httpSrv.EXPECT().GetPrometheusAPIRouter().Return(gin.New().Group("/prometheus"))
				newHTTPServer = func(_ config.HTTP, _ bool, _ *linmetric.Registry) httppkg.Server {
					return httpSrv
				}
				httpSrv.EXPECT().Run().Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				getHostIP = hostutil.GetHostIP
				hostName = os.Hostname
				newGRPCServer = rpc.NewGRPCServer
				newTaskClientFactory = rpc.NewTaskClientFactory
				newStateManager = brokerpkg.NewStateManager
				newChannelManager = replica.NewChannelManager
				newMasterController = coordinator.NewMasterController
				newRegistry = discovery.NewRegistry
				serveGRPCFn = serveGRPC
				newHTTPServer = httppkg.NewServer
				newStateMachineFactory = brokerpkg.NewStateMachineFactory
			}()

			r := &runtime{
				ctx:                 context.TODO(),
				enableSystemMonitor: true,
				config:              &cfg,
				repoFactory:         repoFct,
				logger:              logger.GetLogger("Runtime", "Test"),
			}
			resetNewDepsMock()
			if tt.prepare != nil {
				tt.prepare()
			}
			err := r.Run()
			time.Sleep(50 * time.Millisecond)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				assert.Equal(t, server.Failed, r.State())
			} else {
				assert.Equal(t, server.Running, r.State())
			}
		})
	}
}

func TestBrokerRuntime_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpServer := httppkg.NewMockServer(ctrl)
	registry := discovery.NewMockRegistry(ctrl)
	mc := coordinator.NewMockMasterController(ctrl)
	smFct := discovery.NewMockStateMachineFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	stateMgr := brokerpkg.NewMockStateManager(ctrl)
	connectionMgr := rpc.NewMockConnectionManager(ctrl)
	channelMgr := replica.NewMockChannelManager(ctrl)
	grpcServer := rpc.NewMockGRPCServer(ctrl)
	registry.EXPECT().Deregister().Return(fmt.Errorf("err")).AnyTimes()

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "stop failure",
			prepare: func() {
				httpServer.EXPECT().Close(gomock.Any()).Return(fmt.Errorf("err"))
				registry.EXPECT().Close().Return(fmt.Errorf("err"))
				mc.EXPECT().Stop()
				smFct.EXPECT().Stop()
				repo.EXPECT().Close().Return(fmt.Errorf("err"))
				stateMgr.EXPECT().Close()
				connectionMgr.EXPECT().Close().Return(fmt.Errorf("err"))
				channelMgr.EXPECT().Close()
				grpcServer.EXPECT().Stop()
			},
		},
		{
			name: "stop successfully",
			prepare: func() {
				httpServer.EXPECT().Close(gomock.Any()).Return(nil)
				registry.EXPECT().Close().Return(nil)
				mc.EXPECT().Stop()
				smFct.EXPECT().Stop()
				repo.EXPECT().Close().Return(nil)
				stateMgr.EXPECT().Close()
				connectionMgr.EXPECT().Close().Return(nil)
				channelMgr.EXPECT().Close()
				grpcServer.EXPECT().Stop()
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			r := &runtime{
				ctx:                 ctx,
				cancel:              cancel,
				httpServer:          httpServer,
				registry:            registry,
				master:              mc,
				stateMachineFactory: smFct,
				repo:                repo,
				stateMgr:            stateMgr,
				srv: srv{
					channelManager: channelMgr,
				},
				factory: factory{
					connectionMgr: connectionMgr,
				},
				grpcServer: grpcServer,
				logger:     logger.GetLogger("Runtime", "Test"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			r.Stop()
			assert.Equal(t, server.Terminated, r.State())
		})
	}
}

func TestBrokerRuntime_startGrpcServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpcServer := rpc.NewMockGRPCServer(ctrl)
	grpcServer.EXPECT().Start().Return(fmt.Errorf("err"))
	assert.Panics(t, func() {
		serveGRPC(grpcServer)
	})
}

func TestBrokerRuntime_RunHTTPServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := httppkg.NewMockServer(ctrl)
	s.EXPECT().Run().Return(fmt.Errorf("err"))
	r := &runtime{
		config: &config.Broker{
			BrokerBase: config.BrokerBase{
				HTTP: config.HTTP{
					Port: 8000,
				},
			},
		},
		httpServer: s,
		logger:     logger.GetLogger("Test", "Broker"),
	}
	assert.Panics(t, func() {
		r.runHTTPServer()
	})
}

func resetNewDepsMock() {
	newStateManager = func(ctx context.Context, currentNode models.StatelessNode,
		connectionManager rpc.ConnectionManager,
		taskClientFactory rpc.TaskClientFactory,
	) brokerpkg.StateManager {
		return nil
	}
	newChannelManager = func(ctx context.Context, fct rpc.ClientStreamFactory,
		stateMgr brokerpkg.StateManager,
	) replica.ChannelManager {
		return nil
	}
	newMasterController = func(cfg *coordinator.MasterCfg) coordinator.MasterController {
		return nil
	}
	newRegistry = func(repo state.Repository, path string, node models.Node, ttl time.Duration) discovery.Registry {
		return nil
	}
	serveGRPCFn = func(grpc rpc.GRPCServer) {
	}
	newStateMachineFactory = func(ctx context.Context, discoveryFactory discovery.Factory,
		stateMgr brokerpkg.StateManager,
	) discovery.StateMachineFactory {
		return nil
	}
}

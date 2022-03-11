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
	"google.golang.org/grpc"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	brokerpkg "github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	brokerQuery "github.com/lindb/lindb/query/broker"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/tag"
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
			Port: 9999,
		},
		GRPC: config.GRPC{
			Port: 2881,
		},
	}}

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
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create master controller failure",
			prepare: func() {
				mockGrpcServer(ctrl)
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				newMasterController = func(cfg *coordinator.MasterCfg) (coordinator.MasterController, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "registry alive node failure",
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register(gomock.Any()).Return(fmt.Errorf("err"))
				newRegistry = func(repo state.Repository, prefixPath string, ttl time.Duration) discovery.Registry {
					return registry
				}
			},
			wantErr: true,
		},
		{
			name: "broker state machine start failure, after master election",
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register(gomock.Any()).Return(nil)
				newRegistry = func(repo state.Repository, prefixPath string, ttl time.Duration) discovery.Registry {
					return registry
				}
				mc := coordinator.NewMockMasterController(ctrl)
				newMasterController = func(cfg *coordinator.MasterCfg) (coordinator.MasterController, error) {
					return mc, nil
				}
				mc.EXPECT().WatchMasterElected(gomock.Any()).DoAndReturn(func(fn func(_ *models.Master)) {
					fn(&models.Master{})
				})
				mc.EXPECT().Start()
				smFct := discovery.NewMockStateMachineFactory(ctrl)
				smFct.EXPECT().Start().Return(fmt.Errorf("err"))
				newStateMachineFactory = func(ctx context.Context, discoveryFactory discovery.Factory,
					stateMgr brokerpkg.StateManager) discovery.StateMachineFactory {
					return smFct
				}
			},
			wantErr: true,
		},
		{
			name: "broker successfully",
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				registry := discovery.NewMockRegistry(ctrl)
				registry.EXPECT().Register(gomock.Any()).Return(nil)
				newRegistry = func(repo state.Repository, prefixPath string, ttl time.Duration) discovery.Registry {
					return registry
				}
				mc := coordinator.NewMockMasterController(ctrl)
				newMasterController = func(cfg *coordinator.MasterCfg) (coordinator.MasterController, error) {
					return mc, nil
				}
				mc.EXPECT().WatchMasterElected(gomock.Any()).DoAndReturn(func(fn func(_ *models.Master)) {
					fn(&models.Master{})
				})
				mc.EXPECT().Start()
				smFct := discovery.NewMockStateMachineFactory(ctrl)
				smFct.EXPECT().Start().Return(nil)
				newStateMachineFactory = func(ctx context.Context, discoveryFactory discovery.Factory,
					stateMgr brokerpkg.StateManager) discovery.StateMachineFactory {
					return smFct
				}
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
				newTaskManager = brokerQuery.NewTaskManager
				newMasterController = coordinator.NewMasterController
				newRegistry = discovery.NewRegistry
				serveGRPCFn = serveGRPC

				newNativeProtoPusher = monitoring.NewNativeProtoPusher
				newStateMachineFactory = brokerpkg.NewStateMachineFactory
			}()

			r := &runtime{
				ctx:                 context.TODO(),
				enableSystemMonitor: true,
				config:              &cfg,
				repoFactory:         repoFct,
				log:                 logger.GetLogger("runtime", "test"),
			}
			resetNewDepsMock()
			if tt.prepare != nil {
				tt.prepare()
			}
			err := r.Run()
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

	pusher := monitoring.NewMockNativePusher(ctrl)
	httpServer := http.NewMockServer(ctrl)
	registry := discovery.NewMockRegistry(ctrl)
	mc := coordinator.NewMockMasterController(ctrl)
	smFct := discovery.NewMockStateMachineFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	stateMgr := brokerpkg.NewMockStateManager(ctrl)
	connectionMgr := rpc.NewMockConnectionManager(ctrl)
	channelMgr := replica.NewMockChannelManager(ctrl)
	grpcServer := rpc.NewMockGRPCServer(ctrl)

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "stop failure",
			prepare: func() {
				pusher.EXPECT().Stop()
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
				pusher.EXPECT().Stop()
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
				pusher:              pusher,
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
				log:        logger.GetLogger("runtime", "test"),
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

func TestBrokerRuntime_push_metric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newNativeProtoPusher = monitoring.NewNativeProtoPusher
		ctrl.Finish()
	}()

	pusher := monitoring.NewMockNativePusher(ctrl)
	newNativeProtoPusher = func(ctx context.Context, endpoint string,
		interval time.Duration, pushTimeout time.Duration,
		globalKeyValues tag.Tags) monitoring.NativePusher {
		return pusher
	}
	r := &runtime{
		log: logger.GetLogger("runtime", "test"),
		config: &config.Broker{Monitor: config.Monitor{
			ReportInterval: 10,
		}},
	}
	pusher.EXPECT().Start().AnyTimes()
	r.nativePusher()
	time.Sleep(500 * time.Millisecond)
}

func resetNewDepsMock() {
	newStateManager = func(ctx context.Context, currentNode models.StatelessNode,
		connectionManager rpc.ConnectionManager,
		taskClientFactory rpc.TaskClientFactory) brokerpkg.StateManager {
		return nil
	}
	newChannelManager = func(ctx context.Context, fct rpc.ClientStreamFactory,
		stateMgr brokerpkg.StateManager) replica.ChannelManager {
		return nil
	}
	newTaskManager = func(ctx context.Context, currentNode models.Node,
		taskClientFactory rpc.TaskClientFactory, taskServerFactory rpc.TaskServerFactory,
		taskPool concurrent.Pool, ttl time.Duration) brokerQuery.TaskManager {
		return nil
	}
	newMasterController = func(cfg *coordinator.MasterCfg) (coordinator.MasterController, error) {
		return nil, nil
	}
	newRegistry = func(repo state.Repository, prefixPath string, ttl time.Duration) discovery.Registry {
		return nil
	}
	serveGRPCFn = func(grpc rpc.GRPCServer) {
	}
	newStateMachineFactory = func(ctx context.Context, discoveryFactory discovery.Factory,
		stateMgr brokerpkg.StateManager) discovery.StateMachineFactory {
		return nil
	}
	newNativeProtoPusher = func(ctx context.Context, endpoint string, interval time.Duration,
		pushTimeout time.Duration, globalKeyValues tag.Tags) monitoring.NativePusher {
		return nil
	}
}

func mockGrpcServer(ctrl *gomock.Controller) {
	grpcServer := rpc.NewMockGRPCServer(ctrl)
	grpcServer.EXPECT().GetServer().Return(grpc.NewServer())
	newGRPCServer = func(cfg config.GRPC) rpc.GRPCServer {
		return grpcServer
	}
}

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

package streaming

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/app"
	rpchandler "github.com/lindb/lindb/app/streaming/rpc"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/streaming"
	"github.com/lindb/lindb/internal/api"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/state"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/tag"
	streamingpkg "github.com/lindb/lindb/streaming"
)

// rpcHandler represents all dependency rpc handlers
type rpcHandler struct {
	observer protoReplicaV1.ReplicaServiceServer
}

// just for testing
var (
	getHostIP   = hostutil.GetHostIP
	hostName    = os.Hostname
	newRegistry = discovery.NewRegistry
)

// runtime represents streaming runtime dependency
type runtime struct {
	app.BaseRuntime

	httpServer          httppkg.Server
	ctx                 context.Context
	log                 logger.Logger
	repoFactory         state.RepositoryFactory
	repo                state.Repository
	server              rpc.GRPCServer
	registry            discovery.Registry
	stateMachineFactory discovery.StateMachineFactory
	stateMgr            streaming.StateManager
	cancel              context.CancelFunc
	node                *models.StatelessNode
	config              *config.Streaming
	rpcHandler          *rpcHandler
	version             string
	globalKeyValues     tag.Tags
	state               server.State
}

// NewStreamingRuntime creates streaming runtime.
func NewStreamingRuntime(version string, cfg *config.Streaming) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:       server.New,
		repoFactory: state.NewRepositoryFactory("streaming"),
		version:     version,
		config:      cfg,
		ctx:         ctx,
		cancel:      cancel,
		log:         logger.GetLogger("Streaming", "Runtime"),
	}
}

// Config returns the configure of storage.
func (r *runtime) Config() any {
	return r.config
}

// Name returns the name of streaming service.
func (r *runtime) Name() string {
	return "streaming"
}

// Run runs streaming server.
func (r *runtime) Run() error {
	ip, err := getHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("failed to get server ip address, error: %s", err)
	}

	hostName, err := hostName()
	if err != nil {
		r.log.Error("failed to get host name", logger.Error(err))
		hostName = "unknown"
	}
	r.node = &models.StatelessNode{
		HostIP:     ip,
		GRPCPort:   r.config.StreamingBase.GRPC.Port,
		HostName:   hostName,
		HTTPPort:   r.config.StreamingBase.HTTP.Port,
		OnlineTime: timeutil.Now(),
		Version:    config.Version,
	}

	r.globalKeyValues = tag.Tags{
		{Key: []byte("node"), Value: []byte(r.node.Indicator())},
		{Key: []byte("role"), Value: []byte(constants.StreamingRole)},
		{Key: []byte("namespace"), Value: []byte(r.config.Coordinator.Namespace)},
	}
	r.BaseRuntime = app.NewBaseRuntimeFn(r.ctx, r.config.Monitor, linmetric.StreamingRegistry, r.globalKeyValues)

	// set current observer namespace
	meta.SetCurrentObserver(r.config.StreamingBase.Namespace)

	// start state repo
	if err = r.startStateRepo(); err != nil {
		r.log.Error("start state repo failure", logger.Error(err))
		r.state = server.Failed
		return err
	}

	// register observer node info
	r.registry = newRegistry(r.repo,
		constants.GetObserverLiveNodePath(r.config.StreamingBase.Namespace, r.node.Indicator()),
		r.node, r.config.Coordinator.LeaseTTL.Duration())
	if err = r.registry.Register(); err != nil {
		r.state = server.Failed
		return fmt.Errorf("register broker node error:%s", err)
	}
	discoveryFactory := discovery.NewFactory(r.repo)

	r.stateMgr = streaming.NewStateManager(r.ctx)

	r.stateMgr.RegisterWatcher(streamingpkg.NewCoordinator())

	r.stateMachineFactory = streaming.NewStateMachineFactory(r.ctx, discoveryFactory, r.stateMgr)
	if err := r.stateMachineFactory.Start(); err != nil {
		r.state = server.Failed
		return fmt.Errorf("start state machine failure:%s", err)
	}

	// start tcp server
	r.startTCPServer()
	// start http server
	r.startHTTPServer()

	// start system collector
	r.SystemCollector()
	// start stat monitoring
	r.NativePusher()

	r.state = server.Running
	return nil
}

// State returns current storage server state
func (r *runtime) State() server.State {
	return r.state
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	repo, err := r.repoFactory.CreateNormalRepo(&r.config.Coordinator)
	if err != nil {
		return fmt.Errorf("start storage state repository error:%s", err)
	}
	r.repo = repo
	r.log.Info("start storage state repository successfully")
	return nil
}

// Stop stops storage server
func (r *runtime) Stop() {
	r.log.Info("stopping storage server...")
	defer r.cancel()

	r.Shutdown()

	// close registry, deregister broker node from active list
	if r.registry != nil {
		r.log.Info("closing discovery-registry...")
		if err := r.registry.Deregister(); err != nil {
			r.log.Error("unregister storage node error", logger.Error(err))
		}
		if err := r.registry.Close(); err != nil {
			r.log.Error("unregister storage node error", logger.Error(err))
		} else {
			r.log.Info("closed discovery-registry successfully")
		}
	}

	if r.stateMachineFactory != nil {
		r.stateMachineFactory.Stop()
	}

	// close state repo if exist
	if r.repo != nil {
		r.log.Info("closing state repo...")
		if err := r.repo.Close(); err != nil {
			r.log.Error("close state repo error, when storage stop", logger.Error(err))
		} else {
			r.log.Info("closed state repo successfully")
		}
	}

	if r.stateMgr != nil {
		r.stateMgr.Close()
	}

	if r.httpServer != nil {
		r.log.Info("stopping http server...")
		if err := r.httpServer.Close(r.ctx); err != nil {
			r.log.Error("stopped http server with error", logger.Error(err))
		} else {
			r.log.Info("stopped http server successfully")
		}
	}

	// finally, shutdown rpc server
	if r.server != nil {
		r.log.Info("stopping GRPC server...")
		r.server.Stop()
		r.log.Info("stopped GRPC server")
	}

	r.log.Info("stopped storage server successfully")
	r.state = server.Terminated
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	if r.config.StreamingBase.HTTP.Port <= 0 {
		r.log.Info("http server is disabled as http-port is 0")
		return
	}

	r.httpServer = httppkg.NewServer(r.config.StreamingBase.HTTP, false, linmetric.StreamingRegistry)
	exploreAPI := api.NewExploreAPI(r.globalKeyValues, linmetric.StreamingRegistry)
	v1 := r.httpServer.GetAPIRouter().Group(constants.APIVersion1)
	exploreAPI.Register(v1)

	logAPI := api.NewLoggerAPI(r.config.Logging.Dir)
	logAPI.Register(v1)

	configAPI := api.NewConfigAPI(r.node, r.config)
	configAPI.Register(v1)

	envAPI := api.NewEnvAPI(config.ToEnvs(r.config, config.NewDefaultStreaming()))
	envAPI.Register(v1)

	go r.runHTTPServer()
}

func (r *runtime) runHTTPServer() {
	if err := r.httpServer.Run(); err != http.ErrServerClosed {
		panic(fmt.Sprintf("start http server with error: %s", err))
	}
	r.log.Info("http server stopped successfully")
}

// startTCPServer starts tcp server
func (r *runtime) startTCPServer() {
	r.server = rpc.NewGRPCServer(r.config.StreamingBase.GRPC, linmetric.StreamingRegistry)

	// bind rpc handlers
	r.bindRPCHandlers()

	go r.startRPCServer()
}

func (r *runtime) startRPCServer() {
	if err := r.server.Start(); err != nil {
		panic(err)
	}
}

// bindRPCHandlers binds rpc handlers, registers task into grpc server
func (r *runtime) bindRPCHandlers() {
	// FIXME: (stone1100) need close
	r.rpcHandler = &rpcHandler{
		observer: rpchandler.NewObserverHandler(),
	}

	protoReplicaV1.RegisterReplicaServiceServer(r.server.GetServer(), r.rpcHandler.observer)
}

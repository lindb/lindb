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
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/app"
	"github.com/lindb/lindb/app/broker/api"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/tag"
)

// just for testing
var (
	getHostIP              = hostutil.GetHostIP
	hostName               = os.Hostname
	newStateMachineFactory = broker.NewStateMachineFactory
	newRegistry            = discovery.NewRegistry
	newRepositoryFactory   = state.NewRepositoryFactory
	newGRPCServer          = rpc.NewGRPCServer
	newTaskClientFactory   = rpc.NewTaskClientFactory
	newStateManager        = broker.NewStateManager
	newChannelManager      = replica.NewChannelManager
	newMasterController    = coordinator.NewMasterController
	newHTTPServer          = httppkg.NewServer
	serveGRPCFn            = serveGRPC
)

// srv represents all services for broker
type srv struct {
	channelManager   replica.ChannelManager
	taskManager      query.TaskManager
	transportManager rpc.TransportManager
}

// factory represents all factories for broker
type factory struct {
	taskClient    rpc.TaskClientFactory
	taskServer    rpc.TaskServerFactory
	connectionMgr rpc.ConnectionManager
}

type rpcHandler struct {
	handler *query.TaskHandler
}

// runtime represents broker runtime dependency
type runtime struct {
	app.BaseRuntime
	version string
	state   server.State
	config  *config.Broker
	node    *models.StatelessNode
	// init value when runtime
	repo                state.Repository
	repoFactory         state.RepositoryFactory
	srv                 srv
	factory             factory
	httpServer          httppkg.Server
	master              coordinator.MasterController
	registry            discovery.Registry
	stateMachineFactory discovery.StateMachineFactory
	stateMgr            broker.StateManager

	grpcServer rpc.GRPCServer
	rpcHandler *rpcHandler
	queryPool  concurrent.Pool

	ctx                 context.Context
	cancel              context.CancelFunc
	globalKeyValues     tag.Tags
	enableSystemMonitor bool

	logger *logger.Logger
}

// NewBrokerRuntime creates broker runtime
func NewBrokerRuntime(version string, cfg *config.Broker, enableSystemMonitor bool) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		version:     version,
		state:       server.New,
		config:      cfg,
		repoFactory: newRepositoryFactory("broker"),
		ctx:         ctx,
		cancel:      cancel,
		queryPool: concurrent.NewPool(
			"task-pool",
			cfg.Query.QueryConcurrency,
			cfg.Query.IdleTimeout.Duration(),
			metrics.NewConcurrentStatistics("broker-query", linmetric.BrokerRegistry),
		),
		enableSystemMonitor: enableSystemMonitor,
		logger:              logger.GetLogger("Broker", "Runtime"),
	}
}

// Name returns the broker service's name
func (r *runtime) Name() string {
	return "broker"
}

// Run runs broker server based on config file
func (r *runtime) Run() error {
	ip, err := getHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("cannot get server's ip address, error: %s", err)
	}

	hostName, err := hostName()
	if err != nil {
		r.logger.Error("get host name with error", logger.Error(err))
		hostName = "unknown"
	}
	r.node = &models.StatelessNode{
		HostIP:     ip,
		HostName:   hostName,
		GRPCPort:   r.config.BrokerBase.GRPC.Port,
		HTTPPort:   r.config.BrokerBase.HTTP.Port,
		OnlineTime: timeutil.Now(),
		Version:    config.Version,
	}

	r.logger.Info("starting broker", logger.String("host", hostName), logger.String("ip", ip),
		logger.Uint16("http", r.node.HTTPPort), logger.Uint16("grpc", r.node.GRPCPort))

	// start state repository
	err = r.startStateRepo()
	if err != nil {
		r.logger.Error("failed to startStateRepo", logger.Error(err))
		r.state = server.Failed
		return err
	}
	r.globalKeyValues = tag.Tags{
		{Key: []byte("node"), Value: []byte(r.node.Indicator())},
		{Key: []byte("role"), Value: []byte(constants.BrokerRole)},
	}
	r.BaseRuntime = app.NewBaseRuntimeFn(r.ctx, r.config.Monitor, linmetric.BrokerRegistry, r.globalKeyValues)

	tackClientFct := newTaskClientFactory(r.ctx, r.node, rpc.GetBrokerClientConnFactory())
	r.factory = factory{
		taskClient:    tackClientFct,
		taskServer:    rpc.NewTaskServerFactory(),
		connectionMgr: rpc.NewConnectionManager(tackClientFct),
	}

	r.stateMgr = newStateManager(
		r.ctx,
		*r.node,
		r.factory.connectionMgr,
		r.factory.taskClient)

	r.buildServiceDependency()

	// start tcp server
	r.startGRPCServer()

	discoveryFactory := discovery.NewFactory(r.repo)

	masterCfg := &coordinator.MasterCfg{
		Ctx:              r.ctx,
		Repo:             r.repo,
		Node:             r.node,
		TTL:              int64(r.config.Coordinator.LeaseTTL.Duration().Seconds()),
		DiscoveryFactory: discoveryFactory,
		RepoFactory:      r.repoFactory,
	}
	r.master = newMasterController(masterCfg)

	// register broker node info
	r.registry = newRegistry(r.repo, constants.LiveNodesPath, r.config.Coordinator.LeaseTTL.Duration())
	err = r.registry.Register(r.node)
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("register broker node error:%s", err)
	}

	var wait sync.WaitGroup
	wait.Add(1)
	var errStore atomic.Value
	var stateMachineStarted atomic.Bool

	r.master.WatchMasterElected(func(_ *models.Master) {
		if stateMachineStarted.CAS(false, true) {
			// if state machine is not started, after 5 second when master elected, wait master state sync.
			time.AfterFunc(5*time.Second, func() {
				defer wait.Done()
				// finally, start all state machine
				r.stateMachineFactory = newStateMachineFactory(r.ctx, discoveryFactory, r.stateMgr)
				if err0 := r.stateMachineFactory.Start(); err0 != nil {
					errStore.Store(err0)
				}
			})
		}
	})

	err = r.master.Start()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("start master controller error:%s", err)
	}

	r.logger.Info("waiting broker state machine start")
	// waiting broker state machine started
	wait.Wait()
	// check if it has error when start state machine
	if errVal := errStore.Load(); errVal != nil {
		r.state = server.Failed
		return fmt.Errorf("start state machines error: %v", errVal)
	}
	r.logger.Info("broker state machine started successfully")

	// start http server
	r.startHTTPServer()

	if r.enableSystemMonitor {
		// start system collector
		r.SystemCollector()
	}
	// start stat monitoring
	r.NativePusher()

	r.state = server.Running
	return nil
}

// Config returns the configure of broker.
func (r *runtime) Config() any {
	return r.config
}

// State returns current broker server state
func (r *runtime) State() server.State {
	return r.state
}

// Stop stops broker server,
func (r *runtime) Stop() {
	r.logger.Info("stopping broker server...")
	defer r.cancel()

	r.Shutdown()

	if r.httpServer != nil {
		r.logger.Info("stopping http server...")
		if err := r.httpServer.Close(r.ctx); err != nil {
			r.logger.Error("shutdown http server error", logger.Error(err))
		} else {
			r.logger.Info("stopped http server successfully")
		}
	}

	// close registry, deregister broker node from active list
	if r.registry != nil {
		r.logger.Info("closing discovery-registry...")
		if err := r.registry.Deregister(r.node); err != nil {
			r.logger.Error("unregister broker node error", logger.Error(err))
		}
		if err := r.registry.Close(); err != nil {
			r.logger.Error("unregister broker node error", logger.Error(err))
		} else {
			r.logger.Info("closed discovery-registry successfully")
		}
	}

	if r.master != nil {
		r.logger.Info("stopping master...")
		r.master.Stop()
	}

	if r.stateMachineFactory != nil {
		r.stateMachineFactory.Stop()
	}

	if r.repo != nil {
		r.logger.Info("closing state repo...")
		if err := r.repo.Close(); err != nil {
			r.logger.Error("close state repo error, when broker stop", logger.Error(err))
		} else {
			r.logger.Info("closed state repo successfully")
		}
	}
	if r.stateMgr != nil {
		r.stateMgr.Close()
	}
	if r.srv.channelManager != nil {
		r.logger.Info("closing write channel manager...")
		r.srv.channelManager.Close()
		r.logger.Info("closed write channel successfully")
	}

	if r.factory.connectionMgr != nil {
		if err := r.factory.connectionMgr.Close(); err != nil {
			r.logger.Error("close connection manager error, when broker stop", logger.Error(err))
		} else {
			r.logger.Info("closed connection manager successfully")
		}
	}
	r.logger.Info("close connections successfully")

	// finally, shutdown rpc server
	if r.grpcServer != nil {
		r.logger.Info("stopping grpc server...")
		r.grpcServer.Stop()
		r.logger.Info("stopped grpc server successfully")
	}

	r.state = server.Terminated
	r.logger.Info("stopped broker server successfully")
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	r.logger.Info("starting HTTP server")
	r.httpServer = newHTTPServer(r.config.BrokerBase.HTTP, true, linmetric.BrokerRegistry)
	// TODO login api is not registered
	httpAPI := api.NewAPI(&deps.HTTPDeps{
		Ctx:          r.ctx,
		Node:         r.node,
		BrokerCfg:    r.config,
		Master:       r.master,
		Repo:         r.repo,
		RepoFactory:  r.repoFactory,
		StateMgr:     r.stateMgr,
		TaskMgr:      r.srv.taskManager,
		TransportMgr: r.srv.transportManager,
		CM:           r.srv.channelManager,
		IngestLimiter: concurrent.NewLimiter(
			r.ctx,
			r.config.BrokerBase.Ingestion.MaxConcurrency,
			r.config.BrokerBase.Ingestion.IngestTimeout.Duration(),
			metrics.NewLimitStatistics("ingestion", linmetric.BrokerRegistry),
		),
		QueryLimiter: concurrent.NewLimiter(
			r.ctx,
			r.config.Query.QueryConcurrency,
			r.config.Query.Timeout.Duration(),
			metrics.NewLimitStatistics("query", linmetric.BrokerRegistry),
		),
		GlobalKeyValues: r.globalKeyValues,
	})
	httpAPI.RegisterRouter(r.httpServer.GetAPIRouter())
	go r.runHTTPServer()
}

// runHTTPServer runs http server.
func (r *runtime) runHTTPServer() {
	if err := r.httpServer.Run(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("start http server with error: %s", err))
	}
	r.logger.Info("http server stopped successfully")
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	// set a sub namespace
	repo, err := r.repoFactory.CreateBrokerRepo(&r.config.Coordinator)
	if err != nil {
		return fmt.Errorf("start broker state repository error:%s", err)
	}
	r.repo = repo
	r.logger.Info("start broker state repository successfully")
	return nil
}

// buildServiceDependency builds broker service dependency
func (r *runtime) buildServiceDependency() {
	// create replica channel mgr.
	cm := newChannelManager(r.ctx, rpc.NewClientStreamFactory(r.ctx, r.node, rpc.GetBrokerClientConnFactory()), r.stateMgr)

	taskMgr := query.NewTaskManager(r.queryPool, linmetric.BrokerRegistry)
	// close connections in connection-manager
	r.factory.taskClient.SetTaskReceiver(taskMgr)

	s := srv{
		channelManager:   cm,
		taskManager:      taskMgr,
		transportManager: query.NewTransportManager(r.factory.taskClient, r.factory.taskServer, linmetric.BrokerRegistry),
	}
	r.srv = s
}

// startGRPCServer starts the GRPC server
func (r *runtime) startGRPCServer() {
	r.logger.Info("starting GRPC server")
	r.grpcServer = newGRPCServer(r.config.BrokerBase.GRPC, linmetric.BrokerRegistry)

	// bind grpc handlers
	r.rpcHandler = &rpcHandler{
		handler: query.NewTaskHandler(
			r.config.Query,
			r.factory.taskServer,
			query.NewIntermediateTaskProcessor(*r.node, r.config.Query.Timeout.Duration(),
				r.stateMgr, r.srv.taskManager, r.srv.transportManager),
			r.queryPool,
		),
	}

	protoCommonV1.RegisterTaskServiceServer(r.grpcServer.GetServer(), r.rpcHandler.handler)

	go serveGRPCFn(r.grpcServer)
}

func serveGRPC(grpc rpc.GRPCServer) {
	if err := grpc.Start(); err != nil {
		panic(err)
	}
}

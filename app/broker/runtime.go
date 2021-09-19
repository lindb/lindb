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

	"github.com/lindb/lindb/app/broker/api"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query"
	brokerQuery "github.com/lindb/lindb/query/broker"
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
)

// srv represents all services for broker
type srv struct {
	channelManager replica.ChannelManager
	taskManager    brokerQuery.TaskManager
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
	version string
	state   server.State
	config  *config.Broker
	node    *models.StatelessNode
	// init value when runtime
	repo                state.Repository
	repoFactory         state.RepositoryFactory
	srv                 srv
	factory             factory
	httpServer          *HTTPServer
	master              coordinator.Master
	registry            discovery.Registry
	stateMachineFactory discovery.StateMachineFactory
	stateMgr            broker.StateManager

	grpcServer rpc.GRPCServer
	rpcHandler *rpcHandler
	queryPool  concurrent.Pool

	ctx    context.Context
	cancel context.CancelFunc

	pusher monitoring.NativePusher

	log *logger.Logger
}

// NewBrokerRuntime creates broker runtime
func NewBrokerRuntime(version string, config *config.Broker) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		version:     version,
		state:       server.New,
		config:      config,
		repoFactory: state.NewRepositoryFactory("broker"),
		ctx:         ctx,
		cancel:      cancel,
		queryPool: concurrent.NewPool(
			"task-pool",
			config.Query.QueryConcurrency,
			config.Query.IdleTimeout.Duration(),
			linmetric.NewScope("lindb.concurrent", "pool_name", "broker-query"),
		),
		log: logger.GetLogger("broker", "Runtime"),
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
		r.log.Error("get host name with error", logger.Error(err))
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

	// start state repository
	if err := r.startStateRepo(); err != nil {
		r.log.Error("failed to startStateRepo", logger.Error(err))
		r.state = server.Failed
		return err
	}

	tackClientFct := rpc.NewTaskClientFactory(r.node)
	r.factory = factory{
		taskClient:    tackClientFct,
		taskServer:    rpc.NewTaskServerFactory(),
		connectionMgr: rpc.NewConnectionManager(tackClientFct), //TODO close connections
	}

	r.buildServiceDependency()

	// start tcp server
	r.startGRPCServer()

	discoveryFactory := discovery.NewFactory(r.repo)

	r.stateMgr = broker.NewStateManager(r.ctx, *r.node,
		r.factory.connectionMgr, r.factory.taskClient,
		r.srv.channelManager)
	// finally start all state machine
	r.stateMachineFactory = newStateMachineFactory(r.ctx, discoveryFactory, r.stateMgr)

	if err := r.stateMachineFactory.Start(); err != nil {
		return fmt.Errorf("start state machines error: %s", err)
	}

	masterCfg := &coordinator.MasterCfg{
		Ctx:               r.ctx,
		Repo:              r.repo,
		Node:              r.node,
		TTL:               1, //TODO need config
		DiscoveryFactory:  discoveryFactory,
		ControllerFactory: task.NewControllerFactory(),
		RepoFactory:       r.repoFactory,
	}
	r.master = coordinator.NewMaster(masterCfg)

	// register broker node info
	//TODO TTL default value???
	r.registry = newRegistry(r.repo, 1)
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register broker node error:%s", err)
	}
	r.master.Start()

	// start http server
	r.startHTTPServer()

	// start system collector
	r.systemCollector()
	// start stat monitoring
	r.nativePusher()

	r.state = server.Running
	return nil
}

// State returns current broker server state
func (r *runtime) State() server.State {
	return r.state
}

// Stop stops broker server,
func (r *runtime) Stop() {
	r.log.Info("stopping broker server...")
	defer r.cancel()

	if r.pusher != nil {
		r.pusher.Stop()
		r.log.Info("stopped native metric pusher successfully")
	}

	if r.httpServer != nil {
		r.log.Info("stopping http server...")
		if err := r.httpServer.Close(r.ctx); err != nil {
			r.log.Error("shutdown http server error", logger.Error(err))
		} else {
			r.log.Info("stopped http server successfully")
		}
	}

	// close registry, deregister broker node from active list
	if r.registry != nil {
		r.log.Info("closing discovery-registry...")
		if err := r.registry.Close(); err != nil {
			r.log.Error("unregister broker node error", logger.Error(err))
		} else {
			r.log.Info("closed discovery-registry successfully")
		}
	}

	if r.master != nil {
		r.log.Info("stopping master...")
		r.master.Stop()
	}

	if r.stateMachineFactory != nil {
		r.stateMachineFactory.Stop()
	}

	if r.repo != nil {
		r.log.Info("closing state repo...")
		if err := r.repo.Close(); err != nil {
			r.log.Error("close state repo error, when broker stop", logger.Error(err))
		} else {
			r.log.Info("closed state repo successfully")
		}
	}
	if r.stateMgr != nil {
		r.stateMgr.Close()
	}
	if r.srv.channelManager != nil {
		r.log.Info("closing write channel manager...")
		r.srv.channelManager.Close()
		r.log.Info("closed write channel successfully")
	}

	// finally shutdown rpc server
	if r.grpcServer != nil {
		r.log.Info("stopping grpc server...")
		r.grpcServer.Stop()
		r.log.Info("stopped grpc server successfully")
	}

	r.log.Info("stopped broker server successfully")
	r.state = server.Terminated
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	r.log.Info("starting HTTP server")
	r.httpServer = NewHTTPServer(r.config.BrokerBase.HTTP)
	// TODO login api is not registered
	httpAPI := api.NewAPI(&deps.HTTPDeps{
		Ctx:       r.ctx,
		BrokerCfg: r.config,
		Master:    r.master,
		Repo:      r.repo,
		StateMgr:  r.stateMgr,
		CM:        r.srv.channelManager,
		IngestLimiter: concurrent.NewLimiter(
			r.ctx,
			r.config.BrokerBase.Ingestion.MaxConcurrency,
			r.config.BrokerBase.Ingestion.IngestTimeout.Duration(),
			linmetric.NewScope("lindb.broker.ingestion_limiter"),
		),
		QueryLimiter: concurrent.NewLimiter(
			r.ctx,
			r.config.Query.QueryConcurrency,
			r.config.Query.Timeout.Duration(),
			linmetric.NewScope("lindb.broker.query_limiter"),
		),
		QueryFactory: brokerQuery.NewQueryFactory(
			r.stateMgr,
			r.srv.taskManager,
		),
	})
	httpAPI.RegisterRouter(r.httpServer.GetAPIRouter())
	go func() {
		if err := r.httpServer.Run(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("start http server with error: %s", err))
		}
		r.log.Info("http server stopped successfully")
	}()
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	// set a sub namespace
	repo, err := r.repoFactory.CreateBrokerRepo(r.config.Coordinator)
	if err != nil {
		return fmt.Errorf("start broker state repository error:%s", err)
	}
	r.repo = repo
	r.log.Info("start broker state repository successfully")
	return nil
}

// buildServiceDependency builds broker service dependency
func (r *runtime) buildServiceDependency() {
	// todo watch stateMachine states change.

	// hard code create channel first.
	cm := replica.NewChannelManager(r.ctx, rpc.NewClientStreamFactory(r.node))

	taskManager := brokerQuery.NewTaskManager(
		r.ctx,
		r.node,
		r.factory.taskClient,
		r.factory.taskServer,
		r.queryPool,
		r.config.Query.Timeout.Duration(),
	)

	//FIXME (stone100)close it????
	r.factory.taskClient.SetTaskReceiver(taskManager)

	srv := srv{
		channelManager: cm,
		taskManager:    taskManager,
	}
	r.srv = srv
}

// startGRPCServer starts the GRPC server
func (r *runtime) startGRPCServer() {
	r.log.Info("starting GRPC server")
	r.grpcServer = rpc.NewGRPCServer(r.config.BrokerBase.GRPC)

	// bind grpc handlers
	r.bindGRPCHandlers()

	go func() {
		if err := r.grpcServer.Start(); err != nil {
			panic(err)
		}
	}()
}

// bindGRPCHandlers binds rpc handlers, registers rpcHandler into grpc server
func (r *runtime) bindGRPCHandlers() {
	intermediateTaskProcessor := brokerQuery.NewIntermediateTaskProcessor(
		r.node,
		r.factory.taskClient,
		r.factory.taskServer,
		r.srv.taskManager,
	)
	r.rpcHandler = &rpcHandler{
		handler: query.NewTaskHandler(
			r.config.Query,
			r.factory.taskServer,
			intermediateTaskProcessor,
			r.queryPool,
		),
	}

	protoCommonV1.RegisterTaskServiceServer(r.grpcServer.GetServer(), r.rpcHandler.handler)
}

func (r *runtime) nativePusher() {
	monitorEnabled := r.config.Monitor.ReportInterval > 0
	if !monitorEnabled {
		r.log.Info("pusher won't start because report-interval is 0")
		return
	}
	r.log.Info("pusher is running",
		logger.String("interval", r.config.Monitor.ReportInterval.String()))

	r.pusher = monitoring.NewNativeProtoPusher(
		r.ctx,
		r.config.Monitor.URL,
		r.config.Monitor.ReportInterval.Duration(),
		r.config.Monitor.PushTimeout.Duration(),
		tag.Tags{
			{Key: []byte("node"), Value: []byte(r.node.Indicator())},
			{Key: []byte("role"), Value: []byte("broker")},
		},
	)
	go r.pusher.Start()
}

func (r *runtime) systemCollector() {
	r.log.Info("system collector is running")

	go monitoring.NewSystemCollector(
		r.ctx,
		"",
		r.node, "broker").Run()
}

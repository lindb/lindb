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
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/felixge/fgprof"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/storage/handler"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
	"github.com/lindb/lindb/query"
	storageQuery "github.com/lindb/lindb/query/storage"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
)

// factory represents all factories for storage
type factory struct {
	taskServer rpc.TaskServerFactory
}

// rpcHandler represents all dependency rpc handlers
type rpcHandler struct {
	replica *handler.ReplicaHandler
	write   *handler.WriteHandler
	task    *query.TaskHandler
}

// just for testing
var (
	getHostIP              = hostutil.GetHostIP
	hostName               = os.Hostname
	newStateMachineFactory = storage.NewStateMachineFactory
)

// runtime represents storage runtime dependency
type runtime struct {
	state   server.State
	version string
	config  *config.Storage

	ctx    context.Context
	cancel context.CancelFunc

	stateMachineFactory discovery.StateMachineFactory
	stateMgr            storage.StateManager

	node         *models.StatefulNode
	server       rpc.GRPCServer
	repoFactory  state.RepositoryFactory
	repo         state.Repository
	taskExecutor *storage.TaskExecutor
	factory      factory
	engine       tsdb.Engine
	rpcHandler   *rpcHandler
	httpServer   *http.Server
	queryPool    concurrent.Pool
	pusher       monitoring.NativePusher
	log          *logger.Logger
}

// NewStorageRuntime creates storage runtime
func NewStorageRuntime(version string, config *config.Storage) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:       server.New,
		repoFactory: state.NewRepositoryFactory("storage"),
		version:     version,
		config:      config,
		ctx:         ctx,
		cancel:      cancel,
		queryPool: concurrent.NewPool(
			"task-pool",
			config.Query.QueryConcurrency,
			config.Query.IdleTimeout.Duration(),
			linmetric.NewScope("lindb.concurrent.pool", "pool", "storage-query")),
		log: logger.GetLogger("storage", "Runtime"),
	}
}

// Name returns the storage service's name
func (r *runtime) Name() string {
	return "storage"
}

// Run runs storage server
func (r *runtime) Run() error {
	if r.config.StorageBase.Indicator <= 0 {
		r.state = server.Failed
		return errors.New("storage indicator must be > 0")
	}
	ip, err := getHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("failed to get server ip address, error: %s", err)
	}

	// start tsdb engine for storage server
	engine, err := tsdb.NewEngine()
	if err != nil {
		r.state = server.Failed
		return err
	}
	r.engine = engine

	hostName, err := hostName()
	if err != nil {
		r.log.Error("failed to get host name", logger.Error(err))
		hostName = "unknown"
	}
	r.node = &models.StatefulNode{
		ID: models.NodeID(r.config.StorageBase.Indicator),
		StatelessNode: models.StatelessNode{
			HostIP:     ip,
			GRPCPort:   r.config.StorageBase.GRPC.Port,
			HostName:   hostName,
			HTTPPort:   r.config.StorageBase.GRPC.Port + 1,
			OnlineTime: timeutil.Now(),
			Version:    config.Version,
		},
	}

	r.factory = factory{taskServer: rpc.NewTaskServerFactory()}
	r.stateMgr = storage.NewStateManager(r.ctx, r.node, engine)

	// start tcp server
	r.startTCPServer()
	// start http server
	r.startHTTPServer()

	// start state repo
	if err := r.startStateRepo(); err != nil {
		r.log.Error("start state repo failure", logger.Error(err))
		r.state = server.Failed
		return err
	}

	// Use Leader election mechanism to ensure the uniqueness of stateful node id
	if err := r.MustRegisterStateFulNode(); err != nil {
		return err
	}
	discoveryFactory := discovery.NewFactory(r.repo)
	// finally start all state machine
	r.stateMachineFactory = newStateMachineFactory(r.ctx, discoveryFactory, r.stateMgr)

	if err := r.stateMachineFactory.Start(); err != nil {
		return fmt.Errorf("start state machines error: %s", err)
	}

	r.taskExecutor = storage.NewTaskExecutor(r.ctx, r.node, r.repo, r.engine)
	r.taskExecutor.Run()

	// start system collector
	r.systemCollector()
	// start stat monitoring
	r.nativePusher()

	r.state = server.Running
	return nil
}

// MustRegisterStateFulNode make sure that state node is registered to etcd
func (r *runtime) MustRegisterStateFulNode() error {
	r.log.Info("registering stateful storage node...",
		logger.Int("indicator", int(r.node.ID)),
		logger.Int64("lease-ttl", r.config.Coordinator.LeaseTTL),
	)
	var (
		ok            bool
		err           error
		maxRetries    = 20
		retryInterval = time.Second
	)
	// sometimes lease isn't expired when storage restarts, retry registering is necessary
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ok, _, err = r.repo.Elect(
			r.ctx,
			constants.GetLiveNodePath(strconv.Itoa(int(r.node.ID))),
			encoding.JSONMarshal(r.node),
			r.config.Coordinator.LeaseTTL)
		if ok {
			r.log.Info("registered state node successfully",
				logger.Int("indicator", int(r.node.ID)),
				logger.Int64("lease-ttl", r.config.Coordinator.LeaseTTL),
			)
			return nil
		}
		if err != nil {
			r.log.Error("failed to register state node",
				logger.Int("indicator", int(r.node.ID)),
				logger.Int("attempt", attempt),
				logger.Error(err),
			)
		}
		if !ok {
			r.log.Error("stateful node is already registered",
				logger.Int("indicator", int(r.node.ID)),
				logger.Int("attempt", attempt),
			)
		}
		time.Sleep(retryInterval)
	}
	r.state = server.Failed
	if err != nil {
		// stateful node register err
		return err
	}
	// stateful node already exist
	r.state = server.Failed
	return constants.ErrStatefulNodeExist
}

// State returns current storage server state
func (r *runtime) State() server.State {
	return r.state
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	repo, err := r.repoFactory.CreateStorageRepo(r.config.Coordinator)
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

	if r.pusher != nil {
		r.pusher.Stop()
		r.log.Info("stopped native linmetric pusher successfully")
	}

	if r.taskExecutor != nil {
		r.log.Info("stopping task executor")
		if err := r.taskExecutor.Close(); err != nil {
			r.log.Error("stopped task executor with error", logger.Error(err))
		} else {
			r.log.Info("stooped task executor successfully")
		}
	}

	// close state repo if exist
	if r.repo != nil {
		r.log.Info("closing state repo...")
		if err := r.repo.Delete(r.ctx, constants.GetLiveNodePath(strconv.Itoa(int(r.node.ID)))); err != nil {
			r.log.Warn("delete storage node register info")
		}
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
		if err := r.httpServer.Shutdown(r.ctx); err != nil {
			r.log.Error("stopped http server with error", logger.Error(err))
		} else {
			r.log.Info("stopped http server successfully")
		}
	}

	// finally shutdown rpc server
	if r.server != nil {
		r.log.Info("stopping GRPC server...")
		r.server.Stop()
		r.log.Info("stopped GRPC server")
	}

	// close the storage engine
	if r.engine != nil {
		r.log.Info("stopping tsdb engine...")
		r.engine.Close()
		r.log.Info("stopped tsdb engine")
	}

	r.log.Info("stopped storage server successfully")
	r.state = server.Terminated
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	if !logger.IsDebug() {
		return
	}
	port := r.node.GRPCPort + 1 //TODO need remove
	r.log.Info("starting http server", logger.Uint16("port", port))

	// add prometheus metric report
	g := gin.New()
	pprof.Register(g)
	r.log.Info("/debug/pprof is enabled")
	g.GET("/debug/fgprof", gin.WrapH(fgprof.Handler()))
	r.log.Info("/debug/fgprof is enabled")

	r.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 120,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      g,
	}
	go func() {
		if err := r.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("start http server with error: %s", err))
		}
		r.log.Info("http server stopped successfully")
	}()
}

// startTCPServer starts tcp server
func (r *runtime) startTCPServer() {
	r.server = rpc.NewGRPCServer(r.config.StorageBase.GRPC)

	// bind rpc handlers
	r.bindRPCHandlers()

	go func() {
		if err := r.server.Start(); err != nil {
			panic(err)
		}
	}()
}

// bindRPCHandlers binds rpc handlers, registers task into grpc server
func (r *runtime) bindRPCHandlers() {
	//FIXME: (stone1100) need close
	leafTaskProcessor := storageQuery.NewLeafTaskProcessor(
		r.node,
		r.engine,
		r.factory.taskServer,
	)

	//TODO modify
	walMgr := replica.NewWriteAheadLogManager(
		r.ctx,
		r.config.StorageBase.WAL,
		r.node.ID, r.engine,
		rpc.NewClientStreamFactory(r.node),
		r.stateMgr,
	)
	r.rpcHandler = &rpcHandler{
		replica: handler.NewReplicaHandler(walMgr),
		write:   handler.NewWriteHandler(walMgr),
		task: query.NewTaskHandler(
			r.config.Query,
			r.factory.taskServer,
			leafTaskProcessor,
			r.queryPool,
		),
	}

	//TODO add task service ??????
	protoReplicaV1.RegisterReplicaServiceServer(r.server.GetServer(), r.rpcHandler.replica)
	protoWriteV1.RegisterWriteServiceServer(r.server.GetServer(), r.rpcHandler.write)
	protoCommonV1.RegisterTaskServiceServer(r.server.GetServer(), r.rpcHandler.task)
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
			{Key: []byte("role"), Value: []byte("storage")},
		},
	)
	go r.pusher.Start()
}

func (r *runtime) systemCollector() {
	r.log.Info("system collector is running")

	go monitoring.NewSystemCollector(
		r.ctx,
		r.config.StorageBase.TSDB.Dir,
		&r.node.StatelessNode,
		"storage").Run()
}

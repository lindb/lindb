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

	stateapi "github.com/lindb/lindb/app/storage/api/state"
	rpchandler "github.com/lindb/lindb/app/storage/rpc"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/hostutil"
	httppkg "github.com/lindb/lindb/pkg/http"
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
	replica *rpchandler.ReplicaHandler
	write   *rpchandler.WriteHandler
	task    *query.TaskHandler
}

// just for testing
var (
	getHostIP              = hostutil.GetHostIP
	hostName               = os.Hostname
	newStateMachineFactory = storage.NewStateMachineFactory
	newDatabaseLifecycleFn = NewDatabaseLifecycle
)

// runtime represents storage runtime dependency
type runtime struct {
	state   server.State
	version string
	config  *config.Storage

	delayInit   time.Duration
	initializer bootstrap.ClusterInitializer

	ctx    context.Context
	cancel context.CancelFunc

	jobScheduler kv.JobScheduler

	stateMachineFactory discovery.StateMachineFactory
	stateMgr            storage.StateManager
	walMgr              replica.WriteAheadLogManager
	dbLifecycle         DatabaseLifecycle

	node            *models.StatefulNode
	server          rpc.GRPCServer
	repoFactory     state.RepositoryFactory
	repo            state.Repository
	factory         factory
	engine          tsdb.Engine
	rpcHandler      *rpcHandler
	httpServer      httppkg.Server
	queryPool       concurrent.Pool
	pusher          monitoring.NativePusher
	globalKeyValues tag.Tags

	log *logger.Logger
}

// NewStorageRuntime creates storage runtime
func NewStorageRuntime(version string, cfg *config.Storage) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:       server.New,
		repoFactory: state.NewRepositoryFactory("storage"),
		version:     version,
		config:      cfg,
		ctx:         ctx,
		cancel:      cancel,
		queryPool: concurrent.NewPool(
			"task-pool",
			cfg.Query.QueryConcurrency,
			cfg.Query.IdleTimeout.Duration(),
			metrics.NewConcurrentStatistics("storage-query", linmetric.StorageRegistry)),
		delayInit:   time.Second,
		initializer: bootstrap.NewClusterInitializer(cfg.StorageBase.BrokerEndpoint),
		log:         logger.GetLogger("storage", "Runtime"),
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

	opt := kv.StoreOptions{
		Dir: config.GlobalStorageConfig().TSDB.Dir,
	}
	kv.Options.Store(&opt)
	r.jobScheduler = kv.NewJobScheduler(r.ctx, opt)
	r.jobScheduler.Startup() // startup kv compact job scheduler

	// start TSDB engine for storage server
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
			HTTPPort:   r.config.StorageBase.HTTP.Port,
			OnlineTime: timeutil.Now(),
			Version:    config.Version,
		},
	}
	r.globalKeyValues = tag.Tags{
		{Key: []byte("node"), Value: []byte(r.node.Indicator())},
		{Key: []byte("role"), Value: []byte(constants.StorageRole)},
		{Key: []byte("namespace"), Value: []byte(r.config.Coordinator.Namespace)},
	}

	r.factory = factory{taskServer: rpc.NewTaskServerFactory()}
	r.stateMgr = storage.NewStateManager(r.ctx, r.node, engine)

	walMgr := replica.NewWriteAheadLogManager(
		r.ctx,
		r.config.StorageBase.WAL,
		r.node.ID, r.engine,
		rpc.NewClientStreamFactory(r.ctx, r.node, rpc.GetStorageClientConnFactory()),
		r.stateMgr,
	)
	if err = walMgr.Recovery(); err != nil {
		r.state = server.Failed
		return err
	}
	r.walMgr = walMgr

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

	r.dbLifecycle = newDatabaseLifecycleFn(r.ctx, r.repo, r.walMgr, r.engine)
	r.dbLifecycle.Startup()

	// Use Leader election mechanism to ensure the uniqueness of stateful node id
	if err := r.MustRegisterStateFulNode(); err != nil {
		return err
	}
	discoveryFactory := discovery.NewFactory(r.repo)
	// finally, start all state machine
	r.stateMachineFactory = newStateMachineFactory(r.ctx, discoveryFactory, r.stateMgr)

	if err := r.stateMachineFactory.Start(); err != nil {
		return fmt.Errorf("start state machines error: %s", err)
	}

	// start system collector
	r.systemCollector()
	// start stat monitoring
	r.nativePusher()

	r.state = server.Running

	time.AfterFunc(r.delayInit, func() {
		r.log.Info("starting register storage cluster in broker")
		if err := r.initializer.InitStorageCluster(config.StorageCluster{Config: &r.config.Coordinator}); err != nil {
			r.log.Error("register storage cluster with error", logger.Error(err))
		} else {
			r.log.Info("register storage cluster successfully")
		}
	})
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
		select {
		case <-r.ctx.Done(): // no more retries when context is done
			return nil
		default:
		}
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
	repo, err := r.repoFactory.CreateStorageRepo(&r.config.Coordinator)
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

	if r.jobScheduler != nil {
		r.jobScheduler.Shutdown()
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

	if r.dbLifecycle != nil {
		r.dbLifecycle.Shutdown()
	}

	r.log.Info("stopped storage server successfully")
	r.state = server.Terminated
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	if r.config.StorageBase.HTTP.Port <= 0 {
		r.log.Info("http server is disabled as http-port is 0")
		return
	}

	r.httpServer = httppkg.NewServer(r.config.StorageBase.HTTP, false, linmetric.StorageRegistry)
	exploreAPI := monitoring.NewExploreAPI(r.globalKeyValues, linmetric.StorageRegistry)
	exploreAPI.Register(r.httpServer.GetAPIRouter())
	replicaAPI := stateapi.NewReplicaAPI(r.walMgr)
	replicaAPI.Register(r.httpServer.GetAPIRouter())
	stateMachineAPI := stateapi.NewStorageStateMachineAPI(r.stateMgr)
	stateMachineAPI.Register(r.httpServer.GetAPIRouter())
	logAPI := monitoring.NewLoggerAPI(r.config.Logging.Dir)
	logAPI.Register(r.httpServer.GetAPIRouter())
	configAPI := monitoring.NewConfigAPI(r.node, r.config)
	configAPI.Register(r.httpServer.GetAPIRouter())

	go func() {
		if err := r.httpServer.Run(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("start http server with error: %s", err))
		}
		r.log.Info("http server stopped successfully")
	}()
}

// startTCPServer starts tcp server
func (r *runtime) startTCPServer() {
	r.server = rpc.NewGRPCServer(r.config.StorageBase.GRPC, linmetric.StorageRegistry)

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

	r.rpcHandler = &rpcHandler{
		replica: rpchandler.NewReplicaHandler(r.walMgr),
		write:   rpchandler.NewWriteHandler(r.walMgr),
		task: query.NewTaskHandler(
			r.config.Query,
			r.factory.taskServer,
			leafTaskProcessor,
			r.queryPool,
		),
	}

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
		linmetric.StorageRegistry,
		r.globalKeyValues,
	)
	go r.pusher.Start()
}

func (r *runtime) systemCollector() {
	r.log.Info("system collector is running")

	go monitoring.NewSystemCollector(
		r.ctx,
		r.config.StorageBase.TSDB.Dir,
		linmetric.StorageRegistry).Run()
}

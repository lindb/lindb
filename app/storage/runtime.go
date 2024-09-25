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
	"path/filepath"
	"strconv"
	"time"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/app"
	stateapi "github.com/lindb/lindb/app/storage/api/state"
	rpchandler "github.com/lindb/lindb/app/storage/rpc"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/api"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	internalrpc "github.com/lindb/lindb/internal/rpc"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/state"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/table/metric"
	"github.com/lindb/lindb/sql/execution"
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

var (
	maxRetries    = 20
	retryInterval = time.Second
)

// just for testing
var (
	getHostIP                 = hostutil.GetHostIP
	hostName                  = os.Hostname
	newRegistry               = discovery.NewRegistry
	newStateMachineFactory    = storage.NewStateMachineFactory
	newDatabaseLifecycleFn    = NewDatabaseLifecycle
	newEngineFn               = tsdb.NewEngine
	newWriteAheadLogManagerFn = replica.NewWriteAheadLogManager
	mkDirIfNotExistFn         = fileutil.MkDirIfNotExist
	readFileFn                = os.ReadFile
	writeFileFn               = os.WriteFile

	atoiFn  = strconv.Atoi
	existFn = fileutil.Exist
)

// runtime represents storage runtime dependency
type runtime struct {
	factory             factory
	stateMachineFactory discovery.StateMachineFactory
	queryPool           concurrent.Pool
	httpServer          httppkg.Server
	engine              tsdb.Engine
	ctx                 context.Context
	log                 logger.Logger
	jobScheduler        kv.JobScheduler
	repoFactory         state.RepositoryFactory
	stateMgr            storage.StateManager
	walMgr              replica.WriteAheadLogManager
	dbLifecycle         DatabaseLifecycle
	repo                state.Repository
	server              rpc.GRPCServer
	registry            discovery.Registry
	cancel              context.CancelFunc
	node                *models.StatefulNode
	config              *config.Storage
	rpcHandler          *rpcHandler
	version             string
	app.BaseRuntime
	globalKeyValues tag.Tags
	state           server.State
	myID            int
}

// NewStorageRuntime creates storage runtime
func NewStorageRuntime(version string, myID int, cfg *config.Storage) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		myID:        myID,
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
		log: logger.GetLogger("Storage", "Runtime"),
	}
}

// Config returns the configure of storage.
func (r *runtime) Config() any {
	return r.config
}

// Name returns the storage service's name.
func (r *runtime) Name() string {
	return "storage"
}

// Run runs storage server.
func (r *runtime) Run() error {
	myID, err := r.initMyID()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("init myid failure, err: %s", err)
	}

	if myID <= 0 {
		r.state = server.Failed
		return errors.New("myid of storage server must be > 0")
	}
	ip, err := getHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("failed to get server ip address, error: %s", err)
	}

	r.jobScheduler = kv.NewJobScheduler(r.ctx, kv.DefaultCompactCheckInterval)
	r.jobScheduler.Startup() // startup kv compact job scheduler

	// start TSDB engine for storage server
	engine, err := newEngineFn()
	if err != nil {
		r.state = server.Failed
		return err
	}
	r.engine = engine

	spi.RegisterSplitSourceProvider(&metric.TableHandle{}, metric.NewSplitSourceProvider(engine))
	spi.RegisterPageSourceProvider(&metric.TableHandle{}, metric.NewPageSourceProvider())

	hostName, err := hostName()
	if err != nil {
		r.log.Error("failed to get host name", logger.Error(err))
		hostName = "unknown"
	}
	r.node = &models.StatefulNode{
		ID: models.NodeID(r.myID),
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
	r.BaseRuntime = app.NewBaseRuntimeFn(r.ctx, r.config.Monitor, linmetric.StorageRegistry, r.globalKeyValues)

	walMgr := newWriteAheadLogManagerFn(
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

	// start state repo
	if err := r.startStateRepo(); err != nil {
		r.log.Error("start state repo failure", logger.Error(err))
		r.state = server.Failed
		return err
	}

	r.factory = factory{taskServer: rpc.NewTaskServerFactory()}
	r.stateMgr = storage.NewStateManager(r.ctx, r.repo, r.node, engine)

	// start tcp server
	r.startTCPServer()
	// start http server
	r.startHTTPServer()

	discoveryFactory := discovery.NewFactory(r.repo)
	r.stateMachineFactory = newStateMachineFactory(r.ctx, discoveryFactory, r.stateMgr)
	r.dbLifecycle = newDatabaseLifecycleFn(r.ctx, r.repo, r.walMgr, r.engine)
	r.dbLifecycle.Startup()

	if err := r.startStorageState(); err != nil {
		r.state = server.Failed
		return err
	}
	// start system collector
	r.SystemCollector()
	// start stat monitoring
	r.NativePusher()

	r.state = server.Running
	return nil
}

func (r *runtime) startStorageState() error {
	// Use Leader election mechanism to ensure the uniqueness of stateful node id
	if err := r.MustRegisterStatefulNode(); err != nil {
		return err
	}
	// finally, start all state machine
	if err := r.stateMachineFactory.Start(); err != nil {
		return fmt.Errorf("start state machines error: %s", err)
	}
	return nil
}

// MustRegisterStatefulNode make sure that state node is registered to etcd
func (r *runtime) MustRegisterStatefulNode() error {
	r.log.Info("registering stateful storage node...",
		logger.Int("indicator", int(r.node.ID)),
		logger.String("lease-ttl", r.config.Coordinator.LeaseTTL.String()),
	)
	var err error
	// sometimes lease isn't expired when storage restarts, retry registering is necessary
	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-r.ctx.Done(): // no more retries when context is done
			return r.ctx.Err()
		default:
		}
		// register storage node info
		r.registry = newRegistry(r.repo, constants.GetStorageLiveNodePath(strconv.Itoa(int(r.node.ID))),
			r.node, r.config.Coordinator.LeaseTTL.Duration())
		err = r.registry.Register()
		if err != nil {
			r.log.Error("failed to register state node",
				logger.Int("indicator", int(r.node.ID)),
				logger.Int("attempt", attempt),
				logger.Error(err),
			)
			time.Sleep(retryInterval)
			continue
		}
		r.log.Info("registered state node successfully",
			logger.Int("indicator", int(r.node.ID)),
			logger.String("lease-ttl", r.config.Coordinator.LeaseTTL.String()),
		)
		return nil
	}
	r.state = server.Failed
	// stateful node register err
	return err
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

	if r.jobScheduler != nil {
		r.jobScheduler.Shutdown()
	}

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

	// close state repo if exist
	if r.repo != nil {
		r.log.Info("closing state repo...")
		if err := r.repo.Delete(r.ctx, constants.GetStorageLiveNodePath(strconv.Itoa(int(r.node.ID)))); err != nil {
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
	exploreAPI := api.NewExploreAPI(r.globalKeyValues, linmetric.StorageRegistry)
	v1 := r.httpServer.GetAPIRouter().Group(constants.APIVersion1)
	exploreAPI.Register(v1)
	replicaAPI := stateapi.NewReplicaAPI(r.walMgr)
	replicaAPI.Register(v1)
	tsdbStateAPI := stateapi.NewTSDBAPI()
	tsdbStateAPI.Register(v1)
	stateMachineAPI := stateapi.NewStorageStateMachineAPI(r.stateMgr)
	stateMachineAPI.Register(v1)
	logAPI := api.NewLoggerAPI(r.config.Logging.Dir)
	logAPI.Register(v1)
	configAPI := api.NewConfigAPI(r.node, r.config)
	configAPI.Register(v1)
	requestAPI := stateapi.NewRequestAPI()
	requestAPI.Register(v1)
	metadataAPI := stateapi.NewMetadataAPI(r.engine)
	metadataAPI.Register(v1)

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
	r.server = rpc.NewGRPCServer(r.config.StorageBase.GRPC, linmetric.StorageRegistry)

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
	leafTaskProcessor := query.NewLeafTaskProcessor(
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

	protoMetaV1.RegisterMetaServiceServer(r.server.GetServer(), rpchandler.NewMetaService(r.engine))
	protoCommandV1.RegisterCommandServiceServer(r.server.GetServer(), internalrpc.NewCommandService(execution.NewTaskManager()))
}

// initMyID initializes myid for storage server.
func (r *runtime) initMyID() (int, error) {
	dataPath := config.GlobalStorageConfig().TSDB.Dir
	if err := mkDirIfNotExistFn(dataPath); err != nil {
		return 0, err
	}
	myIDPath := filepath.Join(dataPath, "myid")
	var myID int
	if !existFn(myIDPath) {
		// if myid file not exist, use default value from start cmd and write myid file
		myID = r.myID
		if err := r.writeMyID(myIDPath, myID); err != nil {
			return 0, err
		}
	} else {
		myID0, err := r.readMyID(myIDPath)
		if err != nil {
			return 0, err
		}
		myID = myID0
	}
	return myID, nil
}

// readMyID reads myid from file.
func (r *runtime) readMyID(path string) (int, error) {
	myIDStr, err := readFileFn(path)
	if err != nil {
		return 0, err
	}
	myID, err := atoiFn(string(myIDStr))
	if err != nil {
		return 0, err
	}
	return myID, nil
}

// writeMyID writes myid into file.
func (r *runtime) writeMyID(path string, myID int) error {
	return writeFileFn(path, []byte(fmt.Sprintf("%d", myID)), 0644)
}

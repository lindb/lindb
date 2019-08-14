package broker

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/broker/api/admin"
	masterAPI "github.com/lindb/lindb/broker/api/cluster"
	queryAPI "github.com/lindb/lindb/broker/api/query"
	stateAPI "github.com/lindb/lindb/broker/api/state"
	"github.com/lindb/lindb/broker/handler"
	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/util"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	brokerpb "github.com/lindb/lindb/rpc/proto/broker"
	commonpb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
)

// srv represents all services for broker
type srv struct {
	storageClusterService service.StorageClusterService
	storageStateService   service.StorageStateService
	shardAssignService    service.ShardAssignService
	databaseService       service.DatabaseService
	channelManager        replication.ChannelManager
}

// apiHandler represents all api handlers for broker
type apiHandler struct {
	storageClusterAPI *admin.StorageClusterAPI
	databaseAPI       *admin.DatabaseAPI
	loginAPI          *api.LoginAPI
	storageStateAPI   *stateAPI.StorageAPI
	brokerStateAPI    *stateAPI.BrokerAPI
	masterAPI         *masterAPI.MasterAPI
	metricAPI         *queryAPI.MetricAPI
}

type rpcHandler struct {
	writer *handler.Writer
	task   *parallel.TaskHandler
}

type middlewareHandler struct {
	authentication middleware.Authentication
}

// runtime represents broker runtime dependency
type runtime struct {
	state  server.State
	config config.Broker
	node   models.Node
	// init value when runtime
	repo          state.Repository
	repoFactory   state.RepositoryFactory
	srv           srv
	httpServer    *http.Server
	master        coordinator.Master
	registry      discovery.Registry
	stateMachines *coordinator.BrokerStateMachines

	server     rpc.TCPServer
	handler    *rpcHandler
	middleware *middlewareHandler

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewBrokerRuntime creates broker runtime
func NewBrokerRuntime(config config.Broker) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:  server.New,
		config: config,
		ctx:    ctx,
		cancel: cancel,
		log:    logger.GetLogger("broker/runtime"),
	}
}

// Name returns the broker service's name
func (r *runtime) Name() string {
	return "broker"
}

// Run runs broker server based on config file
func (r *runtime) Run() error {
	ip, err := util.GetHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("cannot get server ip address, error:%s", err)
	}

	r.node = models.Node{IP: ip, Port: r.config.Server.Port, HostName: util.GetHostName()}

	// start state repository
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	r.buildServiceDependency()
	discoveryFactory := discovery.NewFactory(r.repo)

	smFactory := coordinator.NewStateMachineFactory(&coordinator.StateMachineCfg{
		Ctx:                 r.ctx,
		CurrentNode:         r.node,
		DiscoveryFactory:    discoveryFactory,
		ClientStreamFactory: rpc.NewClientStreamFactory(r.node),
	})

	// finally start all state machine
	r.stateMachines = coordinator.NewBrokerStateMachines(smFactory)
	if err := r.stateMachines.Start(); err != nil {
		return fmt.Errorf("start state machines error:%s", err)
	}

	r.buildMiddlewareDependency()
	r.buildAPIDependency()
	// start tcp server
	r.startTCPServer()

	// register storage node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, 1)
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register storage node error:%s", err)
	}
	masterCfg := &coordinator.MasterCfg{
		Ctx:                 r.ctx,
		Repo:                r.repo,
		Node:                r.node,
		TTL:                 1, //TODO need config
		DiscoveryFactory:    discoveryFactory,
		ControllerFactory:   task.NewControllerFactory(),
		ClusterFactory:      storage.NewClusterFactory(),
		RepoFactory:         r.repoFactory,
		StorageStateService: r.srv.storageStateService,
		ShardAssignService:  r.srv.shardAssignService,
	}

	r.master = coordinator.NewMaster(masterCfg)
	r.master.Start()

	// start http server
	r.startHTTPServer()

	r.state = server.Running
	return nil
}

// State returns current broker server state
func (r *runtime) State() server.State {
	return r.state
}

// Stop stops broker server,
func (r *runtime) Stop() error {
	r.log.Info("stopping broker server.....")
	defer r.cancel()

	if r.httpServer != nil {
		r.log.Info("starting shutdown http server")
		if err := r.httpServer.Shutdown(r.ctx); err != nil {
			r.log.Error("shutdown http server error", logger.Error(err))
		}
	}

	if r.master != nil {
		r.master.Stop()
	}

	if r.stateMachines != nil {
		r.stateMachines.Stop()
	}

	if r.repo != nil {
		r.log.Info("closing state repo")
		if err := r.repo.Close(); err != nil {
			r.log.Error("close state repo error, when broker stop", logger.Error(err))
		}
	}

	// finally shutdown rpc server
	if r.server != nil {
		r.log.Info("stopping grpc server")
		r.server.Stop()
	}

	r.log.Info("broker server stop complete")
	r.state = server.Terminated
	return nil
}

// startHTTPServer starts http server for api handler
func (r *runtime) startHTTPServer() {
	port := r.config.HTTP.Port

	r.log.Info("starting http server", logger.Uint16("port", port))
	router := api.NewRouter()
	//TODO add timeout config???
	r.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	go func() {
		if err := r.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("start http server error:%s", err))
		}
		r.log.Info("http server stop complete")
	}()
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	r.repoFactory = state.NewRepositoryFactory()
	repo, err := r.repoFactory.CreateRepo(r.config.Coordinator)
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
	cm := replication.NewChannelManager(r.config.ReplicationChannel, rpc.NewClientStreamFactory(r.node))

	srv := srv{
		storageClusterService: service.NewStorageClusterService(r.repo),
		databaseService:       service.NewDatabaseService(r.repo),
		storageStateService:   service.NewStorageStateService(r.repo),
		shardAssignService:    service.NewShardAssignService(r.repo),
		channelManager:        cm,
	}
	r.srv = srv
}

// buildAPIDependency builds broker api dependency
func (r *runtime) buildAPIDependency() {
	senderManager := parallel.NewTaskSenderManager()
	taskManager := parallel.NewTaskManager(r.node, senderManager)
	jobManager := parallel.NewJobManager(taskManager)

	handlers := apiHandler{
		storageClusterAPI: admin.NewStorageClusterAPI(r.srv.storageClusterService),
		databaseAPI:       admin.NewDatabaseAPI(r.srv.databaseService),
		loginAPI:          api.NewLoginAPI(r.config.User, r.middleware.authentication),
		storageStateAPI:   stateAPI.NewStorageAPI(r.stateMachines.StorageSM),
		brokerStateAPI:    stateAPI.NewBrokerAPI(r.stateMachines.NodeSM),
		masterAPI:         masterAPI.NewMasterAPI(r.master),
		metricAPI: queryAPI.NewMetricAPI(r.stateMachines.ReplicaStatusSM,
			r.stateMachines.NodeSM, query.NewExectorFactory(), jobManager),
	}

	api.AddRoutes("Login", http.MethodPost, "/login", handlers.loginAPI.Login)
	api.AddRoutes("Check", http.MethodGet, "/check/1", handlers.loginAPI.Check)

	api.AddRoutes("SaveStorageCluster", http.MethodPost, "/storage/cluster", handlers.storageClusterAPI.Create)
	api.AddRoutes("GetStorageCluster", http.MethodGet, "/storage/cluster", handlers.storageClusterAPI.GetByName)
	api.AddRoutes("DeleteStorageCluster", http.MethodDelete, "/storage/cluster", handlers.storageClusterAPI.DeleteByName)
	api.AddRoutes("ListStorageClusters", http.MethodGet, "/storage/cluster/list", handlers.storageClusterAPI.List)

	api.AddRoutes("CreateOrUpdateDatabase", http.MethodPost, "/database", handlers.databaseAPI.Save)
	api.AddRoutes("GetDatabase", http.MethodGet, "/database", handlers.databaseAPI.GetByName)
	api.AddRoutes("ListDatabase", http.MethodGet, "/database/list", handlers.databaseAPI.List)

	api.AddRoutes("ListStorageClusterState", http.MethodGet, "/storage/state/list", handlers.storageStateAPI.ListStorageCluster)
	api.AddRoutes("ListBrokerNodesState", http.MethodGet, "/broker/node/state", handlers.brokerStateAPI.ListBrokerNodes)

	api.AddRoutes("GetMasterState", http.MethodGet, "/cluster/master", handlers.masterAPI.GetMaster)

	api.AddRoutes("QueryMetric", http.MethodGet, "/query/metric", handlers.metricAPI.Search)
}

// buildMiddlewareDependency builds middleware dependency
// pattern support regexp matching
func (r *runtime) buildMiddlewareDependency() {
	r.middleware = &middlewareHandler{
		authentication: middleware.NewAuthentication(r.config.User),
	}
	validate, err := regexp.Compile("/check/*")
	if err == nil {
		api.AddMiddleware(r.middleware.authentication.Validate, validate)
	}
}

func (r *runtime) startTCPServer() {
	r.server = rpc.NewTCPServer(fmt.Sprintf(":%d", r.config.Server.Port))

	// bind rpc handlers
	r.bindRPCHandlers()

	go func() {
		if err := r.server.Start(); err != nil {
			panic(err)
		}
	}()
}

// bindRPCHandlers binds rpc handlers, registers handler into grpc server
func (r *runtime) bindRPCHandlers() {
	r.handler = &rpcHandler{
		writer: handler.NewWriter(r.srv.channelManager),
		task:   parallel.NewTaskHandler(rpc.GetServerStreamFactory(), nil),
	}

	brokerpb.RegisterBrokerServiceServer(r.server.GetServer(), r.handler.writer)
	commonpb.RegisterTaskServiceServer(r.server.GetServer(), r.handler.task)
}

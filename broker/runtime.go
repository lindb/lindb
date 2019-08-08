package broker

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/broker/api/admin"
	stateAPI "github.com/lindb/lindb/broker/api/state"
	"github.com/lindb/lindb/broker/handler"
	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/util"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	brokerpb "github.com/lindb/lindb/rpc/proto/broker"
	"github.com/lindb/lindb/service"
)

const (
	cfgName = "broker.toml"
	// DefaultBrokerCfgFile defines broker default config file path
	DefaultBrokerCfgFile = "./" + cfgName
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
}

type rpcHandler struct {
	writer *handler.Writer
}

// stateMachine represents all state machines for broker
type stateMachine struct {
	storageState broker.StorageStateMachine
	nodeState    broker.NodeStateMachine
}

type middlewareHandler struct {
	authentication *middleware.UserAuthentication
}

// runtime represents broker runtime dependency
type runtime struct {
	state   server.State
	cfgPath string
	config  config.Broker
	node    models.Node
	// init value when runtime
	repo         state.Repository
	repoFactory  state.RepositoryFactory
	srv          srv
	httpServer   *http.Server
	master       coordinator.Master
	registry     discovery.Registry
	stateMachine *stateMachine

	server  rpc.TCPServer
	handler *rpcHandler

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewBrokerRuntime creates broker runtime
func NewBrokerRuntime(cfgPath string) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:   server.New,
		cfgPath: cfgPath,
		ctx:     ctx,
		cancel:  cancel,
		log:     logger.GetLogger("broker/runtime"),
	}
}

// Run runs broker server based on config file
func (r *runtime) Run() error {
	if r.cfgPath == "" {
		r.cfgPath = DefaultBrokerCfgFile
	}
	if !fileutil.Exist(r.cfgPath) {
		r.state = server.Failed
		return fmt.Errorf("config file doesn't exist, see how to initialize the config by `lind broker -h`")
	}

	r.config = config.Broker{}
	if err := fileutil.DecodeToml(r.cfgPath, &r.config); err != nil {
		r.state = server.Failed
		return fmt.Errorf("decode config file error:%s", err)
	}
	r.log.Info("load broker config from file successfully", logger.String("config", r.cfgPath))

	ip, err := util.GetHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("cannot get server ip address, error:%s", err)
	}

	r.node = models.Node{IP: ip, Port: r.config.HTTP.Port}

	// start state repository
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	r.buildServiceDependency()

	// finally start all state machine
	if err := r.startStateMachine(); err != nil {
		return fmt.Errorf("start state machine error:%s", err)
	}

	r.buildMiddlewareDependency()
	r.buildAPIDependency()

	// start tcp server
	r.startTCPServer()

	// start http server
	r.startHTTPServer()

	// register storage node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, 1)
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register storage node error:%s", err)
	}

	taskController := task.NewController(r.ctx, r.repo)
	discoveryFactory := discovery.NewFactory(r.repo)
	clusterFactory := storage.NewClusterFactory()

	//TODO config ttl
	r.master = coordinator.NewMaster(r.repo, r.node, 1, taskController,
		discoveryFactory, r.repoFactory, clusterFactory, r.srv.storageStateService, r.srv.shardAssignService)
	if err := r.master.Start(); err != nil {
		return fmt.Errorf("start master error:%s", err)
	}

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

	if r.stateMachine != nil {
		r.stopStateMachine()
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
	handlers := apiHandler{
		storageClusterAPI: admin.NewStorageClusterAPI(r.srv.storageClusterService),
		databaseAPI:       admin.NewDatabaseAPI(r.srv.databaseService),
		loginAPI:          api.NewLoginAPI(r.config.User),
		storageStateAPI:   stateAPI.NewStorageAPI(r.stateMachine.storageState),
		brokerStateAPI:    stateAPI.NewBrokerAPI(r.stateMachine.nodeState),
	}

	api.AddRoutes("Login", http.MethodPost, "/login", handlers.loginAPI.Login)
	api.AddRoutes("Check", http.MethodGet, "/check/1", handlers.loginAPI.Check)

	api.AddRoutes("SaveStorageCluster", http.MethodPost, "/storage/cluster", handlers.storageClusterAPI.Create)
	api.AddRoutes("GetStorageCluster", http.MethodGet, "/storage/cluster", handlers.storageClusterAPI.GetByName)
	api.AddRoutes("DeleteStorageCluster", http.MethodDelete, "/storage/cluster", handlers.storageClusterAPI.DeleteByName)
	api.AddRoutes("ListStorageClusters", http.MethodGet, "/storage/cluster/list", handlers.storageClusterAPI.List)

	api.AddRoutes("CreateOrUpdateDatabase", http.MethodPost, "/database", handlers.databaseAPI.Save)
	api.AddRoutes("GetDatabase", http.MethodGet, "/database", handlers.databaseAPI.GetByName)

	api.AddRoutes("ListStorageClusterState", http.MethodGet, "/storage/state/list", handlers.storageStateAPI.ListStorageCluster)
	api.AddRoutes("ListBrokerNodesState", http.MethodGet, "/broker/node/state", handlers.brokerStateAPI.ListBrokerNodes)
}

// buildMiddlewareDependency builds middleware dependency
// pattern support regexp matching
func (r *runtime) buildMiddlewareDependency() {
	middlewareHandler := middlewareHandler{
		authentication: middleware.NewUserAuthentication(r.config.User),
	}
	validate, err := regexp.Compile("/check/*")
	if err == nil {
		api.AddMiddleware(middlewareHandler.authentication.ValidateMiddleware, validate)
	}

}

// startStateMachine starts related state machines for broker
func (r *runtime) startStateMachine() error {
	r.stateMachine = &stateMachine{}
	discoveryFactory := discovery.NewFactory(r.repo)
	storageStateMachine, err := broker.NewStorageStateMachine(r.ctx, r.repo, discoveryFactory)
	if err != nil {
		return err
	}
	r.stateMachine.storageState = storageStateMachine

	nodeStateMachine, err := broker.NewNodeStateMachine(r.ctx, r.node, discoveryFactory)
	if err != nil {
		return err
	}
	r.stateMachine.nodeState = nodeStateMachine

	return nil
}

// stopStateMachine stops broker's state machines
func (r *runtime) stopStateMachine() {
	if r.stateMachine.storageState != nil {
		if err := r.stateMachine.storageState.Close(); err != nil {
			r.log.Error("close storage state state machine error", logger.Error(err))
		}
	}
	if r.stateMachine.nodeState != nil {
		if err := r.stateMachine.nodeState.Close(); err != nil {
			r.log.Error("close node state state machine error", logger.Error(err))
		}
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
	}

	brokerpb.RegisterBrokerServiceServer(r.server.GetServer(), r.handler.writer)
}

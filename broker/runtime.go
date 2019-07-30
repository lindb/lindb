package broker

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/eleme/lindb/broker/api"
	"github.com/eleme/lindb/broker/api/admin"
	stateAPI "github.com/eleme/lindb/broker/api/state"
	"github.com/eleme/lindb/broker/middleware"
	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/coordinator"
	"github.com/eleme/lindb/coordinator/broker"
	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/fileutil"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/server"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/service"
)

const (
	cfgName = "broker.toml"
	// DefaultBrokerCfgFile defines broker default config file path
	DefaultBrokerCfgFile = "./" + cfgName
)

// srv represents all services for broker
type srv struct {
	storageClusterService service.StorageClusterService
	databaseService       service.DatabaseService
}

// apiHandler represents all api handlers for broker
type apiHandler struct {
	storageClusterAPI *admin.StorageClusterAPI
	databaseAPI       *admin.DatabaseAPI
	loginAPI          *api.LoginAPI
	storageStateAPI   *stateAPI.StorageAPI
	brokerStateAPI    *stateAPI.BrokerAPI
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
	srv          srv
	httpServer   *http.Server
	master       coordinator.Master
	registry     discovery.Registry
	stateMachine *stateMachine

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

	// start http server
	r.startHTTPServer()

	// register storage node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, 1)
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register storage node error:%s", err)
	}

	//TODO config ttl
	r.master = coordinator.NewMaster(r.repo, r.node, 1)
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
	repo, err := state.NewRepo(r.config.Coordinator)
	if err != nil {
		return fmt.Errorf("start broker state repository error:%s", err)
	}
	r.repo = repo
	r.log.Info("start broker state repository successfully")
	return nil
}

// buildServiceDependency builds broker service dependency
func (r *runtime) buildServiceDependency() {
	srv := srv{
		storageClusterService: service.NewStorageClusterService(r.repo),
		databaseService:       service.NewDatabaseService(r.repo),
	}
	r.srv = srv
}

// buildAPIDependency builds broker api dependency
func (r *runtime) buildAPIDependency() {
	handler := apiHandler{
		storageClusterAPI: admin.NewStorageClusterAPI(r.srv.storageClusterService),
		databaseAPI:       admin.NewDatabaseAPI(r.srv.databaseService),
		loginAPI:          api.NewLoginAPI(r.config.User),
		storageStateAPI:   stateAPI.NewStorageAPI(r.stateMachine.storageState),
		brokerStateAPI:    stateAPI.NewBrokerAPI(r.stateMachine.nodeState),
	}

	api.AddRoutes("Login", http.MethodPost, "/login", handler.loginAPI.Login)
	api.AddRoutes("Check", http.MethodGet, "/check/1", handler.loginAPI.Check)

	api.AddRoutes("SaveStorageCluster", http.MethodPost, "/storage/cluster", handler.storageClusterAPI.Create)
	api.AddRoutes("GetStorageCluster", http.MethodGet, "/storage/cluster", handler.storageClusterAPI.GetByName)
	api.AddRoutes("DeleteStorageCluster", http.MethodDelete, "/storage/cluster", handler.storageClusterAPI.DeleteByName)
	api.AddRoutes("ListStorageClusters", http.MethodGet, "/storage/cluster/list", handler.storageClusterAPI.List)

	api.AddRoutes("CreateOrUpdateDatabase", http.MethodPost, "/database", handler.databaseAPI.Save)
	api.AddRoutes("GetDatabase", http.MethodGet, "/database", handler.databaseAPI.GetByName)

	api.AddRoutes("ListStorageClusterState", http.MethodGet, "/storage/state/list", handler.storageStateAPI.ListStorageCluster)
	api.AddRoutes("ListBrokerNodesState", http.MethodGet, "/broker/node/state", handler.brokerStateAPI.ListBrokerNodes)
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
	storageStateMachine, err := broker.NewStorageStateMachine(r.ctx, r.repo)
	if err != nil {
		return err
	}
	r.stateMachine.storageState = storageStateMachine

	nodeStateMachine, err := broker.NewNodeStateMachine(r.ctx, r.repo)
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

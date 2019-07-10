package broker

import (
	"context"
	"fmt"
	"regexp"

	"github.com/eleme/lindb/broker/middleware"

	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/eleme/lindb/broker/api"
	"github.com/eleme/lindb/broker/api/admin"
	"github.com/eleme/lindb/config"
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

type srv struct {
	storageClusterService service.StorageClusterService
	databaseService       service.DatabaseService
}

type apiHandler struct {
	storageClusterAPI *admin.StorageClusterAPI
	databaseAPI       *admin.DatabaseAPI
	loginAPI          *api.LoginAPI
}

type middlewareHandler struct {
	authentication *middleware.UserAuthentication
}

// runtime represents broker runtime dependency
type runtime struct {
	state   server.State
	cfgPath string
	config  config.Broker
	ctx     context.Context
	cancel  context.CancelFunc

	// init value when runtime
	repo       state.Repository
	srv        srv
	httpServer *http.Server

	log *zap.Logger
}

// NewBrokerRuntime creates broker runtime
func NewBrokerRuntime(cfgPath string) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:   server.New,
		cfgPath: cfgPath,
		ctx:     ctx,
		cancel:  cancel,
		log:     logger.GetLogger(),
	}
}

// Run runs broker server based on config file
func (r *runtime) Run() error {
	if r.cfgPath == "" {
		r.cfgPath = DefaultBrokerCfgFile
	}
	if !util.Exist(r.cfgPath) {
		r.state = server.Failed
		return fmt.Errorf("config file doesn't exist, see how to initialize the config by `lind broker -h`")
	}

	r.config = config.Broker{}
	if err := util.DecodeToml(r.cfgPath, &r.config); err != nil {
		r.state = server.Failed
		return fmt.Errorf("decode config file error:%s", err)
	}
	r.log.Info("load broker config from file successfully", zap.String("config", r.cfgPath))

	// start state repository
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	r.buildServiceDependency()
	r.buildMiddlewareDependency()
	r.buildAPIDependency()

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
		r.log.Info("shutdowning http server")
		if err := r.httpServer.Shutdown(r.ctx); err != nil {
			r.log.Error("shutdown http server error", zap.Error(err))
		}
	}

	if r.repo != nil {
		r.log.Info("closing state repo")
		if err := r.repo.Close(); err != nil {
			r.log.Error("close state repo error, when broker stop", zap.Error(err))
		}
	}

	r.log.Info("broker server stop complete")
	r.state = server.Terminated
	return nil
}

// startHTTPServer starts http server for api handler
func (r *runtime) startHTTPServer() {
	port := r.config.HTTP.Port

	r.log.Info("starting http server", zap.Uint16("port", port))
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
	}

	api.AddRoutes("Login", http.MethodPost, "/login", handler.loginAPI.Login)
	api.AddRoutes("Check", http.MethodGet, "/check/1", handler.loginAPI.Check)

	api.AddRoutes("SaveStorageCluster", http.MethodPost, "/storage/cluster", handler.storageClusterAPI.Create)
	api.AddRoutes("GetStorageCluster", http.MethodGet, "/storage/cluster", handler.storageClusterAPI.GetByName)
	api.AddRoutes("DeleteStorageCluster", http.MethodDelete, "/storage/cluster", handler.storageClusterAPI.DeleteByName)
	api.AddRoutes("ListStorageClusters", http.MethodGet, "/storage/cluster/list", handler.storageClusterAPI.List)

	api.AddRoutes("CreateOrUpdateDatabase", http.MethodPost, "/database", handler.databaseAPI.Save)
	api.AddRoutes("GetDatabase", http.MethodGet, "/database", handler.databaseAPI.GetByName)
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

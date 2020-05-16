package storage

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	promreporter "github.com/uber-go/tally/prometheus"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	task "github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/monitoring"
	taskHandler "github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/rpc/proto/storage"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/storage/handler"
	"github.com/lindb/lindb/tsdb"
)

// srv represents all dependency services
type srv struct {
	storageService service.StorageService
}

// factory represents all factories for storage
type factory struct {
	taskServer rpc.TaskServerFactory
}

// rpcHandler represents all dependency rpc handlers
type rpcHandler struct {
	writer *handler.Writer
	task   *taskHandler.TaskHandler
}

// just for testing
var getHostIP = hostutil.GetHostIP
var hostName = os.Hostname

// runtime represents storage runtime dependency
type runtime struct {
	state   server.State
	version string
	config  config.Storage

	ctx    context.Context
	cancel context.CancelFunc

	node         models.Node
	server       rpc.GRPCServer
	repoFactory  state.RepositoryFactory
	repo         state.Repository
	registry     discovery.Registry
	taskExecutor *task.TaskExecutor
	factory      factory
	srv          srv
	handler      *rpcHandler
	httpServer   *http.Server

	pusher monitoring.PrometheusPusher

	log *logger.Logger
}

// NewStorageRuntime creates storage runtime
func NewStorageRuntime(version string, config config.Storage) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:       server.New,
		repoFactory: state.NewRepositoryFactory("storage"),
		version:     version,
		config:      config,
		ctx:         ctx,
		cancel:      cancel,

		log: logger.GetLogger("storage", "Runtime"),
	}
}

// Name returns the storage service's name
func (r *runtime) Name() string {
	return "storage"
}

// Run runs storage server
func (r *runtime) Run() error {
	ip, err := getHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("cannot get server ip address, error:%s", err)
	}

	// build service dependency for storage server
	if err := r.buildServiceDependency(); err != nil {
		r.state = server.Failed
		return err
	}
	hostName, err := hostName()
	if err != nil {
		r.log.Error("get host name with error", logger.Error(err))
		hostName = "unknown"
	}
	r.node = models.Node{
		IP:       ip,
		Port:     r.config.StorageBase.GRPC.Port,
		HostName: hostName,
		HTTPPort: r.config.StorageBase.GRPC.Port + 1,
	}

	r.factory = factory{taskServer: rpc.NewTaskServerFactory()}

	// start tcp server
	r.startTCPServer()
	// start http server
	r.startHTTPServer()

	// start state repo
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	// register storage node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, r.config.StorageBase.GRPC.TTL.Duration())
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register storage node error:%s", err)
	}

	r.taskExecutor = task.NewTaskExecutor(r.ctx, &r.node, r.repo, r.srv.storageService)
	r.taskExecutor.Run()

	// start stat monitoring
	r.monitoring()

	r.state = server.Running
	return nil
}

// State returns current storage server state
func (r *runtime) State() server.State {
	return r.state
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	repo, err := r.repoFactory.CreateRepo(r.config.StorageBase.Coordinator)
	if err != nil {
		return fmt.Errorf("start storage state repository error:%s", err)
	}
	r.repo = repo
	r.log.Info("start storage state repository successfully")
	return nil
}

// Stop stops storage server
func (r *runtime) Stop() error {
	defer r.cancel()

	if r.pusher != nil {
		r.pusher.Stop()
	}

	if r.taskExecutor != nil {
		if err := r.taskExecutor.Close(); err != nil {
			r.log.Error("close task executor error", logger.Error(err))
		}
	}

	// close registry, deregister storage node from active list
	if r.registry != nil {
		if err := r.registry.Close(); err != nil {
			r.log.Error("unregister storage node error", logger.Error(err))
		}
	}

	// close state repo if exist
	if r.repo != nil {
		r.log.Info("closing state repo")
		if err := r.repo.Close(); err != nil {
			r.log.Error("close state repo error, when storage stop", logger.Error(err))
		}
	}

	if r.httpServer != nil {
		r.log.Info("starting shutdown http server")
		if err := r.httpServer.Shutdown(r.ctx); err != nil {
			r.log.Error("shutdown http server error", logger.Error(err))
		}
	}

	// finally shutdown rpc server
	if r.server != nil {
		r.log.Info("stopping grpc server")
		r.server.Stop()
	}

	// close the storage engine
	if r.srv.storageService != nil {
		r.srv.storageService.Close()
	}

	r.log.Info("storage server stop complete")
	r.state = server.Terminated
	return nil
}

// buildServiceDependency builds broker service dependency
func (r *runtime) buildServiceDependency() error {
	engine, err := tsdb.NewEngine(r.config.StorageBase.TSDB)
	if err != nil {
		return err
	}
	srv := srv{
		storageService: service.NewStorageService(engine),
	}
	r.srv = srv
	return nil
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	port := r.node.Port + 1
	r.log.Info("starting http server", logger.Uint16("port", port))

	// add prometheus metric report
	reporter := promreporter.NewReporter(promreporter.Options{})
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/metrics", reporter.HTTPHandler())

	r.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
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
	r.server = rpc.NewGRPCServer(fmt.Sprintf(":%d", r.node.Port))

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
	//FIXME: (stone1100) need close
	dispatcher := taskHandler.NewLeafTaskDispatcher(r.node, r.srv.storageService,
		query.NewExecutorFactory(), r.factory.taskServer)

	r.handler = &rpcHandler{
		writer: handler.NewWriter(r.srv.storageService),
		task:   taskHandler.NewTaskHandler(r.config.StorageBase.Query, r.factory.taskServer, dispatcher),
	}

	//TODO add task service ??????
	storage.RegisterWriteServiceServer(r.server.GetServer(), r.handler.writer)
	common.RegisterTaskServiceServer(r.server.GetServer(), r.handler.task)
}

func (r *runtime) monitoring() {
	systemStatMonitorEnabled := r.config.Monitor.SystemReportInterval > 0
	if systemStatMonitorEnabled {
		r.log.Info("SystemStatMonitor is running")
		go monitoring.NewSystemCollector(
			r.ctx,
			r.config.Monitor.SystemReportInterval.Duration(),
			r.config.StorageBase.TSDB.Dir,
			r.repo,
			constants.GetNodeMonitoringStatPath(r.node.Indicator()),
			models.ActiveNode{
				Version:    r.version,
				Node:       r.node,
				OnlineTime: timeutil.Now(),
			}).Run()
	}

	if !config.StandaloneMode {
		// disable if run as standalone, because broker start same monitoring
		runtimeStatMonitorEnabled := r.config.Monitor.RuntimeReportInterval > 0
		if runtimeStatMonitorEnabled {
			r.log.Info("RuntimeStatMonitor is running")
			go monitoring.NewRunTimeCollector(
				r.ctx,
				r.config.Monitor.RuntimeReportInterval.Duration(),
				map[string]string{"role": "broker", "version": r.version},
			)
		}
	}

	r.pusher = monitoring.NewPrometheusPusher(
		r.ctx,
		r.config.Monitor.URL,
		r.config.Monitor.RuntimeReportInterval.Duration(),
		prometheus.Gatherers{monitoring.StorageGatherer, prometheus.DefaultGatherer},
		[]*dto.LabelPair{
			{
				Name:  proto.String("role"),
				Value: proto.String("storage"),
			},
			{
				Name:  proto.String("namespace"),
				Value: proto.String(r.config.StorageBase.Coordinator.Namespace),
			},
			{
				Name:  proto.String("node"),
				Value: proto.String(r.node.Indicator()),
			},
		},
	)
	go r.pusher.Start()
}

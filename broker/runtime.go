package broker

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	promreporter "github.com/uber-go/tally/prometheus"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/broker/api/admin"
	masterAPI "github.com/lindb/lindb/broker/api/cluster"
	writeAPI "github.com/lindb/lindb/broker/api/metric"
	queryAPI "github.com/lindb/lindb/broker/api/query"
	stateAPI "github.com/lindb/lindb/broker/api/state"
	"github.com/lindb/lindb/broker/api/write"
	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	commonpb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
)

// just for testing
var getHostIP = hostutil.GetHostIP
var hostName = os.Hostname

// srv represents all services for broker
type srv struct {
	storageClusterService service.StorageClusterService
	storageStateService   service.StorageStateService
	shardAssignService    service.ShardAssignService
	databaseService       service.DatabaseService
	replicatorStateReport replication.ReplicatorStateReport
	channelManager        replication.ChannelManager
	taskManager           parallel.TaskManager
	jobManager            parallel.JobManager
}

// factory represents all factories for broker
type factory struct {
	taskClient rpc.TaskClientFactory
	taskServer rpc.TaskServerFactory
}

// apiHandler represents all api handlers for broker
type apiHandler struct {
	storageClusterAPI  *admin.StorageClusterAPI
	databaseAPI        *admin.DatabaseAPI
	databaseFlusherAPI *admin.DatabaseFlusherAPI
	loginAPI           *api.LoginAPI
	storageStateAPI    *stateAPI.StorageAPI
	brokerStateAPI     *stateAPI.BrokerAPI
	masterAPI          *masterAPI.MasterAPI
	metricAPI          *queryAPI.MetricAPI
	metadataAPI        *queryAPI.MetadataAPI
	writeAPI           *writeAPI.WriteAPI
	prometheusWriter   *write.PrometheusWrite
}

type rpcHandler struct {
	task *parallel.TaskHandler
}

type middlewareHandler struct {
	authentication middleware.Authentication
}

// runtime represents broker runtime dependency
type runtime struct {
	version string
	state   server.State
	config  config.Broker
	node    models.Node
	// init value when runtime
	repo          state.Repository
	repoFactory   state.RepositoryFactory
	srv           srv
	factory       factory
	httpServer    *http.Server
	master        coordinator.Master
	registry      discovery.Registry
	stateMachines *coordinator.BrokerStateMachines

	grpcServer rpc.GRPCServer
	rpcHandler *rpcHandler

	middleware *middlewareHandler

	ctx    context.Context
	cancel context.CancelFunc

	pusher monitoring.PrometheusPusher

	log *logger.Logger
}

// NewBrokerRuntime creates broker runtime
func NewBrokerRuntime(version string, config config.Broker) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		version:     version,
		state:       server.New,
		config:      config,
		repoFactory: state.NewRepositoryFactory("broker"),
		ctx:         ctx,
		cancel:      cancel,
		log:         logger.GetLogger("broker", "Runtime"),
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
	r.node = models.Node{
		IP:       ip,
		Port:     r.config.BrokerBase.GRPC.Port,
		HostName: hostName,
		HTTPPort: r.config.BrokerBase.HTTP.Port,
	}

	// start state repository
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	r.factory = factory{
		taskClient: rpc.NewTaskClientFactory(r.node),
		taskServer: rpc.NewTaskServerFactory(),
	}

	r.buildServiceDependency()
	discoveryFactory := discovery.NewFactory(r.repo)

	smFactory := coordinator.NewStateMachineFactory(&coordinator.StateMachineCfg{
		Ctx:               r.ctx,
		CurrentNode:       r.node,
		ChannelManager:    r.srv.channelManager,
		ShardAssignSRV:    r.srv.shardAssignService,
		DiscoveryFactory:  discoveryFactory,
		TaskClientFactory: r.factory.taskClient,
	})

	// finally start all state machine
	r.stateMachines = coordinator.NewBrokerStateMachines(smFactory)
	if err := r.stateMachines.Start(); err != nil {
		return fmt.Errorf("start state machines error: %s", err)
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
		BrokerSM:            r.stateMachines,
	}
	r.master = coordinator.NewMaster(masterCfg)

	r.buildMiddlewareDependency()
	r.buildAPIDependency()
	// start tcp server
	r.startGRPCServer()

	// register broker node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, 1)
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register storage node error:%s", err)
	}
	r.master.Start()

	// start http server
	r.startHTTPServer()

	// start stat monitoring
	r.monitoring()

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

	if r.pusher != nil {
		r.pusher.Stop()
	}

	if r.httpServer != nil {
		r.log.Info("starting shutdown http server")
		if err := r.httpServer.Shutdown(r.ctx); err != nil {
			r.log.Error("shutdown http server error", logger.Error(err))
		}
	}

	// close registry, deregister broker node from active list
	if r.registry != nil {
		if err := r.registry.Close(); err != nil {
			r.log.Error("unregister broker node error", logger.Error(err))
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
	if r.grpcServer != nil {
		r.log.Info("stopping grpc server")
		r.grpcServer.Stop()
	}

	r.log.Info("broker server stop complete")
	r.state = server.Terminated
	return nil
}

// startHTTPServer starts http server for api rpcHandler
func (r *runtime) startHTTPServer() {
	port := r.config.BrokerBase.HTTP.Port

	r.log.Info("starting http server", logger.Uint16("port", port))
	router := api.NewRouter()

	// add prometheus metric report
	reporter := promreporter.NewReporter(promreporter.Options{})
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

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	repo, err := r.repoFactory.CreateRepo(r.config.BrokerBase.Coordinator)
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

	replicatorStateReport := replication.NewReplicatorStateReport(r.node, r.repo)

	// hard code create channel first.
	cm := replication.NewChannelManager(r.config.BrokerBase.ReplicationChannel, rpc.NewClientStreamFactory(r.node), replicatorStateReport)
	taskManager := parallel.NewTaskManager(r.node, r.factory.taskClient, r.factory.taskServer)
	jobManager := parallel.NewJobManager(taskManager)

	//FIXME (stone100)close it????
	taskReceiver := parallel.NewTaskReceiver(jobManager)
	r.factory.taskClient.SetTaskReceiver(taskReceiver)

	srv := srv{
		storageClusterService: service.NewStorageClusterService(r.repo),
		databaseService:       service.NewDatabaseService(r.repo),
		storageStateService:   service.NewStorageStateService(r.repo),
		shardAssignService:    service.NewShardAssignService(r.repo),
		replicatorStateReport: replicatorStateReport,
		channelManager:        cm,
		taskManager:           taskManager,
		jobManager:            jobManager,
	}
	r.srv = srv
}

// buildAPIDependency builds broker api dependency
func (r *runtime) buildAPIDependency() {
	handlers := apiHandler{
		storageClusterAPI:  admin.NewStorageClusterAPI(r.srv.storageClusterService),
		databaseAPI:        admin.NewDatabaseAPI(r.srv.databaseService),
		databaseFlusherAPI: admin.NewDatabaseFlusherAPI(r.master),
		loginAPI:           api.NewLoginAPI(r.config.BrokerBase.User, r.middleware.authentication),
		storageStateAPI:    stateAPI.NewStorageAPI(r.ctx, r.repo, r.stateMachines.StorageSM, r.srv.shardAssignService, r.srv.databaseService),
		brokerStateAPI:     stateAPI.NewBrokerAPI(r.ctx, r.repo, r.stateMachines.NodeSM),
		masterAPI:          masterAPI.NewMasterAPI(r.master),
		metricAPI: queryAPI.NewMetricAPI(r.stateMachines.ReplicaStatusSM,
			r.stateMachines.NodeSM, r.stateMachines.DatabaseSM, query.NewExecutorFactory(), r.srv.jobManager),
		metadataAPI: queryAPI.NewMetadataAPI(r.srv.databaseService, r.stateMachines.ReplicaStatusSM,
			r.stateMachines.NodeSM, query.NewExecutorFactory(), r.srv.jobManager),
		writeAPI:         writeAPI.NewWriteAPI(r.srv.channelManager),
		prometheusWriter: write.NewPrometheusWrite(r.srv.channelManager),
	}

	api.AddRoute("Login", http.MethodPost, "/login", handlers.loginAPI.Login)
	api.AddRoute("Check", http.MethodGet, "/check/1", handlers.loginAPI.Check)

	api.AddRoute("SaveStorageCluster", http.MethodPost, "/storage/cluster", handlers.storageClusterAPI.Create)
	api.AddRoute("GetStorageCluster", http.MethodGet, "/storage/cluster", handlers.storageClusterAPI.GetByName)
	api.AddRoute("DeleteStorageCluster", http.MethodDelete, "/storage/cluster", handlers.storageClusterAPI.DeleteByName)
	api.AddRoute("ListStorageClusters", http.MethodGet, "/storage/cluster/list", handlers.storageClusterAPI.List)

	api.AddRoute("CreateOrUpdateDatabase", http.MethodPost, "/database", handlers.databaseAPI.Save)
	api.AddRoute("GetDatabase", http.MethodGet, "/database", handlers.databaseAPI.GetByName)
	api.AddRoute("ListDatabase", http.MethodGet, "/database/list", handlers.databaseAPI.List)
	api.AddRoute("FLushDatabase", http.MethodGet, "/database/flush", handlers.databaseFlusherAPI.SubmitFlushTask)

	api.AddRoute("ListStorageClusterNodesState", http.MethodGet, "/storage/cluster/state", handlers.storageStateAPI.GetStorageClusterState)
	api.AddRoute("ListStorageClusterState", http.MethodGet, "/storage/cluster/state/list", handlers.storageStateAPI.ListStorageClusterState)
	api.AddRoute("ListBrokerClusterState", http.MethodGet, "/broker/cluster/state", handlers.brokerStateAPI.ListBrokersStat)

	api.AddRoute("GetMasterState", http.MethodGet, "/cluster/master", handlers.masterAPI.GetMaster)

	api.AddRoute("QueryMetric", http.MethodGet, "/query/metric", handlers.metricAPI.Search)
	api.AddRoute("QueryMetadata", http.MethodGet, "/query/metadata", handlers.metadataAPI.Handle)

	api.AddRoute("WriteSumMetric", http.MethodPut, "/metric/sum", handlers.writeAPI.Sum)
	api.AddRoute("PrometheusWriter", http.MethodPut, "/metric/prometheus", handlers.prometheusWriter.Write)
}

// buildMiddlewareDependency builds middleware dependency
// pattern support regexp matching
func (r *runtime) buildMiddlewareDependency() {
	r.middleware = &middlewareHandler{
		authentication: middleware.NewAuthentication(r.config.BrokerBase.User),
	}
	httpAPI, err := regexp.Compile("/*")
	if err == nil {
		api.AddMiddleware(middleware.AccessLogMiddleware, httpAPI)
	}
	validate, err := regexp.Compile("/check/*")
	if err == nil {
		api.AddMiddleware(r.middleware.authentication.Validate, validate)
	}
}

// startGRPCServer starts the GRPC server
func (r *runtime) startGRPCServer() {
	r.grpcServer = rpc.NewGRPCServer(fmt.Sprintf(":%d", r.config.BrokerBase.GRPC.Port))

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
	//FIXME: (stone1100) need close
	dispatcher := parallel.NewIntermediateTaskDispatcher()
	r.rpcHandler = &rpcHandler{
		task: parallel.NewTaskHandler(r.config.BrokerBase.Query, r.factory.taskServer, dispatcher),
	}

	commonpb.RegisterTaskServiceServer(r.grpcServer.GetServer(), r.rpcHandler.task)
}

func (r *runtime) monitoring() {
	systemStatMonitorEnabled := r.config.Monitor.SystemReportInterval > 0
	node := models.ActiveNode{
		Version:    r.version,
		Node:       r.node,
		OnlineTime: timeutil.Now(),
	}
	if systemStatMonitorEnabled {
		r.log.Info("SystemStatMonitor is running")
		go monitoring.NewSystemCollector(
			r.ctx,
			r.config.Monitor.SystemReportInterval.Duration(),
			r.config.BrokerBase.ReplicationChannel.Dir,
			r.repo,
			constants.GetNodeMonitoringStatPath(r.node.Indicator()),
			node).Run()
	}

	runtimeStatMonitorEnabled := r.config.Monitor.RuntimeReportInterval > 0
	if runtimeStatMonitorEnabled {
		r.log.Info("RuntimeStatMonitor is running")
		go monitoring.NewRunTimeCollector(
			r.ctx,
			r.config.Monitor.RuntimeReportInterval.Duration(),
			map[string]string{"role": "broker", "version": r.version},
		)
	}

	r.pusher = monitoring.NewPrometheusPusher(
		r.ctx,
		r.config.Monitor.URL,
		r.config.Monitor.RuntimeReportInterval.Duration(),
		prometheus.Gatherers{monitoring.BrokerGatherer, prometheus.DefaultGatherer},
		[]*dto.LabelPair{
			{
				Name:  proto.String("role"),
				Value: proto.String("broker"),
			},
			{
				Name:  proto.String("node"),
				Value: proto.String(r.node.Indicator()),
			},
		},
	)
	go r.pusher.Start()
}

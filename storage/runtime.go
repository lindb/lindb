package storage

import (
	"context"
	"fmt"
	"os"

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
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/rpc/proto/storage"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/storage/handler"
	"github.com/lindb/lindb/tsdb"
)

// srv represents all dependency services
type srv struct {
	storageService  service.StorageService
	sequenceManager replication.SequenceManager
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
	r.node = models.Node{IP: ip, Port: r.config.GRPC.Port, HostName: hostName}

	r.factory = factory{taskServer: rpc.NewTaskServerFactory()}

	// start tcp server
	r.startTCPServer()

	// start state repo
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	// register storage node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, r.config.GRPC.TTL)
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
	repo, err := r.repoFactory.CreateRepo(r.config.Coordinator)
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

	// finally shutdown rpc server
	if r.server != nil {
		r.log.Info("stopping grpc server")
		r.server.Stop()
	}
	r.log.Info("storage server stop complete")
	r.state = server.Terminated
	return nil
}

// buildServiceDependency builds broker service dependency
func (r *runtime) buildServiceDependency() error {
	sm, err := replication.NewSequenceManager(r.config.Replication.Dir)
	if err != nil {
		return err
	}
	engine, err := tsdb.NewEngine(r.config.Engine)
	if err != nil {
		return err
	}
	srv := srv{
		storageService:  service.NewStorageService(engine),
		sequenceManager: sm,
	}
	r.srv = srv
	return nil
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
		writer: handler.NewWriter(r.srv.storageService, r.srv.sequenceManager),
		task:   taskHandler.NewTaskHandler(r.config.Query, r.factory.taskServer, dispatcher),
	}

	//TODO add task service ??????
	storage.RegisterWriteServiceServer(r.server.GetServer(), r.handler.writer)
	common.RegisterTaskServiceServer(r.server.GetServer(), r.handler.task)
}

func (r *runtime) monitoring() {
	systemStatMonitorEnabled := r.config.Monitor.SystemReportIntervalInSeconds > 0
	if systemStatMonitorEnabled {
		go monitoring.NewSystemCollector(
			r.ctx,
			r.config.Monitor.SystemReportInterval(),
			r.config.Engine.Dir,
			r.repo,
			constants.GetNodeMonitoringStatPath(r.node.Indicator()),
			models.ActiveNode{
				Version:    r.version,
				Node:       r.node,
				OnlineTime: timeutil.Now(),
			}).Run()
	}

	// todo: @stone1100, how to retrieve the broker port?
	runtimeStatMonitorEnabled := r.config.Monitor.RuntimeReportIntervalInSeconds > 0
	if runtimeStatMonitorEnabled {
		go monitoring.NewRunTimeCollector(
			r.ctx,
			fmt.Sprintf("http://localhost:%d/", r.config.GRPC),
			r.config.Monitor.RuntimeReportInterval(),
			map[string]string{"role": "broker", "version": r.version},
		)
	}
}

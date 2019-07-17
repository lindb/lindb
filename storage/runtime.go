package storage

import (
	"context"
	"fmt"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/coordinator/discovery"
	task "github.com/eleme/lindb/coordinator/storage"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/server"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/storage"
	"github.com/eleme/lindb/service"
	"github.com/eleme/lindb/storage/handler"
)

const (
	storageCfgName = "storage.toml"
	// DefaultStorageCfgFile defines storage default config file path
	DefaultStorageCfgFile = "./" + storageCfgName
)

// srv represents all dependency services
type srv struct {
	storageService service.StorageService
}

// rpcHandler represents all dependency rpc handlers
type rpcHandler struct {
	writer *handler.Writer
}

// runtime represents storage runtime dependency
type runtime struct {
	state   server.State
	cfgPath string
	config  config.Storage

	ctx    context.Context
	cancel context.CancelFunc

	node         models.Node
	server       rpc.TCPServer
	repo         state.Repository
	registry     discovery.Registry
	taskExecutor *task.TaskExecutor
	srv          srv

	log *logger.Logger
}

// NewStorageRuntime creates storage runtime
func NewStorageRuntime(cfgPath string) server.Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		state:   server.New,
		cfgPath: cfgPath,
		ctx:     ctx,
		cancel:  cancel,

		log: logger.GetLogger("storage/runtime"),
	}
}

// Run runs storage server
func (r *runtime) Run() error {
	if r.cfgPath == "" {
		r.cfgPath = DefaultStorageCfgFile
	}
	if !util.Exist(r.cfgPath) {
		r.state = server.Failed
		return fmt.Errorf("config file doesn't exist, see how to initialize the config by `lind storage -h`")
	}
	r.config = config.Storage{}
	if err := util.DecodeToml(r.cfgPath, &r.config); err != nil {
		r.state = server.Failed
		return fmt.Errorf("decode config file error:%s", err)
	}

	ip, err := util.GetHostIP()
	if err != nil {
		r.state = server.Failed
		return fmt.Errorf("cannot get server ip address, error:%s", err)
	}

	// build service dependency for storage server
	r.buildServiceDependency()

	r.node = models.Node{IP: ip, Port: r.config.Server.Port}
	// start tcp server
	r.startTCPServer()

	// start state repo
	if err := r.startStateRepo(); err != nil {
		r.state = server.Failed
		return err
	}

	// register storage node info
	//TODO TTL default value???
	r.registry = discovery.NewRegistry(r.repo, constants.ActiveNodesPath, r.config.Server.TTL)
	if err := r.registry.Register(r.node); err != nil {
		return fmt.Errorf("register storage node error:%s", err)
	}

	r.taskExecutor = task.NewTaskExecutor(r.ctx, &r.node, r.repo, r.srv.storageService)
	r.taskExecutor.Run()

	r.state = server.Running
	return nil
}

// State returns current storage server state
func (r *runtime) State() server.State {
	return r.state
}

// startStateRepo starts state repository
func (r *runtime) startStateRepo() error {
	repo, err := state.NewRepo(r.config.Coordinator)
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
			r.log.Error("unregister storage error", logger.Error(err))
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
func (r *runtime) buildServiceDependency() {
	srv := srv{
		storageService: service.NewStorageService(r.config.Engine),
	}
	r.srv = srv
}

// startTCPServer starts tcp server
func (r *runtime) startTCPServer() {
	r.server = rpc.NewTCPServer(fmt.Sprintf("%s:%d", r.node.IP, r.node.Port))

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
	handlers := rpcHandler{
		writer: handler.NewWriter(r.srv.storageService),
	}

	storage.RegisterWriteServiceServer(r.server.GetServer(), handlers.writer)
}

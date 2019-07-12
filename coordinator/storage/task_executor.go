package storage

import (
	"context"

	"go.uber.org/zap"

	"github.com/eleme/lindb/coordinator/task"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/service"
)

// TaskExecutor represents storage node task executor.
// NOTICE: need implements task processor and register it.
type TaskExecutor struct {
	executor       *task.Executor
	storageService service.StorageService
	repo           state.Repository
	ctx            context.Context

	log *zap.Logger
}

// NewTaskExecutor creates task executor
func NewTaskExecutor(ctx context.Context,
	node *models.Node,
	repo state.Repository,
	storageService service.StorageService) *TaskExecutor {
	executor := task.NewExecutor(ctx, node, repo)

	// register task processor
	executor.Register(newCreateShardProcessor(storageService))
	return &TaskExecutor{
		ctx:            ctx,
		repo:           repo,
		executor:       executor,
		storageService: storageService,
		log:            logger.GetLogger(),
	}
}

// Run runs task executor, watches task assign and runs task process based on task kind in background
func (e *TaskExecutor) Run() {
	//TODO refactor
	go e.executor.Run()
	e.log.Info("task executor started")
}

// Close closes task executor
func (e *TaskExecutor) Close() error {
	return e.executor.Close()
}

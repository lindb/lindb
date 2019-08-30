package storage

import (
	"context"

	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

// TaskExecutor represents storage node task executor.
// NOTICE: need implements task processor and register it.
type TaskExecutor struct {
	executor       *task.Executor
	storageService service.StorageService
	repo           state.Repository
	ctx            context.Context

	log *logger.Logger
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
		log:            logger.GetLogger("coordinator", "StorageTaskExecutor"),
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

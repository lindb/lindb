package storage

import (
	"context"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/service"
)

// databaseFlushProcessor represents flush data of memory database for all shards
type databaseFlushProcessor struct {
	storageService service.StorageService
}

// newCreateShardProcessor returns create shard processor instance
func newDatabaseFlushProcessor(storageService service.StorageService) task.Processor {
	return &databaseFlushProcessor{
		storageService: storageService,
	}
}

func (p *databaseFlushProcessor) Kind() task.Kind             { return constants.FlushDatabase }
func (p *databaseFlushProcessor) RetryCount() int             { return 0 }
func (p *databaseFlushProcessor) RetryBackOff() time.Duration { return 0 }
func (p *databaseFlushProcessor) Concurrency() int            { return 1 }

// Process creates shard for storing time series data
func (p *databaseFlushProcessor) Process(ctx context.Context, task task.Task) error {
	param := models.DatabaseFlushTask{}
	if err := encoding.JSONUnmarshal(task.Params, &param); err != nil {
		return err
	}
	r := p.storageService.FlushDatabase(ctx, param.DatabaseName)
	logger.GetLogger("coordinator", "StorageFlushDBProcessor").
		Info("process flush memory database data task",
			logger.String("params", string(task.Params)),
			logger.Any("result", r),
		)
	return nil
}

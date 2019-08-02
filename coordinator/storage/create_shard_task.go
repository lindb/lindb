package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/service"
)

// createShardProcessor represents create shard when receive task.
// create shard if it not exist
type createShardProcessor struct {
	storageService service.StorageService
}

// newCreateShardProcessor returns create shard processor instance
func newCreateShardProcessor(storageService service.StorageService) task.Processor {
	return &createShardProcessor{
		storageService: storageService,
	}
}

func (p *createShardProcessor) Kind() task.Kind             { return constants.CreateShard }
func (p *createShardProcessor) RetryCount() int             { return 0 }
func (p *createShardProcessor) RetryBackOff() time.Duration { return 0 }
func (p *createShardProcessor) Concurrency() int            { return 1 }

// Process creates shard for storing time series data
func (p *createShardProcessor) Process(ctx context.Context, task task.Task) error {
	param := models.CreateShardTask{}
	if err := json.Unmarshal(task.Params, &param); err != nil {
		return err
	}
	logger.GetLogger("create_shard/task").
		Info("process create shard task", logger.String("params", string(task.Params)))
	if err := p.storageService.CreateShards(param.Database, param.ShardOption, param.ShardIDs...); err != nil {
		return err
	}
	return nil
}

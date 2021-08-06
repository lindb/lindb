// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package storage

import (
	"context"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/tsdb"
)

// databaseFlushProcessor represents flush data of memory database for all shards
type databaseFlushProcessor struct {
	engine tsdb.Engine
}

// newCreateShardProcessor returns create shard processor instance
func newDatabaseFlushProcessor(engine tsdb.Engine) task.Processor {
	return &databaseFlushProcessor{
		engine: engine,
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
	r := p.engine.FlushDatabase(ctx, param.DatabaseName)
	logger.GetLogger("coordinator", "StorageFlushDBProcessor").
		Info("process flush memory database data task",
			logger.String("params", string(task.Params)),
			logger.Any("result", r),
		)
	return nil
}

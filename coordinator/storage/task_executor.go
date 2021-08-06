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

	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/tsdb"
)

// TaskExecutor represents storage node task executor.
// NOTICE: need implements task processor and register it.
type TaskExecutor struct {
	executor *task.Executor
	engine   tsdb.Engine
	repo     state.Repository
	ctx      context.Context

	log *logger.Logger
}

// NewTaskExecutor creates task executor
func NewTaskExecutor(ctx context.Context,
	node *models.Node,
	repo state.Repository,
	engine tsdb.Engine,
) *TaskExecutor {
	executor := task.NewExecutor(ctx, node, repo)
	// register task processor
	executor.Register(newCreateShardProcessor(engine))
	executor.Register(newDatabaseFlushProcessor(engine))
	return &TaskExecutor{
		ctx:      ctx,
		repo:     repo,
		executor: executor,
		engine:   engine,
		log:      logger.GetLogger("coordinator", "StorageTaskExecutor"),
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

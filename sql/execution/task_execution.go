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

package execution

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi/types"
	sqlContext "github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner"
)

type TaskExecutionFactory struct {
	logger logger.Logger
}

func NewTaskExecutionFactory() *TaskExecutionFactory {
	return &TaskExecutionFactory{
		logger: logger.GetLogger("SQL", "execution"),
	}
}

func (fct *TaskExecutionFactory) Create(task *SQLTask) *TaskExecution {
	taskPlanner := planner.NewTaskExecutionPlanner()

	ctx := &sqlContext.TaskContext{
		// FIXME: set context
		Context:    context.WithValue(context.TODO(), constants.ContextKeyCurrentTime, task.currentTime),
		TaskID:     task.id,
		Fragment:   task.fragment,
		Partitions: task.partitions,
	}
	plan := taskPlanner.Plan(ctx, task.fragment.Root)

	return &TaskExecution{
		taskCtx: ctx,
		plan:    plan,
		logger:  fct.logger,
	}
}

type TaskExecution struct {
	taskCtx *sqlContext.TaskContext
	plan    *planner.TaskExecutionPlan

	logger logger.Logger
}

func (exe *TaskExecution) Execute(output chan<- *types.Page) error {
	pipelines := exe.plan.GetPipelines()
	var wait sync.WaitGroup
	wait.Add(len(pipelines))
	errs := make([]error, len(pipelines))
	for i := range pipelines {
		pipeline := pipelines[i]
		go func() {
			defer func() {
				wait.Done()

				if r := recover(); r != nil {
					// FIXME: handle not found
					errs[i] = fmt.Errorf("%v", r)
					exe.logger.Warn("task execution pipeline error", logger.Any("error", r), logger.Stack())
				}
			}()
			pipeline.Run(output)
		}()
	}

	wait.Wait()
	return errors.Join(errs...)
}

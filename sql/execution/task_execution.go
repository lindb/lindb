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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
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

func (fct *TaskExecutionFactory) Create(ctx context.Context, task *SQLTask) *TaskExecution {
	taskPlanner := planner.NewTaskExecutionPlanner()

	taskCtx := &sqlContext.TaskContext{
		Context:      context.WithValue(ctx, constants.ContextKeyCurrentTime, task.CurrentTime),
		TaskID:       task.ID,
		Fragment:     task.Fragment,
		Partitions:   task.Partitions,
		Database:     task.Database,
		OutputStream: task.OutputStream,
	}
	plan := taskPlanner.Plan(taskCtx, task.Fragment.Root)

	return &TaskExecution{
		taskCtx: taskCtx,
		plan:    plan,
		logger:  fct.logger,
	}
}

type TaskExecution struct {
	taskCtx *sqlContext.TaskContext
	plan    *planner.TaskExecutionPlan

	logger logger.Logger
}

func (exe *TaskExecution) Execute(output chan<- arrow.RecordBatch) error {
	pipelines := exe.plan.GetPipelines()
	var wait sync.WaitGroup
	wait.Add(len(pipelines))
	errs := make([]error, len(pipelines))
	for i := range pipelines {
		pipeline := pipelines[i]
		go func() {
			defer func() {
				// recover() MUST be called before wait.Done().
				// If wait.Done() ran first and this was the last goroutine, wait.Wait()
				// would unblock and errors.Join(errs) would read a nil slot — the error
				// would be silently dropped and never reach the client.
				r := recover()
				if r != nil {
					errs[i] = fmt.Errorf("%v", r)
					exe.logger.Error("task execution pipeline error", logger.Any("error", r), logger.Stack())
				}
				wait.Done()
			}()
			// pipeline.Run returns errors from async child operators (e.g. panics
			// caught by RunAsync). Merge with any panic caught above so all errors
			// are surfaced via errors.Join.
			if err := pipeline.Run(output); err != nil && errs[i] == nil {
				errs[i] = err
			}
		}()
	}

	wait.Wait()
	// errs[i] already captures panics from the root operator (set by the
	// goroutine's recover). pipeline.Run also returns child-operator panics
	// (from RunAsync) as an error, which is stored in the same slot so that
	// errors.Join surfaces them all to the caller.
	return errors.Join(errs...)
}

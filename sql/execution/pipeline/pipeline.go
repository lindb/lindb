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

package pipeline

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"

	sqlContext "github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/execution/operator"
)

type Pipeline struct {
	taskCtx *sqlContext.TaskContext
	root    operator.Operator
}

func NewPipeline(taskCtx *sqlContext.TaskContext, root operator.Operator) *Pipeline {
	return &Pipeline{
		taskCtx: taskCtx,
		root:    root,
	}
}

// Run executes the pipeline and blocks until the root operator finishes.
// It returns the first error (panic) captured from any async child operator.
func (p *Pipeline) Run(output chan<- arrow.RecordBatch) error {
	fmt.Printf("run pipeline, root=>\n%s\n", renderText(p.root))

	// errCh collects panics from async child operators; buffer = number of async ops
	// so no goroutine ever blocks on send.
	errCh := make(chan error, p.countAsyncOps(p.root))

	p.execOperator(p.taskCtx.Context, p.root, true, output, errCh)

	// Run the root operator synchronously. Any panic here is caught by
	// TaskExecution.Execute's own recover() and returned directly.
	p.root.Run(p.taskCtx.Context, output)

	// Collect the first child-operator error (if any).
	close(errCh)
	var firstErr error
	for err := range errCh {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// countAsyncOps returns the number of non-root operators that will be launched
// via RunAsync, used to size errCh so sends never block.
func (p *Pipeline) countAsyncOps(op operator.Operator) int {
	count := 0
	for _, child := range op.Children() {
		count += 1 + p.countAsyncOps(child)
	}
	return count
}

func (p *Pipeline) execOperator(ctx context.Context,
	op operator.Operator,
	exclude bool,
	output chan<- arrow.RecordBatch,
	errCh chan<- error,
) {
	children := op.Children()
	inbounds := op.GetInbounds()
	for i, child := range children {
		p.execOperator(ctx, child, false, inbounds[i], errCh)
	}

	if !exclude {
		operator.RunAsync(ctx, op, output, errCh)
	}
}

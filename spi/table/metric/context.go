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

package metric

import (
	"context"

	"go.uber.org/atomic"
)

type ExecutionContext struct {
	ctx context.Context

	inflight atomic.Int32 // task inflight count

	completed chan struct{}
}

func NewExecutionContext(ctx context.Context) *ExecutionContext {
	return &ExecutionContext{
		ctx:       ctx,
		inflight:  *atomic.NewInt32(0),
		completed: make(chan struct{}, 1),
	}
}

func (c *ExecutionContext) GetTaskContext() context.Context {
	c.inflight.Inc()
	return c.ctx
}

func (c *ExecutionContext) CompleteTask() {
	val := c.inflight.Dec()
	if val == 0 {
		c.completed <- struct{}{}
	}
}

func (c *ExecutionContext) Waiting() {
	select {
	case <-c.completed:
	case <-c.ctx.Done():
	}
}

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

package context

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./task_context.go -destination=./task_context_mock.go -package=context

// TaskContext represents the task context for distribution query and computing.
type TaskContext interface {
	// Context returns the context.
	Context() context.Context
	// MakePlan executes search logic in compute level.
	// 1) get metadata based on params
	// 2) build execute plan
	MakePlan() error
	// HandleResponse handles task response.
	HandleResponse(resp *protoCommonV1.TaskResponse, fromNode string)
	// SendRequest sends the task request to target node.
	SendRequest(targetNodeID string, req *protoCommonV1.TaskRequest) error
	// GetRequests returns the request list which send to target node.
	GetRequests() map[string]*protoCommonV1.TaskRequest
	// Complete completes the task with error(if execute failure).
	Complete(err error)
	// WaitResponse waits task complete and returns the response.
	WaitResponse() (any, error)
	// SetTracker sets stage tracker.
	SetTracker(stageTracker *tracker.StageTracker)
}

// baseTaskContext implements TaskContext interface, implements some common logic.
type baseTaskContext struct {
	ctx          context.Context
	requests     map[string]*protoCommonV1.TaskRequest
	state        map[string]models.TaskState
	sendTime     time.Time
	sent         int
	transportMgr rpc.TransportManager

	stageTracker *tracker.StageTracker

	// handle response
	doneCh        chan struct{}
	expectResults int
	completed     atomic.Bool
	err           error
	mutex         sync.Mutex
	// tolerantNotFounds keeps the number of how many not found errors can be returned
	// if all nodes return not-found errors, it will be treated as a error
	// other error will be returned immediately
	tolerantNotFounds int32
}

// newBaseTaskContext creates the base task context.
func newBaseTaskContext(ctx context.Context, transportMgr rpc.TransportManager) baseTaskContext {
	return baseTaskContext{
		ctx:          ctx,
		transportMgr: transportMgr,
		doneCh:       make(chan struct{}),
		requests:     make(map[string]*protoCommonV1.TaskRequest),
		state:        make(map[string]models.TaskState),
	}
}

// Context returns the context.
func (ctx *baseTaskContext) Context() context.Context {
	return ctx.ctx
}

// SendRequest sends the task request to target node.
func (ctx *baseTaskContext) SendRequest(targetNodeID string, req *protoCommonV1.TaskRequest) error {
	ctx.mutex.Lock()
	if ctx.sent == 0 {
		// track the start time for first request
		ctx.sendTime = time.Now()
		ctx.sent++
	}
	ctx.state[targetNodeID] = models.Send
	ctx.mutex.Unlock()
	return ctx.transportMgr.SendRequest(targetNodeID, req)
}

// GetRequests returns the request list which send to target node.
func (ctx *baseTaskContext) GetRequests() map[string]*protoCommonV1.TaskRequest {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	return ctx.requests
}

// addRequests adds the task requests based on physical plan.
func (ctx *baseTaskContext) addRequests(req *protoCommonV1.TaskRequest, physicalPlan *models.PhysicalPlan) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	for _, target := range physicalPlan.Targets {
		ctx.expectResults++
		ctx.tolerantNotFounds++

		// add all targets then send task request, for fail fast
		ctx.requests[target.Indicator] = req
		ctx.state[target.Indicator] = models.Init
	}
}

// Complete completes the task with error(if execute failure).
func (ctx *baseTaskContext) Complete(err error) {
	ctx.mutex.Lock()
	ctx.err = err
	ctx.mutex.Unlock()

	ctx.tryClose()
}

// SetTracker sets stage tracker.
func (ctx *baseTaskContext) SetTracker(stageTracker *tracker.StageTracker) {
	ctx.stageTracker = stageTracker
}

// tryClose tries to complete the task.
func (ctx *baseTaskContext) tryClose() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	if ctx.expectResults <= 0 || ctx.err != nil {
		if ctx.completed.CompareAndSwap(false, true) {
			ctx.stageTracker.Complete()
			close(ctx.doneCh)
		}
	}
}

// handleTaskState handles task state based on task response.
func (ctx *baseTaskContext) handleTaskState(resp *protoCommonV1.TaskResponse, fromNode string) {
	if resp.Completed {
		ctx.state[fromNode] = models.Complete
	} else {
		ctx.state[fromNode] = models.Receive
	}
}

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
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	trackerpkg "github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

var (
	leafExecuteCtxLogger = logger.GetLogger("Query", "LeafContext")
)

// LeafExecuteContext represents leaf node execution context.
type LeafExecuteContext struct {
	TaskCtx  *flow.TaskContext
	Tracker  *trackerpkg.StageTracker
	LeafNode *models.Leaf

	StorageExecuteCtx *flow.StorageExecuteContext
	Database          tsdb.Database

	ServerFactory rpc.TaskServerFactory
	Req           *protoCommonV1.TaskRequest

	GroupingCtx *LeafGroupingContext
	ReduceCtx   *LeafReduceContext

	completed atomic.Bool
}

// NewLeafExecuteContext creates a LeafExecuteContext instance.
func NewLeafExecuteContext(taskCtx *flow.TaskContext,
	tracker *trackerpkg.StageTracker,
	queryStmt *stmt.Query,
	req *protoCommonV1.TaskRequest,
	serverFactory rpc.TaskServerFactory,
	leafNode *models.Leaf,
	database tsdb.Database,
) *LeafExecuteContext {
	storageExecuteCtx := &flow.StorageExecuteContext{
		TaskCtx:  taskCtx,
		Query:    queryStmt,
		ShardIDs: leafNode.ShardIDs,
	}
	ctx := &LeafExecuteContext{
		TaskCtx:           taskCtx,
		Tracker:           tracker,
		LeafNode:          leafNode,
		StorageExecuteCtx: storageExecuteCtx,
		Database:          database,
		ServerFactory:     serverFactory,
		Req:               req,
	}
	ctx.GroupingCtx = NewLeafGroupingContext(ctx) // for group by query
	ctx.ReduceCtx = NewLeafReduceContext(ctx.StorageExecuteCtx, ctx.GroupingCtx)
	return ctx
}

// waitCollectGroupingTagsCompleted waits collect grouping tag value tasks completed.
func (ctx *LeafExecuteContext) waitCollectGroupingTagsCompleted() (err error) {
	if ctx.StorageExecuteCtx.Query.HasGroupBy() {
		defer func() {
			ctx.Tracker.SetGroupingCollectStageValues(func(stageStats *models.StageStats) {
				stageStats.End = time.Now().UnixNano()
				stageStats.Cost = stageStats.End - stageStats.Start
				if err != nil {
					stageStats.ErrMsg = err.Error()
					stageStats.State = trackerpkg.ErrorState.String()
				}
				if stageStats.ErrMsg == "" {
					stageStats.State = trackerpkg.CompleteState.String()
				}
			})
		}()
	}
	if ctx.StorageExecuteCtx.HasGroupingTagValueIDs() {
		// if it has grouping tag value ids, need wait collect group by tag values completed
		select {
		case <-ctx.TaskCtx.Ctx.Done():
			err = ctx.TaskCtx.Ctx.Err()
			return
		case <-ctx.GroupingCtx.collectGroupingTagsCompleted:
		}
	}
	return
}

// SendResponse sends lead node execute response, if with err sends error msg, else sends result set.
func (ctx *LeafExecuteContext) SendResponse(err error) {
	if ctx.completed.CAS(false, true) {
		defer ctx.StorageExecuteCtx.Release()

		if err != nil {
			// send error msg
			ctx.sendResponse(nil, err)
			return
		}
		// wait collect tasks completed
		if err := ctx.waitCollectGroupingTagsCompleted(); err != nil {
			ctx.sendResponse(nil, err)
			return
		}

		// build result set
		resultSet := ctx.ReduceCtx.BuildResultSet(ctx.LeafNode)
		// complete stats track
		ctx.Tracker.Complete()

		ctx.sendResponse(resultSet, nil)
	}
}

// sendResponse sends result set based on receivers.
func (ctx *LeafExecuteContext) sendResponse(resultData [][]byte, err error) {
	var stats []byte
	var errMsg string
	if ctx.StorageExecuteCtx.Query.Explain {
		stats = encoding.JSONMarshal(ctx.Tracker.GetStats())
	}
	if err != nil {
		errMsg = err.Error()
	}
	// send result to upstream receivers
	for idx, receiver := range ctx.LeafNode.Receivers {
		stream := ctx.ServerFactory.GetStream(receiver.Indicator())
		if stream == nil {
			leafExecuteCtxLogger.Error("unable to get stream for write response, ignore result",
				logger.String("target", receiver.Indicator()))
			break
		}
		var payload []byte
		if resultData != nil {
			payload = resultData[idx]
		}
		resp := &protoCommonV1.TaskResponse{
			TaskID:    ctx.Req.ParentTaskID,
			Type:      protoCommonV1.TaskType_Leaf,
			Completed: true,
			SendTime:  timeutil.NowNano(),
			Payload:   payload,
			Stats:     stats,
			ErrMsg:    errMsg,
		}
		if err0 := stream.Send(resp); err0 != nil {
			leafExecuteCtxLogger.Error("send storage query result, ignore result",
				logger.String("target", receiver.Indicator()), logger.Error(err0))
		}
	}
}

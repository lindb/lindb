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

package query

import (
	"errors"
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/stage"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// for testing
var (
	newExecutePipelineFn = NewExecutePipeline
)

// leafTaskProcessor represents the leaf node's task, the leaf node is always storage node
// 1. receives the task request, and searches the data from time seres engine
// 2. sends the result to the parent node(root or intermediate)
type leafTaskProcessor struct {
	currentNode       models.Node
	currentNodeID     string
	engine            tsdb.Engine
	taskServerFactory rpc.TaskServerFactory

	statistics *metrics.StorageQueryStatistics
	logger     *logger.Logger
}

// NewLeafTaskProcessor creates the leaf task
func NewLeafTaskProcessor(
	currentNode models.Node,
	engine tsdb.Engine,
	taskServerFactory rpc.TaskServerFactory,
) TaskProcessor {
	return &leafTaskProcessor{
		currentNode:       currentNode,
		currentNodeID:     currentNode.Indicator(),
		engine:            engine,
		taskServerFactory: taskServerFactory,
		statistics:        metrics.NewStorageQueryStatistics(),
		logger:            logger.GetLogger("Query", "LeafTaskDispatcher"),
	}
}

// Process dispatches the request to storage engine query processor
func (p *leafTaskProcessor) Process(
	ctx *flow.TaskContext,
	stream protoCommonV1.TaskService_HandleServer,
	req *protoCommonV1.TaskRequest,
) {
	err := p.process(ctx, req)
	if err != nil {
		// if process fail, need send response with err
		if sendError := stream.Send(&protoCommonV1.TaskResponse{
			TaskID:    req.ParentTaskID,
			Type:      protoCommonV1.TaskType_Leaf,
			Completed: true,
			ErrMsg:    err.Error(),
			SendTime:  timeutil.NowNano(),
		}); sendError != nil {
			p.logger.Error("failed to send error message to target stream",
				logger.String("taskID", req.ParentTaskID),
				logger.Error(err),
			)
		}
	}
}

// Process processes the task request, searches the data of metric from time series engine
func (p *leafTaskProcessor) process(
	ctx *flow.TaskContext,
	req *protoCommonV1.TaskRequest,
) error {
	physicalPlan := models.PhysicalPlan{}
	if err := encoding.JSONUnmarshal(req.PhysicalPlan, &physicalPlan); err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPlan, err)
	}

	foundTask := false
	var curLeaf *models.Leaf
	for _, leaf := range physicalPlan.Leaves {
		if leaf.Indicator == p.currentNodeID {
			foundTask = true
			curLeaf = leaf
			break
		}
	}
	if !foundTask {
		p.statistics.OmitRequest.Incr()
		return fmt.Errorf("%w, i: %s am not a leaf node", ErrBadPhysicalPlan, p.currentNodeID)
	}
	db, ok := p.engine.GetDatabase(physicalPlan.Database)
	if !ok {
		p.statistics.OmitRequest.Incr()
		return fmt.Errorf("%w: %s", ErrNoDatabase, physicalPlan.Database)
	}
	stream := p.taskServerFactory.GetStream(curLeaf.Parent)
	if stream == nil {
		p.statistics.OmitRequest.Incr()
		return fmt.Errorf("%w: %s", ErrNoSendStream, curLeaf.Parent)
	}

	switch req.RequestType {
	case protoCommonV1.RequestType_Data:
		if err := p.processDataSearch(ctx, db, req, curLeaf); err != nil {
			p.statistics.MetricQueryFailures.Incr()
			return err
		}
		p.statistics.MetricQuery.Incr()
	case protoCommonV1.RequestType_Metadata:
		if err := p.processMetadataSuggest(ctx, db, curLeaf.ShardIDs, req, stream); err != nil {
			p.statistics.MetaQueryFailures.Incr()
			return err
		}
		p.statistics.MetaQuery.Incr()
	default:
		p.statistics.OmitRequest.Incr()
		return nil
	}
	return nil
}

func (p *leafTaskProcessor) processMetadataSuggest(
	ctx *flow.TaskContext,
	db tsdb.Database,
	shardIDs []models.ShardID,
	req *protoCommonV1.TaskRequest,
	stream protoCommonV1.TaskService_HandleServer,
) error {
	defer ctx.Release()
	var stmtQuery = &stmt.MetricMetadata{}
	if err := stmtQuery.UnmarshalJSON(req.Payload); err != nil {
		return ErrUnmarshalSuggest
	}
	leafExecuteCtx := context.NewLeafMetadataContext(stmtQuery, db, shardIDs)
	pipeline := newExecutePipelineFn(false, func(err error) {
		var errMsg string
		var payload []byte
		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			errMsg = err.Error()
			p.statistics.MetaQueryFailures.Incr()
		} else {
			payload = encoding.JSONMarshal(&models.SuggestResult{Values: leafExecuteCtx.ResultSet})
		}
		// send result to upstream
		if err := stream.Send(&protoCommonV1.TaskResponse{
			Type:      protoCommonV1.TaskType_Leaf,
			TaskID:    req.ParentTaskID,
			Completed: true,
			ErrMsg:    errMsg,
			SendTime:  timeutil.NowNano(),
			Payload:   payload,
		}); err != nil {
			p.logger.Error("failed to send error message to target stream",
				logger.String("taskID", req.ParentTaskID),
				logger.Error(err),
			)
		}
	})
	pipeline.Execute(stage.NewMetadataSuggestStage(leafExecuteCtx))
	return nil
}

func (p *leafTaskProcessor) processDataSearch(
	ctx *flow.TaskContext,
	db tsdb.Database,
	req *protoCommonV1.TaskRequest,
	leafNode *models.Leaf,
) error {
	stmtQuery := stmt.Query{}
	if err := stmtQuery.UnmarshalJSON(req.Payload); err != nil {
		return ErrUnmarshalQuery
	}

	// execute leaf pipeline
	leafExecuteCtx := context.NewLeafExecuteContext(ctx, &stmtQuery, req, p.taskServerFactory, leafNode, db)

	pipeline := newExecutePipelineFn(stmtQuery.Explain, func(err error) {
		leafExecuteCtx.SendResponse(err)
	})
	pipeline.Execute(stage.NewMetadataLookupStage(leafExecuteCtx))
	return nil
}

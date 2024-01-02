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
	"fmt"
	"time"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/tree"
)

// for testing
var (
	execFn                 = exec
	metricMetadataSearchFn = MetricMetadataSearch
)

// intermediateTaskProcessor represents the intermediate node's task, the intermediate node is always broker node.
// 1. created for group by query
// 2. exchanges leaf task
// 3. receives root task's request
type intermediateTaskProcessor struct {
	timeout      time.Duration
	curNode      models.StatelessNode
	stateMgr     broker.StateManager
	taskMgr      TaskManager
	transportMgr rpc.TransportManager

	logger logger.Logger
}

// NewIntermediateTaskProcessor creates a intermediate task processor.
func NewIntermediateTaskProcessor(
	curNode models.StatelessNode,
	timeout time.Duration,
	stateMgr broker.StateManager,
	taskMgr TaskManager,
	transportMgr rpc.TransportManager,
) TaskProcessor {
	return &intermediateTaskProcessor{
		curNode:      curNode,
		timeout:      timeout,
		stateMgr:     stateMgr,
		taskMgr:      taskMgr,
		transportMgr: transportMgr,
		logger:       logger.GetLogger("Query", "IntermediateTaskProcessor"),
	}
}

// Process processes the intermediate task request.
// if current node is only receive task response need ignore search execute.
func (p *intermediateTaskProcessor) Process(ctx *flow.TaskContext,
	stream protoCommonV1.TaskService_HandleServer, req *protoCommonV1.TaskRequest) error {
	physicalPlan := &models.PhysicalPlan{}
	if err := encoding.JSONUnmarshal(req.PhysicalPlan, physicalPlan); err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPlan, err)
	}
	foundTask := false
	var curTarget *models.Target
	for _, target := range physicalPlan.Targets {
		if target.Indicator == p.curNode.Indicator() {
			foundTask = true
			curTarget = target
			break
		}
	}
	if !foundTask {
		return fmt.Errorf("%w, i: %s am not a target node", ErrBadPhysicalPlan, p.curNode.Indicator())
	}
	if curTarget.ReceiveOnly {
		// if target is receive node, do nothing
		return nil
	}
	switch req.RequestType {
	case protoCommonV1.RequestType_Metadata:
		return p.processMetadataSearch(ctx, stream, req, physicalPlan)
	case protoCommonV1.RequestType_Data:
		return p.processDataSearch(ctx, stream, req, physicalPlan)
	}
	return nil
}

// processDataSearch executes metric data search.
func (p *intermediateTaskProcessor) processDataSearch(
	ctx *flow.TaskContext, stream protoCommonV1.TaskService_HandleServer,
	req *protoCommonV1.TaskRequest,
	physicalPlan *models.PhysicalPlan,
) error {
	var stmtQuery = &tree.Query1{}
	if err := stmtQuery.UnmarshalJSON(req.Payload); err != nil {
		return ErrUnmarshalQuery
	}
	// use intermediate task's targets as leaf's receivers
	var receivers []string
	for _, target := range physicalPlan.Targets {
		receivers = append(receivers, target.Indicator)
	}
	rs, err := execFn(
		context.NewIntermediateMetricContext(ctx.Ctx,
			p.transportMgr, p.stateMgr, req, p.curNode,
			physicalPlan, stmtQuery,
			receivers),
		&models.Request{
			DB: physicalPlan.Database,
		}, &SearchMgr{
			Timeout:      p.timeout,
			RequestID:    req.RequestID,
			CurNode:      p.curNode,
			Choose:       p.stateMgr,
			TaskMgr:      p.taskMgr,
			TransportMgr: p.transportMgr,
		})
	if err != nil {
		return err
	}
	// send result to upstream
	p.sendResponse(stream, req, rs.(*protoCommonV1.TaskResponse))
	return nil
}

// processMetadataSearch executes metric metadata search.
func (p *intermediateTaskProcessor) processMetadataSearch(
	ctx *flow.TaskContext, stream protoCommonV1.TaskService_HandleServer,
	req *protoCommonV1.TaskRequest,
	physicalPlan *models.PhysicalPlan,
) error {
	var stmtQuery = &tree.MetricMetadata{}
	if err := stmtQuery.UnmarshalJSON(req.Payload); err != nil {
		return ErrUnmarshalSuggest
	}
	rs, err := metricMetadataSearchFn(ctx.Ctx, &models.ExecuteParam{
		Database: physicalPlan.Database,
	}, stmtQuery, &SearchMgr{
		Timeout:      p.timeout,
		RequestID:    req.RequestID,
		CurNode:      p.curNode,
		Choose:       p.stateMgr,
		TaskMgr:      p.taskMgr,
		TransportMgr: p.transportMgr,
	})
	if err != nil {
		return err
	}
	payload := encoding.JSONMarshal(&models.SuggestResult{Values: rs.([]string)})
	// send result to upstream
	p.sendResponse(stream, req, &protoCommonV1.TaskResponse{
		RequestID: req.RequestID,
		Completed: true,
		SendTime:  timeutil.NowNano(),
		Payload:   payload,
	})
	return nil
}

// sendResponse sends task response to client.
func (p *intermediateTaskProcessor) sendResponse(
	stream protoCommonV1.TaskService_HandleServer,
	req *protoCommonV1.TaskRequest,
	resp *protoCommonV1.TaskResponse,
) {
	if err := stream.Send(resp); err != nil {
		p.logger.Error("failed to send error message to target stream",
			logger.String("requestID", req.RequestID),
			logger.String("RequestType", req.RequestType.String()),
			logger.Error(err),
		)
	}
}

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

package brokerquery

import (
	"context"
	"fmt"
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

// intermediateTaskProcessor represents intermediate task dispatcher for broker
type intermediateTaskProcessor struct {
	currentNode       models.Node
	currentNodeID     string
	taskClientFactory rpc.TaskClientFactory
	taskServerFactory rpc.TaskServerFactory
	taskManager       TaskManager
	logger            *logger.Logger
}

// NewIntermediateTaskProcessor create an intermediate task processor
// 1. only created for group by query
// 2. exchanges leaf task
// 3. receives leaf task's result
func NewIntermediateTaskProcessor(
	currentNode models.Node,
	taskClientFactory rpc.TaskClientFactory,
	taskServerFactory rpc.TaskServerFactory,
	taskManager TaskManager,
) query.TaskProcessor {
	return &intermediateTaskProcessor{
		currentNode:       currentNode,
		currentNodeID:     currentNode.Indicator(),
		taskClientFactory: taskClientFactory,
		taskServerFactory: taskServerFactory,
		taskManager:       taskManager,
		logger:            logger.GetLogger("query", "IntermediateTaskProcessor"),
	}
}

// Process dispatches the request to distribution query processor, merges the results
func (p *intermediateTaskProcessor) Process(
	ctx context.Context,
	stream protoCommonV1.TaskService_HandleServer,
	req *protoCommonV1.TaskRequest,
) {
	var err error
	if req.RequestType != protoCommonV1.RequestType_Data {
		err = query.ErrOnlySupportIntermediateTask
	}
	switch req.Type {
	case protoCommonV1.TaskType_UNKNOWN, protoCommonV1.TaskType_Leaf:
		err = query.ErrOnlySupportIntermediateTask
	}
	if err != nil {
		goto ErrToRoot
	}
	if err = p.processIntermediateTask(ctx, req); err == nil {
		return
	}

ErrToRoot:
	if streamErr := stream.Send(&protoCommonV1.TaskResponse{
		TaskID:    req.ParentTaskID,
		Completed: true,
		ErrMsg:    err.Error(),
		SendTime:  timeutil.NowNano(),
	}); streamErr != nil {
		p.logger.Error("failed to send error message to target stream",
			logger.String("taskID", req.ParentTaskID),
			logger.Error(streamErr),
		)
	}
}

// processIntermediateTask processes the task request, sends task request to leaf nodes based on physical plan,
// and tracks the task state
func (p *intermediateTaskProcessor) processIntermediateTask(ctx context.Context, req *protoCommonV1.TaskRequest) error {
	startTime := time.Now()
	stmtQuery := stmt.Query{}
	if err := stmtQuery.UnmarshalJSON(req.Payload); err != nil {
		return query.ErrUnmarshalQuery
	}
	physicalPlan, intermediate, err := p.decodePhysicalPlan(req)
	if err != nil {
		return err
	}
	eventCh := p.taskManager.SubmitIntermediateMetricTask(
		physicalPlan,
		&stmtQuery,
		req.ParentTaskID,
	)
	select {
	case event, ok := <-eventCh:
		if !ok {
			return fmt.Errorf("missing response from sent tasks")
		}
		if event.Err != nil {
			return event.Err
		}
		if event.Stats != nil {
			event.Stats.WaitCost = ltoml.Duration(time.Since(startTime))
		}
		taskResponse := p.makeTaskResponse(req, event)
		return p.taskManager.SendResponse(intermediate.Parent, taskResponse)
	case <-ctx.Done():
		// ignore timeout case, as the caller is already timed out
		return nil
	}
}

func (p *intermediateTaskProcessor) decodePhysicalPlan(
	req *protoCommonV1.TaskRequest,
) (
	physicalPlan *models.PhysicalPlan,
	intermediate *models.Intermediate,
	err error,
) {
	physicalPlan = new(models.PhysicalPlan)
	if err := encoding.JSONUnmarshal(req.PhysicalPlan, physicalPlan); err != nil {
		return nil, nil, query.ErrUnmarshalPlan
	}

	var whoAmI *models.Intermediate
	for _, intermediate := range physicalPlan.Intermediates {
		intermediate := intermediate
		if intermediate.Indicator == p.currentNodeID {
			whoAmI = &intermediate
			break
		}
	}
	if whoAmI == nil {
		return nil, nil,
			fmt.Errorf("%w, i: %s amd not a intermediate node",
				query.ErrBadPhysicalPlan, p.currentNode.Indicator())
	}
	return physicalPlan, whoAmI, nil
}

func (p *intermediateTaskProcessor) makeTaskResponse(
	req *protoCommonV1.TaskRequest,
	event *series.TimeSeriesEvent,
) *protoCommonV1.TaskResponse {
	var stats []byte
	if event.Stats != nil {
		stats = encoding.JSONMarshal(event.Stats)
	}
	var timeSeriesList []*protoCommonV1.TimeSeries
	for _, itr := range event.SeriesList {
		fields := make(map[string][]byte)
		for itr.HasNext() {
			fieldItr := itr.Next()
			data, err := fieldItr.MarshalBinary()
			if err != nil || len(data) == 0 {
				continue
			}
			fields[string(fieldItr.FieldName())] = data
		}
		if len(fields) > 0 {
			// always have group by
			timeSeriesList = append(timeSeriesList, &protoCommonV1.TimeSeries{
				Tags:   itr.Tags(),
				Fields: fields,
			})
		}
	}

	var aggregatorSpecs []*protoCommonV1.AggregatorSpec
	for _, spec := range event.AggregatorSpecs {
		aggregatorSpecs = append(aggregatorSpecs, spec)
	}
	seriesList := protoCommonV1.TimeSeriesList{
		TimeSeriesList: timeSeriesList,
		FieldAggSpecs:  aggregatorSpecs,
	}
	data, _ := seriesList.Marshal()
	return &protoCommonV1.TaskResponse{
		TaskID:    req.ParentTaskID,
		Type:      protoCommonV1.TaskType_Intermediate,
		Completed: true,
		SendTime:  timeutil.NowNano(),
		Stats:     stats,
		Payload:   data,
	}
}

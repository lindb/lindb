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
	"time"

	"github.com/lindb/common/pkg/encoding"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/tree"
)

// IntermediateMetricContext represents intermediate metric data search context.
type IntermediateMetricContext struct {
	MetricContext
	req             *protoCommonV1.TaskRequest
	stateMgr        broker.StateManager
	rawPhysicalPlan *models.PhysicalPlan
	statement       *tree.Query1
	currentNode     models.StatelessNode
	receivers       []string

	responseCh chan *protoCommonV1.TaskResponse
}

// NewIntermediateMetricContext creates intermediate metric data search context.
func NewIntermediateMetricContext(ctx context.Context,
	transportMgr rpc.TransportManager, stateMgr broker.StateManager,
	req *protoCommonV1.TaskRequest, curNode models.StatelessNode,
	physicalPlan *models.PhysicalPlan, statement *tree.Query1, receivers []string,
) *IntermediateMetricContext {
	return &IntermediateMetricContext{
		MetricContext:   newMetricContext(ctx, transportMgr),
		stateMgr:        stateMgr,
		req:             req,
		rawPhysicalPlan: physicalPlan,
		statement:       statement,
		currentNode:     curNode,
		receivers:       receivers,
		responseCh:      make(chan *protoCommonV1.TaskResponse),
	}
}

// WaitResponse waits the task completed, then returns the result set.
func (ctx *IntermediateMetricContext) WaitResponse() (any, error) {
	err := ctx.waitResponse()
	if err != nil {
		return nil, err
	}
	return ctx.makeTaskResponse(), nil
}

// MakePlan makes the metric data physical plan.
func (ctx *IntermediateMetricContext) MakePlan() error {
	database := ctx.rawPhysicalPlan.Database
	// TODO: root=intermediate node
	physicalPlans, err := ctx.stateMgr.Choose(database, 1)
	if err != nil {
		return err
	}
	if len(physicalPlans) == 0 {
		return constants.ErrReplicaNotFound
	}
	databaseCfg, ok := ctx.stateMgr.GetDatabaseCfg(database)
	if !ok {
		return constants.ErrDatabaseNotFound
	}

	calcTimeRangeAndInterval(ctx.statement, databaseCfg)

	payload, _ := ctx.statement.MarshalJSON()
	for _, physicalPlan := range physicalPlans {
		for _, receiver := range ctx.receivers {
			physicalPlan.AddReceiver(receiver)
		}
		if err := physicalPlan.Validate(); err != nil {
			return err
		}
		ctx.addRequests(
			&protoCommonV1.TaskRequest{
				RequestID:    ctx.req.RequestID,
				RequestType:  protoCommonV1.RequestType_Data,
				PhysicalPlan: encoding.JSONMarshal(physicalPlan),
				Payload:      payload,
			}, physicalPlan)
	}
	return nil
}

// makeTaskResponse builds task response.
func (ctx *IntermediateMetricContext) makeTaskResponse() *protoCommonV1.TaskResponse {
	var stats []byte
	if ctx.stats != nil {
		end := time.Now()
		ctx.stats.End = end.UnixNano()
		ctx.stats.TotalCost = end.Sub(ctx.startTime).Nanoseconds()
		stats = encoding.JSONMarshal(ctx.stats)
	}
	var timeSeriesList []*protoCommonV1.TimeSeries
	if ctx.groupAgg != nil {
		groupIts := ctx.groupAgg.ResultSet()
		for _, itr := range groupIts {
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
	}

	var aggregatorSpecs []*protoCommonV1.AggregatorSpec
	for _, spec := range ctx.aggregatorSpecs {
		aggregatorSpecs = append(aggregatorSpecs, spec)
	}
	seriesList := protoCommonV1.TimeSeriesList{
		Start:          ctx.timeRange.Start,
		End:            ctx.timeRange.End,
		Interval:       ctx.interval,
		TimeSeriesList: timeSeriesList,
		FieldAggSpecs:  aggregatorSpecs,
	}
	data, _ := seriesList.Marshal()
	return &protoCommonV1.TaskResponse{
		RequestID:   ctx.req.RequestID,
		RequestType: ctx.req.RequestType,
		Completed:   true,
		SendTime:    commontimeutil.NowNano(),
		Stats:       stats,
		Payload:     data,
	}
}

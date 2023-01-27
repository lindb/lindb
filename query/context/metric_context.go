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
	"errors"
	"strings"
	"time"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// MetricContext represents metric data search context.
type MetricContext struct {
	baseTaskContext

	groupAgg aggregation.GroupingAggregator
	stats    *models.NodeStats
	// field name -> aggregator spec
	// we will use it during intermediate tasks
	aggregatorSpecs map[string]*protoCommonV1.AggregatorSpec
	timeRange       timeutil.TimeRange
	interval        int64
	startTime       time.Time // task start time
}

// newMetricContext creates metric data search context.
func newMetricContext(ctx context.Context, transportMgr rpc.TransportManager) MetricContext {
	return MetricContext{
		baseTaskContext: newBaseTaskContext(ctx, transportMgr),
		aggregatorSpecs: make(map[string]*protoCommonV1.AggregatorSpec),
		startTime:       time.Now(),
	}
}

// HandleResponse handles metric data search task response.
func (ctx *MetricContext) HandleResponse(resp *protoCommonV1.TaskResponse, fromNode string) {
	ctx.handleResponse(resp, fromNode)
	ctx.tryClose()
}

// waitResponse waits metric data search task completed.
func (ctx *MetricContext) waitResponse() error {
	select {
	case <-ctx.doneCh:
		if ctx.err != nil {
			return ctx.err
		}
		return nil
	case <-ctx.ctx.Done():
		return constants.ErrTimeout
	}
}

// handleResponse hanles task response.
func (ctx *MetricContext) handleResponse(resp *protoCommonV1.TaskResponse, fromNode string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.handleTaskState(resp, fromNode)
	ctx.expectResults--

	ctx.handleStats(resp, fromNode)

	ignoreResponse, err := ctx.checkError(resp.ErrMsg)
	if err != nil {
		ctx.err = err
		return
	}
	// partial not-found errors
	if ignoreResponse {
		return
	}

	tsList := &protoCommonV1.TimeSeriesList{}
	if err := tsList.Unmarshal(resp.Payload); err != nil {
		ctx.err = err
		return
	}

	if len(tsList.FieldAggSpecs) == 0 {
		// if it gets empty aggregator spec(empty response), need ignore response.
		// if not ignore, will build empty group aggregator, and cannot aggregate real response data.
		return
	}
	ctx.timeRange = timeutil.TimeRange{
		Start: tsList.Start,
		End:   tsList.End,
	}
	ctx.interval = tsList.Interval

	for _, spec := range tsList.FieldAggSpecs {
		ctx.aggregatorSpecs[spec.FieldName] = spec
	}

	if ctx.groupAgg == nil {
		AggregatorSpecs := make(aggregation.AggregatorSpecs, len(tsList.FieldAggSpecs))
		for idx, aggSpec := range tsList.FieldAggSpecs {
			AggregatorSpecs[idx] = aggregation.NewAggregatorSpec(
				field.Name(aggSpec.FieldName),
				field.Type(aggSpec.FieldType),
			)
			for _, funcType := range aggSpec.FuncTypeList {
				AggregatorSpecs[idx].AddFunctionType(function.FuncType(funcType))
			}
		}
		ctx.groupAgg = newGroupingAgg(
			timeutil.Interval(ctx.interval),
			1, // interval ratio is 1 when do merge result.
			ctx.timeRange,
			AggregatorSpecs,
		)
	}

	for _, ts := range tsList.TimeSeriesList {
		// if no field data, ignore this response
		if len(ts.Fields) == 0 {
			continue
		}
		fields := make(map[field.Name][]byte)
		for k, v := range ts.Fields {
			fields[field.Name(k)] = v
		}
		ctx.groupAgg.Aggregate(series.NewGroupedIterator(ts.Tags, fields))
	}
}

// checkError checks if it has an error should be returned.
// node of the cluster may return not found error,
// ignoreResponse=true symbols that the response should be ignored
func (ctx *MetricContext) checkError(errMsg string) (ignoreResponse bool, err error) {
	if errMsg == "" {
		return false, nil
	}
	// real error
	if !strings.Contains(errMsg, "not found") {
		goto ReturnError
	}
	ctx.tolerantNotFounds--
	// not found, but there may be still more responses not reached
	if ctx.tolerantNotFounds > 0 {
		return true, nil
	}
	// fallthrough, all node returns not found errors
ReturnError:
	return true, errors.New(errMsg)
}

// handleStats handles the node stats of query task.
func (ctx *MetricContext) handleStats(resp *protoCommonV1.TaskResponse, fromNode string) {
	if len(resp.Stats) == 0 {
		return
	}
	// if has query stats, need merge task query stats
	if ctx.stats == nil {
		ctx.stats = &models.NodeStats{}
		ctx.stats.Stages = ctx.stageTracker.GetStages()
		ctx.stats.Start = ctx.startTime.UnixNano()
		ctx.stats.WaitEnd = time.Now().UnixNano()
		ctx.stats.WaitStart = ctx.sendTime.UnixNano()
		ctx.stats.WaitCost = ctx.stats.WaitEnd - ctx.stats.WaitStart
	}
	nodeStats := &models.NodeStats{}
	_ = encoding.JSONUnmarshal(resp.Stats, nodeStats)
	nodeStats.Node = fromNode
	nodeStats.NetPayload = int64(len(resp.Stats) + len(resp.Payload))
	ctx.stats.Children = append(ctx.stats.Children, nodeStats)
}

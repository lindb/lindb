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
	"sync"

	"github.com/cespare/xxhash/v2"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
)

// LeafReduceContext represents reduce the result after down sampling aggregate.
type LeafReduceContext struct {
	storageExecuteCtx *flow.StorageExecuteContext
	leafGroupingCtx   *LeafGroupingContext
	reduceAgg         aggregation.GroupingAggregator
	lock              sync.Mutex
}

// NewLeafReduceContext creates a LeafReduceContext instance.
func NewLeafReduceContext(storageExecuteCtx *flow.StorageExecuteContext, leafGroupingCtx *LeafGroupingContext) *LeafReduceContext {
	return &LeafReduceContext{
		storageExecuteCtx: storageExecuteCtx,
		leafGroupingCtx:   leafGroupingCtx,
	}
}

// Reduce reduces the down sampling aggregator's result.
func (ctx *LeafReduceContext) Reduce(it series.GroupedIterator) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if ctx.reduceAgg == nil {
		storageExecuteCtx := ctx.storageExecuteCtx
		ctx.reduceAgg = aggregation.NewGroupingAggregator(storageExecuteCtx.Query.Interval,
			storageExecuteCtx.Query.IntervalRatio, storageExecuteCtx.Query.TimeRange, storageExecuteCtx.AggregatorSpecs)
	}

	ctx.reduceAgg.Aggregate(it)
}

// BuildResultSet returns the result set from reduce aggregator based on receivers.
func (ctx *LeafReduceContext) BuildResultSet(leafNode *models.Leaf) [][]byte {
	aggSpecs := ctx.storageExecuteCtx.AggregatorSpecs
	aggregatorSpecs := make([]*protoCommonV1.AggregatorSpec, len(aggSpecs))
	for idx, spec := range aggSpecs {
		aggregatorSpecs[idx] = &protoCommonV1.AggregatorSpec{
			FieldName: string(spec.FieldName()),
			FieldType: uint32(spec.GetFieldType()),
		}
		for _, funcType := range spec.Functions() {
			aggregatorSpecs[idx].FuncTypeList = append(aggregatorSpecs[idx].FuncTypeList, uint32(funcType))
		}
	}
	numOfReceivers := len(leafNode.Receivers)
	resultSet := make([][]byte, numOfReceivers)
	timeSeriesList := ctx.makeTimeSeriesList()
	// root -> leaf task, return the raw total series
	if numOfReceivers == 1 {
		leaf2RootSeries := protoCommonV1.TimeSeriesList{
			TimeSeriesList: timeSeriesList,
			FieldAggSpecs:  aggregatorSpecs,
		}
		leaf2RootSeriesPayload, _ := leaf2RootSeries.Marshal()
		resultSet[0] = leaf2RootSeriesPayload
	} else {
		// during intermediate task, time series will be grouped by hash
		// and send to multi intermediate receiver
		// hash mod -> series list
		var timeSeriesHashGroups = make([][]*protoCommonV1.TimeSeries, numOfReceivers)
		for _, ts := range timeSeriesList {
			h := xxhash.Sum64String(ts.Tags)
			index := int(h % uint64(numOfReceivers))
			timeSeriesHashGroups[index] = append(timeSeriesHashGroups[index], ts)
		}
		for idx, timeSeriesHashGroup := range timeSeriesHashGroups {
			leaf2IntermediateSeries := protoCommonV1.TimeSeriesList{
				TimeSeriesList: timeSeriesHashGroup,
				FieldAggSpecs:  aggregatorSpecs,
			}
			leaf2IntermediatePayload, _ := leaf2IntermediateSeries.Marshal()
			resultSet[idx] = leaf2IntermediatePayload
		}
	}
	return resultSet
}

// makeTimeSeriesList returns the time series data from reduce aggregator.
func (ctx *LeafReduceContext) makeTimeSeriesList() []*protoCommonV1.TimeSeries {
	if ctx.reduceAgg == nil {
		// if no data found or do aggregate
		return nil
	}

	hasGroupBy := ctx.storageExecuteCtx.Query.HasGroupBy()
	// 1. get reduce aggregator result set
	groupedSeriesList := ctx.reduceAgg.ResultSet()
	// 2. build rpc response data
	var timeSeriesList []*protoCommonV1.TimeSeries
	for _, groupedSeriesItr := range groupedSeriesList {
		fields := make(map[string][]byte)
		for groupedSeriesItr.HasNext() {
			seriesItr := groupedSeriesItr.Next()
			data, err := seriesItr.MarshalBinary()
			if err != nil || len(data) == 0 {
				if err != nil {
					leafExecuteCtxLogger.Error("marshal series data, ignore it.", logger.Error(err))
				}
				continue
			}
			fields[string(seriesItr.FieldName())] = data
		}

		if len(fields) > 0 {
			tags := ""
			if hasGroupBy {
				tagValueIDs := groupedSeriesItr.Tags() // returns tag value ids string value under leaf node.
				tags = ctx.leafGroupingCtx.getTagValues(tagValueIDs)
			}
			timeSeriesList = append(timeSeriesList, &protoCommonV1.TimeSeries{
				Tags:   tags,
				Fields: fields,
			})
		}
	}
	return timeSeriesList
}

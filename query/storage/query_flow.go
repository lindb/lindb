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

package storagequery

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"

	xxhash "github.com/cespare/xxhash/v2"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	querypkg "github.com/lindb/lindb/query"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
)

const (
	tagValueNotFound = "tag_value_not_found"
)

var (
	storageQueryFlowLogger = logger.GetLogger("query", "StorageQueryFlow")
)

// storageQueryFlow represents the storage engine query execute flow
type storageQueryFlow struct {
	storageExecuteCtx *flow.StorageExecuteContext
	pendingTasks      map[int32]flow.Stage // pending task ref counter for each stage
	taskIDSeq         atomic.Int32         // task id gen sequence
	executorPool      *tsdb.ExecutorPool
	reduceAgg         aggregation.GroupingAggregator
	leafNode          *models.Leaf
	req               *protoCommonV1.TaskRequest
	ctx               context.Context
	serverFactory     rpc.TaskServerFactory

	aggregatorSpecs []*protoCommonV1.AggregatorSpec

	tagsMap                    map[string]string   // tag value ids => tag values
	tagValuesMap               []map[uint32]string // tag value id=> tag value for each group by tag key
	tagValues                  []string
	waitingForCollectTagValues atomic.Int32
	signal                     chan struct{}

	mux       sync.Mutex
	completed atomic.Bool
}

// NewStorageQueryFlow creates a storage query flow.
func NewStorageQueryFlow(
	storageExecuteCtx *flow.StorageExecuteContext,
	req *protoCommonV1.TaskRequest,
	serverFactory rpc.TaskServerFactory,
	leafNode *models.Leaf,
	executorPool *tsdb.ExecutorPool,
) flow.StorageQueryFlow {
	return &storageQueryFlow{
		ctx:               storageExecuteCtx.TaskCtx.Ctx,
		storageExecuteCtx: storageExecuteCtx,
		req:               req,
		leafNode:          leafNode,
		serverFactory:     serverFactory,
		executorPool:      executorPool,
		pendingTasks:      make(map[int32]flow.Stage),
		signal:            make(chan struct{}),
	}
}

func (qf *storageQueryFlow) Prepare() {
	aggregatorSpecs := qf.storageExecuteCtx.AggregatorSpecs
	qf.reduceAgg = aggregation.NewGroupingAggregator(qf.storageExecuteCtx.Query.Interval,
		qf.storageExecuteCtx.Query.IntervalRatio, qf.storageExecuteCtx.Query.TimeRange, aggregatorSpecs)
	qf.aggregatorSpecs = make([]*protoCommonV1.AggregatorSpec, len(aggregatorSpecs))
	for idx, spec := range aggregatorSpecs {
		qf.aggregatorSpecs[idx] = &protoCommonV1.AggregatorSpec{
			FieldName: string(spec.FieldName()),
			FieldType: uint32(spec.GetFieldType()),
		}
		for _, funcType := range spec.Functions() {
			qf.aggregatorSpecs[idx].FuncTypeList = append(qf.aggregatorSpecs[idx].FuncTypeList, uint32(funcType))
		}
	}

	// for group by
	groupByKenLen := len(qf.storageExecuteCtx.Query.GroupBy)
	if groupByKenLen > 0 {
		qf.tagValuesMap = make([]map[uint32]string, groupByKenLen)
		qf.tagsMap = make(map[string]string)
		qf.tagValues = make([]string, groupByKenLen)
		qf.waitingForCollectTagValues.Add(int32(groupByKenLen))
	}
}

// Complete completes the query flow with error
func (qf *storageQueryFlow) Complete(err error) {
	if qf.completed.CAS(false, true) {
		// if complete with err, need send err msg directly and mark task completed
		qf.sendResponse(nil, err)
	}
}

// Reduce reduces the down sampling aggregator's result.
func (qf *storageQueryFlow) Reduce(it series.GroupedIterator) {
	if qf.completed.Load() {
		storageQueryFlowLogger.Warn("reduce the aggregator data after storage query flow completed")
		return
	}

	qf.mux.Lock()
	defer qf.mux.Unlock()

	qf.reduceAgg.Aggregate(it)
}

// ReduceTagValues reduces the group by tag values
func (qf *storageQueryFlow) ReduceTagValues(tagKeyIndex int, tagValues map[uint32]string) {
	qf.mux.Lock()
	defer qf.mux.Unlock()
	qf.tagValuesMap[tagKeyIndex] = tagValues
	if qf.waitingForCollectTagValues.Sub(1) == 0 {
		close(qf.signal)
	}
}

func (qf *storageQueryFlow) getTagValues(tags string) string {
	qf.mux.Lock()
	defer qf.mux.Unlock()

	tagValues, ok := qf.tagsMap[tags]
	if ok {
		return tagValues
	}
	tagsData := []byte(tags)
	for idx := range qf.tagValues {
		tagValuesForKey := qf.tagValuesMap[idx]
		offset := idx * 4
		tagValueID := binary.LittleEndian.Uint32(tagsData[offset:])
		tagValue, ok := tagValuesForKey[tagValueID]
		if ok {
			qf.tagValues[idx] = tagValue
		} else {
			qf.tagValues[idx] = tagValueNotFound
		}
	}
	tagsOfStr := tag.ConcatTagValues(qf.tagValues)
	qf.tagsMap[tags] = tagsOfStr
	return tagsOfStr
}

func (qf *storageQueryFlow) completeTask(taskID int32) {
	completed := false
	// get complete execute result
	qf.mux.Lock()
	delete(qf.pendingTasks, taskID)
	completed = len(qf.pendingTasks) == 0 || qf.ctx.Err() != nil
	qf.mux.Unlock()

	if !completed || !qf.completed.CAS(false, true) {
		return
	}

	defer qf.storageExecuteCtx.Release()

	resultSet := make([][]byte, len(qf.leafNode.Receivers))
	if qf.reduceAgg != nil {
		if qf.storageExecuteCtx.HasGroupingTagValueIDs() {
			// if it has grouping tag value ids, need wait collect group by tag value complete
			select {
			case <-qf.ctx.Done():
				qf.sendResponse(nil, qf.ctx.Err())
				return
			case <-qf.signal:
			}
		}

		timeSeriesList := qf.makeTimeSeriesList()
		// root -> leaf task, return the raw total series
		if len(qf.leafNode.Receivers) == 1 {
			leaf2RootSeries := protoCommonV1.TimeSeriesList{
				TimeSeriesList: timeSeriesList,
				FieldAggSpecs:  qf.aggregatorSpecs,
			}
			leaf2RootSeriesPayload, _ := leaf2RootSeries.Marshal()
			resultSet[0] = leaf2RootSeriesPayload
		} else {
			// during intermediate task, time series will be grouped by hash
			// and send to multi intermediate receiver
			// hash mod -> series list
			var timeSeriesHashGroups = make([][]*protoCommonV1.TimeSeries, len(qf.leafNode.Receivers))
			for _, ts := range timeSeriesList {
				h := xxhash.Sum64String(ts.Tags)
				index := int(h % uint64(len(qf.leafNode.Receivers)))
				timeSeriesHashGroups[index] = append(timeSeriesHashGroups[index], ts)
			}
			for idx, timeSeriesHashGroup := range timeSeriesHashGroups {
				leaf2IntermediateSeries := protoCommonV1.TimeSeriesList{
					TimeSeriesList: timeSeriesHashGroup,
					FieldAggSpecs:  qf.aggregatorSpecs,
				}
				leaf2IntermediatePayload, _ := leaf2IntermediateSeries.Marshal()
				resultSet[idx] = leaf2IntermediatePayload
			}
		}
	}
	qf.sendResponse(resultSet, nil)
}

func (qf *storageQueryFlow) sendResponse(resultData [][]byte, err error) {
	var stats []byte
	var errMsg string
	if err == nil {
		executeStats := qf.storageExecuteCtx.QueryStats()
		if executeStats != nil {
			stats = encoding.JSONMarshal(executeStats)
		}
	} else {
		errMsg = err.Error()
	}
	// send result to upstream receivers
	for idx, receiver := range qf.leafNode.Receivers {
		stream := qf.serverFactory.GetStream(receiver.Indicator())
		if stream == nil {
			storageQueryFlowLogger.Error("unable to get stream for write response",
				logger.String("target", receiver.Indicator()))
			qf.Complete(querypkg.ErrNoSendStream)
			break
		}
		var payload []byte
		if resultData != nil {
			payload = resultData[idx]
		}
		resp := &protoCommonV1.TaskResponse{
			TaskID:    qf.req.ParentTaskID,
			Type:      protoCommonV1.TaskType_Leaf,
			Completed: true,
			SendTime:  timeutil.NowNano(),
			Payload:   payload,
			Stats:     stats,
			ErrMsg:    errMsg,
		}
		if err0 := stream.Send(resp); err0 != nil {
			storageQueryFlowLogger.Error("send storage query result", logger.Error(err0))
		}
	}
}

func (qf *storageQueryFlow) makeTimeSeriesList() []*protoCommonV1.TimeSeries {
	hasGroupBy := qf.storageExecuteCtx.Query.HasGroupBy()
	// 1. get reduce aggregator result set
	groupedSeriesList := qf.reduceAgg.ResultSet()
	// 2. build rpc response data
	var timeSeriesList []*protoCommonV1.TimeSeries
	for _, groupedSeriesItr := range groupedSeriesList {
		fields := make(map[string][]byte)
		for groupedSeriesItr.HasNext() {
			seriesItr := groupedSeriesItr.Next()
			data, err := seriesItr.MarshalBinary()
			if err != nil || len(data) == 0 {
				if err != nil {
					storageQueryFlowLogger.Error("marshal seriesItr data", logger.Error(err))
				}
				continue
			}
			fields[string(seriesItr.FieldName())] = data
		}

		if len(fields) > 0 {
			tags := ""
			if hasGroupBy {
				tags = qf.getTagValues(groupedSeriesItr.Tags())
			}
			timeSeriesList = append(timeSeriesList, &protoCommonV1.TimeSeries{
				Tags:   tags,
				Fields: fields,
			})
		}
	}
	return timeSeriesList
}

func (qf *storageQueryFlow) isCompleted() bool {
	if qf.completed.Load() {
		return true
	}

	err := qf.ctx.Err()
	if err != nil {
		// context is done, complete query flow.
		qf.Complete(err)
		return true
	}
	return false
}

// Submit submits an async task when do query pipeline.
func (qf *storageQueryFlow) Submit(stage flow.Stage, task func()) {
	if qf.isCompleted() {
		// query flow is completed, reject new task execute
		return
	}
	var executePool concurrent.Pool
	switch stage {
	case flow.FilteringStage:
		executePool = qf.executorPool.Filtering
	case flow.GroupingStage:
		executePool = qf.executorPool.Grouping
	case flow.ScannerStage:
		executePool = qf.executorPool.Scanner
	}
	if executePool == nil {
		qf.Complete(fmt.Errorf("execute pool not found for stage:%s", stage))
		return
	}
	// 1. retain the task pending count before submit task
	qf.mux.Lock()
	taskID := qf.taskIDSeq.Inc()
	qf.pendingTasks[taskID] = stage
	qf.mux.Unlock()

	executePool.Submit(qf.ctx, concurrent.NewTask(func() {
		defer func() {
			// 3. complete task and dec task pending after task handle
			qf.completeTask(taskID)
		}()
		if !qf.isCompleted() {
			// 2. handle task logic in background goroutine, if it's not completed.
			task()
		}
	}, func(err error) {
		qf.Complete(err)
	}))
}

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
	"errors"
	"sync"

	"github.com/cespare/xxhash/v2"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
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
	storageExecuteCtx StorageExecuteContext
	query             *stmt.Query
	pendingTasks      map[int32]Stage // pending task ref counter for each stage
	taskIDSeq         atomic.Int32    // task id gen sequence
	executorPool      *tsdb.ExecutorPool
	reduceAgg         aggregation.GroupingAggregator
	leafNode          *models.Leaf
	req               *protoCommonV1.TaskRequest
	ctx               context.Context
	serverFactory     rpc.TaskServerFactory

	aggregatorSpecs []*protoCommonV1.AggregatorSpec

	tagsMap      map[string]string   // tag value ids => tag values
	tagValuesMap []map[uint32]string // tag value id=> tag value for each group by tag key
	tagValues    []string
	signal       sync.WaitGroup

	mux       sync.Mutex
	completed atomic.Bool
}

func NewStorageQueryFlow(
	ctx context.Context,
	storageExecuteCtx StorageExecuteContext,
	query *stmt.Query,
	req *protoCommonV1.TaskRequest,
	serverFactory rpc.TaskServerFactory,
	leafNode *models.Leaf,
	executorPool *tsdb.ExecutorPool,
) flow.StorageQueryFlow {
	return &storageQueryFlow{
		ctx:               ctx,
		storageExecuteCtx: storageExecuteCtx,
		query:             query,
		req:               req,
		leafNode:          leafNode,
		serverFactory:     serverFactory,
		executorPool:      executorPool,
		pendingTasks:      make(map[int32]Stage),
	}
}

func (qf *storageQueryFlow) Prepare(
	interval timeutil.Interval,
	intervalRatio int,
	timeRange timeutil.TimeRange,
	aggregatorSpecs aggregation.AggregatorSpecs,
) {
	qf.reduceAgg = aggregation.NewGroupingAggregator(interval, intervalRatio, timeRange, aggregatorSpecs)
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
	groupByKenLen := len(qf.query.GroupBy)
	if groupByKenLen > 0 {
		qf.tagValuesMap = make([]map[uint32]string, groupByKenLen)
		qf.tagsMap = make(map[string]string)
		qf.tagValues = make([]string, groupByKenLen)
		qf.signal.Add(groupByKenLen)
	}
}

// Complete completes the query flow with error
func (qf *storageQueryFlow) Complete(err error) {
	if err != nil && qf.completed.CAS(false, true) {
		// if complete with err, need send err msg directly and mark task completed
		for _, receiver := range qf.leafNode.Receivers {
			stream := qf.serverFactory.GetStream(receiver.Indicator())
			if stream == nil {
				storageQueryFlowLogger.Error("unable to get stream for answering error",
					logger.String("target", receiver.Indicator()))
				continue
			}
			if err := stream.Send(&protoCommonV1.TaskResponse{
				TaskID:    qf.req.ParentTaskID,
				Type:      protoCommonV1.TaskType_Leaf,
				Completed: true,
				ErrMsg:    err.Error(),
			}); err != nil {
				storageQueryFlowLogger.Error("send storage execute result", logger.Error(err))
			}
		}
	}
}

func (qf *storageQueryFlow) Load(task concurrent.Task) {
	qf.execute(Scanner, task)
}

func (qf *storageQueryFlow) Grouping(task concurrent.Task) {
	qf.execute(Grouping, task)
}

func (qf *storageQueryFlow) Filtering(task concurrent.Task) {
	qf.execute(Filtering, task)
}

func (qf *storageQueryFlow) Reduce(_ string, it series.GroupedIterator) {
	if qf.completed.Load() {
		storageQueryFlowLogger.Warn("reduce the aggregator data after storage query flow completed")
		return
	}

	qf.mux.Lock()
	defer qf.mux.Unlock()

	//TODO impl
	qf.reduceAgg.Aggregate(it)
}

// ReduceTagValues reduces the group by tag values
func (qf *storageQueryFlow) ReduceTagValues(tagKeyIndex int, tagValues map[uint32]string) {
	qf.mux.Lock()
	defer qf.mux.Unlock()
	qf.tagValuesMap[tagKeyIndex] = tagValues
	qf.signal.Done()
}

func (qf *storageQueryFlow) getTagValues(tags string) string {
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
	completed = len(qf.pendingTasks) == 0
	qf.mux.Unlock()

	if !completed || !qf.completed.CAS(false, true) {
		return
	}

	hashGroupData := make([][]byte, len(qf.leafNode.Receivers))
	if qf.reduceAgg != nil {
		hasGroupBy := qf.query.HasGroupBy()
		if hasGroupBy {
			qf.signal.Wait() // wait collect group by tag value complete
		}
		timeSeriesList := qf.makeTimeSeriesList()
		// root -> leaf task, return the raw total series
		if len(qf.leafNode.Receivers) == 1 {
			leaf2RootSeries := protoCommonV1.TimeSeriesList{
				TimeSeriesList: timeSeriesList,
				FieldAggSpecs:  qf.aggregatorSpecs,
			}
			leaf2RootSeriesPayload, _ := leaf2RootSeries.Marshal()
			hashGroupData[0] = leaf2RootSeriesPayload
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
				hashGroupData[idx] = leaf2IntermediatePayload
			}
		}
	}
	qf.sendResponse(hashGroupData)
}

func (qf *storageQueryFlow) sendResponse(hashGroupData [][]byte) {
	var stats []byte
	if qf.storageExecuteCtx.QueryStats() != nil {
		stats = encoding.JSONMarshal(qf.storageExecuteCtx.QueryStats())
	}
	// send result to upstream receivers
	for idx, receiver := range qf.leafNode.Receivers {
		stream := qf.serverFactory.GetStream(receiver.Indicator())
		if stream == nil {
			storageQueryFlowLogger.Error("unable to get stream for write response",
				logger.String("target", receiver.Indicator()))
			qf.Complete(query.ErrNoSendStream)
			break
		}
		if err := stream.Send(&protoCommonV1.TaskResponse{
			TaskID:    qf.req.ParentTaskID,
			Type:      protoCommonV1.TaskType_Leaf,
			Completed: true,
			SendTime:  timeutil.NowNano(),
			Payload:   hashGroupData[idx],
			Stats:     stats,
		}); err != nil {
			storageQueryFlowLogger.Error("send storage query result", logger.Error(err))
		}
	}
}

func (qf *storageQueryFlow) makeTimeSeriesList() []*protoCommonV1.TimeSeries {
	hasGroupBy := qf.query.HasGroupBy()
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

// execute executes the query task by stage
func (qf *storageQueryFlow) execute(stage Stage, task concurrent.Task) {
	if qf.completed.Load() {
		// query flow is completed, reject new task execute
		return
	}
	var executePool concurrent.Pool
	switch stage {
	case Filtering:
		executePool = qf.executorPool.Filtering
	case Grouping:
		executePool = qf.executorPool.Grouping
	case Scanner:
		executePool = qf.executorPool.Scanner
	}
	if executePool != nil {
		// 1. retain the task pending count before submit task
		qf.mux.Lock()
		taskID := qf.taskIDSeq.Inc()
		qf.pendingTasks[taskID] = stage
		qf.mux.Unlock()

		executePool.Submit(func() {
			defer func() {
				// 3. complete task and dec task pending after task handle
				qf.completeTask(taskID)
				var err error
				r := recover()
				if r != nil {
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("unknown error")
					}
					storageQueryFlowLogger.Error("do task fail when execute storage query flow",
						logger.Error(err), logger.Stack())
					qf.Complete(err)
				}
			}()

			// 2. handle task logic in background goroutine
			task()
		})
	}
}

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

package parallel

import (
	"context"
	"encoding/binary"
	"errors"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

const (
	tagValueNotFound = "tag_value_not_found"
)

var storageQueryFlowLogger = logger.GetLogger("parallel", "storageQueryFlow")

// storageQueryFlow represents the storage engine query execute flow
type storageQueryFlow struct {
	storageExecuteCtx StorageExecuteContext
	query             *stmt.Query
	pendingTasks      map[int32]Stage // pending task ref counter for each stage
	taskIDSeq         atomic.Int32    // task id gen sequence
	executorPool      *tsdb.ExecutorPool
	reduceAgg         aggregation.GroupingAggregator
	stream            pb.TaskService_HandleServer
	req               *pb.TaskRequest
	ctx               context.Context

	aggregatorSpecs []*pb.AggregatorSpec

	tagsMap      map[string]string   // tag value ids => tag values
	tagValuesMap []map[uint32]string // tag value id=> tag value for each group by tag key
	tagValues    []string
	signal       sync.WaitGroup

	mux       sync.Mutex
	completed atomic.Bool
}

func NewStorageQueryFlow(ctx context.Context,
	storageExecuteCtx StorageExecuteContext,
	query *stmt.Query,
	req *pb.TaskRequest,
	stream pb.TaskService_HandleServer,
	executorPool *tsdb.ExecutorPool,
) flow.StorageQueryFlow {
	return &storageQueryFlow{
		ctx:               ctx,
		storageExecuteCtx: storageExecuteCtx,
		query:             query,
		req:               req,
		stream:            stream,
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
	qf.aggregatorSpecs = make([]*pb.AggregatorSpec, len(aggregatorSpecs))
	for idx, spec := range aggregatorSpecs {
		qf.aggregatorSpecs[idx] = &pb.AggregatorSpec{
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
		if err := qf.stream.Send(&pb.TaskResponse{
			JobID:     qf.req.JobID,
			TaskID:    qf.req.ParentTaskID,
			Completed: true,
			ErrMsg:    err.Error(),
		}); err != nil {
			storageQueryFlowLogger.Error("send storage execute result", logger.Error(err))
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

func (qf *storageQueryFlow) Reduce(tags string, it series.GroupedIterator) {
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

	if completed && qf.completed.CAS(false, true) {
		// if all tasks of all stages completed
		var data []byte
		if qf.reduceAgg != nil {
			hasGroupBy := qf.query.HasGroupBy()
			if hasGroupBy {
				qf.signal.Wait() // wait collect group by tag value complete
			}
			// 1. get reduce aggregator result set
			groupedSeriesList := qf.reduceAgg.ResultSet()
			// 2. build rpc response data
			var timeSeriesList []*pb.TimeSeries
			for _, ts := range groupedSeriesList {
				fields := make(map[string][]byte)
				for ts.HasNext() {
					fieldIt := ts.Next()
					data, err := fieldIt.MarshalBinary()
					if err != nil || len(data) == 0 {
						if err != nil {
							storageQueryFlowLogger.Error("marshal iterator data", logger.Error(err))
						}
						continue
					}

					fields[string(fieldIt.FieldName())] = data
				}
				if len(fields) > 0 {
					tags := ""
					if hasGroupBy {
						tags = qf.getTagValues(ts.Tags())
					}
					timeSeriesList = append(timeSeriesList, &pb.TimeSeries{
						Tags:   tags,
						Fields: fields,
					})
				}
			}

			seriesList := pb.TimeSeriesList{
				TimeSeriesList: timeSeriesList,
				FieldAggSpecs:  qf.aggregatorSpecs,
			}
			// no error
			data, _ = seriesList.Marshal()
		}

		var stats []byte
		if qf.storageExecuteCtx.QueryStats() != nil {
			stats = encoding.JSONMarshal(qf.storageExecuteCtx.QueryStats())
		}
		// send result to upstream
		if err := qf.stream.Send(&pb.TaskResponse{
			JobID:     qf.req.JobID,
			TaskID:    qf.req.ParentTaskID,
			Completed: true,
			SendTime:  timeutil.NowNano(),
			Payload:   data,
			Stats:     stats,
		}); err != nil {
			storageQueryFlowLogger.Error("send storage execute result", logger.Error(err))
		}
	}
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

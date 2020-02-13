package parallel

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/tsdb"
)

type allocAgg func(aggSpecs aggregation.AggregatorSpecs) aggregation.FieldAggregates

var storageQueryFlowLogger = logger.GetLogger("parallel", "storageQueryFlow")

// storageQueryFlow represents the storage engine query execute flow
type storageQueryFlow struct {
	aggPool      chan aggregation.FieldAggregates // use aggregator for request scope
	pendingTasks map[int32]Stage                  // pending task ref counter for each stage
	taskIDSeq    atomic.Int32                     // task id gen sequence
	executorPool *tsdb.ExecutorPool
	reduceAgg    aggregation.GroupingAggregator
	stream       pb.TaskService_HandleServer
	req          *pb.TaskRequest
	ctx          context.Context
	allocAgg     allocAgg

	queryTimeRange     timeutil.TimeRange
	queryInterval      timeutil.Interval
	queryIntervalRatio int
	downSamplingSpecs  aggregation.AggregatorSpecs

	mux       sync.Mutex
	completed atomic.Bool

	err error
}

func NewStorageQueryFlow(ctx context.Context,
	req *pb.TaskRequest,
	stream pb.TaskService_HandleServer,
	executorPool *tsdb.ExecutorPool,
	queryTimeRange timeutil.TimeRange,
	queryInterval timeutil.Interval,
	queryIntervalRatio int,
) flow.StorageQueryFlow {
	return &storageQueryFlow{
		ctx:                ctx,
		req:                req,
		stream:             stream,
		executorPool:       executorPool,
		pendingTasks:       make(map[int32]Stage),
		queryTimeRange:     queryTimeRange,
		queryInterval:      queryInterval,
		queryIntervalRatio: queryIntervalRatio,
	}
}

func (qf *storageQueryFlow) Prepare(downSamplingSpecs aggregation.AggregatorSpecs) {
	qf.reduceAgg = aggregation.NewGroupingAggregator(qf.queryInterval, qf.queryTimeRange, downSamplingSpecs)
	qf.aggPool = make(chan aggregation.FieldAggregates, 64)
	qf.downSamplingSpecs = downSamplingSpecs
	qf.allocAgg = func(aggSpecs aggregation.AggregatorSpecs) aggregation.FieldAggregates {
		return aggregation.NewFieldAggregates(qf.queryInterval, qf.queryIntervalRatio, qf.queryTimeRange,
			true, aggSpecs)
	}
}

func (qf *storageQueryFlow) GetAggregator() (agg aggregation.FieldAggregates) {
	select {
	case agg = <-qf.aggPool:
		// reuse existing aggregator
	default:
		// create new field aggregator
		agg = qf.allocAgg(qf.downSamplingSpecs)
	}
	return agg
}

// releaseAgg releases the current field aggregator for reuse in request scope
func (qf *storageQueryFlow) releaseAgg(agg aggregation.FieldAggregates) {
	// 1. reset aggregator context
	agg.Reset()

	// 2. try put it back into pool
	select {
	case qf.aggPool <- agg:
		// aggregator went back into pool
	default:
		// aggregator didn't go back into pool, just discard
	}
}

func (qf *storageQueryFlow) Complete(err error) {
	if err != nil {
		qf.mux.Lock()
		defer qf.mux.Unlock()
		qf.err = err
		qf.completed.Store(true)
	}
}

func (qf *storageQueryFlow) Scanner(task concurrent.Task) {
	qf.execute(Scanner, task)
}

func (qf *storageQueryFlow) Grouping(task concurrent.Task) {
	qf.execute(Grouping, task)
}

func (qf *storageQueryFlow) Filtering(task concurrent.Task) {
	qf.execute(Filtering, task)
}

func (qf *storageQueryFlow) Reduce(tags string, agg aggregation.FieldAggregates) {
	//NOTICE: don't do reduce operator in other goroutine, because big overhead when goroutine schedule
	defer func() {
		qf.releaseAgg(agg)
	}()

	if qf.completed.Load() {
		storageQueryFlowLogger.Warn("reduce the aggregator data after storage query flow completed")
		return
	}

	qf.mux.Lock()
	defer qf.mux.Unlock()

	qf.reduceAgg.Aggregate(agg.ResultSet(tags))
}

func (qf *storageQueryFlow) completeTask(taskID int32) {
	qf.mux.Lock()
	defer qf.mux.Unlock()

	delete(qf.pendingTasks, taskID)

	if len(qf.pendingTasks) == 0 {
		// if all tasks of all stages completed
		qf.completed.Store(true)

		errMsg := ""
		var data []byte
		if qf.err != nil {
			errMsg = qf.err.Error()
		} else if qf.reduceAgg != nil {
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

					fields[fieldIt.FieldName()] = data
				}
				if len(fields) > 0 {
					timeSeriesList = append(timeSeriesList, &pb.TimeSeries{
						Tags:   ts.Tags(),
						Fields: fields,
					})
				}
			}

			seriesList := pb.TimeSeriesList{
				TimeSeriesList: timeSeriesList,
			}
			// no error
			data, _ = seriesList.Marshal()
		}

		// send result to upstream
		if err := qf.stream.Send(&pb.TaskResponse{
			JobID:     qf.req.JobID,
			TaskID:    qf.req.ParentTaskID,
			Completed: true,
			Payload:   data,
			ErrMsg:    errMsg,
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
						err = errors.New("UnKnow ERROR")
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

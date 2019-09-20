package query

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

var log = logger.GetLogger("query", "StorageExecutor")

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Scanning
// 3) Grouping if need
// 4) Down sampling
// 5) Sample aggregation
type storageExecutor struct {
	engine   tsdb.Engine
	query    *stmt.Query
	shardIDs []int32

	shards []tsdb.Shard

	metricID uint32

	fieldIDs           []uint16
	storageExecutePlan *storageExecutePlan
	intervalType       interval.Type

	resultCh    chan *series.TimeSeriesEvent
	executorCtx *storageExecutorContext

	ctx context.Context
	err error
}

// newStorageExecutor creates the execution which queries the data of storage engine
func newStorageExecutor(ctx context.Context, engine tsdb.Engine, shardIDs []int32, query *stmt.Query) parallel.Executor {
	return &storageExecutor{
		engine:   engine,
		shardIDs: shardIDs,
		query:    query,
		ctx:      ctx,
	}
}

// Execute executes search logic in storage level,
// 1) validation input params
// 2) build execute plan
// 3) build execute pipeline
// 4) run pipeline
func (e *storageExecutor) Execute() <-chan *series.TimeSeriesEvent {
	// do query validation
	if err := e.validation(); err != nil {
		e.err = err
		return nil
	}

	// get shard by given query shard id list
	for _, shardID := range e.shardIDs {
		shard := e.engine.GetShard(shardID)
		// if shard exist, add shard to query list
		if shard != nil {
			e.shards = append(e.shards, shard)
		}
	}

	// check got shards if valid
	if err := e.checkShards(); err != nil {
		e.err = err
		return nil
	}

	plan := newStorageExecutePlan(e.engine.GetIDGetter(), e.query)
	if err := plan.Plan(); err != nil {
		e.err = err
		return nil
	}
	storageExecutePlan := plan.(*storageExecutePlan)

	e.metricID = storageExecutePlan.metricID
	e.intervalType = interval.CalcIntervalType(e.query.Interval)

	//TODO set size
	e.resultCh = make(chan *series.TimeSeriesEvent, 10)
	e.executorCtx = newStorageExecutorContext(e.resultCh)

	e.fieldIDs = storageExecutePlan.getFieldIDs()
	e.storageExecutePlan = storageExecutePlan

	// need retain total memory and shard search
	e.executorCtx.retainTask(int32(len(e.shards) * 2))
	for _, shard := range e.shards {
		go e.memoryDBSearch(shard)
		e.shardLevelSearch(shard)
	}
	return e.resultCh
}

// Statement returns the query statement
func (e *storageExecutor) Statement() *stmt.Query {
	return e.query
}

// Error returns the execution error
func (e *storageExecutor) Error() error {
	return e.err
}

// memoryDBSearch searches data from memory database
func (e *storageExecutor) memoryDBSearch(shard tsdb.Shard) {
	defer e.executorCtx.completeTask()

	memoryDB := shard.GetMemoryDatabase()
	seriesIDSet := e.searchSeriesIDs(memoryDB)
	if seriesIDSet == nil || seriesIDSet.IsEmpty() {
		return
	}

	aggWorker := createAggWorker(e.query.Interval, &e.query.TimeRange, e.storageExecutePlan.getFields(), e.resultCh)
	worker := createScanWorker(e.ctx, e.metricID, e.query.GroupBy, memoryDB, aggWorker)
	defer worker.Close()
	memoryDB.Scan(&series.ScanContext{
		MetricID:    e.metricID,
		FieldIDs:    e.fieldIDs,
		TimeRange:   e.query.TimeRange,
		SeriesIDSet: seriesIDSet,
		Worker:      worker,
	})
}

// searchSeriesIDs searches series ids from index
func (e *storageExecutor) searchSeriesIDs(filter series.Filter) (seriesIDSet *series.MultiVerSeriesIDSet) {
	condition := e.query.Condition
	metricID := e.metricID
	if condition != nil {
		seriesSearch := newSeriesSearch(metricID, filter, e.query)
		idSet, err := seriesSearch.Search()
		if err != nil {
			if err != series.ErrNotFound {
				e.err = err
			}
			return
		}
		seriesIDSet = idSet
	}
	//TODO add metric level search for no condition
	return
}

// shardLevelSearch searches data from shard
func (e *storageExecutor) shardLevelSearch(shard tsdb.Shard) {
	// must complete task
	defer e.executorCtx.completeTask()

	// find data family
	families := shard.GetDataFamilies(e.intervalType, e.query.TimeRange)
	if len(families) == 0 {
		return
	}

	seriesIDSet := e.searchSeriesIDs(shard.GetSeriesIDsFilter())
	if seriesIDSet == nil || seriesIDSet.IsEmpty() {
		return
	}
	// retain family task first
	e.executorCtx.retainTask(int32(2 * len(families)))
	aggWorker := createAggWorker(e.query.Interval, &e.query.TimeRange, e.storageExecutePlan.getFields(), e.resultCh)
	worker := createScanWorker(e.ctx, e.metricID, e.query.GroupBy, shard.GetMetaGetter(), aggWorker)
	defer worker.Close()
	for _, family := range families {
		go e.familyLevelSearch(worker, family, seriesIDSet)
	}
}

// familyLevelSearch searches data from data family, do down sampling and aggregation
func (e *storageExecutor) familyLevelSearch(worker series.ScanWorker, family tsdb.DataFamily,
	seriesIDSet *series.MultiVerSeriesIDSet) {
	// must complete task
	defer e.executorCtx.completeTask()

	family.Scan(&series.ScanContext{
		MetricID:    e.metricID,
		FieldIDs:    e.fieldIDs,
		TimeRange:   e.query.TimeRange,
		SeriesIDSet: seriesIDSet,
		Worker:      worker,
	})
}

// validation validates query input params are valid
func (e *storageExecutor) validation() error {
	// check input shardIDs if empty
	if len(e.shardIDs) == 0 {
		return fmt.Errorf("there is no shard id in search condition")
	}
	numOfShards := e.engine.NumOfShards()
	// check engine has shard
	if numOfShards == 0 {
		return fmt.Errorf("tsdb engine[%s] hasn't shard", e.engine.Name())
	}
	if numOfShards != len(e.shardIDs) {
		return fmt.Errorf("storage's num. of shard not match search condition")
	}
	return nil
}

// checkShards checks got shards if valid
func (e *storageExecutor) checkShards() error {
	numOfShards := len(e.shards)
	if numOfShards == 0 {
		return fmt.Errorf("cannot find shard by given shard id")
	}
	numOfShardIDs := len(e.shardIDs)
	if numOfShards != numOfShardIDs {
		return fmt.Errorf("got shard size[%d] not eqauls input shard size[%d]", numOfShards, numOfShardIDs)
	}
	return nil
}

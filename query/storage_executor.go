package query

import (
	"fmt"
	"sync"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Scanning
// 3) Grouping if need
// 4) Down sampling
// 5) Sample aggregation
type storageExecutor struct {
	database tsdb.Database
	query    *stmt.Query
	shardIDs []int32

	shards []tsdb.Shard

	metricID uint32

	fieldIDs           []uint16
	storageExecutePlan *storageExecutePlan
	intervalType       timeutil.IntervalType

	executorPool *tsdb.ExecutorPool

	executeCtx parallel.ExecuteContext
}

// newStorageExecutor creates the execution which queries the data of storage engine
func newStorageExecutor(
	ctx parallel.ExecuteContext,
	database tsdb.Database,
	shardIDs []int32,
	query *stmt.Query,
) parallel.Executor {
	return &storageExecutor{
		database:     database,
		shardIDs:     shardIDs,
		query:        query,
		executorPool: database.ExecutorPool(),
		executeCtx:   ctx,
	}
}

// Execute executes search logic in storage level,
// 1) validation input params
// 2) build execute plan
// 3) build execute pipeline
// 4) run pipeline
func (e *storageExecutor) Execute() {
	// do query validation
	if err := e.validation(); err != nil {
		e.executeCtx.Complete(err)
		return
	}

	// get shard by given query shard id list
	for _, shardID := range e.shardIDs {
		shard, ok := e.database.GetShard(shardID)
		// if shard exist, add shard to query list
		if ok {
			e.shards = append(e.shards, shard)
		}
	}

	// check got shards if valid
	if err := e.checkShards(); err != nil {
		e.executeCtx.Complete(err)
		return
	}

	plan := newStorageExecutePlan(e.database.IDGetter(), e.query)
	if err := plan.Plan(); err != nil {
		e.executeCtx.Complete(err)
		return
	}
	storageExecutePlan := plan.(*storageExecutePlan)

	e.metricID = storageExecutePlan.metricID
	e.intervalType = timeutil.Interval(e.query.Interval).Type()

	e.fieldIDs = storageExecutePlan.getFieldIDs()
	e.storageExecutePlan = storageExecutePlan

	// need retain total memory and shard search
	e.executeCtx.RetainTask(1)
	for idx := range e.shards {
		shard := e.shards[idx]
		// execute memory db search in background goroutine
		e.executeCtx.RetainTask(1)
		e.executorPool.Scanners.Submit(func() {
			e.memoryDBSearch(shard)
		})

		e.executeCtx.RetainTask(1)
		e.shardLevelSearch(shard)
	}
	e.executeCtx.Complete(nil)
}

// memoryDBSearch searches data from memory database
func (e *storageExecutor) memoryDBSearch(shard tsdb.Shard) {
	memoryDB := shard.MemoryDatabase()
	seriesIDSet := e.searchSeriesIDs(memoryDB)
	if seriesIDSet == nil || seriesIDSet.IsEmpty() {
		// if series ids not found, complete the search task
		e.executeCtx.Complete(nil)
		return
	}

	timeRange, intervalRatio, queryInterval := downSamplingTimeRange(e.query.Interval, memoryDB.Interval(), e.query.TimeRange)
	aggSpecs := e.storageExecutePlan.getDownSamplingAggSpecs()
	groupAgg := aggregation.NewGroupingAggregator(queryInterval, timeRange, aggSpecs)

	// scan data and complete task in scan worker after scan worker completed
	worker := createScanWorker(e.executeCtx, e.metricID, e.query.GroupBy, memoryDB, groupAgg, e.executorPool)
	defer worker.Close()
	memoryDB.Scan(&series.ScanContext{
		MetricID:    e.metricID,
		FieldIDs:    e.fieldIDs,
		SeriesIDSet: seriesIDSet,
		HasGroupBy:  e.storageExecutePlan.hasGroupBy(),
		Worker:      worker,
		Aggregators: e.getAggregatorPool(queryInterval, intervalRatio, timeRange),
	})
}

// getAggregatorPool returns aggregator pool
func (e *storageExecutor) getAggregatorPool(
	queryInterval timeutil.Interval,
	intervalRatio int,
	timeRange timeutil.TimeRange,
) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return aggregation.NewFieldAggregates(queryInterval, intervalRatio, timeRange, true,
				e.storageExecutePlan.getDownSamplingAggSpecs())
		},
	}
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
				e.executeCtx.Complete(err)
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
	// find data family
	families := shard.GetDataFamilies(e.intervalType, e.query.TimeRange)
	if len(families) == 0 {
		e.executeCtx.Complete(nil)
		return
	}

	seriesIDSet := e.searchSeriesIDs(shard.IndexFilter())
	if seriesIDSet == nil || seriesIDSet.IsEmpty() {
		e.executeCtx.Complete(nil)
		return
	}
	// retain family task first
	e.executeCtx.RetainTask(int32(2 * len(families)))
	//FIXME get interval
	timeRange, _, queryInterval := downSamplingTimeRange(e.query.Interval, 10, e.query.TimeRange)
	aggSpecs := e.storageExecutePlan.getDownSamplingAggSpecs()
	groupAgg := aggregation.NewGroupingAggregator(queryInterval, timeRange, aggSpecs)

	worker := createScanWorker(
		e.executeCtx,
		e.metricID,
		e.query.GroupBy,
		shard.IndexMetaGetter(),
		groupAgg,
		e.executorPool,
	)
	for _, family := range families {
		go e.familyLevelSearch(worker, family, seriesIDSet)
	}
}

// familyLevelSearch searches data from data family, do down sampling and aggregation
func (e *storageExecutor) familyLevelSearch(worker series.ScanWorker, family tsdb.DataFamily,
	seriesIDSet *series.MultiVerSeriesIDSet) {
	// must complete task
	defer e.executeCtx.Complete(nil)

	family.Scan(&series.ScanContext{
		MetricID:    e.metricID,
		FieldIDs:    e.fieldIDs,
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
	numOfShards := e.database.NumOfShards()
	// check engine has shard
	if numOfShards == 0 {
		return fmt.Errorf("tsdb database[%s] hasn't shard", e.database.Name())
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

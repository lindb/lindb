package query

import (
	"fmt"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/series"
)

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Scanning
// 3) Grouping if need
// 4) down sampling
// 5) Sample aggregation
type storageExecutor struct {
	engine   tsdb.Engine
	query    *stmt.Query
	shardIDs []int32

	shards []tsdb.Shard

	metricID uint32

	fieldIDs      []uint16
	aggregations  map[uint16]*aggregation.AggregatorSpec
	intervalRatio int
	interval      int64

	resultCh chan field.GroupedTimeSeries

	err error
}

// newStorageExecutor creates the execution which queries the data of storage engine
func newStorageExecutor(engine tsdb.Engine, shardIDs []int32, query *stmt.Query) parallel.Executor {
	interval := query.Interval
	if interval <= 0 {
		//TODO use storage interval
		interval = 10 * timeutil.OneSecond
	}
	return &storageExecutor{
		engine:   engine,
		shardIDs: shardIDs,
		query:    query,
		interval: interval,
	}
}

// Execute executes search logic in storage level,
// 1) validation input params
// 2) build execute plan
// 3) build execute pipeline
// 4) run pipeline
func (e *storageExecutor) Execute() <-chan field.GroupedTimeSeries {
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

	//TODO need modify
	e.intervalRatio = timeutil.CalIntervalRatio(100, 100)

	plan := newStorageExecutePlan(e.engine.GetIDGetter(), e.query)
	if err := plan.Plan(); err != nil {
		e.err = err
		return nil
	}
	storageExecutePlan, ok := plan.(*storageExecutePlan)
	if !ok {
		e.err = fmt.Errorf("cannot get storage execute plan")
		return nil
	}

	//TODO set size
	e.resultCh = make(chan field.GroupedTimeSeries, 10)

	e.fieldIDs = storageExecutePlan.getFieldIDs()
	e.aggregations = storageExecutePlan.fields

	for _, shard := range e.shards {
		e.shardLevelSearch(shard)
	}
	return e.resultCh
}

// Error returns the execution error
func (e *storageExecutor) Error() error {
	return e.err
}

// shardLevelSearch searches data from shard
func (e *storageExecutor) shardLevelSearch(shard tsdb.Shard) {
	condition := e.query.Condition
	metricID := e.metricID
	var seriesIDSet *series.MultiVerSeriesIDSet
	if condition != nil {
		seriesSearch := newSeriesSearch(metricID, shard.GetSeriesIDsFilter(), e.query)
		idSet, err := seriesSearch.Search()
		if err != nil {
			//TODO
			return
		}
		if idSet == nil || idSet.IsEmpty() {
			return
		}
		seriesIDSet = idSet
	}
	//TODO need group by
	timeRange := e.query.TimeRange
	segments := shard.GetSegments(e.query.IntervalType, timeRange)
	for _, segment := range segments {
		families := segment.GetDataFamilyScanners(timeRange)
		for _, family := range families {
			e.familyLevelSearch(family, seriesIDSet)
		}
	}
}

// familyLevelSearch searches data from data family, do down sampling and aggregation
func (e *storageExecutor) familyLevelSearch(scanner series.DataFamilyScanner, seriesIDSet *series.MultiVerSeriesIDSet) {
	scanItr := scanner.Scan(
		series.ScanContext{
			MetricID:    e.metricID,
			FieldIDs:    e.fieldIDs,
			TimeRange:   e.query.TimeRange,
			SeriesIDSet: seriesIDSet,
		})

	if scanItr == nil {
		return
	}
	defer scanItr.Close()
	for scanItr.HasNext() {
		timeSeries := scanItr.Next()
		if timeSeries == nil {
			break
		}
		for timeSeries.HasNext() {
			it := timeSeries.Next()
			//TODO use family time range
			agg := aggregation.NewFieldAggregator(1, e.interval, e.query.TimeRange.Start, e.query.TimeRange.End, e.intervalRatio, e.aggregations[it.ID()])
			agg.Aggregate(it)
		}
	}
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

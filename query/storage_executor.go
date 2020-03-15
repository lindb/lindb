package query

import (
	"errors"
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

// for testing
var (
	newTagSearchFunc          = newTagSearch
	newStorageExecutePlanFunc = newStorageExecutePlan
	newSeriesSearchFunc       = newSeriesSearch
)

var (
	errNoShardID         = errors.New("there is no shard id in search condition")
	errNoShardInDatabase = errors.New("there is no shard in database storage engine")
	errShardNotFound     = errors.New("shard not found in database storage engine")
	errShardNotMatch     = errors.New("storage's num. of shard not match search condition")
	errShardNumNotMatch  = errors.New("got shard size not equals input shard size")
)

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Grouping if need
// 3) Scanning and Loading
// 4) Down sampling
// 5) Simple aggregation
type storageExecutor struct {
	database  tsdb.Database
	namespace string
	query     *stmt.Query

	shardIDs []int32
	shards   []tsdb.Shard

	metricID           uint32
	fieldIDs           []field.ID
	storageExecutePlan *storageExecutePlan
	intervalType       timeutil.IntervalType

	filterResult map[string]*filterResult

	queryFlow flow.StorageQueryFlow

	// group by query need
	mutex              sync.Mutex
	groupByTagKeyIDs   []uint32
	tagValueIDs        []*roaring.Bitmap // for group by query store tag value ids for each group tag key
	pendingForShard    atomic.Int32
	pendingForGrouping atomic.Int32
	collecting         atomic.Bool
}

// newStorageExecutor creates the execution which queries the data of storage engine
func newStorageExecutor(
	queryFlow flow.StorageQueryFlow,
	database tsdb.Database,
	namespace string,
	shardIDs []int32,
	query *stmt.Query,
) parallel.Executor {
	return &storageExecutor{
		database:  database,
		namespace: namespace,
		shardIDs:  shardIDs,
		query:     query,
		queryFlow: queryFlow,
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
		e.queryFlow.Complete(err)
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
		e.queryFlow.Complete(err)
		return
	}

	plan := newStorageExecutePlanFunc(e.namespace, e.database.Metadata(), e.query)
	if err := plan.Plan(); err != nil {
		e.queryFlow.Complete(err)
		return
	}

	condition := e.query.Condition
	var err error
	if condition != nil {
		tagSearch := newTagSearchFunc(e.namespace, e.query, e.database.Metadata())
		e.filterResult, err = tagSearch.Filter()
		if err != nil {
			e.queryFlow.Complete(err)
			return
		}
		if len(e.filterResult) == 0 {
			// filter not match, return not found
			e.queryFlow.Complete(constants.ErrNotFound)
			return
		}
	}

	storageExecutePlan := plan.(*storageExecutePlan)
	e.intervalType = e.query.Interval.Type()

	// prepare storage query flow
	e.queryFlow.Prepare(storageExecutePlan.getDownSamplingAggSpecs())

	e.metricID = storageExecutePlan.metricID
	e.fieldIDs = storageExecutePlan.getFieldIDs()
	e.storageExecutePlan = storageExecutePlan

	// execute query flow
	e.executeQuery()
}

// executeQuery executes query flow for each shard
func (e *storageExecutor) executeQuery() {
	e.pendingForShard.Store(int32(len(e.shards)))
	for idx := range e.shards {
		shard := e.shards[idx]
		e.queryFlow.Filtering(func() {
			defer func() {
				// finish shard query
				e.pendingForShard.Dec()
				// try start collect tag values
				e.collectGroupByTagValues()
			}()
			// 1. get series ids by query condition
			seriesIDs := e.searchSeriesIDs(shard.IndexDatabase())
			// if series ids not found
			if seriesIDs == nil || seriesIDs.IsEmpty() {
				return
			}
			var rs []flow.FilterResultSet
			// 2. filter data in memory database
			resultSet, err := shard.MemoryDatabase().Filter(e.metricID, e.fieldIDs, seriesIDs, e.query.TimeRange)
			if err != nil && err != constants.ErrNotFound {
				// maybe data not exist in memory database, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			rs = append(rs, resultSet...)
			// 3. filter data each data family in shard
			resultSet, err = e.filterForShard(shard, seriesIDs)
			if err != nil && err != constants.ErrNotFound {
				// maybe data not exist in shard, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			rs = append(rs, resultSet...)
			if len(rs) == 0 {
				// data not found
				return
			}
			// 4. merge all series ids after filtering => final series ids
			seriesIDsAfterFilter := roaring.New()
			for _, result := range rs {
				seriesIDsAfterFilter.Or(result.SeriesIDs())
			}
			// 5. execute group by
			e.pendingForGrouping.Inc()
			e.queryFlow.Grouping(func() {
				defer func() {
					e.pendingForGrouping.Dec()
					// try start collect tag values
					e.collectGroupByTagValues()
				}()
				e.executeGroupBy(shard.IndexDatabase(), rs, seriesIDs)
			})
		})
	}
}

// searchSeriesIDs searches series ids from index
func (e *storageExecutor) searchSeriesIDs(filter series.Filter) (seriesIDs *roaring.Bitmap) {
	condition := e.query.Condition
	var err error
	if condition != nil {
		// if get tag filter result do series ids searching
		seriesSearch := newSeriesSearchFunc(filter, e.filterResult, e.query)
		seriesIDs, err = seriesSearch.Search()
	} else {
		// get series ids for metric level
		seriesIDs, err = filter.GetSeriesIDsForMetric(e.namespace, e.query.MetricName)
		if err == nil && !e.query.HasGroupBy() {
			// add series id without tags, maybe one metric has too many series, but one series without tags
			seriesIDs.Add(constants.SeriesIDWithoutTags)
		}
	}
	if err != nil && err != constants.ErrNotFound {
		// maybe series ids not found in shard, so ignore not found err
		e.queryFlow.Complete(err)
	}
	return
}

// filterForShard filtering data in shard
func (e *storageExecutor) filterForShard(shard tsdb.Shard, seriesIDs *roaring.Bitmap) (rs []flow.FilterResultSet, err error) {
	families := shard.GetDataFamilies(e.intervalType, e.query.TimeRange)
	if len(families) == 0 {
		return nil, nil
	}
	for idx := range families {
		family := families[idx]
		// execute data family search in background goroutine
		resultSet, err := family.Filter(e.metricID, e.fieldIDs, seriesIDs, e.query.TimeRange)
		if err != nil {
			return nil, err
		}
		rs = append(rs, resultSet...)
	}
	return rs, nil
}

// executeGroupBy executes the query flow, step as below:
// 1. grouping
// 2. loading
func (e *storageExecutor) executeGroupBy(indexDB indexdb.IndexDatabase, rs []flow.FilterResultSet, seriesIDs *roaring.Bitmap) {
	var groupingCtx series.GroupingContext
	// 1. grouping, if has group by, do group by tag keys, else just split series ids as batch first,
	// get grouping context if need
	if e.query.HasGroupBy() {
		e.groupByTagKeyIDs = e.storageExecutePlan.groupByKeyIDs()
		gCtx, err := indexDB.GetGroupingContext(e.storageExecutePlan.groupByKeyIDs(), seriesIDs)
		if err != nil && err != constants.ErrNotFound {
			// maybe group by not found, so ignore not found
			e.queryFlow.Complete(err)
			return
		}
		if gCtx == nil {
			return
		}
		groupingCtx = gCtx
	}
	keys := seriesIDs.GetHighKeys()
	e.pendingForGrouping.Add(int32(len(keys)))
	var groupWait atomic.Int32
	groupWait.Add(int32(len(keys)))

	for idx, key := range keys {
		// be carefully, need use new variable for variable scope problem
		highKey := key
		container := seriesIDs.GetContainerAtIndex(idx)
		// grouping based on group by tag keys for each container
		e.queryFlow.Grouping(func() {
			defer func() {
				groupWait.Dec()
				if groupingCtx != nil && groupWait.Load() == 0 {
					// current group by query completed, need merge group by tag value ids
					e.mergeGroupByTagValueIDs(groupingCtx.GetGroupByTagValueIDs())
				}
				e.pendingForGrouping.Dec()
				// try start collect tag values for group by query
				e.collectGroupByTagValues()
			}()
			var groupedSeries map[string][]uint16
			if groupingCtx != nil {
				// build group by data, grouped series: tags => series IDs
				groupedSeries = groupingCtx.BuildGroup(highKey, container)
			} else {
				groupedSeries = map[string][]uint16{"": container.ToArray()}
			}
			for _, resultSet := range rs {
				// 3.load data by grouped seriesIDs
				filteringRS := resultSet
				e.queryFlow.Scanner(func() {
					filteringRS.Load(e.queryFlow, e.fieldIDs, highKey, groupedSeries)
				})
			}
		})
	}
}

// mergeGroupByTagValueIDs merges group by tag value ids for each shard
func (e *storageExecutor) mergeGroupByTagValueIDs(tagValueIDs []*roaring.Bitmap) {
	if tagValueIDs == nil {
		return
	}
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.tagValueIDs == nil {
		e.tagValueIDs = make([]*roaring.Bitmap, len(e.groupByTagKeyIDs))
	}

	for idx, tagVIDs := range e.tagValueIDs {
		if tagVIDs == nil {
			e.tagValueIDs[idx] = tagValueIDs[idx]
		} else {
			tagVIDs.Or(tagValueIDs[idx])
		}
	}
}

// collectGroupByTagValues collects group tag values
func (e *storageExecutor) collectGroupByTagValues() {
	// all shard pending query tasks and grouping task completed, start collect tag values
	if e.pendingForShard.Load() == 0 && e.pendingForGrouping.Load() == 0 {
		if e.collecting.CAS(false, true) {
			for idx, tagKeyID := range e.groupByTagKeyIDs {
				tagKey := tagKeyID
				tagValueIDs := e.tagValueIDs[idx]
				tagIndex := idx
				e.queryFlow.Scanner(func() {
					//FIXME need check group by tag value ids is nil???
					tagValues := make(map[uint32]string)
					if err := e.database.Metadata().TagMetadata().CollectTagValues(tagKey, tagValueIDs, tagValues); err != nil {
						e.queryFlow.Complete(err)
						return
					}
					e.queryFlow.ReduceTagValues(tagIndex, tagValues)
				})
			}
		}
	}
}

// validation validates query input params are valid
func (e *storageExecutor) validation() error {
	// check input shardIDs if empty
	if len(e.shardIDs) == 0 {
		return errNoShardID
	}
	numOfShards := e.database.NumOfShards()
	// check engine has shard
	if numOfShards == 0 {
		return errNoShardInDatabase
	}
	if numOfShards != len(e.shardIDs) {
		return errShardNotMatch
	}
	return nil
}

// checkShards checks got shards if valid
func (e *storageExecutor) checkShards() error {
	numOfShards := len(e.shards)
	if numOfShards == 0 {
		return errShardNotFound
	}
	numOfShardIDs := len(e.shardIDs)
	if numOfShards != numOfShardIDs {
		return errShardNumNotMatch
	}
	return nil
}

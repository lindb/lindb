package query

import (
	"errors"

	"github.com/lindb/roaring"

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
	database tsdb.Database
	query    *stmt.Query

	shardIDs []int32
	shards   []tsdb.Shard

	metricID           uint32
	fieldIDs           []field.ID
	storageExecutePlan *storageExecutePlan
	intervalType       timeutil.IntervalType

	filterResult map[string]*filterResult

	queryFlow flow.StorageQueryFlow
}

// newStorageExecutor creates the execution which queries the data of storage engine
func newStorageExecutor(
	queryFlow flow.StorageQueryFlow,
	database tsdb.Database,
	shardIDs []int32,
	query *stmt.Query,
) parallel.Executor {
	return &storageExecutor{
		database:  database,
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

	plan := newStorageExecutePlanFunc(e.database.Metadata(), e.query)
	if err := plan.Plan(); err != nil {
		e.queryFlow.Complete(err)
		return
	}

	condition := e.query.Condition
	var err error
	if condition != nil {
		tagSearch := newTagSearchFunc(e.query, e.database.Metadata())
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

	for idx := range e.shards {
		shard := e.shards[idx]
		seriesIDs := e.searchSeriesIDs(shard.IndexDatabase())
		// if series ids not found
		if seriesIDs == nil || seriesIDs.IsEmpty() {
			continue
		}

		// execute memory db search in background goroutine
		e.queryFlow.Filtering(func() {
			e.executeQueryFlow(shard.IndexDatabase(), shard.MemoryDatabase(), seriesIDs)
		})

		e.queryFlow.Filtering(func() {
			e.shardLevelSearch(shard, seriesIDs)
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
		//FIXME
		//} else {
		//	seriesIDs, err = filter.(metricID, e.query.TimeRange)
	}
	if err != nil {
		e.queryFlow.Complete(err)
	}
	return
}

// shardLevelSearch searches data from shard
func (e *storageExecutor) shardLevelSearch(shard tsdb.Shard, seriesIDs *roaring.Bitmap) {
	families := shard.GetDataFamilies(e.intervalType, e.query.TimeRange)
	if len(families) == 0 {
		return
	}
	for idx := range families {
		family := families[idx]
		// execute data family search in background goroutine
		e.queryFlow.Filtering(func() {
			e.executeQueryFlow(shard.IndexDatabase(), family, seriesIDs)
		})
	}
}

// executeQueryFlow executes the query flow, step as below:
// 1. filtering
// 2. grouping
// 3. loading
func (e *storageExecutor) executeQueryFlow(indexDB indexdb.IndexDatabase, filter flow.DataFilter, seriesIDs *roaring.Bitmap) {
	hasGroupBy := e.query.HasGroupBy()
	// 1. filtering, check series ids if exist in storage
	e.queryFlow.Filtering(func() {
		resultSet, err := filter.Filter(e.metricID, e.fieldIDs, seriesIDs, e.query.TimeRange)
		if err != nil {
			e.queryFlow.Complete(err)
			return
		}
		if len(resultSet) == 0 {
			// not found in storage, return it
			return
		}
		// 2. grouping, if has group by, do group by tag keys, else just split series ids as batch first,
		// get grouping context if need
		var groupingCtx series.GroupingContext
		if hasGroupBy {
			//FIXME
			gCtx, err := indexDB.GetGroupingContext(nil)
			if err != nil {
				e.queryFlow.Complete(err)
				return
			}
			if gCtx == nil {
				return
			}
			groupingCtx = gCtx
		}
		keys := seriesIDs.GetHighKeys()

		for idx, key := range keys {
			// be carefully, need use new variable for variable scope problem
			highKey := key
			container := seriesIDs.GetContainerAtIndex(idx)
			// grouping based on group by tag keys for each container
			e.queryFlow.Grouping(func() {
				var groupedSeries map[string][]uint16
				if hasGroupBy {
					// build group by data, grouped series: tags => series IDs
					groupedSeries = groupingCtx.BuildGroup(highKey, container)
				} else {
					groupedSeries = map[string][]uint16{"": container.ToArray()}
				}
				for _, rs := range resultSet {
					// 3.load data by grouped seriesIDs
					filteringRS := rs
					e.queryFlow.Scanner(func() {
						filteringRS.Load(e.queryFlow, e.fieldIDs, highKey, groupedSeries)
					})
				}
			})
		}
	})
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

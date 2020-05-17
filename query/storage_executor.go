package query

import (
	"errors"
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
)

// for testing
var (
	newTagSearchFunc          = newTagSearch
	newStorageExecutePlanFunc = newStorageExecutePlan
	newSeriesSearchFunc       = newSeriesSearch
	newBuildGroupTaskFunc     = newBuildGroupTask
	newDataLoadTaskFunc       = newDataLoadTask
)

var (
	errNoShardID         = errors.New("there is no shard id in search condition")
	errNoShardInDatabase = errors.New("there is no shard in database storage engine")
	errShardNotFound     = errors.New("shard not found in database storage engine")
	errShardNotMatch     = errors.New("storage's num. of shard not match search condition")
	errShardNumNotMatch  = errors.New("got shard size not equals input shard size")
)

// filterResultSet represents data filter result set
type filterResultSet struct {
	rs []flow.FilterResultSet
}

// groupingResult represents the grouping context result
type groupingResult struct {
	groupingCtx series.GroupingContext
}

// groupedSeriesResult represents grouped series for group by query
type groupedSeriesResult struct {
	groupedSeries map[string][]uint16
}

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Grouping if need
// 3) Scanning and Loading
// 4) Down sampling
// 5) Simple aggregation
type storageExecutor struct {
	database tsdb.Database
	ctx      *storageExecuteContext
	shards   []tsdb.Shard

	metricID           uint32
	fieldIDs           []field.ID
	storageExecutePlan *storageExecutePlan

	queryFlow flow.StorageQueryFlow

	// group by query need
	mutex              sync.Mutex
	groupByTagKeyIDs   []tag.Meta
	tagValueIDs        []*roaring.Bitmap // for group by query store tag value ids for each group tag key
	pendingForShard    atomic.Int32
	pendingForGrouping atomic.Int32
	collecting         atomic.Bool
}

// newStorageExecutor creates the execution which queries the data of storage engine
func newStorageExecutor(
	queryFlow flow.StorageQueryFlow,
	database tsdb.Database,
	storageExecuteCtx parallel.StorageExecuteContext,
) parallel.Executor {
	ctx := storageExecuteCtx.(*storageExecuteContext)
	return &storageExecutor{
		database:  database,
		ctx:       ctx,
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
	for _, shardID := range e.ctx.shardIDs {
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

	plan := newStorageExecutePlanFunc(e.ctx.query.Namespace, e.database.Metadata(), e.ctx.query)
	t := newStoragePlanTask(e.ctx, plan)

	if err := t.Run(); err != nil {
		e.queryFlow.Complete(err)
		return
	}
	condition := e.ctx.query.Condition
	if condition != nil {
		tagSearch := newTagSearchFunc(e.ctx.query.Namespace, e.ctx.query, e.database.Metadata())
		t = newTagFilterTask(e.ctx, tagSearch)
		if err := t.Run(); err != nil {
			e.queryFlow.Complete(err)
			return
		}
	}

	storageExecutePlan := plan.(*storageExecutePlan)

	// prepare storage query flow
	e.queryFlow.Prepare(storageExecutePlan.getDownSamplingAggSpecs())

	e.metricID = storageExecutePlan.metricID
	e.fieldIDs = storageExecutePlan.getFieldIDs()
	e.storageExecutePlan = storageExecutePlan
	if e.ctx.query.HasGroupBy() {
		e.groupByTagKeyIDs = e.storageExecutePlan.groupByKeyIDs()
		e.tagValueIDs = make([]*roaring.Bitmap, len(e.groupByTagKeyIDs))
	}

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
			seriesIDs := roaring.New()
			t := newSeriesIDsSearchTask(e.ctx, shard, seriesIDs)
			err := t.Run()
			if err != nil && err != constants.ErrNotFound {
				// maybe series ids not found in shard, so ignore not found err
				e.queryFlow.Complete(err)
			}
			// if series ids not found
			if seriesIDs.IsEmpty() {
				return
			}

			rs := &filterResultSet{}
			// 2. filter data in memory database
			t = newMemoryDataFilterTask(e.ctx, shard, e.metricID, e.fieldIDs, seriesIDs, rs)
			err = t.Run()
			if err != nil && err != constants.ErrNotFound {
				// maybe data not exist in memory database, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			// 3. filter data each data family in shard
			t = newFileDataFilterTask(e.ctx, shard, e.metricID, e.fieldIDs, seriesIDs, rs)
			err = t.Run()
			if err != nil && err != constants.ErrNotFound {
				// maybe data not exist in shard, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			if len(rs.rs) == 0 {
				// data not found
				return
			}
			// 4. merge all series ids after filtering => final series ids
			seriesIDsAfterFilter := roaring.New()
			for _, result := range rs.rs {
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
				e.executeGroupBy(shard, rs.rs, seriesIDs)
			})
		})
	}
}

// executeGroupBy executes the query flow, step as below:
// 1. grouping
// 2. loading
func (e *storageExecutor) executeGroupBy(shard tsdb.Shard, rs []flow.FilterResultSet, seriesIDs *roaring.Bitmap) {
	groupingResult := &groupingResult{}
	var groupingCtx series.GroupingContext
	// 1. grouping, if has group by, do group by tag keys, else just split series ids as batch first,
	// get grouping context if need
	if e.ctx.query.HasGroupBy() {
		tagKeys := make([]uint32, len(e.groupByTagKeyIDs))
		for idx, tagKeyID := range e.groupByTagKeyIDs {
			tagKeys[idx] = tagKeyID.ID
		}
		t := newGroupingContextFindTask(e.ctx, shard, tagKeys, seriesIDs, groupingResult)
		err := t.Run()
		if err != nil && err != constants.ErrNotFound {
			// maybe group by not found, so ignore not found
			e.queryFlow.Complete(err)
			return
		}
		if groupingResult.groupingCtx == nil {
			return
		}
		groupingCtx = groupingResult.groupingCtx
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
			groupedResult := &groupedSeriesResult{}
			t := newBuildGroupTaskFunc(e.ctx, shard, groupingCtx, highKey, container, groupedResult)
			if err := t.Run(); err != nil {
				e.queryFlow.Complete(err)
				return
			}
			for _, resultSet := range rs {
				// 3.load data by grouped seriesIDs
				filteringRS := resultSet
				e.queryFlow.Scanner(func() {
					t := newDataLoadTaskFunc(e.ctx, shard, e.queryFlow, filteringRS, e.fieldIDs, highKey, groupedResult.groupedSeries)
					if err := t.Run(); err != nil {
						e.queryFlow.Complete(err)
						return
					}
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
				if tagValueIDs == nil || tagValueIDs.IsEmpty() {
					e.queryFlow.ReduceTagValues(tagIndex, nil)
					continue
				}
				e.queryFlow.Scanner(func() {
					tagValues := make(map[uint32]string)
					t := newCollectTagValuesTask(e.ctx, e.database.Metadata(), tagKey, tagValueIDs, tagValues)
					if err := t.Run(); err != nil {
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
	if len(e.ctx.shardIDs) == 0 {
		return errNoShardID
	}
	numOfShards := e.database.NumOfShards()
	// check engine has shard
	if numOfShards == 0 {
		return errNoShardInDatabase
	}
	if numOfShards != len(e.ctx.shardIDs) {
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
	numOfShardIDs := len(e.ctx.shardIDs)
	if numOfShards != numOfShardIDs {
		return errShardNumNotMatch
	}
	return nil
}

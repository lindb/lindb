package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestBaseQueryTask_Run(t *testing.T) {
	task := &baseQueryTask{}
	assert.Nil(t, task.Run())
}

func TestQueryStatTask_Run(t *testing.T) {
	task := &queryStatTask{}
	task.BeforeRun()
	task.AfterRun()
}

func TestStoragePlanTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	plan := NewMockPlan(ctrl)
	plan.EXPECT().Plan().Return(nil).AnyTimes()
	// case 1: normal
	task := newStoragePlanTask(newStorageExecuteContext("ns", nil, &stmt.Query{}), plan)
	err := task.Run()
	assert.NoError(t, err)
	// case 2: explain track stats
	task = newStoragePlanTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}), plan)
	err = task.Run()
	assert.NoError(t, err)
}

func TestTagFilterTask_AfterRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagSearch := NewMockTagSearch(ctrl)
	task := newTagFilterTask(newStorageExecuteContext("ns", nil, &stmt.Query{}), tagSearch)
	// case 1: tag filter err
	tagSearch.EXPECT().Filter().Return(nil, fmt.Errorf("err"))
	err := task.Run()
	assert.Error(t, err)
	// case 2: not found
	tagSearch.EXPECT().Filter().Return(nil, nil)
	err = task.Run()
	assert.Equal(t, err, constants.ErrNotFound)
	// case 3: normal
	tagSearch.EXPECT().Filter().Return(map[string]*tagFilterResult{"test": nil}, nil)
	err = task.Run()
	assert.NoError(t, err)
	// case 4: explain case
	task = newTagFilterTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}), tagSearch)
	tagSearch.EXPECT().Filter().Return(map[string]*tagFilterResult{"test": nil}, nil)
	err = task.Run()
	assert.NoError(t, err)
}

func TestSeriesIDsSearchTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newSeriesSearchFunc = newSeriesSearch
		ctrl.Finish()
	}()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	result := roaring.New()
	task := newSeriesIDsSearchTask(newStorageExecuteContext("ns", nil, &stmt.Query{}), shard, result)
	// case 1: search err
	indexDB.EXPECT().GetSeriesIDsForMetric("ns", gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := task.Run()
	assert.Error(t, err)
	// case 2: no group by add series ids without tags
	indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(roaring.New(), nil)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(constants.SeriesIDWithoutTags), result)
	result.Clear()
	// case 3: group by tag
	indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(roaring.New(), nil)
	task = newSeriesIDsSearchTask(newStorageExecuteContext("ns", nil, &stmt.Query{GroupBy: []string{"host"}}), shard, result)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.GetCardinality())
	// case 4: has condition, return err
	q, _ := sql.Parse("select f from cpu where ip<>'1.1.1.1'")
	query := q.(*stmt.Query)
	seriesSearch := NewMockSeriesSearch(ctrl)
	newSeriesSearchFunc = func(filter series.Filter, filterResult map[string]*tagFilterResult, query *stmt.Query) SeriesSearch {
		return seriesSearch
	}
	seriesSearch.EXPECT().Search().Return(nil, fmt.Errorf("err"))
	task = newSeriesIDsSearchTask(newStorageExecuteContext("ns", nil, query), shard, result)
	err = task.Run()
	assert.Error(t, err)
	// case 5: has condition, return series ids
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), result)
	result.Clear()
	// case 6: explain
	q, _ = sql.Parse("explain select f from cpu where ip<>'1.1.1.1'")
	query = q.(*stmt.Query)
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil)
	shard.EXPECT().ShardID().Return(int32(10))
	task = newSeriesIDsSearchTask(newStorageExecuteContext("ns", nil, query), shard, result)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), result)
}

func TestMemoryDataFilterTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	shard.EXPECT().MemoryDatabase().Return(memDB).AnyTimes()
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	result := &filterResultSet{}
	task := newMemoryDataFilterTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		shard, 1, []field.ID{10}, seriesIDs, result)
	// case 1: filter err
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := task.Run()
	assert.Error(t, err)
	// case 2: filter success
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	err = task.Run()
	assert.NoError(t, err)
	// case 4: explain
	task = newMemoryDataFilterTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}),
		shard, 1, []field.ID{10}, seriesIDs, result)
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	shard.EXPECT().ShardID().Return(int32(10))
	err = task.Run()
	assert.NoError(t, err)
}

func TestFileDataFilterTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	result := &filterResultSet{}
	task := newFileDataFilterTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		shard, 1, []field.ID{10}, seriesIDs, result)
	// case 1: get empty family
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil)
	err := task.Run()
	assert.NoError(t, err)
	assert.Nil(t, result.rs)
	// case 2: family filter err
	family := tsdb.NewMockDataFamily(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family}).AnyTimes()
	family.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = task.Run()
	assert.Error(t, err)
	// case 3: get data
	family.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{flow.NewMockFilterResultSet(ctrl)}, nil)
	err = task.Run()
	assert.NoError(t, err)
	assert.NotNil(t, result.rs)
	// case 4: explain
	task = newFileDataFilterTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}),
		shard, 1, []field.ID{10}, seriesIDs, result)
	family.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{flow.NewMockFilterResultSet(ctrl)}, nil)
	shard.EXPECT().ShardID().Return(int32(10))
	err = task.Run()
	assert.NoError(t, err)
	assert.NotNil(t, result.rs)
}

func TestGroupingContextFindTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	result := &groupingResult{}
	task := newGroupingContextFindTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		shard, nil, seriesIDs, result)
	// case 1: get grouping context err
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := task.Run()
	assert.Error(t, err)
	// case 2: get grouping context
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, nil)
	err = task.Run()
	assert.NoError(t, err)
	// case 3: explain
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, nil)
	task = newGroupingContextFindTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}),
		shard, nil, seriesIDs, result)
	shard.EXPECT().ShardID().Return(int32(10))
	err = task.Run()
	assert.NoError(t, err)
}

func TestBuildGroupTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	result := &groupedSeriesResult{}
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	task := newBuildGroupTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		shard, nil, 0, seriesIDs.GetContainer(0), result)
	// case 1: no group
	err := task.Run()
	assert.NoError(t, err)
	// case 2: has grouping
	groupingCtx := series.NewMockGroupingContext(ctrl)
	groupingCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	task = newBuildGroupTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		shard, groupingCtx, 0, seriesIDs.GetContainer(0), result)
	err = task.Run()
	assert.NoError(t, err)
	// case 3: explain
	task = newBuildGroupTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}),
		shard, groupingCtx, 0, seriesIDs.GetContainer(0), result)
	shard.EXPECT().ShardID().Return(int32(10))
	err = task.Run()
	assert.NoError(t, err)
}

func TestDataLoadTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	qf := flow.NewMockStorageQueryFlow(ctrl)
	rs := flow.NewMockFilterResultSet(ctrl)
	task := newDataLoadTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		shard, qf, rs, nil, 1, nil)
	rs.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	// case 1: load data
	err := task.Run()
	assert.NoError(t, err)
	// case 2: explain
	task = newDataLoadTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}),
		shard, qf, rs, nil, 1, nil)
	shard.EXPECT().ShardID().Return(int32(10)).AnyTimes()
	rs.EXPECT().Identifier().Return("memory")
	err = task.Run()
	assert.NoError(t, err)
	rs.EXPECT().Identifier().Return("shard/10/segment/day/20190202/10/1.sst")
	err = task.Run()
	assert.NoError(t, err)
}

func TestCollectTagValuesTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	meta := metadb.NewMockMetadata(ctrl)
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	meta.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	task := newCollectTagValuesTask(newStorageExecuteContext("ns", nil, &stmt.Query{}),
		meta, tag.Meta{ID: 10}, roaring.BitmapOf(1, 2), nil)
	// case 1: collect tag values
	tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err := task.Run()
	assert.NoError(t, err)
	// case 2: explain
	task = newCollectTagValuesTask(newStorageExecuteContext("ns", nil, &stmt.Query{Explain: true}),
		meta, tag.Meta{ID: 10}, roaring.BitmapOf(1, 2), nil)
	err = task.Run()
	assert.NoError(t, err)

}

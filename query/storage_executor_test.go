package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

type mockQueryFlow struct {
}

func (m *mockQueryFlow) Prepare(downSamplingSpecs aggregation.AggregatorSpecs) {
}

func (m *mockQueryFlow) Filtering(task concurrent.Task) {
	task()
}

func (m *mockQueryFlow) Grouping(task concurrent.Task) {
	task()
}

func (m *mockQueryFlow) Scanner(task concurrent.Task) {
	task()
}

func (m *mockQueryFlow) GetAggregator() (agg aggregation.FieldAggregates) {
	return nil
}

func (m *mockQueryFlow) Reduce(tags string, agg aggregation.FieldAggregates) {
}

func (m *mockQueryFlow) Complete(err error) {
}

func newMockQueryFlow() flow.StorageQueryFlow {
	return &mockQueryFlow{}
}

func TestStorageExecute_validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)
	exeCtx := parallel.NewMockExecuteContext(ctrl)
	exeCtx.EXPECT().Complete(gomock.Any()).AnyTimes()

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.Interval(timeutil.OneSecond)}

	// query shards is empty
	exec := newStorageExecutor(queryFlow, mockDatabase, nil, query)
	queryFlow.EXPECT().Complete(errNoShardID)
	exec.Execute()

	// shards of engine is empty
	mockDatabase.EXPECT().NumOfShards().Return(0)
	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1, 2, 3}, query)
	queryFlow.EXPECT().Complete(errNoShardInDatabase)
	exec.Execute()

	// num. of shard not match
	mockDatabase.EXPECT().NumOfShards().Return(2)
	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1, 2, 3}, query)
	queryFlow.EXPECT().Complete(errShardNotMatch)
	exec.Execute()

	mockDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, false).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1, 2, 3}, query)
	queryFlow.EXPECT().Complete(errShardNotFound)
	exec.Execute()

	// normal case
	query, _ = sql.Parse("select f from cpu")
	mockDB1 := newMockDatabase(ctrl)
	exec = newStorageExecutor(queryFlow, mockDB1, []int32{1, 2, 3}, query)
	gomock.InOrder(
		queryFlow.EXPECT().Prepare(gomock.Any()),
		queryFlow.EXPECT().Filtering(gomock.Any()).MaxTimes(3*2), //memory db and shard
	)
	exec.Execute()
}

func TestStorageExecute_Plan_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)
	queryFlow.EXPECT().Complete(gomock.Any()).AnyTimes()

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	mockDatabase.EXPECT().GetShard(gomock.Any()).Return(shard, true).MaxTimes(3)
	mockDatabase.EXPECT().NumOfShards().Return(3)
	idGetter := metadb.NewMockIDGetter(ctrl)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), fmt.Errorf("err"))
	mockDatabase.EXPECT().IDGetter().Return(idGetter).AnyTimes()

	// find metric name err
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(queryFlow, mockDatabase, []int32{1, 2, 3}, query)
	exec.Execute()
}

func TestStorageExecute_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := newMockQueryFlow()

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	idGetter := metadb.NewMockIDGetter(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().Interval().Return(int64(10)).AnyTimes()

	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(3)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().GetShard(int32(2)).Return(shard, true)
	mockDatabase.EXPECT().GetShard(int32(3)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB).MaxTimes(3)
	shard.EXPECT().IndexMetaGetter().Return(nil).MaxTimes(3)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil).MaxTimes(3)
	filterRS := flow.NewMockFilterResultSet(ctrl)
	filterRS.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(3)
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil).MaxTimes(3)

	// normal case with filter
	query, _ := sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(queryFlow, mockDatabase, []int32{1, 2, 3}, query)
	exec.Execute()

	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()

	// normal case without filter
	query, _ = sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")

	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().GetSeriesIDsForMetric(uint32(10), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()
}

func TestStorageExecutor_Execute_GroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := newMockQueryFlow()

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	idGetter := metadb.NewMockIDGetter(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().Interval().Return(int64(10)).AnyTimes()

	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(20), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	groupingCtx := series.NewMockGroupingContext(ctrl)
	groupingCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(map[string][]uint16{"1.1.1.1": {1, 2, 3}})
	memDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(groupingCtx, nil)
	filterRS := flow.NewMockFilterResultSet(ctrl)
	filterRS.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil)

	// normal case
	query, _ := sql.Parse("select f from cpu where host='1.1.1.1' " +
		"and time>'20190729 11:00:00' and time<'20190729 12:00:00' group by type")
	exec := newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()

	// get grouping context err
	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(20), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	memDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil)

	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()

	// get grouping context nil
	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(20), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	memDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	memDB.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil)

	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()
}

func TestStorageExecutor_Execute_Find_Series_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := newMockQueryFlow()

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	idGetter := metadb.NewMockIDGetter(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().Interval().Return(int64(10)).AnyTimes()

	// find series err
	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(nil, series.ErrNotFound)

	query, _ := sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()

	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(1)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	mockDatabase.EXPECT().IDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().MemoryDatabase().Return(memDB)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("err"))
	exec = newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	exec.Execute()
}

func TestStorageExecutor_shardLevelSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := newMockQueryFlow()
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	query, _ := sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(queryFlow, mockDatabase, []int32{1}, query)
	sExec := exec.(*storageExecutor)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil)
	sExec.shardLevelSearch(shard)

	family := tsdb.NewMockDataFamily(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family})
	filter := series.NewMockFilter(ctrl)
	shard.EXPECT().IndexFilter().Return(filter)
	filter.EXPECT().FindSeriesIDsByExpr(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.NewVersion(), roaring.BitmapOf(1, 2)), nil)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB)
	family.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	sExec.shardLevelSearch(shard)

	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family})
	shard.EXPECT().IndexFilter().Return(filter)
	filter.EXPECT().FindSeriesIDsByExpr(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	sExec.shardLevelSearch(shard)
}

func TestStorageExecutor_checkShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)
	queryFlow.EXPECT().Complete(gomock.Any()).AnyTimes()
	queryFlow.EXPECT().Prepare(gomock.Any()).AnyTimes()
	queryFlow.EXPECT().Filtering(gomock.Any()).AnyTimes()

	mockDatabase := newMockDatabase(ctrl)
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(queryFlow, mockDatabase, []int32{1, 2, 3}, query)
	exec.Execute()

	execImpl := exec.(*storageExecutor)
	// check shards error
	execImpl.shardIDs = nil
	assert.NotNil(t, execImpl.checkShards())
}

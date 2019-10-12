package query

import (
	"fmt"
	"testing"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/memdb"
)

func TestStorageExecute_validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exeCtx := parallel.NewMockExecuteContext(ctrl)
	exeCtx.EXPECT().Complete(gomock.Any()).AnyTimes()
	exeCtx.EXPECT().RetainTask(gomock.Any()).AnyTimes()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().GetExecutePool().Return(execPool).AnyTimes()
	engine.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.OneSecond}

	// query shards is empty
	exec := newStorageExecutor(exeCtx, engine, nil, query)
	exec.Execute()

	// shards of engine is empty
	engine.EXPECT().NumOfShards().Return(0)
	exec = newStorageExecutor(exeCtx, engine, []int32{1, 2, 3}, query)
	exec.Execute()

	// num. of shard not match
	engine.EXPECT().NumOfShards().Return(2)
	exec = newStorageExecutor(exeCtx, engine, []int32{1, 2, 3}, query)
	exec.Execute()

	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	engine.EXPECT().GetShard(gomock.Any()).Return(nil).MaxTimes(3)
	exec = newStorageExecutor(exeCtx, engine, []int32{1, 2, 3}, query)
	exec.Execute()

	// normal case
	query, _ = sql.Parse("select f from cpu")
	engine1 := MockTSDBEngine(ctrl)
	engine1.EXPECT().GetExecutePool().Return(execPool)

	exec = newStorageExecutor(exeCtx, engine1, []int32{1, 2, 3}, query)
	exec.Execute()
}

func TestStorageExecute_Plan_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exeCtx := parallel.NewMockExecuteContext(ctrl)
	exeCtx.EXPECT().Complete(gomock.Any()).AnyTimes()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().GetExecutePool().Return(execPool).AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	engine.EXPECT().GetShard(gomock.Any()).Return(shard).MaxTimes(3)
	engine.EXPECT().NumOfShards().Return(3)
	idGetter := diskdb.NewMockIDGetter(ctrl)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), fmt.Errorf("err"))
	engine.EXPECT().GetIDGetter().Return(idGetter).AnyTimes()

	// find metric name err
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(exeCtx, engine, []int32{1, 2, 3}, query)
	exec.Execute()
}

func TestStorageExecute_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exeCtx := parallel.NewMockExecuteContext(ctrl)
	exeCtx.EXPECT().Complete(gomock.Any()).AnyTimes()
	exeCtx.EXPECT().RetainTask(gomock.Any()).AnyTimes()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().GetExecutePool().Return(execPool).AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	idGetter := diskdb.NewMockIDGetter(ctrl)
	family := tsdb.NewMockDataFamily(ctrl)
	filter := series.NewMockFilter(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().Interval().Return(int64(10)).AnyTimes()

	// mock data
	engine.EXPECT().NumOfShards().Return(3)
	engine.EXPECT().GetShard(int32(1)).Return(shard)
	engine.EXPECT().GetShard(int32(2)).Return(shard)
	engine.EXPECT().GetShard(int32(3)).Return(shard)
	engine.EXPECT().GetIDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family, family}).MaxTimes(3)
	shard.EXPECT().GetMemoryDatabase().Return(memDB).MaxTimes(3)
	shard.EXPECT().GetSeriesIDsFilter().Return(filter).MaxTimes(3)
	shard.EXPECT().GetMetaGetter().Return(nil).MaxTimes(3)
	filter.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	filter.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf()), nil)
	filter.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).Return(nil, nil)
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil).MaxTimes(3)
	memDB.EXPECT().Scan(gomock.Any()).MaxTimes(3)
	family.EXPECT().Scan(gomock.Any()).MaxTimes(2 * 3)

	// normal case
	query, _ := sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(exeCtx, engine, []int32{1, 2, 3}, query)
	exec.Execute()
	time.Sleep(100 * time.Millisecond)
	e := exec.(*storageExecutor)
	pool := e.getAggregatorPool(10, 1, &query.TimeRange)
	assert.NotNil(t, pool.Get())

	// find series err
	// mock data
	engine.EXPECT().NumOfShards().Return(1)
	engine.EXPECT().GetShard(int32(1)).Return(shard)
	engine.EXPECT().GetIDGetter().Return(idGetter)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), nil)
	idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family, family})
	shard.EXPECT().GetMemoryDatabase().Return(memDB)
	shard.EXPECT().GetSeriesIDsFilter().Return(filter)
	filter.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("err"))
	memDB.EXPECT().FindSeriesIDsByExpr(uint32(10), gomock.Any(), gomock.Any()).
		Return(nil, series.ErrNotFound)
	exec = newStorageExecutor(exeCtx, engine, []int32{1}, query)
	exec.Execute()
	time.Sleep(100 * time.Millisecond)
}

func TestStorageExecutor_checkShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exeCtx := parallel.NewMockExecuteContext(ctrl)
	exeCtx.EXPECT().Complete(gomock.Any()).AnyTimes()
	exeCtx.EXPECT().RetainTask(gomock.Any()).AnyTimes()

	engine := MockTSDBEngine(ctrl)
	engine.EXPECT().GetExecutePool().Return(execPool).AnyTimes()
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(exeCtx, engine, []int32{1, 2, 3}, query)
	exec.Execute()

	execImpl := exec.(*storageExecutor)
	// check shards error
	execImpl.shardIDs = nil
	assert.NotNil(t, execImpl.checkShards())
}

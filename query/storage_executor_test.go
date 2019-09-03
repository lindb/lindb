package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.OneSecond}

	// query shards is empty
	exec := newStorageExecutor(context.TODO(), engine, nil, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// shards of engine is empty
	engine.EXPECT().NumOfShards().Return(0)
	exec = newStorageExecutor(context.TODO(), engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// num. of shard not match
	engine.EXPECT().NumOfShards().Return(2)
	exec = newStorageExecutor(context.TODO(), engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	engine.EXPECT().GetShard(gomock.Any()).Return(nil).MaxTimes(3)
	exec = newStorageExecutor(context.TODO(), engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// normal case
	query, _ = sql.Parse("select f from cpu")
	engine1 := MockTSDBEngine(ctrl)

	exec = newStorageExecutor(context.TODO(), engine1, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())
}

func TestStorageExecute_Plan_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := tsdb.NewMockEngine(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	engine.EXPECT().GetShard(gomock.Any()).Return(shard).MaxTimes(3)
	engine.EXPECT().NumOfShards().Return(3)
	idGetter := diskdb.NewMockIDGetter(ctrl)
	idGetter.EXPECT().GetMetricID("cpu").Return(uint32(10), fmt.Errorf("err"))
	engine.EXPECT().GetIDGetter().Return(idGetter).AnyTimes()

	// find metric name err
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(context.TODO(), engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

}

func TestStorageExecute_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := tsdb.NewMockEngine(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	idGetter := diskdb.NewMockIDGetter(ctrl)
	family := tsdb.NewMockDataFamily(ctrl)
	filter := series.NewMockFilter(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)

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
	exec := newStorageExecutor(context.TODO(), engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())

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
	exec = newStorageExecutor(context.TODO(), engine, []int32{1}, query)
	rs := exec.Execute()
	assert.NotNil(t, exec.Error())
	count := 0
	for range rs {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestStorageExecutor_checkShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := MockTSDBEngine(ctrl)
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(context.TODO(), engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())

	execImpl := exec.(*storageExecutor)
	// check shards error
	execImpl.shardIDs = nil
	assert.NotNil(t, execImpl.checkShards())
}

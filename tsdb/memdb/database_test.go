package memdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

var cfg = MemoryDatabaseCfg{
	TimeWindow: 32,
	Interval:   timeutil.Interval(10 * timeutil.OneSecond),
}

func Test_NewMemoryDatabase(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	assert.NotNil(t, mdINTF)
	assert.Equal(t, int64(10*1000), mdINTF.Interval())
}

func Test_MemoryDatabase_Write(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock generator
	mockGen := metadb.NewMockIDGenerator(ctrl)
	mockGen.EXPECT().GenMetricID("test1").Return(uint32(1)).AnyTimes()
	mockIndex := indexdb.NewMockIndexDatabase(ctrl)
	// build memory-database
	cfg.Generator = mockGen
	cfg.Index = mockIndex
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	// write error
	gomock.InOrder(
		mockIndex.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint32(10), nil),
		mockMStore.EXPECT().Write(uint32(10), gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("error")),
	)
	// load mock
	md.mStores.put(uint32(1), mockMStore)
	// write error
	err := md.Write(&pb.Metric{Name: "test1", Timestamp: 1564300800000})
	assert.Error(t, err)
	assert.Nil(t, md.Families())
	// write ok
	gomock.InOrder(
		mockIndex.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint32(10), nil),
		mockMStore.EXPECT().Write(uint32(10), gomock.Any(), gomock.Any()).Return(20, nil).AnyTimes(),
	)
	err = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564300800000})
	assert.NoError(t, err)
	// test families
	gomock.InOrder(
		mockIndex.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint32(10), nil),
		mockMStore.EXPECT().Write(uint32(10), gomock.Any(), gomock.Any()).Return(20, nil).AnyTimes(),
	)
	_ = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564297200000})
	gomock.InOrder(
		mockIndex.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint32(10), nil),
		mockMStore.EXPECT().Write(uint32(10), gomock.Any(), gomock.Any()).Return(20, nil).AnyTimes(),
	)
	_ = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564308000000})
	assert.NotNil(t, md.Families())
	assert.Len(t, md.Families(), 3)
}

func Test_MemoryDatabase_FlushFamilyTo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fluster := metricsdata.NewMockFlusher(ctrl)
	fluster.EXPECT().Commit().AnyTimes()
	mdINTF := NewMemoryDatabase(ctx, cfg)
	_ = mdINTF.FlushFamilyTo(fluster, 10)
	_ = mdINTF.FlushFamilyTo(fluster, 10)
	_ = mdINTF.FlushFamilyTo(fluster, 10)
	time.Sleep(time.Millisecond * 10)
}

func Test_MemoryDatabase_flushFamilyTo_ok(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fluster := metricsdata.NewMockFlusher(ctrl)
	fluster.EXPECT().Commit().AnyTimes()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	mockMStore := NewMockmStoreINTF(ctrl)

	returnNil := mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(nil)
	returnError := mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	gomock.InOrder(returnNil, returnError)

	md.mStores.put(uint32(1), mockMStore)
	assert.Nil(t, md.FlushFamilyTo(fluster, 10))
	assert.NotNil(t, md.FlushFamilyTo(fluster, 10))
}

func Test_MemoryDatabase_Filter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	// not found
	_, _ = md.Filter(0, []uint16{1}, 1, nil)

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	md.mStores.put(uint32(3333), mockMStore)

	_, _ = md.Filter(uint32(3333), []uint16{1}, 1, nil)
}

func Test_MemoryDatabase_MemSize(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	assert.Zero(t, md.MemSize())
}

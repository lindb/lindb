package memdb

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"

	"github.com/cespare/xxhash"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	count := uint32(0)
	mockGen.EXPECT().GenMetricID("test1").
		Do(func() {
			count++
		}).Return(count).AnyTimes()

	// build memory-database
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	md.generator = mockGen

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().GetMetricID().Return(uint32(1)).AnyTimes()
	errCall1 := mockMStore.EXPECT().Write(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("error"))
	okCall2 := mockMStore.EXPECT().Write(gomock.Any(), gomock.Any()).Return(20, nil).AnyTimes()
	gomock.InOrder(errCall1, okCall2)
	// load mock
	hash := xxhash.Sum64String("test1")
	md.getBucket(hash).hash2MStore[hash] = mockMStore
	// write error
	err := md.Write(&pb.Metric{Name: "test1", Timestamp: 1564300800000})
	assert.NotNil(t, err)
	assert.Nil(t, md.Families())
	// write ok
	err = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564300800000})
	assert.Nil(t, err)
	// test families
	_ = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564297200000})
	_ = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564308000000})
	assert.NotNil(t, md.Families())
	assert.Len(t, md.Families(), 3)
}

func Test_MemoryDatabase_setLimitations_countTags_countMetrics_resetMStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	md.generator = makeMockIDGenerator(ctrl)
	// count metrics
	assert.Equal(t, 0, md.CountMetrics())

	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().SetMaxTagsLimit(gomock.Any()).Return().AnyTimes()
	mockMStore.EXPECT().GetTagsUsed().Return(1).AnyTimes()
	mockMStore.EXPECT().ResetVersion().Return(100, nil).AnyTimes()
	// setLimitations
	limitations := map[string]uint32{"cpu.load": 10, "memory": 100}
	hash := xxhash.Sum64String("cpu.load")
	md.getOrCreateMStore("cpu.load", hash)
	md.getBucket(hash).hash2MStore[hash] = mockMStore
	md.setLimitations(limitations)

	// countTags
	assert.Equal(t, -1, md.CountTags("cpu.load1"))
	assert.Equal(t, 1, md.CountTags("cpu.load"))

	// count metrics
	assert.Equal(t, 1, md.CountMetrics())

	// reset mStore
	assert.NotNil(t, md.ResetMetricStore("cpu.load2"))
	assert.Nil(t, md.ResetMetricStore("cpu.load"))
}

func Test_MemoryDatabase_WithMaxTagsLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	limitationCh := make(chan map[string]uint32)
	md.WithMaxTagsLimit(limitationCh)
	md.WithMaxTagsLimit(limitationCh)

	limitationCh <- nil
	limitationCh <- map[string]uint32{"cpu.load": 10}
	time.Sleep(time.Millisecond * 10)

	close(limitationCh)
	time.Sleep(time.Millisecond * 10)
}

func Test_MemoryDatabase_WithMaxTagsLimit_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	mdINTF := NewMemoryDatabase(ctx, cfg)
	limitationCh := make(chan map[string]uint32)
	mdINTF.WithMaxTagsLimit(limitationCh)
	cancel()
	time.Sleep(time.Millisecond * 10)
}

func Test_MemoryDatabase_evict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// mock generator
	mockGen := metadb.NewMockIDGenerator(ctrl)
	for i := 0; i < 1000; i++ {
		mockGen.EXPECT().GenMetricID(strconv.Itoa(i)).Return(uint32(i)).AnyTimes()
	}
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	md.generator = mockGen
	// prepare mStores
	for i := 0; i < 1000; i++ {
		md.getOrCreateMStore(strconv.Itoa(i), xxhash.Sum64String(strconv.Itoa(i)))
	}
	// evict all
	for _, store := range md.mStoresList {
		md.evict(store)
	}
}

func Test_MemoryDatabase_evictor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	md.evictNotifier <- struct{}{}
	md.evictNotifier <- struct{}{}
	md.evictNotifier <- struct{}{}
	time.Sleep(time.Millisecond * 10)
}

func Test_FindSeriesIDsByExpr_GetSeriesIDsForTag(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().FindSeriesIDsByExpr(gomock.Any()).Return(nil, nil).AnyTimes()
	mockMStore.EXPECT().GetSeriesIDsForTag("").Return(nil, nil).AnyTimes()
	// not exist
	_, err := md.FindSeriesIDsByExpr(1, nil, timeutil.TimeRange{})
	assert.NotNil(t, err)
	_, err = md.GetSeriesIDsForTag(1, "", timeutil.TimeRange{})
	assert.NotNil(t, err)
	// exist
	md.getBucket(3333).hash2MStore[3333] = mockMStore
	md.metricID2Hash.Store(uint32(1), uint64(3333))
	_, err = md.FindSeriesIDsByExpr(1, nil, timeutil.TimeRange{})
	assert.Nil(t, err)
	_, err = md.GetSeriesIDsForTag(1, "", timeutil.TimeRange{})
	assert.Nil(t, err)
}

func Test_MemoryDatabase_FlushFamilyTo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	_ = mdINTF.FlushFamilyTo(nil, 10)
	_ = mdINTF.FlushFamilyTo(nil, 10)
	_ = mdINTF.FlushFamilyTo(nil, 10)
	time.Sleep(time.Millisecond * 10)
}

func Test_MemoryDatabase_flushFamilyTo_ok(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().GetMetricID().Return(uint32(1)).AnyTimes()
	mockMStore.EXPECT().Evict().Return(100).AnyTimes()
	mockMStore.EXPECT().IsEmpty().Return(false).AnyTimes()

	returnNil := mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(100, nil)
	returnError := mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("error"))
	gomock.InOrder(returnNil, returnError)

	md.getBucket(4).hash2MStore[1] = mockMStore
	md.getBucket(4).familyTimes = map[int64]struct{}{33: {}}
	assert.Nil(t, md.FlushFamilyTo(nil, 10))
	assert.NotNil(t, md.FlushFamilyTo(nil, 10))
}

func Test_MemoryDatabase_flushIndexTo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	// test FlushIndexTo
	assert.Nil(t, mdINTF.FlushInvertedIndexTo(nil))
	assert.Nil(t, mdINTF.FlushForwardIndexTo(nil))

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	gomock.InOrder(
		mockMStore.EXPECT().FlushInvertedIndexTo(gomock.Any(), gomock.Any()).Return(nil),
		mockMStore.EXPECT().FlushInvertedIndexTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")),
		mockMStore.EXPECT().FlushForwardIndexTo(gomock.Any()).Return(nil),
		mockMStore.EXPECT().FlushForwardIndexTo(gomock.Any()).Return(fmt.Errorf("error")),
	)
	// insert to bucket
	md.getBucket(4).hash2MStore[1] = mockMStore
	// test flushInvertedIndexTo
	assert.Nil(t, md.FlushInvertedIndexTo(nil))
	assert.NotNil(t, md.FlushInvertedIndexTo(nil))
	// test flushForwardIndexTo
	assert.Nil(t, md.FlushForwardIndexTo(nil))
	assert.NotNil(t, md.FlushForwardIndexTo(nil))
}

func Test_MemoryDatabase_GetTagValues(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().GetTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	md.getBucket(3333).hash2MStore[3333] = mockMStore
	md.metricID2Hash.Store(uint32(3333), uint64(3333))

	// existed metricID
	_, err := mdINTF.GetTagValues(3333, nil, 1, nil)
	assert.Nil(t, err)
	// inexisted metricID
	_, err = mdINTF.GetTagValues(3334, nil, 1, nil)
	assert.NotNil(t, err)

}

func Test_MemoryDatabase_Suggset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	assert.Nil(t, md.SuggestMetrics("", 100))
	assert.Nil(t, md.SuggestTagKeys("", "", 100))
	assert.Nil(t, md.SuggestTagValues("", "", "", 100))

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().SuggestTagKeys(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockMStore.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	md.getBucket(xxhash.Sum64String("test")).hash2MStore[xxhash.Sum64String("test")] = mockMStore

	assert.Nil(t, md.SuggestTagKeys("test", "", 100))
	assert.Nil(t, md.SuggestTagValues("test", "", "", 100))
}

func Test_MemoryDatabase_Scan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(ctx, cfg)
	md := mdINTF.(*memoryDatabase)

	// not found
	md.Scan(&series.ScanContext{MetricID: 0})

	// mock mStore
	sCtx := &series.ScanContext{MetricID: 3333}
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Scan(sCtx)
	md.metricID2Hash.Store(uint32(3333), xxhash.Sum64String("test"))
	md.getBucket(xxhash.Sum64String("test")).hash2MStore[xxhash.Sum64String("test")] = mockMStore
	md.Scan(sCtx)
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

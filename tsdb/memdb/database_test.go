package memdb

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/indexdb"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/stretchr/testify/assert"
)

func Test_NewMemoryDatabase(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mdINTF, err := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	assert.Nil(t, err)
	assert.NotNil(t, mdINTF)

	mdINTF, err = NewMemoryDatabase(ctx, 32, 10*1000, interval.Type(3232323))
	assert.Nil(t, mdINTF)
	assert.NotNil(t, err)
}

func Test_MemoryDatabase_Write(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock generator
	mockGen := indexdb.NewMockIDGenerator(ctrl)
	count := uint32(0)
	mockGen.EXPECT().GenMetricID("test1").
		Do(func() {
			count++
		}).Return(count).AnyTimes()

	// build memory-database
	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	md := mdINTF.(*memoryDatabase)
	md.generator = mockGen

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().getMetricID().Return(uint32(1)).AnyTimes()
	errCall1 := mockMStore.EXPECT().write(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	okCall2 := mockMStore.EXPECT().write(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	gomock.InOrder(errCall1, okCall2)
	// load mock
	hash := fnv1a.HashString64("test1")
	md.getBucket(hash).hash2MStore[hash] = mockMStore
	// write error
	err := md.Write(&pb.Metric{Name: "test1", Timestamp: 1564300800000})
	assert.NotNil(t, err)
	assert.Nil(t, md.Families())
	// write ok
	err = md.Write(&pb.Metric{Name: "test1", Timestamp: 1564300800000})
	assert.Nil(t, err)
	// test families
	md.Write(&pb.Metric{Name: "test1", Timestamp: 1564297200000})
	md.Write(&pb.Metric{Name: "test1", Timestamp: 1564308000000})
	assert.NotNil(t, md.Families())
	assert.Len(t, md.Families(), 3)
}

func Test_MemoryDatabase_setLimitations_countTags_countMetrics_resetMStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	md := mdINTF.(*memoryDatabase)
	md.generator = makeMockIDGenerator(ctrl)
	// count metrics
	assert.Equal(t, 0, md.CountMetrics())

	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().setMaxTagsLimit(gomock.Any()).Return().AnyTimes()
	mockMStore.EXPECT().getTagsCount().Return(1).AnyTimes()
	mockMStore.EXPECT().resetVersion().Return(nil).AnyTimes()
	// setLimitations
	limitations := map[string]uint32{"cpu.load": 10, "memory": 100}
	hash := fnv1a.HashString64("cpu.load")
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

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
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

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
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
	mockGen := indexdb.NewMockIDGenerator(ctrl)
	for i := 0; i < 1000; i++ {
		mockGen.EXPECT().GenMetricID(strconv.Itoa(i)).Return(uint32(i)).AnyTimes()
	}
	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	md := mdINTF.(*memoryDatabase)
	md.generator = mockGen
	// prepare mStores
	for i := 0; i < 1000; i++ {
		md.getOrCreateMStore(strconv.Itoa(i), fnv1a.HashString64(strconv.Itoa(i)))
	}
	// evict all
	for _, store := range md.mStoresList {
		md.evict(store)
	}
}

func Test_MemoryDatabase_evictor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
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

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	md := mdINTF.(*memoryDatabase)
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().findSeriesIDsByExpr(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockMStore.EXPECT().getSeriesIDsForTag("", gomock.Any()).Return(nil, nil).AnyTimes()
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

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	mdINTF.FlushFamilyTo(nil, 10)
	mdINTF.FlushFamilyTo(nil, 10)
	mdINTF.FlushFamilyTo(nil, 10)
	time.Sleep(time.Millisecond * 10)
}

func Test_MemoryDatabase_flushFamilyTo_ok(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	md := mdINTF.(*memoryDatabase)

	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().getMetricID().Return(uint32(1)).AnyTimes()
	mockMStore.EXPECT().evict().Return().AnyTimes()
	mockMStore.EXPECT().isEmpty().Return(false).AnyTimes()

	returnNil := mockMStore.EXPECT().flushMetricsTo(gomock.Any(), gomock.Any()).Return(nil)
	returnError := mockMStore.EXPECT().flushMetricsTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	gomock.InOrder(returnNil, returnError)

	md.getBucket(4).hash2MStore[1] = mockMStore
	md.getBucket(4).familyTimes = map[int64]struct{}{33: {}}
	assert.Nil(t, md.flushFamilyTo(nil, 10))
	assert.NotNil(t, md.flushFamilyTo(nil, 10))
}

func Test_FlushSeriesIndexTo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mdINTF, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)
	md := mdINTF.(*memoryDatabase)
	// test FlushSeriesIndexTo
	mdINTF.FlushSeriesIndexTo(nil)
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	gomock.InOrder(
		mockMStore.EXPECT().flushIndexesTo(gomock.Any(), gomock.Any()).Return(nil),
		mockMStore.EXPECT().flushIndexesTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")),
	)
	// insert to bucket
	md.getBucket(4).hash2MStore[1] = mockMStore
	// test flushSeriesIndexTo
	assert.Nil(t, md.flushSeriesIndexTo(nil))
	assert.NotNil(t, md.flushSeriesIndexTo(nil))

}

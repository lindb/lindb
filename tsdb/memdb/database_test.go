package memdb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/eleme/lindb/pkg/interval"

	"github.com/stretchr/testify/assert"
)

func Test_NewMemoryDatabase_GetVersion(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := NewMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	assert.NotNil(t, md)
	assert.NotZero(t, md.GetVersion())
}

func Test_getBucket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	for i := 0; i < 1000; i++ {
		assert.NotNil(t, md.getBucket(strconv.Itoa(i)))
	}
}

func Test_getMetricStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	for i := 0; i < 1000; i++ {
		assert.NotNil(t, md.getMetricStore(strconv.Itoa(i)))
	}
}

func Test_memoryDatabase_PrefixSearchMetricNames(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	md.getMetricStore("abc")
	md.getMetricStore("abcd")
	md.getMetricStore("xyz")
	md.getMetricStore("ab")

	assert.Len(t, md.PrefixSearchMetricNames("", 10), 0)
	assert.Len(t, md.PrefixSearchMetricNames("bcd", 10), 0)
	assert.Len(t, md.PrefixSearchMetricNames("abcd", 10), 1)
	assert.Len(t, md.PrefixSearchMetricNames("ab", 10), 3)
	assert.Len(t, md.PrefixSearchMetricNames("ab", 2), 2)
	assert.Len(t, md.PrefixSearchMetricNames("ab", 0), 0)
}

func Test_memoryDatabase_RegexSearchTags(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	md.getMetricStore("cpu.load").getTimeSeries("host=alpha-32,ezone=nj")
	assert.Len(t, md.RegexSearchTags("loadavg", "host=alpha.*"), 0)
	assert.Len(t, md.RegexSearchTags("cpu.load", "host=alpha.*"), 1)
}

func Test_setLimitations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	limitations := map[string]uint32{"cpu.load": 10}
	md.getMetricStore("cpu.load")
	md.getMetricStore("loadavg")

	md.setLimitations(limitations)
	assert.Equal(t, uint32(10), md.getMetricStore("cpu.load").getMaxTagsLimit())
	assert.NotEqual(t, uint32(10), md.getMetricStore("loadavg").getMaxTagsLimit())
}

func Test_WithMaxTagsLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	md.getMetricStore("cpu.load")
	limitationCh := make(chan map[string]uint32)
	md.WithMaxTagsLimit(limitationCh)
	md.WithMaxTagsLimit(limitationCh)

	limitationCh <- nil
	limitationCh <- map[string]uint32{"cpu.load": 10}
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, uint32(10), md.getMetricStore("cpu.load").getMaxTagsLimit())

	close(limitationCh)
}

// func Test_Write(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

// 	assert.NotNil(t, md.Write(nil))

// 	p := models.NewMockPoint(ctrl)
// 	p.EXPECT().Name().Return("cpu.load").AnyTimes()
// 	p.EXPECT().TagsID().Return("idle").AnyTimes()
// 	p.EXPECT().Timestamp().Return(util.Now()).AnyTimes()
// 	p.EXPECT().Fields().Return(nil).Times(1)
// 	assert.NotNil(t, md.Write(p))

// 	fakeFields := map[string]models.Field{"a": nil, "b": nil}
// 	p.EXPECT().Fields().Return(fakeFields).AnyTimes()
// 	assert.Nil(t, md.Write(p))

// 	// assert error
// 	mStore := md.getMetricStore("cpu.load")
// 	for i := 0; i < 20000; i++ {
// 		mStore.getTimeSeries(strconv.Itoa(i))
// 	}
// 	assert.Equal(t, models.ErrTooManyTags, md.Write(p))
// }

func Test_evict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md, _ := newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	setTagsIDTTL(60 * 1000) // 60 s
	for i := 0; i < 1000; i++ {
		md.getMetricStore(strconv.Itoa(i))
	}
	for _, store := range md.mStoresList {
		md.evict(store)
	}
	// purges all
	assert.Equal(t, 0, len(md.PrefixSearchMetricNames("1", 10)))
	assert.Equal(t, 0, len(md.mStoresList[0].m))
}

func Test_evictorRunner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newMemoryDatabase(ctx, 32, 10*1000, interval.Day)

	setTagsIDTTL(1)
	setEvictInterval(10) // 10 ms
	time.Sleep(time.Millisecond * 1000)
}

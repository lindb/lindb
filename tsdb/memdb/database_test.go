package memdb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/eleme/lindb/models"

	"github.com/stretchr/testify/assert"
)

func Test_NewMemoryDatabase_GetVersion(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := NewMemoryDatabase(ctx)

	assert.NotNil(t, md)
	assert.NotZero(t, md.GetVersion())
}

func Test_getBucket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	for i := 0; i < 1000; i++ {
		assert.NotNil(t, md.getBucket(strconv.Itoa(i)))
	}
}

func Test_getMeasurementStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	for i := 0; i < 1000; i++ {
		assert.NotNil(t, md.getMeasurementStore(strconv.Itoa(i)))
	}
}

func Test_memoryDatabase_PrefixSearchMeasurements(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	md.getMeasurementStore("abc")
	md.getMeasurementStore("abcd")
	md.getMeasurementStore("xyz")
	md.getMeasurementStore("ab")

	assert.Len(t, md.PrefixSearchMeasurements("", 10), 0)
	assert.Len(t, md.PrefixSearchMeasurements("bcd", 10), 0)
	assert.Len(t, md.PrefixSearchMeasurements("abcd", 10), 1)
	assert.Len(t, md.PrefixSearchMeasurements("ab", 10), 3)
	assert.Len(t, md.PrefixSearchMeasurements("ab", 2), 2)
	assert.Len(t, md.PrefixSearchMeasurements("ab", 0), 0)
}

func Test_memoryDatabase_RegexSearchTags(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	md.getMeasurementStore("cpu.load").getTimeSeries("host=alpha-32,ezone=nj")
	assert.Len(t, md.RegexSearchTags("loadavg", "host=alpha.*"), 0)
	assert.Len(t, md.RegexSearchTags("cpu.load", "host=alpha.*"), 1)
}

func Test_setLimitations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	limitations := map[string]uint32{"cpu.load": 10}
	md.getMeasurementStore("cpu.load")
	md.getMeasurementStore("loadavg")

	md.setLimitations(limitations)
	assert.Equal(t, uint32(10), md.getMeasurementStore("cpu.load").getMaxTagsLimit())
	assert.NotEqual(t, uint32(10), md.getMeasurementStore("loadavg").getMaxTagsLimit())
}

func Test_WithMaxTagsLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	md.getMeasurementStore("cpu.load")
	limitationCh := make(chan map[string]uint32)
	md.WithMaxTagsLimit(limitationCh)
	md.WithMaxTagsLimit(limitationCh)

	limitationCh <- nil
	limitationCh <- map[string]uint32{"cpu.load": 10}
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, uint32(10), md.getMeasurementStore("cpu.load").getMaxTagsLimit())

	close(limitationCh)
}

func Test_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	assert.NotNil(t, md.Write(nil, 0, 0))

	p := models.NewMockPoint(ctrl)
	p.EXPECT().Name().Return("cpu.load").AnyTimes()
	p.EXPECT().TagsID().Return("idle").AnyTimes()

	p.EXPECT().Fields().Return(nil).Times(1)
	assert.NotNil(t, md.Write(p, 0, 0))

	fakeFields := map[string]models.Field{"a": nil, "b": nil}
	p.EXPECT().Fields().Return(fakeFields).AnyTimes()
	assert.Nil(t, md.Write(p, 0, 0))

	// assert error
	mStore := md.getMeasurementStore("cpu.load")
	for i := 0; i < 20000; i++ {
		mStore.getTimeSeries(strconv.Itoa(i))
	}
	assert.Equal(t, models.ErrTooManyTags, md.Write(p, 0, 0))
}

func Test_evict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	md := newMemoryDatabase(ctx)

	setTagsIDTTL(60 * 1000) // 60 s
	for i := 0; i < 1000; i++ {
		md.getMeasurementStore(strconv.Itoa(i))
	}
	for _, store := range md.mStoresList {
		md.evict(store)
	}
	// purges all
	assert.Equal(t, 0, len(md.PrefixSearchMeasurements("1", 10)))
	assert.Equal(t, 0, len(md.mStoresList[0].m))
}

func Test_evictorRunner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newMemoryDatabase(ctx)

	setTagsIDTTL(1)
	setEvictInterval(10) // 10 ms
	time.Sleep(time.Millisecond * 1000)
}

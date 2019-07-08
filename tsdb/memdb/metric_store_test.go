package memdb

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newMetricStore(t *testing.T) {
	mStore := newMetricStore("cpu.load")
	assert.NotNil(t, mStore)
	assert.NotNil(t, mStore.tsMap)
	assert.NotZero(t, mStore.maxTagsLimit)
}

func Test_metricStore_isEmpty_isFull(t *testing.T) {
	mStore := newMetricStore("cpu.load")
	assert.True(t, mStore.isEmpty())
	assert.False(t, mStore.isFull())

	for i := 0; i < 100; i++ {
		mStore.tsMap[strconv.Itoa(i)] = nil
	}
	assert.False(t, mStore.isFull())
	assert.False(t, mStore.isEmpty())

	for i := 0; i < defaultMaxTagsLimit; i++ {
		mStore.tsMap[strconv.Itoa(i)] = nil
	}
	assert.True(t, mStore.isFull())
	assert.False(t, mStore.isEmpty())
}

func Test_metricStore_regexSearchTags(t *testing.T) {
	mStore := newMetricStore("cpu.load")

	// invalid pattern
	assert.Nil(t, mStore.regexSearchTags(""))
	assert.Nil(t, mStore.regexSearchTags(`[\w-32]`))

	tagsIDFormat := "host=alpha-%d,ezone=nj,ip=192.168.1.1"
	for i := 0; i < 100; i++ {
		mStore.tsMap[fmt.Sprintf(tagsIDFormat, i)] = nil
	}

	assert.Len(t, mStore.regexSearchTags("host=alpha-102.*,ezone=nj"), 0)
	assert.Len(t, mStore.regexSearchTags("host=alpha-10.*,ezone=gz"), 0)
	assert.Len(t, mStore.regexSearchTags("host=alpha-10.*,ezone=nj"), 1)
	assert.Len(t, mStore.regexSearchTags("host=alpha-1.*,ezone=nj"), 11)

	matched := mStore.regexSearchTags("host=alpha-2.*,ezone=nj")
	assert.Contains(t, matched[0], "alpha-2")
	assert.Contains(t, matched[10], "alpha-29")
}

func Test_metricStore_getTimeSeries(t *testing.T) {
	mStore := newMetricStore("cpu.load")

	assert.NotNil(t, mStore.getTimeSeries("host=alpha-1"))
	assert.Equal(t, mStore.getTimeSeries("host=alpha-2"), mStore.getTimeSeries("host=alpha-2"))
	assert.Equal(t, mStore.getTagsCount(), 2)
}

func Test_metricStore_evict(t *testing.T) {
	mStore := newMetricStore("cpu.load")
	mStore.evict()
	assert.True(t, mStore.isEmpty())
	// has not been purged
	for i := 0; i < 2000; i++ {
		mStore.getTimeSeries(strconv.Itoa(i)).getFieldStore("t")
	}
	setTagsIDTTL(60 * 1000) // 1 minute
	mStore.evict()
	assert.Equal(t, 2000, mStore.getTagsCount())

	// purge half
	time.Sleep(time.Millisecond * 20)
	setTagsIDTTL(20) // 20 ms
	for i := 0; i < 1000; i++ {
		mStore.getTimeSeries(strconv.Itoa(i)).getFieldStore("t")
	}
	mStore.evict()
	assert.Equal(t, 1000, mStore.getTagsCount())
	// purge all
	time.Sleep(time.Millisecond * 20)
	setTagsIDTTL(20) // 20 ms
	mStore.evict()
	assert.True(t, mStore.isEmpty())
}

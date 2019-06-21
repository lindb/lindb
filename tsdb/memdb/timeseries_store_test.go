package memdb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newTimeSeriesStore(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	assert.NotNil(t, tsStore)
	assert.Nil(t, tsStore.element)
	assert.NotZero(t, tsStore.lastAccessedAt)
}

func Test_getFieldStore(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	tsStore.lastAccessedAt = 0

	fStore := tsStore.getFieldStore("idle")
	assert.NotNil(t, fStore)
	assert.NotEqual(t, int64(0), tsStore.lastAccessedAt)
}

func Test_shouldBeEvicted(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	assert.False(t, tsStore.shouldBeEvicted())

	setTagsIDTTL(1) // 1 ms
	time.Sleep(time.Millisecond)
	assert.True(t, tsStore.shouldBeEvicted())
}

func Test_getFieldsCount(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	assert.Equal(t, 0, tsStore.getFieldsCount())

	tsStore.getFieldStore("idle")
	assert.Equal(t, 1, tsStore.getFieldsCount())
	tsStore.getFieldStore("idle")
	assert.Equal(t, 1, tsStore.getFieldsCount())
}

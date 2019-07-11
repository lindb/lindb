package memdb

import (
	"testing"
	"time"

	"github.com/eleme/lindb/pkg/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_newTimeSeriesStore(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	assert.NotNil(t, tsStore)
	assert.NotZero(t, tsStore.lastAccessedAt)
}

func Test_mustGetTSID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gen := makeMockIDGenerator(ctrl)
	tsStore := newTimeSeriesStore("host=alpha")

	assert.NotZero(t, tsStore.mustGetTSID(32, gen))
	assert.NotZero(t, tsStore.mustGetTSID(32, gen))
}

func Test_getOrCreateFStore(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	tsStore.lastAccessedAt = 0

	fStore, err := tsStore.getOrCreateFStore("idle", field.MaxField)
	assert.NotNil(t, fStore)
	assert.Nil(t, err)
	assert.NotEqual(t, int64(0), tsStore.lastAccessedAt)

	fStore, err = tsStore.getOrCreateFStore("idle", field.SumField)
	assert.Nil(t, fStore)
	assert.NotNil(t, err)
	assert.NotEqual(t, int64(0), tsStore.lastAccessedAt)
}

func Test_shouldBeEvicted(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	fStore := newFieldStore("sum", field.SumField)

	tsStore.fields[1] = fStore
	assert.False(t, tsStore.shouldBeEvicted())

	fStore.segments[1] = nil
	tsStore.fields[1] = fStore
	assert.False(t, tsStore.shouldBeEvicted())

	delete(tsStore.fields, 1)
	setTagsIDTTL(1) // 1 ms
	time.Sleep(time.Millisecond)
	assert.True(t, tsStore.shouldBeEvicted())
}

func Test_getFieldsCount(t *testing.T) {
	tsStore := newTimeSeriesStore("host=alpha")
	assert.Equal(t, 0, tsStore.getFieldsCount())

	tsStore.getOrCreateFStore("idle", field.MaxField)
	assert.Equal(t, 1, tsStore.getFieldsCount())
	tsStore.getOrCreateFStore("idle", field.MaxField)
	assert.Equal(t, 1, tsStore.getFieldsCount())
}

func Test_flushTSEntryTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tsStore := newTimeSeriesStore("host=alpha")
	tw := makeMockTableWriter(ctrl)
	gen := makeMockIDGenerator(ctrl)

	tsStore.getOrCreateFStore("idle", field.MaxField)
	tsStore.getOrCreateFStore("system", field.MaxField)

	tsStore.flushTSEntryTo(tw, 3, 2, gen)
}

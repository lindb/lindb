package memdb

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/series"
)

func Test_mStore_scan(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	mStore.fieldsMetas = fieldsMetas{
		{"sum3", 3, field.SumField},
		{"sum4", 4, field.SumField},
		{"sum5", 5, field.SumField},
		{"sum6", 6, field.SumField},
	}
	// v1:
	ti1 := newTagIndex().(*tagIndex)
	ti1.version = 1
	ti1.startTime, ti1.endTime = 100, 200
	// v2
	ti2 := newTagIndex().(*tagIndex)
	ti2.version = 2
	ti2.startTime, ti2.endTime = 200, 300
	ts5 := newTimeSeriesStore(55)
	ts6 := newTimeSeriesStore(66)
	ts7 := newTimeSeriesStore(77)
	ts8 := newTimeSeriesStore(88)
	ts5.(*timeSeriesStore).insertFStore(newFieldStore(1))
	ts5.(*timeSeriesStore).insertFStore(newFieldStore(2))
	ts5.(*timeSeriesStore).insertFStore(newFieldStore(3))
	ts5.(*timeSeriesStore).insertFStore(newFieldStore(4))

	ti2.seriesID2TStore = map[uint32]tStoreINTF{
		5: ts5,
		6: ts6,
		7: ts7,
		8: ts8,
	}
	// build id-set
	idset := series.NewMultiVerSeriesIDSet()
	idset.Add(0, roaring.New())
	bitmap := roaring.New()
	bitmap.AddMany([]uint32{1, 2, 3, 4, 5, 7})
	idset.Add(2, bitmap)

	// build mStore
	mStore.immutable = []tagIndexINTF{ti1}
	mStore.mutable = ti2

	// timeRange not overlaps
	itr1 := mStore.scan(series.ScanContext{SeriesIDSet: idset,
		TimeRange: timeutil.TimeRange{Start: 0, End: 0}})
	assert.False(t, itr1.HasNext())
	assert.Nil(t, itr1.Close())
	assert.False(t, itr1.Next().HasNext())
	assert.Zero(t, itr1.Version())
	// timeRange overlaps
	itr2 := mStore.scan(series.ScanContext{
		FieldIDs:    []uint16{1, 5, 6},
		SeriesIDSet: idset,
		TimeRange:   timeutil.TimeRange{Start: 100 * 1000, End: 500 * 1000}})
	// default version value
	assert.Zero(t, itr2.Version())
	assert.True(t, itr2.HasNext())
	tsItr2 := itr2.Next()
	// default tStore value
	assert.Zero(t, tsItr2.SeriesID())
	assert.Equal(t, uint32(2), itr2.Version())
	// next ts
	assert.True(t, tsItr2.HasNext())
	assert.Equal(t, uint32(5), tsItr2.SeriesID())
	assert.NotNil(t, tsItr2.Next())
	// test field-iterator, todo: fixme
	fItr5 := tsItr2.Next()
	assert.False(t, fItr5.HasNext())
	assert.NotNil(t, fItr5.Next())
	assert.Zero(t, fItr5.FieldID())
	assert.Zero(t, fItr5.FieldType())
	assert.Len(t, fItr5.FieldName(), 0)

	// seriesID: 7
	assert.True(t, tsItr2.HasNext())
	assert.Equal(t, uint32(7), tsItr2.SeriesID())
	assert.False(t, tsItr2.HasNext())

}

func Test_primitiveIterator(t *testing.T) {
	// todo: fixme
	pi := newPrimitiveIterator(series.ScanContext{})
	pi.reset(nil)
	assert.Zero(t, pi.FieldID())
	assert.False(t, pi.HasNext())
	timeSlot, value := pi.Next()
	assert.Zero(t, timeSlot)
	assert.Zero(t, value)
}

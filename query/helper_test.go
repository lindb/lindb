package query

import (
	"sort"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/memdb"
)

///////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

func MockTSDBEngine(ctrl *gomock.Controller) tsdb.Engine {
	shard := tsdb.NewMockShard(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	shard.EXPECT().GetMemoryDatabase().Return(memDB).AnyTimes()
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	metadataIndex := diskdb.NewMockIDGetter(ctrl)
	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), gomock.Any()).Return(uint16(10), field.SumField, nil).AnyTimes()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().GetShard(gomock.Any()).Return(shard).AnyTimes()
	engine.EXPECT().GetIDGetter().Return(metadataIndex).AnyTimes()
	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	return engine
}

// MockSumFieldSeries returns mock an iterator of sum field
func MockSumFieldSeries(ctrl *gomock.Controller, fieldID uint16, primitiveFieldID uint16,
	points map[int]interface{}) series.Iterator {

	it := series.NewMockFieldIterator(ctrl)
	//it.EXPECT().ID().Return(fieldID)
	it.EXPECT().HasNext().Return(true)

	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().Next().Return(primitiveIt)

	primitiveIt.EXPECT().FieldID().Return(primitiveFieldID)

	var keys []int
	for timeSlot := range points {
		keys = append(keys, timeSlot)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, timeSlot := range keys {
		primitiveIt.EXPECT().HasNext().Return(true)
		primitiveIt.EXPECT().Next().Return(timeSlot, points[timeSlot])
	}
	// mock nil primitive iterator
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(nil)
	it.EXPECT().FieldMeta().Return(field.Meta{ID: fieldID})
	// return hasNext=>false, finish primitive iterator
	primitiveIt.EXPECT().HasNext().Return(false).AnyTimes()

	// sum field only has one primitive field
	it.EXPECT().HasNext().Return(false).AnyTimes()

	timeSeries := series.NewMockIterator(ctrl)
	timeSeries.EXPECT().HasNext().Return(true)
	timeSeries.EXPECT().Next().Return(it)
	timeSeries.EXPECT().HasNext().Return(false)
	//return it
	return timeSeries
}

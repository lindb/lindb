package query

import (
	"sort"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/index"
)

///////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

func MockTSDBEngine(ctrl *gomock.Controller, families ...tsdb.DataFamily) tsdb.Engine {
	segment := tsdb.NewMockSegment(ctrl)
	if len(families) > 0 {
		for _, f := range families {
			segment.EXPECT().GetDataFamilies(gomock.Any()).Return([]tsdb.DataFamily{f})
		}
	} else {
		segment.EXPECT().GetDataFamilies(gomock.Any()).Return(nil).AnyTimes()
	}

	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().GetSegments(gomock.Any(), gomock.Any()).Return([]tsdb.Segment{segment}).AnyTimes()

	metadataIndex := index.NewMockMetadataGetter(ctrl)
	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), gomock.Any()).Return(uint16(10), field.SumField, nil).AnyTimes()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().GetShard(gomock.Any()).Return(shard).AnyTimes()
	engine.EXPECT().GetMetadataGetter().Return(metadataIndex).AnyTimes()
	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	return engine
}

// MockSumFieldSeries returns mock an iterator of sum field
func MockSumFieldSeries(ctrl *gomock.Controller, fieldID uint16, primitiveFieldID uint16, points map[int]interface{}) field.MultiTimeSeries {
	it := field.NewMockIterator(ctrl)
	//it.EXPECT().ID().Return(fieldID)
	it.EXPECT().HasNext().Return(true)

	primitiveIt := field.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().Next().Return(primitiveIt)

	primitiveIt.EXPECT().ID().Return(primitiveFieldID)

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
	it.EXPECT().ID().Return(fieldID)

	// return hasNext=>false, finish primitive iterator
	primitiveIt.EXPECT().HasNext().Return(false).AnyTimes()

	// sum field only has one primitive field
	it.EXPECT().HasNext().Return(false).AnyTimes()

	timeSeries := field.NewMockMultiTimeSeries(ctrl)
	timeSeries.EXPECT().Iterator().Return(it)
	//return it
	return timeSeries
}

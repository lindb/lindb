package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
)

func TestMetricStore_Filter(t *testing.T) {
	metricStore := mockMetricStore()

	// case 1: field not found
	rs, err := metricStore.Filter([]field.ID{1, 2}, nil, nil)
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 2: family not found
	rs, err = metricStore.Filter([]field.ID{1, 20}, nil, map[familyID]int64{
		familyID(10): 100,
	})
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 3: series ids not found
	rs, err = metricStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 2), map[familyID]int64{
		familyID(20): 100,
	})
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 3: found data
	rs, err = metricStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200), map[familyID]int64{
		familyID(20): 100,
	})
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	mrs := rs[0].(*memFilterResultSet)
	assert.Equal(t, []familyID{20}, mrs.familyIDs)
	assert.Equal(t,
		map[familyID]int64{
			familyID(20): 100,
		}, mrs.familyIDMap)
	assert.Equal(t,
		field.Metas{{
			ID:   20,
			Type: field.SumField,
		}}, mrs.fields)
}

func TestMemFilterResultSet_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	qFlow := flow.NewMockStorageQueryFlow(ctrl)
	mStore := mockMetricStore()

	rs, err := mStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200), map[familyID]int64{
		familyID(1):  100,
		familyID(20): 1000,
	})
	assert.NoError(t, err)
	sAgg := aggregation.NewMockSeriesAggregator(ctrl)
	fAgg := aggregation.NewMockFieldAggregator(ctrl)
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	// case 1: load data success
	gomock.InOrder(
		qFlow.EXPECT().GetAggregator().Return(aggregation.FieldAggregates{sAgg}),
		sAgg.EXPECT().GetAggregator(int64(100)).Return(fAgg, false),
		sAgg.EXPECT().GetAggregator(int64(1000)).Return(fAgg, true),
		fAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg}),
		pAgg.EXPECT().FieldID().Return(field.PrimitiveID(10)),
		qFlow.EXPECT().Reduce("host", gomock.Any()),
	)
	rs[0].Load(qFlow, []field.ID{20, 30}, 0, map[string][]uint16{
		"host": {1, 2},
	})
	// case 2: series ids not found
	gomock.InOrder(
		qFlow.EXPECT().GetAggregator().Return(aggregation.FieldAggregates{sAgg}),
		sAgg.EXPECT().GetAggregator(int64(100)).Return(fAgg, false),
		sAgg.EXPECT().GetAggregator(int64(1000)).Return(fAgg, true),
		fAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg}),
		pAgg.EXPECT().FieldID().Return(field.PrimitiveID(10)),
		qFlow.EXPECT().Reduce("host", gomock.Any()),
	)
	rs[0].Load(qFlow, []field.ID{20, 30}, 0, map[string][]uint16{
		"host": {100, 200},
	})
	// case 3: high key not exist
	rs[0].Load(qFlow, []field.ID{20, 30}, 10, map[string][]uint16{
		"host": {100, 200},
	})
	// case 4: field agg is empty
	gomock.InOrder(
		qFlow.EXPECT().GetAggregator().Return(aggregation.FieldAggregates{sAgg}),
		sAgg.EXPECT().GetAggregator(int64(100)).Return(nil, false),
		sAgg.EXPECT().GetAggregator(int64(1000)).Return(nil, false),
		qFlow.EXPECT().Reduce("host", gomock.Any()),
	)
	rs[0].Load(qFlow, []field.ID{20, 30}, 0, map[string][]uint16{
		"host": {100, 200},
	})
}

func mockMetricStore() *metricStore {
	mStore := newMetricStore()
	mStore.AddField(field.ID(10), field.SumField)
	mStore.AddField(field.ID(20), field.SumField)
	mStore.SetTimestamp(familyID(1), 10)
	mStore.SetTimestamp(familyID(20), 20)
	mStore.GetOrCreateTStore(100)
	mStore.GetOrCreateTStore(120)
	mStore.GetOrCreateTStore(200)
	return mStore.(*metricStore)
}

//
//import (
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"github.com/lindb/roaring"
//	"github.com/stretchr/testify/assert"
//
//	"github.com/lindb/lindb/pkg/timeutil"
//	pb "github.com/lindb/lindb/rpc/proto/field"
//	"github.com/lindb/lindb/series"
//	"github.com/lindb/lindb/series/field"
//	"github.com/lindb/lindb/tsdb/metadb"
//)
//
//func Test_MetricStore_scan(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
//
//	mStoreInterface := newMetricStore()
//	mStore := mStoreInterface.(*metricStore)
//	mStore.fieldsMetas.Store(field.Metas{
//		{ID: 3, Type: field.SumField, Name: "sum3"},
//		{ID: 4, Type: field.SumField, Name: "sum4"},
//		{ID: 5, Type: field.SumField, Name: "sum5"},
//		{ID: 6, Type: field.SumField, Name: "sum6"}})
//	// v1:
//	ti1 := newTagIndex().(*tagIndex)
//	ti1.version = 1
//	// v2
//	ti2 := newTagIndex().(*tagIndex)
//	ti2.version = 2
//	ts5 := newTimeSeriesStore()
//	ts6 := newTimeSeriesStore()
//	ts7 := newTimeSeriesStore()
//	ts8 := newTimeSeriesStore()
//	ts5.(*timeSeriesStore).insertFStore(newFieldStore(1))
//	ts5.(*timeSeriesStore).insertFStore(newFieldStore(2))
//	ts5.(*timeSeriesStore).insertFStore(newFieldStore(3))
//	ts5.(*timeSeriesStore).insertFStore(newFieldStore(4))
//
//	ti2.seriesID2TStore = newMetricMap()
//	ti2.seriesID2TStore.put(5, ts5)
//	ti2.seriesID2TStore.put(6, ts6)
//	ti2.seriesID2TStore.put(7, ts7)
//	ti2.seriesID2TStore.put(8, ts8)
//	// build id-set
//	idset := series.NewMultiVerSeriesIDSet()
//	idset.Add(0, roaring.New())
//	bitmap := roaring.New()
//	bitmap.AddMany([]uint32{1, 2, 3, 4, 5, 7})
//	idset.Add(2, bitmap)
//
//	// build mStore
//	mStore.immutable.Store(ti1)
//	mStore.mutable = ti2
//	fields := []*pb.Field{
//		{Name: "sum3", Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 1.0}}},
//		{Name: "sum4", Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 1.0}}},
//	}
//
//	generator := metadb.NewMockIDGenerator(ctrl)
//	generator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()
//
//	idGet := NewMockmStoreFieldIDGetter(ctrl)
//	idGet.EXPECT().GetFieldIDOrGenerate(gomock.Any(), "sum3", gomock.Any(), gomock.Any()).Return(uint16(3), nil)
//	idGet.EXPECT().GetFieldIDOrGenerate(gomock.Any(), "sum4", gomock.Any(), gomock.Any()).Return(uint16(4), nil)
//	bs := newBlockStore(10)
//	writtenSize, err := mStore.Write(uint32(10), fields,
//		writeContext{
//			generator:           generator,
//			blockStore:          bs,
//			familyTime:          familyTime,
//			slotIndex:           20,
//			metricID:            uint32(10),
//			mStoreFieldIDGetter: idGet,
//		})
//	assert.NotZero(t, writtenSize)
//	if err != nil {
//		t.Fatal(err)
//	}
//	//
//	//worker := &mockScanWorker{}
//	//mStore.Scan(&series.ScanContext{
//	//	SeriesIDSet:  idset,
//	//	FieldIDs:     []uint16{3, 4, 5},
//	//	IntervalCalc: calc,
//	//	Worker:       worker,
//	//})
//	//assert.Equal(t, 1, len(worker.events))
//	//// field not found
//	//mStore.Scan(&series.ScanContext{
//	//	SeriesIDSet: idset,
//	//	FieldIDs:    []uint16{1, 2},
//	//})
//	//// field not match
//	//mStore.Scan(&series.ScanContext{
//	//	SeriesIDSet: idset,
//	//	FieldIDs:    []uint16{1, 2, 3, 4},
//	//})
//}

package metricsdata

//
//import (
//	"testing"
//)
//
//func Test_NewMetricsDataScanner(t *testing.T) {
//	//assert.NotNil(t, NewScanner(nil))
//}
//
//func Test_newMDTVersionBlock(t *testing.T) {
//	// empty block
//	//vb, err := newMDTVersionBlock(series.Version(1), nil, &series.ScanContext{})
//	//assert.NotNil(t, err)
//	//assert.Nil(t, vb)
//	//
//	//vb, err = newMDTVersionBlock(series.Version(1), []byte{
//	//	1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5,
//	//}, &series.ScanContext{})
//	//assert.NotNil(t, err)
//	//assert.Nil(t, vb)
//}
////
////func buildGoodData() []byte {
////	nopKvFlusher := kv.NewNopFlusher()
////	flusherImpl := NewFlusher(nopKvFlusher)
////
////	flusherImpl.FlushFieldMetas([]field.Meta{
////		{ID: 1, Type: field.SumField, Name: "sum"},
////		{ID: 2, Type: field.MinField, Name: "min"},
////		{ID: 3, Type: field.MaxField, Name: "max"},
////	})
////	flusherImpl.FlushField(1, []byte{1, 1, 1, 1})
////	flusherImpl.FlushSeries(1)
////	flusherImpl.FlushField(2, []byte{2, 2, 2, 2})
////	flusherImpl.FlushSeries(2)
////	flusherImpl.FlushVersion(series.Version(100))
////	flusherImpl.FlushField(3, []byte{3, 3, 3, 3})
////	flusherImpl.FlushSeries(3)
////	flusherImpl.FlushVersion(series.Version(101))
////	flusherImpl.FlushField(3, []byte{3, 3, 3, 3})
////	flusherImpl.FlushSeries(4)
////	flusherImpl.FlushVersion(series.Version(102))
////	_ = flusherImpl.FlushMetric(1)
////	return nopKvFlusher.Bytes()
////}
//
//func Test_pickVersion2Blocks(t *testing.T) {
//	//ctrl := gomock.NewController(t)
//	//mockReader := table.NewMockReader(ctrl)
//	//mockReader.EXPECT().Get(uint32(1)).Return(buildGoodData()).AnyTimes()
//	//
//	//idSet := series.NewMultiVerSeriesIDSet()
//	//idSet.Add(series.Version(100), roaring.BitmapOf(1, 2))
//	//idSet.Add(series.Version(101), roaring.BitmapOf(3))
//	//
//	//scanner1 := NewScanner([]table.Reader{mockReader}).(*metricsDataScanner)
//	//m := scanner1.pickVersion2Blocks(&series.ScanContext{
//	//	MetricID:    1,
//	//	FieldIDs:    []uint16{1, 2, 3},
//	//	SeriesIDSet: idSet})
//	//assert.Len(t, m, 2)
//	//
//	//scanner2 := NewScanner([]table.Reader{mockReader, mockReader}).(*metricsDataScanner)
//	//m = scanner2.pickVersion2Blocks(&series.ScanContext{
//	//	MetricID:    1,
//	//	FieldIDs:    []uint16{1, 2, 3},
//	//	SeriesIDSet: idSet})
//	//assert.Len(t, m, 2)
//	//
//	//mdt := m[series.Version(100)][0]
//	//testMdtVersionBlock(t, mdt)
//}
//
////func testMdtVersionBlock(t *testing.T, mdt *mdtVersionBlock) {
//	//assert.NotNil(t, mdt.SeriesIDs())
//	//scanned := mdt.Scan()
//	//assert.True(t, scanned)
////}

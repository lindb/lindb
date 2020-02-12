package memdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
)

// Filter filters the data based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (ms *metricStore) Filter(metricID uint32, fieldIDs []uint16,
	version series.Version, seriesIDs *roaring.Bitmap,
) ([]flow.FilterResultSet, error) {
	//// first need check query's fields is match store's fields, if not return.
	//fmList := ms.fieldsMetas.Load().(field.Metas)
	//_, ok := fmList.Intersects(fieldIDs)
	//if !ok {
	//	return nil, nil
	//}
	//// scan tagIndex when version matches the idSet
	//var resultSet []flow.FilterResultSet
	//scanOnVersionMatch := func(idx tagIndexINTF) {
	//	if idx.Version() == version && idx.filter(seriesIDs) {
	//		resultSet = append(resultSet, &memFilterResultSet{tagIndex: idx})
	//	}
	//}
	//
	//ms.mux.RLock()
	//scanOnVersionMatch(ms.mutable)
	//immutable := ms.atomicGetImmutable()
	//ms.mux.RUnlock()
	//
	//if immutable != nil {
	//	scanOnVersionMatch(immutable)
	//}
	//return resultSet, nil
	return nil, nil
}

//type memFilterResultSet struct {
//	//tagIndex tagIndexINTF
//}
//
//func (rs *memFilterResultSet) Load(flow flow.StorageQueryFlow, fieldIDs []uint16,
//	highKey uint16, groupedSeries map[string][]uint16,
//) {
//	//rs.tagIndex.loadData(flow, fieldIDs, highKey, groupedSeries)
//}

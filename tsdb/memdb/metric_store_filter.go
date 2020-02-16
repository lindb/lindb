package memdb

import (
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

// Filter filters the data based on fieldIDs/seriesIDs/familyIDs,
// if finds data then returns the FilterResultSet, else returns constants.ErrNotFound
func (ms *metricStore) Filter(fieldIDs []field.ID,
	seriesIDs *roaring.Bitmap, familyIDs map[familyID]int64,
) ([]flow.FilterResultSet, error) {
	// first need check query's fields is match store's fields, if not return.
	fields, _ := ms.fields.Intersects(fieldIDs)
	if len(fields) == 0 {
		// field not found
		return nil, constants.ErrNotFound
	}
	resultFamilyIDMap := make(map[familyID]int64)
	var resultFamilyIDs []familyID // sort by family id

	for _, entry := range ms.families {
		fTime, ok := familyIDs[entry.id]
		if ok {
			resultFamilyIDMap[entry.id] = fTime
			resultFamilyIDs = append(resultFamilyIDs, entry.id)
		}
	}
	if len(resultFamilyIDMap) == 0 {
		// family time not found
		return nil, constants.ErrNotFound
	}

	// after and operator, query bitmap is sub of store bitmap
	matchSeriesIDs := roaring.FastAnd(seriesIDs, ms.keys)
	if matchSeriesIDs.IsEmpty() {
		// series id not found
		return nil, constants.ErrNotFound
	}

	// sort by family ids
	sort.Slice(resultFamilyIDs, func(i, j int) bool { return resultFamilyIDs[i] < resultFamilyIDs[j] })

	// returns the filter result set
	return []flow.FilterResultSet{
		&memFilterResultSet{
			store:       ms,
			fields:      fields,
			familyIDs:   resultFamilyIDs,
			familyIDMap: resultFamilyIDMap,
		},
	}, nil
}

// fieldAggregator represents the field aggregator that does memory data scan and aggregates
type fieldAggregator struct {
	familyID   familyID
	fieldMeta  field.Meta
	aggregator aggregation.PrimitiveAggregator

	fieldKey uint32
}

// newFieldAggregator creates a field aggregator
func newFieldAggregator(familyID familyID, fieldMeta field.Meta, aggregator aggregation.PrimitiveAggregator) *fieldAggregator {
	fieldKey := buildFieldKey(familyID, fieldMeta.ID, aggregator.FieldID())
	return &fieldAggregator{
		familyID:   familyID,
		fieldMeta:  fieldMeta,
		aggregator: aggregator,
		fieldKey:   fieldKey,
	}
}

// memFilterResultSet represents memory filter result set for loading data in query flow
type memFilterResultSet struct {
	store       *metricStore
	fields      field.Metas // sort by field id
	familyIDs   []familyID  // sort by family id
	familyIDMap map[familyID]int64
}

// prepare prepares the field aggregator based on query condition
func (rs *memFilterResultSet) prepare(fieldIDs []field.ID, aggregator aggregation.FieldAggregates) (aggs []*fieldAggregator) {
	for _, fID := range rs.familyIDs { // sort by family ids
		familyTime := rs.familyIDMap[fID]
		for idx, fieldID := range fieldIDs { // sort by field ids
			fMeta, ok := rs.fields.GetFromID(fieldID)
			if !ok {
				continue
			}
			fieldAggregator, ok := aggregator[idx].GetAggregator(familyTime)
			if !ok {
				continue
			}
			pAggregators := fieldAggregator.GetAllAggregators() // sort by primitive field ids
			for _, agg := range pAggregators {
				aggs = append(aggs, newFieldAggregator(fID, fMeta, agg))
			}
		}
	}
	return
}

// Load loads the data from storage, then does down sampling, finally reduces the down sampling results.
func (rs *memFilterResultSet) Load(flow flow.StorageQueryFlow, fieldIDs []field.ID,
	highKey uint16, groupedSeries map[string][]uint16,
) {
	//FIXME need add lock?????

	// 1. get high container index by the high key of series ID
	highContainerIdx := rs.store.keys.GetContainerIndex(highKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series ID not exist) return it
		return
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := rs.store.keys.GetContainerAtIndex(highContainerIdx)

	memScanCtx := &memScanContext{
		tsd: encoding.GetTSDDecoder(),
	}
	for groupByTags, lowSeriesIDs := range groupedSeries {
		aggregator := flow.GetAggregator()
		memScanCtx.fieldAggs = rs.prepare(fieldIDs, aggregator)
		if len(memScanCtx.fieldAggs) == 0 {
			// reduce empty aggregator for re-use
			flow.Reduce(groupByTags, aggregator)
			continue
		}
		for _, lowSeriesID := range lowSeriesIDs {
			// check low series id if exist
			if !lowContainer.Contains(lowSeriesID) {
				continue
			}
			// get the index of low series id in container
			idx := lowContainer.Rank(lowSeriesID)
			// scan the data and aggregate the values
			store := rs.store.values[highContainerIdx][idx-1]
			store.scan(memScanCtx)
		}
		flow.Reduce(groupByTags, aggregator)
	}
	encoding.ReleaseTSDDecoder(memScanCtx.tsd)
}

package memdb

import (
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
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
			seriesIDs:   matchSeriesIDs,
		},
	}, nil
}

// fieldAggregator represents the field aggregator that does memory data scan and aggregates
type fieldAggregator struct {
	familyID  familyID
	fieldMeta field.Meta
	block     series.Block

	fieldKey uint32
}

// newFieldAggregator creates a field aggregator
func newFieldAggregator(familyID familyID, fieldMeta field.Meta, block series.Block) *fieldAggregator {
	fieldKey := buildFieldKey(familyID, fieldMeta.ID)
	return &fieldAggregator{
		familyID:  familyID,
		fieldMeta: fieldMeta,
		block:     block,
		fieldKey:  fieldKey,
	}
}

// memFilterResultSet represents memory filter result set for loading data in query flow
type memFilterResultSet struct {
	store       *metricStore
	fields      field.Metas // sort by field id
	familyIDs   []familyID  // sort by family id
	familyIDMap map[familyID]int64

	seriesIDs *roaring.Bitmap
}

// prepare prepares the field aggregator based on query condition
func (rs *memFilterResultSet) prepare(fieldIDs []field.ID, aggregator aggregation.ContainerAggregator) (aggs []*fieldAggregator) {
	for _, fID := range rs.familyIDs { // sort by family ids
		familyTime := rs.familyIDMap[fID]
		for idx, fieldID := range fieldIDs { // sort by field ids
			fMeta, ok := rs.fields.GetFromID(fieldID)
			if !ok {
				continue
			}
			block, ok := aggregator.GetFieldAggregates()[idx].GetAggregateBlock(familyTime)
			if !ok {
				continue
			}
			aggs = append(aggs, newFieldAggregator(fID, fMeta, block))
		}
	}
	return
}

// Identifier identifies the source of result set from memory storage
func (rs *memFilterResultSet) Identifier() string {
	return "memory"
}

// SeriesIDs returns the series ids which matches with query series ids
func (rs *memFilterResultSet) SeriesIDs() *roaring.Bitmap {
	return rs.seriesIDs
}

// Load loads the data from storage, then returns the memory storage metric scanner.
func (rs *memFilterResultSet) Load(flow flow.StorageQueryFlow,
	fieldIDs []field.ID,
	highKey uint16, seriesID roaring.Container,
) flow.Scanner {
	//FIXME need add lock?????

	// 1. get high container index by the high key of series ID
	highContainerIdx := rs.store.keys.GetContainerIndex(highKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series ID not exist) return it
		return nil
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := rs.store.keys.GetContainerAtIndex(highContainerIdx)
	foundSeriesIDs := lowContainer.And(seriesID)
	if foundSeriesIDs.GetCardinality() == 0 {
		return nil
	}

	aggregator := flow.GetAggregator(highKey)
	fieldAggs := rs.prepare(fieldIDs, aggregator)
	if len(fieldAggs) == 0 {
		return nil
	}

	// must use lowContainer from store, because get series index based on container
	return newMetricStoreScanner(lowContainer, rs.store.values[highContainerIdx], fieldAggs)
}

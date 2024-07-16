// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package memdb

import (
	"fmt"
	"sort"

	commontimeutil "github.com/lindb/common/pkg/timeutil"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

// fieldEntry represents field context for searching.
type fieldEntry struct {
	pageBuf     DataPointBuffer
	compressBuf CompressStore // time series compress buffer
	field       field.Meta

	buf []byte // time series current write buffer
}

// Reset resets time series current write buffer.
func (fe *fieldEntry) Reset(buf []byte) {
	fe.buf = buf
}

// GetValue returns value by time slot, if it hasn't, return false.
func (fe *fieldEntry) GetValue(slot uint16) (float64, bool) {
	if len(fe.buf) == 0 {
		return 0, false
	}
	startTime := getStart(fe.buf)
	return getCurrentValue(fe.buf, startTime, slot)
}

// getCompressBuf returns time series compress buffer by memory time series id.
func (fe *fieldEntry) getCompressBuf(memSeriesID uint32) []byte {
	if fe.compressBuf == nil {
		return nil
	}
	return fe.compressBuf.GetCompressBuffer(memSeriesID)
}

// getPage returns current time series write buffer by memory time series id.
func (fe *fieldEntry) getPage(memTimeSeriesID uint32) ([]byte, bool) {
	if fe.pageBuf == nil {
		return nil, false
	}
	return fe.pageBuf.GetPage(memTimeSeriesID)
}

// Filter filters the data based on fields/seriesIDs/family time,
// if it finds data then returns the FilterResultSet, else returns constants.ErrFieldNotFound
func (md *memoryDatabase) filter(shardExecuteContext *flow.ShardExecuteContext,
	memMetricID uint64, slotRange *timeutil.SlotRange,
	timeSeriesIndex TimeSeriesIndex,
) ([]flow.FilterResultSet, error) {
	mStore, ok := md.indexDB.GetMetadataDatabase().GetMetricMeta(memMetricID)
	if !ok {
		// metric meta not found
		return nil, nil
	}
	fields := shardExecuteContext.StorageExecuteCtx.Fields.Clone()
	// NOTE: must re-stort by field name, if not cannot find field from query fields
	sort.Sort(fields)
	// first need check query's fields is match store's fields, if not return.
	foundFields := mStore.FindFields(fields)
	if len(foundFields) == 0 {
		// field not found
		return nil, fmt.Errorf("%w, fields: %s", constants.ErrFieldNotFound, fields.String())
	}

	var fieldEntries []*fieldEntry
	for _, fm := range foundFields {
		fStore, ok := md.fieldWriteStores.Load(fm.Index)
		fcStore, fcOK := md.fieldCompressStore.Load(fm.Index)
		if ok || fcOK {
			queryField, _ := fields.GetFromName(fm.Name)
			fieldEntry := &fieldEntry{
				pageBuf: fStore.(DataPointBuffer), // TEST: add test case
				field:   queryField,
			}
			fieldEntries = append(fieldEntries, fieldEntry)
			if ok {
				fieldEntry.pageBuf = fStore.(DataPointBuffer)
			}
			if fcOK {
				fieldEntry.compressBuf = fcStore.(CompressStore)
			}
		}
	}

	if len(fieldEntries) == 0 {
		// field temp store buffer not found
		return nil, fmt.Errorf("%w, fields: %s", constants.ErrFieldNotFound, fields.String())
	}

	seriesIDs := shardExecuteContext.SeriesIDsAfterFiltering
	familyTime := md.FamilyTime()
	// after and operator, query bitmap is sub of store bitmap
	matchSeriesIDs := roaring.FastAnd(seriesIDs, timeSeriesIndex.TimeSeriesIDs())
	if matchSeriesIDs.IsEmpty() {
		// series id not found
		return nil, fmt.Errorf("%w when Filter, familyTime: %d, fields: %s",
			constants.ErrSeriesIDNotFound, familyTime, fields.String())
	}
	// returns the filter result set
	return []flow.FilterResultSet{
		&memFilterResultSet{
			db:              md,
			timeSeriesIndex: timeSeriesIndex,
			familyTime:      familyTime,
			fields:          fieldEntries,
			slotRange:       slotRange,
			storeSeriesIDs:  seriesIDs,
			seriesIDs:       matchSeriesIDs,
		},
	}, nil
}

// memFilterResultSet represents memory filter result set for loading data in query flow
type memFilterResultSet struct {
	timeSeriesIndex TimeSeriesIndex
	db              *memoryDatabase
	slotRange       *timeutil.SlotRange
	storeSeriesIDs  *roaring.Bitmap
	seriesIDs       *roaring.Bitmap
	fields          []*fieldEntry
	familyTime      int64
}

// Identifier identifies the source of result set from memory storage
func (rs *memFilterResultSet) Identifier() string {
	dbStatus := "readwrite"
	if rs.db.IsReadOnly() {
		dbStatus = "readonly"
	}
	return fmt.Sprintf("%s/memory/%s",
		commontimeutil.FormatTimestamp(rs.familyTime, commontimeutil.DataTimeFormat2), dbStatus)
}

// FamilyTime returns the family time of storage.
func (rs *memFilterResultSet) FamilyTime() int64 {
	return rs.familyTime
}

// SlotRange returns the slot range of storage.
func (rs *memFilterResultSet) SlotRange() timeutil.SlotRange {
	return *rs.slotRange
}

// SeriesIDs returns the series ids which matches with query series ids
func (rs *memFilterResultSet) SeriesIDs() *roaring.Bitmap {
	return rs.seriesIDs
}

// Load loads the data from storage, then returns the memory storage metric scanner.
func (rs *memFilterResultSet) Load(ctx *flow.DataLoadContext) flow.DataLoader {
	// 1. get high container index by the high key of series ID
	highContainerIdx := rs.storeSeriesIDs.GetContainerIndex(ctx.SeriesIDHighKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series ID not exist) return it
		return nil
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := rs.storeSeriesIDs.GetContainerAtIndex(highContainerIdx)
	foundSeriesIDs := lowContainer.And(ctx.LowSeriesIDsContainer)
	if foundSeriesIDs.GetCardinality() == 0 {
		return nil
	}
	// must use lowContainer from store, because get series index based on container
	return NewTimeSeriesLoader(rs.db, rs.timeSeriesIndex, ctx.SeriesIDHighKey, *rs.slotRange, rs.fields)
}

// Close release the resource during doing query operation.
func (rs *memFilterResultSet) Close() {
	// do nothing
}

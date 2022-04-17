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

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

// Filter filters the data based on fields/seriesIDs/family time,
// if it finds data then returns the FilterResultSet, else returns constants.ErrFieldNotFound
func (ms *metricStore) Filter(shardExecuteContext *flow.ShardExecuteContext, db MemoryDatabase) ([]flow.FilterResultSet, error) {
	fields := shardExecuteContext.StorageExecuteCtx.Fields
	// first need check query's fields is match store's fields, if not return.
	foundFields, _ := ms.fields.Intersects(fields)
	if len(foundFields) == 0 {
		// field not found
		return nil, fmt.Errorf("%w, fields: %s", constants.ErrFieldNotFound, fields.String())
	}

	seriesIDs := shardExecuteContext.SeriesIDsAfterFiltering
	familyTime := db.FamilyTime()
	// after and operator, query bitmap is sub of store bitmap
	matchSeriesIDs := roaring.FastAnd(seriesIDs, ms.keys)
	if matchSeriesIDs.IsEmpty() {
		// series id not found
		return nil, fmt.Errorf("%w when Filter, familyTime: %d, fields: %s",
			constants.ErrSeriesIDNotFound, familyTime, fields.String())
	}

	// returns the filter result set
	return []flow.FilterResultSet{
		&memFilterResultSet{
			db:         db,
			familyTime: familyTime,
			store:      ms,
			fields:     fields,
			seriesIDs:  matchSeriesIDs,
		},
	}, nil
}

// memFilterResultSet represents memory filter result set for loading data in query flow
type memFilterResultSet struct {
	db         MemoryDatabase
	familyTime int64
	store      *metricStore
	fields     field.Metas // sort by field id

	seriesIDs *roaring.Bitmap
}

// Identifier identifies the source of result set from memory storage
func (rs *memFilterResultSet) Identifier() string {
	dbStatus := "readwrite"
	if rs.db.IsReadOnly() {
		dbStatus = "readonly"
	}
	return fmt.Sprintf("%s/memory/%s",
		timeutil.FormatTimestamp(rs.familyTime, timeutil.DataTimeFormat2), dbStatus)
}

// FamilyTime returns the family time of storage.
func (rs *memFilterResultSet) FamilyTime() int64 {
	return rs.familyTime
}

// SlotRange returns the slot range of storage.
func (rs *memFilterResultSet) SlotRange() timeutil.SlotRange {
	return *rs.store.slotRange
}

// SeriesIDs returns the series ids which matches with query series ids
func (rs *memFilterResultSet) SeriesIDs() *roaring.Bitmap {
	return rs.seriesIDs
}

// Load loads the data from storage, then returns the memory storage metric scanner.
func (rs *memFilterResultSet) Load(ctx *flow.DataLoadContext) flow.DataLoader {
	// 1. get high container index by the high key of series ID
	highContainerIdx := rs.store.keys.GetContainerIndex(ctx.SeriesIDHighKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series ID not exist) return it
		return nil
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := rs.store.keys.GetContainerAtIndex(highContainerIdx)
	foundSeriesIDs := lowContainer.And(ctx.LowSeriesIDsContainer)
	if foundSeriesIDs.GetCardinality() == 0 {
		return nil
	}

	// must use lowContainer from store, because get series index based on container
	return newMetricStoreLoader(rs.db, lowContainer, rs.store.values[highContainerIdx], *rs.store.slotRange, rs.fields)
}

// Close release the resource during doing query operation.
func (rs *memFilterResultSet) Close() {
	// do nothing
}

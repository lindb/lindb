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
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
)

// Filter filters the data based on fieldIDs/seriesIDs/familyIDs,
// if finds data then returns the FilterResultSet, else returns constants.ErrNotFound
func (ms *metricStore) Filter(fieldIDs []field.ID, seriesIDs *roaring.Bitmap) ([]flow.FilterResultSet, error) {
	// first need check query's fields is match store's fields, if not return.
	fields, _ := ms.fields.Intersects(fieldIDs)
	if len(fields) == 0 {
		// field not found
		return nil, constants.ErrNotFound
	}

	// after and operator, query bitmap is sub of store bitmap
	matchSeriesIDs := roaring.FastAnd(seriesIDs, ms.keys)
	if matchSeriesIDs.IsEmpty() {
		// series id not found
		return nil, constants.ErrNotFound
	}

	// returns the filter result set
	return []flow.FilterResultSet{
		&memFilterResultSet{
			store:     ms,
			fields:    fields,
			seriesIDs: matchSeriesIDs,
		},
	}, nil
}

// memFilterResultSet represents memory filter result set for loading data in query flow
type memFilterResultSet struct {
	store       *metricStore
	fields      field.Metas // sort by field id
	queryFields field.Metas // query fields sort by field id

	seriesIDs *roaring.Bitmap
}

// prepare prepares the field aggregator based on query condition
func (rs *memFilterResultSet) prepare(fieldIDs []field.ID) {
	for _, fieldID := range fieldIDs { // sort by field ids
		fMeta, ok := rs.fields.GetFromID(fieldID)
		if !ok {
			continue
		}
		rs.queryFields = append(rs.queryFields, fMeta)
	}
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
func (rs *memFilterResultSet) Load(highKey uint16, seriesIDs roaring.Container, fieldIDs []field.ID) flow.Scanner {
	//FIXME need add lock?????

	// 1. get high container index by the high key of series ID
	highContainerIdx := rs.store.keys.GetContainerIndex(highKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series ID not exist) return it
		return nil
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := rs.store.keys.GetContainerAtIndex(highContainerIdx)
	foundSeriesIDs := lowContainer.And(seriesIDs)
	if foundSeriesIDs.GetCardinality() == 0 {
		return nil
	}

	rs.prepare(fieldIDs)
	if len(rs.queryFields) == 0 {
		return nil
	}

	// must use lowContainer from store, because get series index based on container
	return newMetricStoreScanner(lowContainer, rs.store.values[highContainerIdx], rs.fields)
}

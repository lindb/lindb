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
	"sort"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./timeseries_store.go -destination=./timeseries_store_mock.go -package memdb

const emptyTimeSeriesStoreSize = 24 // fStores slice

// tStoreINTF abstracts a time-series store
type tStoreINTF interface {
	// Capacity returns the size of tStoreINTF without fields
	Capacity() int
	// GetFStore returns the fStore in field list by field id.
	GetFStore(fieldID field.ID) (fStoreINTF, bool)
	// InsertFStore inserts a new fStore to field list.
	InsertFStore(fStore fStoreINTF)
	// FlushSeriesTo flushes the series data segment.
	FlushSeriesTo(flusher metricsdata.Flusher, flushCtx flushContext)
	// load loads the time series data based on field ids
	load(fields field.Metas, slotRange timeutil.SlotRange) [][]byte
}

// fStoreNodes implements sort.Interface
type fStoreNodes []fStoreINTF

func (f fStoreNodes) Len() int           { return len(f) }
func (f fStoreNodes) Less(i, j int) bool { return f[i].GetFieldID() < f[j].GetFieldID() }
func (f fStoreNodes) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// timeSeriesStore holds a mapping relation of field and fieldStore.
type timeSeriesStore struct {
	fStoreNodes fStoreNodes // key: sorted fStore list by field-name, insert-only
}

// newTimeSeriesStore returns a new tStoreINTF.
func newTimeSeriesStore() tStoreINTF {
	return &timeSeriesStore{}
}

// GetFStore returns the fStore in this list from field-id.
func (ts *timeSeriesStore) GetFStore(fieldID field.ID) (fStoreINTF, bool) {
	fieldLength := len(ts.fStoreNodes)
	if fieldLength == 1 {
		if ts.fStoreNodes[0].GetFieldID() != fieldID {
			return nil, false
		}
		return ts.fStoreNodes[0], true
	}
	// fast path
	if fieldLength < 20 {
		for idx := range ts.fStoreNodes {
			if ts.fStoreNodes[idx].GetFieldID() == fieldID {
				return ts.fStoreNodes[idx], true
			}
		}
		return nil, false
	}
	idx := sort.Search(fieldLength, func(i int) bool {
		return ts.fStoreNodes[i].GetFieldID() >= fieldID
	})
	if idx >= fieldLength || ts.fStoreNodes[idx].GetFieldID() != fieldID {
		return nil, false
	}
	return ts.fStoreNodes[idx], true
}

func (ts *timeSeriesStore) Capacity() int {
	return emptyTimeSeriesStoreSize + 8*cap(ts.fStoreNodes)
}

// InsertFStore inserts a new fStore to field list.
func (ts *timeSeriesStore) InsertFStore(fStore fStoreINTF) {
	ts.fStoreNodes = append(ts.fStoreNodes, fStore)
	if len(ts.fStoreNodes) > 1 {
		sort.Sort(ts.fStoreNodes)
	}
}

// FlushSeriesTo flushes the series data segment.
func (ts *timeSeriesStore) FlushSeriesTo(flusher metricsdata.Flusher, flushCtx flushContext) {
	stores := ts.fStoreNodes
	fStoreLen := len(stores)
	// if no field store under current data family
	if stores == nil || fStoreLen == 0 {
		return
	}
	fieldMetas := flusher.GetFieldMetas()
	idx := 0
	for _, fieldMeta := range fieldMetas {
		if idx < fStoreLen && fieldMeta.ID == stores[idx].GetFieldID() {
			// flush field data
			stores[idx].FlushFieldTo(flusher, fieldMeta, flushCtx)
			idx++
		} else {
			// must flush nil data for metric has multi-field
			flusher.FlushField(nil)
		}
	}
}

// load loads the time series data based on key(family+field).
// NOTICE: field ids and fields aggregator must be in order.
func (ts *timeSeriesStore) load(fields field.Metas, slotRange timeutil.SlotRange) [][]byte {
	fieldLength := len(ts.fStoreNodes)
	fieldCount := len(fields)
	rs := make([][]byte, fieldCount)
	// find equals field id index
	idx := sort.Search(fieldLength, func(i int) bool {
		return ts.fStoreNodes[i].GetFieldID() >= fields[0].ID
	})
	j := 0
	for i := idx; i < fieldLength; i++ {
		fieldStore := ts.fStoreNodes[i]
		queryFieldID := fields[j].ID
		storeFieldID := fieldStore.GetFieldID()
		switch {
		case storeFieldID == queryFieldID:
			// load field data
			rs[j] = fieldStore.Load(fields[j].Type, slotRange)
			j++ // goto next query field id
			// found all query fields return it
			if fieldCount == j {
				return rs
			}
		case storeFieldID > queryFieldID:
			// store field id > query field id, return it
			return rs
		}
	}
	return rs
}

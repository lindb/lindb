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

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock.go -package memdb

// for testing
var (
	flushFunc = flush
)

const (
	emptyMStoreSize = 8 +
		12 + // slot and pointer
		8 + // pointer in metricStore
		24 // slice in metricStore
	// keys in metricStore are ignored
)

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// Filter filters the data based on fields/seriesIDs/family time,
	// if finds data then returns the flow.FilterResultSet, else returns constants.ErrNotFound
	Filter(familyTime int64, seriesIDs *roaring.Bitmap, fields field.Metas) ([]flow.FilterResultSet, error)
	// SetSlot sets the current write slot
	SetSlot(slot uint16)
	// GetSlotRange returns slot range.
	GetSlotRange() *timeutil.SlotRange
	// AddField adds field meta into metric level
	AddField(fieldID field.ID, fieldType field.Type)
	// GetOrCreateTStore constructs the index and return a tStore
	GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, createdSize int)
	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(tableFlusher metricsdata.Flusher, flushCtx flushContext) (err error)
}

// metricStore represents metric level storage, stores all series data, and fields/family times metadata
type metricStore struct {
	MetricStore

	slotRange *timeutil.SlotRange
	fields    field.Metas // field metadata
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	var ms metricStore
	ms.keys = roaring.New() // init keys
	return &ms
}

// SetSlot sets the current write timestamp
func (ms *metricStore) SetSlot(slot uint16) {
	if ms.slotRange == nil {
		slotRange := timeutil.NewSlotRange(slot, slot)
		ms.slotRange = &slotRange
	} else {
		ms.slotRange.SetSlot(slot)
	}
}

// GetSlotRange returns slot range.
func (ms *metricStore) GetSlotRange() *timeutil.SlotRange {
	return ms.slotRange
}

// AddField adds field meta into metric level
func (ms *metricStore) AddField(fieldID field.ID, fieldType field.Type) {
	_, ok := ms.fields.GetFromID(fieldID)
	if !ok {
		ms.fields = ms.fields.Insert(field.Meta{
			ID:   fieldID,
			Type: fieldType,
		})
		// sort by field id
		sort.Slice(ms.fields, func(i, j int) bool { return ms.fields[i].ID < ms.fields[j].ID })
	}
}

// GetOrCreateTStore constructs the index and return a tStore
func (ms *metricStore) GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, createdSize int) {
	tStore, ok := ms.Get(seriesID)
	if !ok {
		tStore = newTimeSeriesStore()
		ms.Put(seriesID, tStore)
		createdSize += emptyTimeSeriesStoreSize + 8 // pointer size
	}
	return tStore, createdSize
}

// FlushMetricsTo Writes metric-data to the table.
func (ms *metricStore) FlushMetricsDataTo(flusher metricsdata.Flusher, flushCtx flushContext) (err error) {
	slotRange := ms.slotRange
	// field not exist, return
	fieldLen := len(ms.fields)
	if fieldLen == 0 {
		return
	}
	// flush field meta info
	flusher.FlushFieldMetas(ms.fields)
	// set current family's slot range
	flushCtx.Start, flushCtx.End = slotRange.GetRange()
	if err := ms.WalkEntry(func(key uint32, value tStoreINTF) error {
		return flushFunc(flusher, flushCtx, key, value)
	}); err != nil {
		return err
	}
	return flusher.FlushMetric(flushCtx.metricID, slotRange.Start, slotRange.End)
}

// flush flushes series data
func flush(flusher metricsdata.Flusher, flushCtx flushContext, key uint32, value tStoreINTF) error {
	value.FlushSeriesTo(flusher, flushCtx)
	flusher.FlushSeries(key)
	return nil
}

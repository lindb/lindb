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

	"go.uber.org/atomic"

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
	emptyMStoreSize = 4 + // metricStore put-count
		8 + // metric store roaring pointer
		4 + // capacity size
		12 // slot range pointer and struct
	fieldMetaSize = 2 + // field id
		1 + // field Type
		8 + // field name, string internal pointer
		4 // field name, string internal len size
)

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// Capacity returns the memory usage of metric-store,
	// without tStores and FieldStores
	Capacity() int
	// Filter filters the data based on fields/seriesIDs/family time,
	// if data founded then returns the flow.FilterResultSet, else returns constants.ErrNotFound
	Filter(shardExecuteContext *flow.ShardExecuteContext, db MemoryDatabase) ([]flow.FilterResultSet, error)
	// SetSlot sets the current write slot
	SetSlot(slot uint16)
	// GetSlotRange returns slot range.
	GetSlotRange() *timeutil.SlotRange
	// AddField adds field meta into metric level
	AddField(fieldID field.ID, fieldType field.Type)
	// GetOrCreateTStore constructs the index and return a tStore
	GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, created bool)
	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(tableFlusher metricsdata.Flusher, flushCtx *flushContext) (err error)
}

// metricStore represents metric level storage, stores all series data, and fields/family times metadata
type metricStore struct {
	capacity atomic.Int32 // memory usage

	MetricStore

	slotRange *timeutil.SlotRange
	fields    field.Metas // field metadata
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	var ms metricStore
	ms.keys = roaring.New() // init keys
	ms.capacity.Store(int32(emptyMStoreSize + cap(ms.fields)*fieldMetaSize))
	return &ms
}

func (ms *metricStore) Capacity() int {
	return int(ms.capacity.Load())
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
	if _, ok := ms.fields.GetFromID(fieldID); !ok {
		fieldsCap := cap(ms.fields)
		ms.fields = ms.fields.Insert(field.Meta{
			ID:   fieldID,
			Type: fieldType,
		})
		ms.capacity.Add(int32((cap(ms.fields) - fieldsCap) * fieldMetaSize))
		if len(ms.fields) <= 1 {
			return
		}
		// sort by field id
		sort.Slice(ms.fields, func(i, j int) bool { return ms.fields[i].ID < ms.fields[j].ID })
	}
}

func (ms *metricStore) mStoreSize() int {
	var size int
	size += cap(ms.MetricStore.values)*24 + 24
	for idx := range ms.MetricStore.values {
		size += cap(ms.MetricStore.values[idx])*8 + 24
	}
	return size
}

// GetOrCreateTStore constructs the index and return a tStore
func (ms *metricStore) GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, created bool) {
	tStore, ok := ms.Get(seriesID)
	if !ok {
		tStore = newTimeSeriesStore()
		beforeMStoreSize := ms.mStoreSize()
		ms.Put(seriesID, tStore)
		ms.capacity.Add(int32(ms.mStoreSize() - beforeMStoreSize))
		created = true
	}
	return tStore, created
}

// FlushMetricsDataTo Writes metric-data to the table.
func (ms *metricStore) FlushMetricsDataTo(flusher metricsdata.Flusher, flushCtx *flushContext) (err error) {
	slotRange := ms.slotRange
	// field not exist, return
	fieldLen := len(ms.fields)
	if fieldLen == 0 {
		return
	}
	// prepare for flushing metric
	flusher.PrepareMetric(flushCtx.metricID, ms.fields)
	// set current family's slot range
	flushCtx.Start, flushCtx.End = slotRange.GetRange()
	if err := ms.WalkEntry(func(key uint32, value tStoreINTF) error {
		return flushFunc(flusher, flushCtx, key, value)
	}); err != nil {
		return err
	}
	return flusher.CommitMetric(flushCtx.SlotRange)
}

// flush series data
func flush(flusher metricsdata.Flusher, flushCtx *flushContext, key uint32, tStore tStoreINTF) error {
	if err := tStore.FlushFieldsTo(flusher, flushCtx); err != nil {
		return err
	}
	return flusher.FlushSeries(key)
}

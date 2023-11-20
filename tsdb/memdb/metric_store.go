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

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock.go -package memdb

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// Capacity returns the memory usage of metric-store,
	// without tStores and FieldStores
	Capacity() int
	// Filter filters the data based on fields/seriesIDs/family time,
	// if data founded then returns the flow.FilterResultSet, else returns constants.ErrNotFound
	Filter(shardExecuteContext *flow.ShardExecuteContext, db *memoryDatabase) ([]flow.FilterResultSet, error)
	// SetSlot sets the current write slot
	SetSlot(slot uint16)
	// GetSlotRange returns slot range.
	GetSlotRange() *timeutil.SlotRange
	// GenField generates field meta under memory database.
	GenField(fieldName field.Name, fieldType field.Type) (f field.Meta, created bool)
	// UpdateFieldMeta updates field meta after metric meta updated.
	UpdateFieldMeta(fieldID field.ID, fm field.Meta)
	// FindFields returns fields from store based on current written fields.
	FindFields(fields field.Metas) (found field.Metas)
	// GenTStore generates memory level time series id under memory database.
	GenTStore(tagHash uint64, create func() uint32) uint32
	// IndexTStore adds memorty store series id -> global series id mapping.
	IndexTStore(seriesID, seriesIdx uint32)
	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(
		tableFlusher metricsdata.Flusher,
		flushCtx *flushContext,
		flushFields func(memSeriesID uint32, fields field.Metas) error,
	) (err error)
}

// metricStore represents metric level storage, stores all series data, and fields/family times metadata
type metricStore struct {
	hashes map[uint64]uint32 // tag hash -> memory time series id
	ids    *imap.IntMap[uint32]

	slotRange *timeutil.SlotRange
	fields    field.Metas // field metadata
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	var ms metricStore
	ms.hashes = make(map[uint64]uint32)
	ms.ids = imap.NewIntMap[uint32]()
	return &ms
}

func (ms *metricStore) Capacity() int {
	return 0
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

// GenField generates field meta under memory database.
func (ms *metricStore) GenField(name field.Name, fType field.Type) (f field.Meta, created bool) {
	fm, ok := ms.fields.GetFromName(name)
	if !ok {
		index := uint8(len(ms.fields))
		fm = field.Meta{
			Type:  fType,
			Name:  name,
			Index: index,
		}
		ms.fields = append(ms.fields, fm)
		// sort by field name
		sort.Slice(ms.fields, func(i, j int) bool { return ms.fields[i].Name < ms.fields[j].Name })
		return fm, true
	}
	return fm, false
}

// UpdateFieldMeta updates field meta after metric meta updated.
func (ms *metricStore) UpdateFieldMeta(fieldID field.ID, fm field.Meta) {
	idx, ok := ms.fields.FindIndexByName(fm.Name)
	if ok {
		ms.fields[idx].ID = fieldID
		ms.fields[idx].Persisted = true
	}
}

// GenTStore generates memory level time series id under memory database.
func (ms *metricStore) GenTStore(tagHash uint64, create func() uint32) uint32 {
	ts, ok := ms.hashes[tagHash]
	if !ok {
		ts = create()
		ms.hashes[tagHash] = ts
	}
	return ts
}

// FindFields returns fields from store based on current written fields.
func (ms *metricStore) FindFields(fields field.Metas) (found field.Metas) {
	for _, f := range fields {
		fm, ok := ms.fields.Find(f.Name)
		if ok {
			found = append(found, fm)
		}
	}
	return
}

// IndexTStore adds memorty store series id -> global series id mapping.
func (ms *metricStore) IndexTStore(seriesID, seriesIdx uint32) {
	ms.ids.Put(seriesID, seriesIdx)
}

// FlushMetricsDataTo Writes metric-data to the table.
func (ms *metricStore) FlushMetricsDataTo(
	flusher metricsdata.Flusher,
	flushCtx *flushContext,
	flushFields func(memSeriesID uint32, fields field.Metas) error,
) (err error) {
	slotRange := ms.slotRange
	var fields field.Metas
	for idx := range ms.fields {
		f := ms.fields[idx]
		if f.Persisted {
			fields = append(fields, f)
		}
	}
	// field not exist, return
	fieldLen := len(fields)
	if fieldLen == 0 {
		return
	}
	// prepare for flushing metric
	flusher.PrepareMetric(flushCtx.metricID, fields)
	// set current family's slot range
	flushCtx.Start, flushCtx.End = slotRange.GetRange()
	if err := ms.ids.WalkEntry(func(seriesID, memSeriesID uint32) error {
		if err := flushFields(memSeriesID, fields); err != nil {
			return err
		}
		return flusher.FlushSeries(seriesID)
	}); err != nil {
		return err
	}
	return flusher.CommitMetric(flushCtx.SlotRange)
}

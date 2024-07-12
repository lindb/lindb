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
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./time_series_index.go -destination=./time_series_index_mock.go -package=memdb

// TimeSeriesIndex represents metric level time series index(shard level).
// 1. tags hash => memory time series id;
// 2. global time series id => memory time series id;
// 3. time range of data family;
type TimeSeriesIndex interface {
	// GenMemTimeSeriesID generates memory time series id based on tags hash.
	GenMemTimeSeriesID(tags uint64, newID func() uint32) (memSeriesID uint32, isNew bool)
	// IndexTimeSeries indexes memory time series id.(link global time series id => memory time series id)
	IndexTimeSeries(seriesID, memSeriesID uint32)
	// MemTimeSeriesIDs returns all memory time series ids.
	MemTimeSeriesIDs() *roaring.Bitmap
	// TimeSeriesIDs returns all global time series ids under current memory time series index.
	TimeSeriesIDs() *roaring.Bitmap
	// StoreTimeRange stores family level time slat.
	StoreTimeRange(familyCreateTime int64, slot uint16)
	// GetTimeRange returns family level time slot range.
	GetTimeRange(familyCreateTime int64) (*timeutil.SlotRange, bool)
	// ClearTimeRange clears family level time slot range.
	ClearTimeRange(familyCreateTime int64)
	// Load loads field data based search context and time series ids.
	Load(
		ctx *flow.DataLoadContext,
		seriesIDHighKey uint16,
		slotRange timeutil.SlotRange,
		fields []*fieldEntry,
	)
	// FlushMetricsDataTo flushes metric data.
	FlushMetricsDataTo(
		tableFlusher metricsdata.Flusher,
		flushFields func(memSeriesID uint32) error,
	) (err error)
}

// timeSeriesIndex implements TimeSeriesIndex interface.
type timeSeriesIndex struct {
	hashes sync.Map             // tag hash => memory time series id(map[uint64]uint32)
	ids    *imap.IntMap[uint32] // global series id => memory time series id

	families sync.Map // family create timestamp(ns) => metric write time range(map[uint64]*timeutil.SlotRange)

	lock sync.RWMutex
}

// NewTimeSeriesIndex creates TimeSeriesIndex instance.
func NewTimeSeriesIndex() TimeSeriesIndex {
	return &timeSeriesIndex{
		ids: imap.NewIntMap[uint32](),
	}
}

// IndexTimeSeries indexes memory time series id.(link global time series id => memory time series id)
func (idx *timeSeriesIndex) IndexTimeSeries(seriesID, memSeriesID uint32) {
	idx.lock.Lock()
	defer idx.lock.Unlock()

	idx.ids.PutIfNotExist(seriesID, memSeriesID)
}

// GenMemTimeSeriesID generates memory time series id based on tags hash.
func (idx *timeSeriesIndex) GenMemTimeSeriesID(tags uint64, newID func() uint32) (memSeriesID uint32, isNew bool) {
	memTimeSeriesID, ok := idx.hashes.Load(tags)
	if ok {
		return memTimeSeriesID.(uint32), false
	}
	idx.lock.Lock()
	defer idx.lock.Unlock()

	return idx.genMemTimeSeriesID(tags, newID)
}

func (idx *timeSeriesIndex) genMemTimeSeriesID(tags uint64, newID func() uint32) (memSeriesID uint32, isNew bool) {
	memTimeSeriesID, ok := idx.hashes.Load(tags)
	if ok {
		return memTimeSeriesID.(uint32), false
	}

	newMemTimeSeriesID := newID()

	// store new time series id mapping
	idx.hashes.Store(tags, newMemTimeSeriesID)
	return newMemTimeSeriesID, true
}

// MemTimeSeriesIDs returns all memory time series ids.
func (idx *timeSeriesIndex) MemTimeSeriesIDs() *roaring.Bitmap {
	ids := roaring.New()
	idx.hashes.Range(func(key, value any) bool {
		ids.Add(value.(uint32))
		return true
	})
	return ids
}

// TimeSeriesIDs returns all global time series ids under current memory time series index.
func (idx *timeSeriesIndex) TimeSeriesIDs() *roaring.Bitmap {
	idx.lock.RLock()
	defer idx.lock.RUnlock()

	return idx.ids.Keys().Clone()
}

// StoreTimeRange stores family level time slat.
func (idx *timeSeriesIndex) StoreTimeRange(familyCreateTime int64, slot uint16) {
	slotRange, ok := idx.families.Load(familyCreateTime)
	if !ok {
		newSlotRange := timeutil.NewSlotRange(slot, slot)
		idx.families.Store(familyCreateTime, &newSlotRange)
	} else {
		(slotRange.(*timeutil.SlotRange)).SetSlot(slot)
	}
}

// ClearTimeRange clears family level time slot range.
func (idx *timeSeriesIndex) ClearTimeRange(familyCreateTime int64) {
	idx.families.Delete(familyCreateTime)
}

// GetTimeRange returns family level time slot range.
func (idx *timeSeriesIndex) GetTimeRange(familyCreateTime int64) (*timeutil.SlotRange, bool) {
	slotRange, ok := idx.families.Load(familyCreateTime)
	if !ok {
		return nil, false
	}
	return slotRange.(*timeutil.SlotRange), true
}

// Load loads field data based search context and time series ids.
func (idx *timeSeriesIndex) Load(
	ctx *flow.DataLoadContext,
	seriesIDHighKey uint16,
	slotRange timeutil.SlotRange,
	fields []*fieldEntry,
) {
	idx.lock.RLock()
	defer idx.lock.RUnlock()

	highContainerIdx := idx.ids.Keys().GetContainerIndex(seriesIDHighKey)
	if highContainerIdx == -1 {
		// not found
		return
	}
	lowContainer := idx.ids.Keys().GetContainerAtIndex(highContainerIdx)
	memTimeSeriesIDs := idx.ids.Values()[highContainerIdx]

	ctx.IterateLowSeriesIDs(lowContainer, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
		memTimeSeriesID := memTimeSeriesIDs[seriesIdxFromStorage]
		for _, fm := range fields {
			// read field compress data
			compress := fm.getCompressBuf(memTimeSeriesID)
			var tsd *encoding.TSDDecoder
			size := len(compress)
			if size > 0 {
				tsd = ctx.Decoder
				tsd.Reset(compress)
				ctx.DownSampling(slotRange, seriesIdxFromQuery, int(fm.field.Index), tsd)
			}
			// read field current write buffer
			buf, ok := fm.getPage(memTimeSeriesID)
			if ok {
				fm.Reset(buf)
				ctx.DownSampling(slotRange, seriesIdxFromQuery, int(fm.field.Index), fm)
			}
		}
	})
}

// FlushMetricsDataTo flushes metric data.
func (idx *timeSeriesIndex) FlushMetricsDataTo(
	tableFlusher metricsdata.Flusher,
	flushFields func(memSeriesID uint32) error,
) (err error) {
	var timeSeriesIDs *roaring.Bitmap
	idx.lock.RLock()
	timeSeriesIDs = idx.ids.Keys().Clone() // NOTE: must copy keys, lock free when flush data
	idx.lock.RUnlock()

	return idx.flushMetricsDataTo(timeSeriesIDs, tableFlusher, flushFields)
}

// flushMetricsDataTo flushes metric data.
func (idx *timeSeriesIndex) flushMetricsDataTo(
	timeSeriesIDs *roaring.Bitmap,
	tableFlusher metricsdata.Flusher,
	flushFields func(memSeriesID uint32) error,
) (err error) {
	timeSeriesIDsIt := timeSeriesIDs.Iterator()
	for timeSeriesIDsIt.HasNext() {
		seriesID := timeSeriesIDsIt.Next()
		memSeriesID, ok := idx.getMemSeriesID(seriesID)
		if !ok {
			continue
		}
		// flush fields of time series
		if err := flushFields(memSeriesID); err != nil {
			return err
		}

		// commit time series
		if err := tableFlusher.FlushSeries(seriesID); err != nil {
			return err
		}
	}
	return nil
}

// getMemSeriesID returns memory time series id based on global series id.
func (idx *timeSeriesIndex) getMemSeriesID(seriesID uint32) (uint32, bool) {
	idx.lock.RLock()
	defer idx.lock.RUnlock()

	return idx.ids.Get(seriesID)
}

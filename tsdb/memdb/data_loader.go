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
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

// timeSeriesLoader represents time series store loader.
type timeSeriesLoader struct {
	db              *memoryDatabase
	mStore          *metricStore
	seriesIDHighKey uint16
	fields          field.Metas        // metric store field meta
	slotRange       timeutil.SlotRange // slot range of metric store
}

// NewTimeSeriesLoader creates a time series store loader.
func NewTimeSeriesLoader(
	db *memoryDatabase,
	mStore *metricStore,
	seriesIDHighKey uint16,
	fields field.Metas,
	slotRange timeutil.SlotRange,
) flow.DataLoader {
	return &timeSeriesLoader{
		db:              db,
		mStore:          mStore,
		seriesIDHighKey: seriesIDHighKey,
		fields:          fields,
		slotRange:       slotRange,
	}
}

// Load implements flow.DataLoader
func (tsl *timeSeriesLoader) Load(ctx *flow.DataLoadContext) {
	release := tsl.db.WithLock()
	defer release()

	keys := tsl.mStore.ids.Keys()
	highContainerIdx := keys.GetContainerIndex(tsl.seriesIDHighKey)
	lowContainer := keys.GetContainerAtIndex(highContainerIdx)
	fStores := tsl.mStore.ids.Values()[highContainerIdx]

	ctx.IterateLowSeriesIDs(lowContainer, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
		tsKey := fStores[seriesIdxFromStorage]
		for idx := range tsl.fields {
			fm := tsl.fields[idx]

			tsStores := tsl.db.timeSeriesStores[fm.Index]
			fStore, ok := tsStores.Get(tsKey)
			if ok {
				// read field data
				fStore.Load(ctx, seriesIdxFromQuery, int(fm.Index), fm.Type, tsl.slotRange)
			}
		}
	})
}

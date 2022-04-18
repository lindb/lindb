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

	"github.com/lindb/roaring"
)

// metricStoreLoader implements flow.DataLoader interface that loads metric data from memory storage.
type metricStoreLoader struct {
	db               MemoryDatabase
	lowContainer     roaring.Container
	timeSeriesStores []tStoreINTF
	slotRange        timeutil.SlotRange // slot range of metric store
	fields           field.Metas        // sort by field id
}

// newMetricStoreLoader creates a memory storage metric loader.
func newMetricStoreLoader(db MemoryDatabase,
	lowContainer roaring.Container,
	timeSeriesStores []tStoreINTF,
	slotRange timeutil.SlotRange,
	fields field.Metas,
) flow.DataLoader {
	return &metricStoreLoader{
		db:               db,
		lowContainer:     lowContainer,
		timeSeriesStores: timeSeriesStores,
		slotRange:        slotRange,
		fields:           fields,
	}
}

// Load loads the metric data by given series id from memory storage.
func (s *metricStoreLoader) Load(loadCtx *flow.DataLoadContext) {
	release := s.db.WithLock()
	defer release()

	loadCtx.IterateLowSeriesIDs(s.lowContainer, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
		store := s.timeSeriesStores[seriesIdxFromStorage]
		// read series data of fields
		store.load(loadCtx, seriesIdxFromQuery, s.fields, s.slotRange)
	})
}

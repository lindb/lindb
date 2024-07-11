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
)

// timeSeriesLoader represents time series store loader.
type timeSeriesLoader struct {
	db              *memoryDatabase
	timeSeriesIndex TimeSeriesIndex
	fields          []*fieldEntry
	slotRange       timeutil.SlotRange // slot range of metric memory store
	seriesIDHighKey uint16
}

// NewTimeSeriesLoader creates a time series store loader.
func NewTimeSeriesLoader(
	db *memoryDatabase,
	timeSeriesIndex TimeSeriesIndex,
	seriesIDHighKey uint16,
	slotRange timeutil.SlotRange,
	fields []*fieldEntry,
) flow.DataLoader {
	return &timeSeriesLoader{
		db:              db,
		timeSeriesIndex: timeSeriesIndex,
		seriesIDHighKey: seriesIDHighKey,
		fields:          fields,
		slotRange:       slotRange,
	}
}

// Load implements flow.DataLoader.
func (tsl *timeSeriesLoader) Load(ctx *flow.DataLoadContext) {
	tsl.timeSeriesIndex.Load(ctx, tsl.seriesIDHighKey, tsl.slotRange, tsl.fields)
}

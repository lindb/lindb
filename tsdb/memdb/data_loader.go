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
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

// timeSeriesLoader represents time series store loader.
type timeSeriesLoader struct {
	db        *memoryDatabase
	seriesIDs *flow.LowSeriesIDs
	decoder   *encoding.TSDDecoder

	memTimeSeriesIDs []uint32
	fields           []*fieldEntry

	slotRange timeutil.SlotRange // slot range of metric memory store
}

// NewTimeSeriesLoader creates a time series store loader.
func NewTimeSeriesLoader(
	db *memoryDatabase,
	seriesIDs *flow.LowSeriesIDs,
	memTimeSeriesIDs []uint32,
	slotRange timeutil.SlotRange,
	fields []*fieldEntry,
) flow.DataLoader {
	return &timeSeriesLoader{
		db:        db,
		fields:    fields,
		slotRange: slotRange,

		seriesIDs:        seriesIDs,
		memTimeSeriesIDs: memTimeSeriesIDs,
		decoder:          encoding.GetTSDDecoder(), // TODO: refact
	}
}

func (tsl *timeSeriesLoader) Load(seriesID uint16, fn func(field field.Meta, geter encoding.TSDValueGetter)) {
	index, ok := tsl.seriesIDs.Find(seriesID)
	// FIXME: add lock????
	if ok {
		memTimeSeriesID := tsl.memTimeSeriesIDs[index]
		for _, fm := range tsl.fields {
			// read field compress data
			compress := fm.getCompressBuf(memTimeSeriesID)
			size := len(compress)
			if size > 0 {
				tsl.decoder.Reset(compress)
				fn(fm.field, tsl.decoder)
			}
			// read current field write buffer
			buf, ok := fm.getPage(memTimeSeriesID)
			if ok {
				fm.Reset(buf)
				fn(fm.field, fm)
			}
		}
	}
}

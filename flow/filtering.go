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

package flow

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./filtering.go -destination=./filtering_mock.go -package=flow

// DataFilter represents the filter ability over memory database and files under data family.
type DataFilter interface {
	// Filter filters the data based on metricID/fields/seriesIDs/timeRange,
	// if finds data then returns filter result set, else returns nil.
	Filter(ctx *MetricScanContext) ([]FilterResultSet, error)
}

// FilterResultSet represents the filter result set, loads data based on this interface.
type FilterResultSet interface {
	// Interval returns the interval of storage.
	Interval() timeutil.Interval
	// Identifier identifies the source of result set(mem/kv etc.).
	Identifier() string
	// FamilyTime returns the family time of storage.
	FamilyTime() int64
	// SlotRange returns the slot range of storage.
	SlotRange() timeutil.SlotRange
	// Load loads the data from storage, then returns the data loader.
	Load(seriesIDHighKey uint16, lowSeriesIDs roaring.Container) DataLoader
	// SeriesIDs returns the series ids which matches with query series ids.
	SeriesIDs() *roaring.Bitmap
	// Close release the resource during doing query operation.
	Close()
}

// LoaderCallback represents the callback function after loading data.
// NOTE: field must be one of query fields.
type LoaderCallback func(fieldOfQuery field.Meta, getter encoding.TSDValueGetter)

// DataLoader represents the loader which load metric data from storage.
type DataLoader interface {
	Load(seriesID uint16, fn LoaderCallback)
}

type LowSeriesIDs struct {
	lowSeriesIDs roaring.Container

	it       roaring.PeekableShortIterator
	index    int
	seriesID uint16
}

func NewLowSeriesIDs(lowSeriesIDs roaring.Container) *LowSeriesIDs {
	return &LowSeriesIDs{lowSeriesIDs: lowSeriesIDs}
}

func (ids *LowSeriesIDs) Find(target uint16) (index int, ok bool) {
	// TODO: get value from container
	if ids.it != nil {
		if ids.seriesID == target {
			// previous series id match
			return ids.index - 1, true
		} else if ids.seriesID > target {
			return -1, false
		}
	}
	if ids.it == nil {
		// initial low series ids iterator
		ids.it = ids.lowSeriesIDs.PeekableIterator()
		ids.index = 0
	}
	for ids.it.HasNext() {
		ids.seriesID = ids.it.Next()
		ids.index++
		if ids.seriesID == target {
			// found series id
			return ids.index - 1, true
		} else if ids.seriesID > target {
			break
		}
		// series id < target series id, loop next series id
	}
	return -1, false
}

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

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./filtering.go -destination=./filtering_mock.go -package=flow

// DataFilter represents the filter ability over memory database and files under data family.
type DataFilter interface {
	// Filter filters the data based on metricIDs/fields/seriesIDs/timeRange,
	// if finds data then returns filter result set, else returns nil.
	Filter(metricID uint32,
		seriesIDs *roaring.Bitmap,
		timeRange timeutil.TimeRange,
		fields field.Metas,
	) ([]FilterResultSet, error)
}

// FilterResultSet represents the filter result set, loads data based on this interface.
type FilterResultSet interface {
	// Identifier identifies the source of result set(mem/kv etc.).
	Identifier() string
	// FamilyTime returns the family time of storage.
	FamilyTime() int64
	// SlotRange returns the slot range of storage.
	SlotRange() timeutil.SlotRange
	// Load loads the data from storage, then returns the data loader.
	Load(highKey uint16, seriesID roaring.Container) DataLoader
	// SeriesIDs returns the series ids which matches with query series ids.
	SeriesIDs() *roaring.Bitmap
}

// DataLoader represents the loader which load metric data from storage.
type DataLoader interface {
	// Load loads the metric data by given low series id.
	Load(lowSeriesID uint16) (timeutil.SlotRange, [][]byte)
}

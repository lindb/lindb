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
	// Filter filters the data based on metricIDs/fieldIDs/seriesIDs/timeRange,
	// if finds data then returns filter result set, else returns nil.
	Filter(metricID uint32, fieldIDs []field.ID,
		seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
	) ([]FilterResultSet, error)
}

// FilterResultSet represents the filter result set, loads data and does down sampling need based on this interface.
type FilterResultSet interface {
	// Identifier identifies the source of result set(mem/kv etc.)
	Identifier() string
	// Load loads the data from storage, then returns the data scanner.
	Load(highKey uint16, seriesID roaring.Container, fieldIDs []field.ID) Scanner
	// SeriesIDs returns the series ids which matches with query series ids
	SeriesIDs() *roaring.Bitmap
}

// Scanner represents the scanner which scan metric data from storage.
type Scanner interface {
	// Scan scans the metric data by given series id.
	Scan(lowSeriesID uint16) [][]byte
}

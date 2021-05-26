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

package metricsdata

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./filter.go -destination=./filter_mock.go -package metricsdata

// Filter implements filtering metrics from sst files.
type Filter interface {
	// Filter filters data under each sst file based on query condition
	Filter(seriesIDs *roaring.Bitmap, fields field.Metas) ([]flow.FilterResultSet, error)
}

// metricsDataFilter represents the sst file data filter
type metricsDataFilter struct {
	familyTime int64
	snapshot   version.Snapshot //FIXME stone1100, need close version snapshot
	readers    []MetricReader
}

// NewFilter creates the sst file data filter
func NewFilter(familyTime int64, snapshot version.Snapshot, readers []MetricReader) Filter {
	return &metricsDataFilter{
		familyTime: familyTime,
		snapshot:   snapshot,
		readers:    readers,
	}
}

// Filter filters the data under each sst file based on metric/version/seriesIDs,
// if finds data then returns the flow.FilterResultSet, else returns nil
func (f *metricsDataFilter) Filter(
	seriesIDs *roaring.Bitmap, fields field.Metas,
) (rs []flow.FilterResultSet, err error) {
	for _, reader := range f.readers {
		//FIXME add time range compare????
		fieldMetas, _ := reader.GetFields().Intersects(fields)
		if len(fieldMetas) == 0 {
			// field not found
			continue
		}
		// after and operator, query bitmap is sub of store bitmap
		matchSeriesIDs := roaring.FastAnd(seriesIDs, reader.GetSeriesIDs())
		if matchSeriesIDs.IsEmpty() {
			// series ids not found
			continue
		}
		rs = append(rs, newFileFilterResultSet(f.familyTime, fields, matchSeriesIDs, reader))
	}
	// not founds
	if len(rs) == 0 {
		return nil, constants.ErrNotFound
	}
	return
}

// fileFilterResultSet represents sst file metricReader for loading file data based on query condition
type fileFilterResultSet struct {
	reader     MetricReader
	familyTime int64
	fields     field.Metas
	seriesIDs  *roaring.Bitmap
}

// newFileFilterResultSet creates the file filter result set
func newFileFilterResultSet(familyTime int64, fields field.Metas,
	seriesIDs *roaring.Bitmap, reader MetricReader,
) flow.FilterResultSet {
	return &fileFilterResultSet{
		familyTime: familyTime,
		reader:     reader,
		fields:     fields,
		seriesIDs:  seriesIDs,
	}
}

// Identifier identifies the source of result set from kv store
func (f *fileFilterResultSet) Identifier() string {
	return f.reader.Path()
}

// SeriesIDs returns the series ids which matches with query series ids
func (f *fileFilterResultSet) SeriesIDs() *roaring.Bitmap {
	return f.seriesIDs
}

// FamilyTime returns the family time of storage.
func (f *fileFilterResultSet) FamilyTime() int64 {
	return f.familyTime
}

// SlotRange returns the slot range of storage.
func (f *fileFilterResultSet) SlotRange() timeutil.SlotRange {
	return f.reader.GetTimeRange()
}

// Load reads data from sst files, then returns the data file scanner.
func (f *fileFilterResultSet) Load(highKey uint16, seriesID roaring.Container) flow.DataLoader {
	return f.reader.Load(highKey, seriesID, f.fields)
}

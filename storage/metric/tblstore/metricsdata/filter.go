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
	slotRange  timeutil.SlotRange
	interval   timeutil.Interval
	snapshot   version.Snapshot
	readers    []MetricReader
}

// NewFilter creates the sst file data filter
func NewFilter(familyTime int64, interval timeutil.Interval, slotRange timeutil.SlotRange,
	snapshot version.Snapshot, readers []MetricReader,
) Filter {
	return &metricsDataFilter{
		familyTime: familyTime,
		slotRange:  slotRange,
		interval:   interval,
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
		if fields.Len() > 0 {
			fieldMetas, _ := reader.GetFields().Intersects(fields)
			if len(fieldMetas) == 0 {
				// field not found
				continue
			}
		}
		// after and operator, query bitmap is sub of store bitmap
		matchSeriesIDs := roaring.FastAnd(seriesIDs, reader.GetSeriesIDs())
		if matchSeriesIDs.IsEmpty() {
			// series ids not found
			continue
		}
		slotRnage := reader.GetTimeRange()
		rs = append(rs, newFileFilterResultSet(f.familyTime, f.interval, (&slotRnage).Intersect(f.slotRange),
			matchSeriesIDs, fields, reader, f.snapshot))
	}
	// not founds
	if len(rs) == 0 {
		return nil, constants.ErrNotFound
	}
	return
}

// fileFilterResultSet represents sst file metricReader for loading file data based on query condition
type fileFilterResultSet struct {
	snapshot   version.Snapshot
	slotRange  timeutil.SlotRange
	reader     MetricReader
	familyTime int64
	interval   timeutil.Interval
	fields     field.Metas
	seriesIDs  *roaring.Bitmap
}

// newFileFilterResultSet creates the file filter result set
func newFileFilterResultSet(
	familyTime int64,
	interval timeutil.Interval,
	slotRange timeutil.SlotRange,
	seriesIDs *roaring.Bitmap,
	fields field.Metas,
	reader MetricReader,
	snapshot version.Snapshot,
) flow.FilterResultSet {
	return &fileFilterResultSet{
		familyTime: familyTime,
		interval:   interval,
		slotRange:  slotRange,
		seriesIDs:  seriesIDs,
		fields:     fields,
		reader:     reader,
		snapshot:   snapshot,
	}
}

// Interval returns the interval of file storage.
func (f *fileFilterResultSet) Interval() timeutil.Interval {
	return f.interval
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
	return f.slotRange
}

// Load reads data from sst files, then returns the data file scanner.
func (f *fileFilterResultSet) Load(seriesIDHighKey uint16, lowSeriesIDs roaring.Container) flow.DataLoader {
	return f.reader.Load(seriesIDHighKey, lowSeriesIDs, f.fields)
}

// Close release the resource during doing query operation.
func (f *fileFilterResultSet) Close() {
	// release kv snapshot
	f.snapshot.Close()
}

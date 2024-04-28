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

package index

import (
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"

	"path"
	"strconv"
	"sync"
	"time"
)

var (
	mkDirFunc            = fileutil.MkDirIfNotExist
	getDirectoryListFunc = fileutil.GetDirectoryList
	newIndexDBFunc       = NewMetricIndexDatabase
)

type metricIndexSegment struct {
	dir     string
	metaDB  MetricMetaDatabase
	indices map[int]MetricIndexDatabase
	lock    sync.Mutex
}

// NewMetricIndexSegment creates a metric index segment store.
func NewMetricIndexSegment(dir string, metaDB MetricMetaDatabase) (segment MetricIndexSegment, err error) {
	if err0 := mkDirFunc(dir); err0 != nil {
		return nil, err0
	}
	dirs, err := getDirectoryListFunc(dir)
	if err != nil {
		return
	}
	var indices = make(map[int]MetricIndexDatabase)
	for _, d := range dirs {
		t, err0 := time.Parse("200601", d)
		if err0 != nil {
			continue
		}
		segment := t.Year()*100 + int(t.Month())
		database, err0 := NewMetricIndexDatabase(path.Join(dir, d), metaDB)
		if err0 != nil {
			return nil, err0
		}
		indices[segment] = database
	}
	return &metricIndexSegment{
		dir:     dir,
		metaDB:  metaDB,
		indices: indices,
	}, nil
}

func (m *metricIndexSegment) Close() error {
	if len(m.indices) == 0 {
		return nil
	}
	var (
		indices = m.indices
		g       = sync.WaitGroup{}
		err     error
	)
	g.Add(len(indices))
	for _, index := range indices {
		index := index
		go func() {
			defer g.Done()
			if err0 := index.Close(); err0 != nil {
				err = err0
			}
		}()
	}
	g.Wait()
	return err
}

func (m *metricIndexSegment) PrepareFlush() {}

func (m *metricIndexSegment) getSegment(timeRange timeutil.TimeRange) (start, end int) {
	start = timeutil.GetSegment(timeRange.Start)
	end = timeutil.GetSegment(timeRange.End)
	return
}

func (m *metricIndexSegment) Flush() error {
	if len(m.indices) == 0 {
		return nil
	}
	var (
		indices = m.indices
		g       = sync.WaitGroup{}
		err     error
	)
	g.Add(len(indices))
	for _, index := range indices {
		index := index
		go func() {
			defer g.Done()
			ch := make(chan error, 1)
			index.Notify(&FlushNotifier{
				Callback: func(err error) {
					ch <- err
				},
			})
			if err0 := <-ch; err0 != nil {
				err = err0
			}
		}()
	}
	g.Wait()
	return err
}

func (m *metricIndexSegment) GetSeriesIDsByTagValueIDs(
	tagKeyID tag.KeyID,
	tagValueIDs *roaring.Bitmap,
	timeRange timeutil.TimeRange,
) (*roaring.Bitmap, error) {
	var (
		seriesID = roaring.NewBitmap()
		indices  = m.indices
	)
	start, end := m.getSegment(timeRange)
	for segment, index := range indices {
		if segment >= start && segment <= end {
			id, err := index.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
			if err != nil {
				return nil, err
			}
			seriesID.Or(id)
		}
	}
	return seriesID, nil
}

func (m *metricIndexSegment) GetSeriesIDsForTag(
	tagKeyID tag.KeyID,
	timeRange timeutil.TimeRange,
) (*roaring.Bitmap, error) {
	var (
		seriesID = roaring.NewBitmap()
		indices  = m.indices
	)
	start, end := m.getSegment(timeRange)
	for segment, index := range indices {
		if segment >= start && segment <= end {
			id, err := index.GetSeriesIDsForTag(tagKeyID)
			if err != nil {
				return nil, err
			}
			seriesID.Or(id)
		}
	}
	return seriesID, nil
}

func (m *metricIndexSegment) GetSeriesIDsForMetric(
	metricID metric.ID,
	timeRange timeutil.TimeRange,
) (*roaring.Bitmap, error) {
	var (
		seriesID = roaring.NewBitmap()
		indices  = m.indices
	)
	start, end := m.getSegment(timeRange)
	for segment, index := range indices {
		if segment >= start && segment <= end {
			id, err := index.GetSeriesIDsForMetric(metricID)
			if err != nil {
				return nil, err
			}
			seriesID.Or(id)
		}
	}
	return seriesID, nil
}

// GetGroupingContext returns the context of group by
func (m *metricIndexSegment) GetGroupingContext(ctx *flow.ShardExecuteContext) error {
	var (
		timeRange  = ctx.StorageExecuteCtx.Query.TimeRange
		scannerMap = make(map[tag.KeyID][]flow.GroupingScanner)
		indices    = m.indices
	)
	start, end := m.getSegment(timeRange)
	for segment, index := range indices {
		if segment >= start && segment <= end {
			scanners, err := index.GetGroupingContext(ctx)
			if err != nil {
				continue
			}
			for id, ss := range scanners {
				scannerMap[id] = append(scannerMap[id], ss...)
			}
		}
	}
	ctx.GroupingContext = flow.NewGroupContext(ctx.StorageExecuteCtx.GroupByTagKeyIDs, scannerMap)
	return nil
}

func (m *metricIndexSegment) GetOrCreateIndex(familyTime int64) (MetricIndexDatabase, error) {
	segment := timeutil.GetSegment(familyTime)
	index, ok := m.indices[segment]
	if ok {
		return index, nil
	}
	m.lock.Lock()
	index, ok = m.indices[segment]
	if ok {
		m.lock.Unlock()
		return index, nil
	}
	defer m.lock.Unlock()

	index, err := newIndexDBFunc(path.Join(m.dir, strconv.Itoa(segment)), m.metaDB)
	if err != nil {
		return nil, err
	}
	newIndices := make(map[int]MetricIndexDatabase, len(m.indices))
	for k, v := range m.indices {
		newIndices[k] = v
	}
	newIndices[segment] = index
	m.indices = newIndices
	return index, nil
}

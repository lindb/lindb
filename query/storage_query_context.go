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

package query

import (
	"sort"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/lindb/roaring"
)

// storageExecuteContext represents storage query execute context
type storageExecuteContext struct {
	query    *stmt.Query
	shardIDs []int32

	tagFilterResult map[string]*tagFilterResult

	stats *models.StorageStats // storage query stats track for explain query
}

// newStorageExecuteContext creates storage execute context
func newStorageExecuteContext(shardIDs []int32, query *stmt.Query) *storageExecuteContext {
	ctx := &storageExecuteContext{
		query:    query,
		shardIDs: shardIDs,
	}
	if query.Explain {
		// if explain query, create storage query stats
		ctx.stats = models.NewStorageStats()
	}
	return ctx
}

// QueryStats returns the storage query stats
func (ctx *storageExecuteContext) QueryStats() *models.StorageStats {
	if ctx.stats != nil {
		ctx.stats.Complete()
	}
	return ctx.stats
}

// setTagFilterResult sets tag filter result
func (ctx *storageExecuteContext) setTagFilterResult(tagFilterResult map[string]*tagFilterResult) {
	ctx.tagFilterResult = tagFilterResult
}

// timeSpans represents the time span slice in query time range.
type timeSpans []*timeSpan

func (f timeSpans) Len() int           { return len(f) }
func (f timeSpans) Less(i, j int) bool { return f[i].familyTime < f[j].familyTime }
func (f timeSpans) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// timeSpan represents a time span in query time range.
type timeSpan struct {
	identifier     string
	familyTime     int64
	source, target timeutil.SlotRange
	interval       timeutil.Interval

	resultSets []flow.FilterResultSet
	loaders    []flow.DataLoader
}

type timeSpanResultSet struct {
	spanMap   map[int64]*timeSpan
	seriesIDs *roaring.Bitmap

	filterRSCount int
}

func newTimeSpanResultSet() *timeSpanResultSet {
	return &timeSpanResultSet{
		spanMap:   make(map[int64]*timeSpan),
		seriesIDs: roaring.New(),
	}
}

func (s *timeSpanResultSet) addFilterResultSet(interval timeutil.Interval, rs flow.FilterResultSet) {
	familyTime := rs.FamilyTime()
	span, ok := s.spanMap[familyTime]
	if !ok {
		span = &timeSpan{
			identifier: rs.Identifier(),
			familyTime: familyTime,
			source:     rs.SlotRange(),
			target:     rs.SlotRange(),
			interval:   interval,
		}
		s.spanMap[familyTime] = span
	} else {
		// calc target slot range
		target := span.target
		source := rs.SlotRange()
		span.target = *((&target).Intersect(&source))
	}

	span.resultSets = append(span.resultSets, rs)

	// increase filter rs
	s.filterRSCount++

	// merge all series ids after filtering => final series ids
	s.seriesIDs.Or(rs.SeriesIDs())
}

func (s *timeSpanResultSet) getFilterRSCount() int {
	return s.filterRSCount
}

func (s *timeSpanResultSet) isEmpty() bool {
	return len(s.spanMap) == 0
}

func (s *timeSpanResultSet) getTimeSpans() timeSpans {
	var timeSpans timeSpans
	for _, span := range s.spanMap {
		timeSpans = append(timeSpans, span)
	}
	sort.Sort(timeSpans)
	return timeSpans
}

// getSeriesIDs returns final series ids after family filtering.
func (s *timeSpanResultSet) getSeriesIDs() *roaring.Bitmap {
	return s.seriesIDs
}

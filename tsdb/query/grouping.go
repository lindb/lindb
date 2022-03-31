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
	"encoding/binary"
	"strings"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
)

// GroupingContext represents the context of group by query for tag keys
// builds tags => series ids mapping, using such as counting sort
// https://en.wikipedia.org/wiki/Counting_sort
type GroupingContext struct {
	tagKeys  []tag.KeyID
	scanners map[tag.KeyID][]flow.GroupingScanner
}

// NewGroupContext creates a GroupingContext
func NewGroupContext(tagKeys []tag.KeyID, scanners map[tag.KeyID][]flow.GroupingScanner) flow.GroupingContext {
	return &GroupingContext{
		tagKeys:  tagKeys,
		scanners: scanners,
	}
}

// ScanTagValueIDs scans grouping context by high key/container of series ids,
// then returns grouped tag value ids for each tag key
func (g *GroupingContext) ScanTagValueIDs(highKey uint16, container roaring.Container) []*roaring.Bitmap {
	result := make([]*roaring.Bitmap, len(g.tagKeys))
	for i, tagKey := range g.tagKeys {
		scanners := g.scanners[tagKey]
		tagValues := roaring.New()
		result[i] = tagValues
		for _, scanner := range scanners {
			// get series ids/tag value ids mapping by high key
			lowContainer, tagValueIDs := scanner.GetSeriesAndTagValue(highKey)
			if lowContainer == nil {
				// high key not exist
				continue
			}
			// iterator all series ids after filtering
			it := lowContainer.PeekableIterator()
			idx := 0
			for it.HasNext() {
				seriesID := it.Next()
				if container.Contains(seriesID) {
					tagValues.Add(tagValueIDs[idx])
				}
				idx++
			}
		}
	}
	return result
}

// BuildGroup builds the grouped series ids by the high key of series id
// and the container includes low keys of series id.
func (g *GroupingContext) BuildGroup(ctx *flow.DataLoadContext) {
	// new tag value ids array for each group by tag key
	groupByTagValueIDs := g.buildTagValueIDs2SeriesIDs(ctx)

	// current group by query completed, need merge group by tag value ids
	ctx.ShardExecuteCtx.StorageExecuteCtx.CollectGroupingTagValueIDs(groupByTagValueIDs)
}

// buildTagValueIDs2SeriesIDs builds tag value id => series id mapping
func (g *GroupingContext) buildTagValueIDs2SeriesIDs(ctx *flow.DataLoadContext) []*roaring.Bitmap {
	// new seriesIDs2Tags array based on range of min ~ max
	seriesIDHighKey := ctx.SeriesIDHighKey
	min := ctx.MinSeriesID
	tagSize := len(g.tagKeys)
	tagValueIDsForTagKeys := make([]*roaring.Bitmap, tagSize)
	tagValueIDsForGrouping := make([][]uint32, tagSize)

	for tagKeyIdx, tagKey := range g.tagKeys {
		scanners := g.scanners[tagKey]
		tagValueIDsForTagKey := roaring.New()
		tagValueIDsForTagKeys[tagKeyIdx] = tagValueIDsForTagKey

		groupingTagValueIDs := make([]uint32, len(ctx.LowSeriesIDs))
		tagValueIDsForGrouping[tagKeyIdx] = groupingTagValueIDs
		for _, scanner := range scanners {
			lowSeriesIDs, tagValueIDs := scanner.GetSeriesAndTagValue(seriesIDHighKey)
			if lowSeriesIDs == nil {
				// high key not exist
				continue
			}
			ctx.IterateLowSeriesIDs(lowSeriesIDs, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
				groupingTagValueIDs[seriesIdxFromQuery] = tagValueIDs[seriesIdxFromStorage]
			})
		}
	}

	var scratch [4]byte
	result := make(map[string]uint16)
	var keyBuilder strings.Builder

	groupingSeriesAggIdx := uint16(0)

	it := ctx.LowSeriesIDsContainer.PeekableIterator()
	for it.HasNext() {
		seriesID := it.Next()
		seriesIdxFromQuery := seriesID - min // series index in query counting sort array
		found := true
		keyBuilder.Reset()
		for tagKeyIdx := range tagValueIDsForGrouping {
			tagValueID := tagValueIDsForGrouping[tagKeyIdx][seriesIdxFromQuery]
			if tagValueID == 0 {
				found = false
				break
			}
			binary.LittleEndian.PutUint32(scratch[:], tagValueID)
			keyBuilder.Write(scratch[:])

			tagValueIDsForTagKeys[tagKeyIdx].Add(tagValueID)
		}
		if found {
			key := keyBuilder.String()
			aggIdx, ok := result[key]
			if !ok {
				ctx.GroupingSeriesAgg = append(ctx.GroupingSeriesAgg, &flow.GroupingSeriesAgg{
					Key:        key,
					Aggregator: ctx.NewSeriesAggregator(),
				})
				aggIdx = groupingSeriesAggIdx
				result[key] = aggIdx
				groupingSeriesAggIdx++
			}
			ctx.GroupingSeriesAggRefs[seriesIdxFromQuery] = aggIdx
		}
	}

	return tagValueIDsForTagKeys
}

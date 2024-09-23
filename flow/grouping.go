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
	"encoding/binary"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source=./grouping.go -destination=./grouping_mock.go -package=flow

// GroupingContext represents the context of group by query for tag keys
type GroupingContext interface {
	// BuildGroup builds the grouped series ids by the high key of series id
	// and the container includes low keys of series id
	BuildGroup(ctx *DataLoadContext)
	// ScanTagValueIDs scans grouping context by high key/container of series ids,
	// then returns grouped tag value ids for each tag key
	ScanTagValueIDs(highKey uint16, container roaring.Container) []*roaring.Bitmap
}

// GroupingScanner represents the scanner which scans the group by data by high key of series id
type GroupingScanner interface {
	// GetSeriesAndTagValue returns group by container and tag value ids
	GetSeriesAndTagValue(highKey uint16) (roaring.Container, []uint32)
	// GetSeriesIDs returns the series ids in current scanner.
	GetSeriesIDs() *roaring.Bitmap
}

// Grouping represents the getter grouping scanners for tag key group by query
type Grouping interface {
	// GetGroupingScanner returns the grouping scanners based on tag key ids and series ids
	GetGroupingScanner(tagKeyID tag.KeyID, seriesIDs *roaring.Bitmap) ([]GroupingScanner, error)
}

// GroupingBuilder represents grouping tag builder.
type GroupingBuilder interface {
	// GetGroupingContext returns the context of group by
	GetGroupingContext(
		groupingTags tag.Metas, seriesIDs *roaring.Bitmap,
	) (*roaring.Bitmap, GroupingContext, error)
}

// groupingContext represents the context of group by query for tag keys
// builds tags => series ids mapping, using such as counting sort
// https://en.wikipedia.org/wiki/Counting_sort
type groupingContext struct {
	tags     tag.Metas // grouping tags
	scanners map[tag.KeyID][]GroupingScanner
}

// NewGroupContext creates a GroupingContext
func NewGroupContext(tags tag.Metas, scanners map[tag.KeyID][]GroupingScanner) GroupingContext {
	return &groupingContext{
		tags:     tags,
		scanners: scanners,
	}
}

// ScanTagValueIDs scans grouping context by high key/container of series ids,
// then returns grouped tag value ids for each tag key
func (g *groupingContext) ScanTagValueIDs(highKey uint16, container roaring.Container) []*roaring.Bitmap {
	result := make([]*roaring.Bitmap, len(g.tags))
	for i, groupingTag := range g.tags {
		scanners := g.scanners[groupingTag.ID]
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
func (g *groupingContext) BuildGroup(ctx *DataLoadContext) {
	if g.tags.Len() == 1 {
		g.buildGroupForSingleTag(ctx)
	} else {
		g.buildGroupForMultiTags(ctx)
	}
}

// buildGroupForMultiTags builds grouping for multi-tags.
func (g *groupingContext) buildGroupForMultiTags(ctx *DataLoadContext) {
	tagSize := g.tags.Len()
	tagValueIDsForGrouping := make([][]byte, len(ctx.LowSeriesIDs))
	groupingTagValueIDs := make([]*roaring.Bitmap, tagSize)
	g.scanGroupingTags(ctx, func(seriesIdxFromQuery uint16, tagKeyIDIdx int, tagValueID uint32) {
		tagValueIDs := tagValueIDsForGrouping[seriesIdxFromQuery]
		if tagValueIDs == nil {
			tagValueIDs = make([]byte, tagSize*4)
			tagValueIDsForGrouping[seriesIdxFromQuery] = tagValueIDs
		}
		tagOffset := tagKeyIDIdx * 4
		binary.LittleEndian.PutUint32(tagValueIDs[tagOffset:], tagValueID)

		// add grouping tag value ids
		if groupingTagValueIDs[tagKeyIDIdx] == nil {
			groupingTagValueIDs[tagKeyIDIdx] = roaring.New()
		}
		groupingTagValueIDs[tagKeyIDIdx].Add(tagValueID)

		if tagKeyIDIdx == tagSize-1 {
			key := strutil.ByteSlice2String(tagValueIDs)
			ctx.CollectGrouping(key, seriesIdxFromQuery)
		}
	})

	// add grouping tag value ids
	ctx.Grouping(groupingTagValueIDs)
}

// buildGroupForMultiTags builds grouping for single-tags.
func (g *groupingContext) buildGroupForSingleTag(ctx *DataLoadContext) {
	// result := make(map[uint32]uint16)
	tagValueIDs := roaring.New()
	var scratch [4]byte
	g.scanGroupingTags(ctx, func(seriesIdxFromQuery uint16, tagKeyIDIdx int, tagValueID uint32) {
		// TODO: opt
		binary.LittleEndian.PutUint32(scratch[:], tagValueID)
		key := string(scratch[:])
		ctx.CollectGrouping(key, seriesIdxFromQuery)
		tagValueIDs.Add(tagValueID)
	})

	// add grouping tag value ids
	ctx.Grouping([]*roaring.Bitmap{tagValueIDs})
}

// scanGroupingTags scans grouping tags(series ids=>tag value ids)
func (g *groupingContext) scanGroupingTags(ctx *DataLoadContext,
	fn func(seriesIdxFromQuery uint16, tagKeyIDIdx int, tagValueID uint32),
) {
	seriesIDHighKey := ctx.SeriesIDHighKey
	for tagKeyIdx, groupingTag := range g.tags {
		scanners := g.scanners[groupingTag.ID]
		for _, scanner := range scanners {
			lowSeriesIDs, tagValueIDs := scanner.GetSeriesAndTagValue(seriesIDHighKey)
			if lowSeriesIDs == nil {
				// high key not exist
				continue
			}
			ctx.IterateLowSeriesIDs(lowSeriesIDs, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
				fn(seriesIdxFromQuery, tagKeyIdx, tagValueIDs[seriesIdxFromStorage])
			})
		}
	}
}

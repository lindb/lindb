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

package metric

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/storage"
)

const (
	tagValueNotFound = "tag_value_not_found"
)

type GroupingKey []uint32

func (key GroupingKey) Clone() *GroupingKey {
	newKey := make(GroupingKey, len(key))
	copy(newKey, key)
	return &newKey
}

type TagsScanner struct {
	scanners    []*flow.TagScanner
	tagValueIDs []*roaring.Bitmap
	key         GroupingKey
}

func NewTagsScanner(scanners []*flow.TagScanner) *TagsScanner {
	return &TagsScanner{
		scanners:    scanners,
		key:         make(GroupingKey, len(scanners)),
		tagValueIDs: make([]*roaring.Bitmap, len(scanners)),
	}
}

func (s *TagsScanner) FindTagValues(lowSeriesID uint16) *GroupingKey {
	for idx, scanner := range s.scanners {
		tagValueID, ok := scanner.FindTagValue(lowSeriesID)
		if s.tagValueIDs[idx] == nil {
			s.tagValueIDs[idx] = roaring.New()
		}
		s.tagValueIDs[idx].Add(tagValueID)
		if ok {
			s.key[idx] = tagValueID
		} else {
			s.key[idx] = tag.EmptyTagValueID
		}
	}
	return &s.key
}

func (s *TagsScanner) GetTagValueIDs() []*roaring.Bitmap {
	return s.tagValueIDs
}

type Grouping struct {
	db   storage.Database
	tags tag.Metas

	// tag value ids for each grouping tag key
	groupingTagValueIDs []*roaring.Bitmap
	tagValuesMap        []map[uint32]string // tag value id=> tag value for each group by tag key
}

func NewGrouping(db storage.Database, tags tag.Metas) *Grouping {
	lenOfTags := tags.Len()
	return &Grouping{
		db:                  db,
		tags:                tags,
		groupingTagValueIDs: make([]*roaring.Bitmap, lenOfTags),
		tagValuesMap:        make([]map[uint32]string, lenOfTags),
	}
}

func (g *Grouping) CollectTagValueIDs(tagValueIDs []*roaring.Bitmap) {
	// TODO: add lock?
	for idx, ids := range tagValueIDs {
		if g.groupingTagValueIDs[idx] == nil {
			g.groupingTagValueIDs[idx] = ids
		} else {
			g.groupingTagValueIDs[idx].Or(ids)
		}
	}
}

func (g *Grouping) CollectTagValues() {
	metaDB := g.db.MetaDB()

	for idx := range g.groupingTagValueIDs {
		tagKey := g.tags[idx]
		tagValueIDs := g.groupingTagValueIDs[idx]

		if tagValueIDs == nil || tagValueIDs.IsEmpty() {
			continue
		}

		tagValues := make(map[uint32]string) // tag value id => tag value
		err := metaDB.CollectTagValues(tagKey.ID, tagValueIDs, tagValues)
		if err != nil {
			panic(err)
		}
		fmt.Printf("collect tag values...%v\n", tagValues)
		g.tagValuesMap[idx] = tagValues
	}
}

func (g *Grouping) GetTagValues(tagValueIDs GroupingKey) []string {
	// TODO: cache grouping tag values
	// if tagValues, ok := ctx.tagsMap[tagValueIDs]; ok {
	// 	return tagValues
	// }

	fmt.Printf("get value values==%v\n", tagValueIDs)
	tagValues := make([]string, g.tags.Len())
	for idx := range g.tagValuesMap {
		tagValuesForKey := g.tagValuesMap[idx]
		tagValueID := tagValueIDs[idx]
		if tagValue, ok := tagValuesForKey[tagValueID]; ok {
			tagValues[idx] = tagValue
		} else {
			fmt.Printf("tag value not found...%v\n", tagValueID)
			tagValues[idx] = tagValueNotFound
		}
	}
	return tagValues
}

type grouping interface {
	GetAggregator(lowSeriesID uint16) []*collections.FloatArray
	ForEach(fn func(tags *GroupingKey, rs []*collections.FloatArray))
}

type groupingWithTags struct {
	aggregators map[*GroupingKey][]*collections.FloatArray
	tagsScanner *TagsScanner
	tableScan   *TableScan
}

func newGroupingWithTags(tagsScanner *TagsScanner, tableScan *TableScan) grouping {
	return &groupingWithTags{
		tagsScanner: tagsScanner,
		tableScan:   tableScan,
		aggregators: make(map[*GroupingKey][]*collections.FloatArray),
	}
}

func (g *groupingWithTags) ForEach(fn func(tags *GroupingKey, rs []*collections.FloatArray)) {
	for tags, aggregator := range g.aggregators {
		fn(tags, aggregator)
	}
}

func (g *groupingWithTags) GetAggregator(lowSeriesID uint16) []*collections.FloatArray {
	key := g.tagsScanner.FindTagValues(lowSeriesID)
	var (
		rs []*collections.FloatArray
		ok bool
	)
	rs, ok = g.aggregators[key]
	if !ok {
		rs = make([]*collections.FloatArray, g.tableScan.numOfAggs)
		g.aggregators[key.Clone()] = rs
	}
	return rs
}

type groupingWithoutTags struct {
	aggregator []*collections.FloatArray // family time => streams of fields
}

func newGroupingWithoutTags(tableScan *TableScan) grouping {
	return &groupingWithoutTags{
		aggregator: make([]*collections.FloatArray, tableScan.numOfAggs),
	}
}

func (g *groupingWithoutTags) ForEach(fn func(tags *GroupingKey, rs []*collections.FloatArray)) {
	fn(nil, g.aggregator)
}

func (g *groupingWithoutTags) GetAggregator(_ uint16) []*collections.FloatArray {
	return g.aggregator
}

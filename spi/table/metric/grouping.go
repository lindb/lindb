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

	"github.com/lindb/common/models"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/storage/metric"
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
	db   *metric.Database
	tags tag.Metas

	// tag value ids for each grouping tag key
	groupingTagValueIDs []*roaring.Bitmap
	tagValuesMap        []map[uint32]string // tag value id=> tag value for each group by tag key
}

func NewGrouping(db *metric.Database, tags tag.Metas) *Grouping {
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
	GetAggregator(lowSeriesID uint16) []Result
	ForEach(fn func(tags *GroupingKey, rs []Result))
}

type groupingWithTags struct {
	aggregators map[*GroupingKey][]Result
	tagsScanner *TagsScanner
	tableScan   *TableScan
}

func newGroupingWithTags(tagsScanner *TagsScanner, tableScan *TableScan) grouping {
	return &groupingWithTags{
		tagsScanner: tagsScanner,
		tableScan:   tableScan,
		aggregators: make(map[*GroupingKey][]Result),
	}
}

func (g *groupingWithTags) ForEach(fn func(tags *GroupingKey, rs []Result)) {
	for tags, aggregator := range g.aggregators {
		fn(tags, aggregator)
	}
}

func (g *groupingWithTags) GetAggregator(lowSeriesID uint16) []Result {
	key := g.tagsScanner.FindTagValues(lowSeriesID)
	var (
		rs []Result
		ok bool
	)
	rs, ok = g.aggregators[key]
	if !ok {
		rs = make([]Result, g.tableScan.numOfAggs)
		g.aggregators[key.Clone()] = rs
	}
	return rs
}

type groupingWithoutTags struct {
	aggregator []Result
}

func newGroupingWithoutTags(tableScan *TableScan) grouping {
	return &groupingWithoutTags{
		aggregator: make([]Result, tableScan.numOfAggs),
	}
}

func (g *groupingWithoutTags) ForEach(fn func(tags *GroupingKey, rs []Result)) {
	fn(nil, g.aggregator)
}

func (g *groupingWithoutTags) GetAggregator(_ uint16) []Result {
	return g.aggregator
}

type Result any

type result[V float64 | *models.Exemplar] struct {
	array *Array[V]
}

func NewResult[V float64 | *models.Exemplar](numOfPoints int) Result {
	return &result[V]{
		array: NewArray[V](numOfPoints),
	}
}

const blockSize = 8

type Array[V float64 | *models.Exemplar] struct {
	marks    []uint8
	values   []V
	capacity int
	size     int
	isSingle bool
}

func NewArray[V float64 | *models.Exemplar](capacity int) *Array[V] {
	markLen := capacity / blockSize
	if capacity%blockSize > 0 {
		markLen++
	}
	return &Array[V]{
		capacity: capacity,
		values:   make([]V, capacity),
		marks:    make([]uint8, markLen),
	}
}

// Values returns the values of array.
func (f *Array[V]) Values() []V {
	return f.values
}

// HasValue returns if has value with pos
func (f *Array[V]) HasValue(pos int) bool {
	if !f.checkPos(pos) {
		return false
	}
	blockIdx := pos / blockSize
	idx := pos % blockSize
	mark := f.marks[blockIdx]
	return mark&(1<<uint64(idx)) != 0
}

// GetValue returns value with pos, if it has not value return 0
func (f *Array[V]) GetValue(pos int) V {
	return f.values[pos]
}

// SetValue sets value with pos, if pos out of bounds, return it
func (f *Array[V]) SetValue(pos int, value V) {
	if !f.checkPos(pos) {
		return
	}
	f.values[pos] = value

	if !f.HasValue(pos) {
		blockIdx := pos / blockSize
		idx := pos - pos/blockSize*blockSize
		mark := f.marks[blockIdx]
		mark |= 1 << uint64(idx)
		f.marks[blockIdx] = mark

		f.size++
	}
}

func (f *Array[V]) Reset() {
	f.size = 0
	f.isSingle = false
	for i := range f.marks {
		f.marks[i] = 0
	}
}

func (f *Array[V]) checkPos(pos int) bool {
	if pos < 0 || pos >= f.capacity {
		return false
	}
	return true
}

type Stream[V any] interface {
	SetAtStep(step int, value V, fn func(a, b V) V)
	GetAtStep(step int) V
	Reset()
}

type stream[V float64 | *models.Exemplar] struct {
	values *Array[V]
}

func NewStream[V float64 | *models.Exemplar](size int) Stream[V] {
	return &stream[V]{
		values: NewArray[V](size),
	}
}

func (s *stream[V]) SetAtStep(step int, value V, fn func(a, b V) V) {
	if s.values.HasValue(step) {
		s.values.SetValue(step, fn(s.values.GetValue(step), value))
		return
	}
	s.values.SetValue(step, value)
}

func (s *stream[V]) GetAtStep(step int) V {
	return s.values.GetValue(step)
}

func (s *stream[V]) Reset() {
	s.values.Reset()
}

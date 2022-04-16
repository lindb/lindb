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

package indexdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb/tblstore/tagindex"
)

//go:generate mockgen -source ./tag_index.go -destination=./tag_index_mock.go -package=indexdb

// TagIndex represents the tag inverted index
type TagIndex interface {
	// GetGroupingScanner returns the grouping scanners based on series ids
	GetGroupingScanner(seriesIDs *roaring.Bitmap) ([]flow.GroupingScanner, error)
	// buildInvertedIndex builds inverted index for tag value id
	buildInvertedIndex(tagValueID uint32, seriesID uint32)
	// getSeriesIDsByTagValueIDs returns series ids by tag value ids
	getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) *roaring.Bitmap
	// getValues returns the all tag values and series ids
	getValues() *InvertedStore
	// getAllSeriesIDs returns all series ids
	getAllSeriesIDs() *roaring.Bitmap
	// flush flushes tag index under spec tag key,
	// write series ids of tag key level with constants.TagValueIDForTag
	flush(tagKeyID uint32, forward tagindex.ForwardFlusher, inverted tagindex.InvertedFlusher) error
}

// memGroupingScanner implements series.GroupingScanner for memory tag index
type memGroupingScanner struct {
	forward *ForwardStore // TODO add read lock
}

// GetSeriesAndTagValue returns group by container and tag value ids
func (g *memGroupingScanner) GetSeriesAndTagValue(highKey uint16) (lowSeriesIDs roaring.Container, tagValueIDs []uint32) {
	index := g.forward.keys.GetContainerIndex(highKey)
	if index < 0 {
		// data not found
		return nil, nil
	}
	return g.forward.keys.GetContainerAtIndex(index), g.forward.values[index]
}

// GetSeriesIDs returns the series ids in current memory scanner.
func (g *memGroupingScanner) GetSeriesIDs() *roaring.Bitmap {
	return g.forward.keys
}

// tagIndex is a inverted mapping relation of tag-value and seriesID group.
type tagIndex struct {
	forward  *ForwardStore  // store forward index, series id=>tag value id, maybe have same tag value id
	inverted *InvertedStore // store all tag value id=>series ids of tag level
}

// newTagKVEntrySet returns a new tagKVEntrySet
func newTagIndex() TagIndex {
	return &tagIndex{
		inverted: NewInvertedStore(),
		forward:  NewForwardStore(),
	}
}

// GetGroupingScanner returns the grouping scanners based on series ids
func (index *tagIndex) GetGroupingScanner(seriesIDs *roaring.Bitmap) ([]flow.GroupingScanner, error) {
	// check reader if it has series ids(after filtering)
	finalSeriesIDs := roaring.FastAnd(seriesIDs, index.forward.Keys())
	if finalSeriesIDs.IsEmpty() {
		// not found
		return nil, nil
	}
	// TODO add lock
	return []flow.GroupingScanner{&memGroupingScanner{forward: index.forward}}, nil
}

// buildInvertedIndex builds inverted index for tag value id
func (index *tagIndex) buildInvertedIndex(tagValueID, seriesID uint32) {
	seriesIDs, ok := index.inverted.Get(tagValueID)
	if !ok {
		// create new series ids for new tag value
		seriesIDs = roaring.NewBitmap()
		index.inverted.Put(tagValueID, seriesIDs)
	}
	seriesIDs.Add(seriesID)

	// build forward index, because series id is a unique id, so just put into forward index
	index.forward.Put(seriesID, tagValueID)
}

// getSeriesIDsByTagValueIDs returns series ids by tag value ids
func (index *tagIndex) getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) *roaring.Bitmap {
	result := roaring.New()
	values := index.inverted.Values()
	keys := index.inverted.Keys()
	// get final tag value ids need to load
	finalTagValueIDs := roaring.And(tagValueIDs, keys)
	highKeys := finalTagValueIDs.GetHighKeys()
	for idx, highKey := range highKeys {
		loadLowContainer := finalTagValueIDs.GetContainerAtIndex(idx)
		lowContainerIdx := keys.GetContainerIndex(highKey)
		lowContainer := keys.GetContainerAtIndex(lowContainerIdx)
		it := loadLowContainer.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// get the index of low tag value id in container
			lowIdx := lowContainer.Rank(lowTagValueID)
			result.Or(values[lowContainerIdx][lowIdx-1])
		}
	}
	return result
}

// getAllSeriesIDs returns all series ids
func (index *tagIndex) getAllSeriesIDs() *roaring.Bitmap {
	return index.forward.keys.Clone()
}

// getValues returns the all tag values and series ids
func (index *tagIndex) getValues() *InvertedStore {
	return index.inverted
}

// flush flushes tag index under spec tag key,
// write series ids of tag key level with constants.TagValueIDForTag
func (index *tagIndex) flush(
	tagKeyID uint32,
	forward tagindex.ForwardFlusher,
	inverted tagindex.InvertedFlusher,
) error {
	forward.PrepareTagKey(tagKeyID)
	inverted.PrepareTagKey(tagKeyID)
	for _, tagValueIDs := range index.forward.values {
		if err := forward.FlushForwardIndex(tagValueIDs); err != nil {
			return err
		}
	}
	if err := forward.CommitTagKey(index.forward.keys); err != nil {
		return err
	}
	// write each tag value series ids
	if err := index.inverted.WalkEntry(inverted.FlushInvertedIndex); err != nil {
		return err
	}
	if err := inverted.CommitTagKey(); err != nil {
		return err
	}
	return nil
}

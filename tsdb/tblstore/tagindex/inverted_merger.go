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

package tagindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

var SeriesInvertedMerger kv.MergerType = "SeriesInvertedMerger"

// init registers series inverted merger create function
func init() {
	kv.RegisterMerger(SeriesInvertedMerger, NewInvertedMerger)
}

// invertedMerger implements kv.Merger for merging inverted index data for each tag key
type invertedMerger struct {
	invertedFlusher InvertedFlusher
	kvFlusher       kv.Flusher
}

// NewInvertedMerger creates a inverted merger
func NewInvertedMerger(flusher kv.Flusher) (kv.Merger, error) {
	iFlusher, err := NewInvertedFlusher(flusher)
	if err != nil {
		return nil, err
	}
	return &invertedMerger{
		kvFlusher:       flusher,
		invertedFlusher: iFlusher,
	}, nil
}

func (m *invertedMerger) Init(_ map[string]interface{}) {}

// Merge merges the multi inverted index data into a inverted index for same tag key id
func (m *invertedMerger) Merge(key uint32, values [][]byte) error {
	var scanners []*tagInvertedScanner
	targetTagValueIDs := roaring.New() // target merged tag value ids
	// 1. prepare tag inverted scanner
	for _, value := range values {
		reader, err := newTagInvertedReader(value)
		if err != nil {
			return err
		}
		targetTagValueIDs.Or(reader.keys)
		newScanner, err := newTagInvertedScanner(reader)
		if err != nil {
			return err
		}
		scanners = append(scanners, newScanner)
	}

	m.invertedFlusher.PrepareTagKey(key)
	// 2. merge inverted index by roaring container
	highKeys := targetTagValueIDs.GetHighKeys()
	seriesIDs := roaring.New()
	for idx, highKey := range highKeys {
		container := targetTagValueIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// scan index data then merge series ids
			for _, scanner := range scanners {
				if err := scanner.scan(highKey, lowTagValueID, seriesIDs); err != nil {
					return err
				}
			}

			hk := uint32(highKey) << 16
			// flush tag value id=>series ids mapping
			if err := m.invertedFlusher.
				FlushInvertedIndex(encoding.ValueWithHighLowBits(hk, lowTagValueID), seriesIDs); err != nil {
				return err
			}
			seriesIDs.Clear() // clear target series ids
		}
	}
	return m.invertedFlusher.CommitTagKey()
}

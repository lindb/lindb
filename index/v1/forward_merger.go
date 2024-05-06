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

package v1

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
)

// for testing
var (
	newForwardIndexFlusher = NewForwardIndexFlusher
	newTagForwardReader    = NewTagForwardReader
)

var ForwardIndexMerger kv.MergerType = "ForwardIndexMergerV1"

func init() {
	// register forward index merger
	kv.RegisterMerger(ForwardIndexMerger, NewForwardIndexMerger)
}

// forwardIndexMerger implements kv.Merger interface.
type forwardIndexMerger struct {
	flusher ForwardIndexFlusher

	seriesIDs   *roaring.Bitmap
	tagValueIDs []uint32
	scanners    []*tagForwardScanner
}

// NewForwardIndexMerger creates a forward index merger.
func NewForwardIndexMerger(kvFlusher kv.Flusher) (kv.Merger, error) {
	flusher, err := newForwardIndexFlusher(kvFlusher)
	if err != nil {
		return nil, err
	}
	return &forwardIndexMerger{
		flusher:   flusher,
		seriesIDs: roaring.New(),
	}, nil
}

func (m *forwardIndexMerger) Init(params map[string]interface{}) {}

// Merge merges series ids -> tag value ids.
func (m *forwardIndexMerger) Merge(tagKeyID uint32, values [][]byte) error {
	m.seriesIDs.Clear() // target merged series ids
	m.scanners = m.scanners[:0]
	// 1. prepare tag forward scanner
	for _, value := range values {
		reader, err := newTagForwardReader(value)
		if err != nil {
			return err
		}
		m.seriesIDs.Or(reader.GetSeriesIDs())
		m.scanners = append(m.scanners, newTagForwardScanner(reader))
	}

	// 2. merge forward index by roaring container
	highKeys := m.seriesIDs.GetHighKeys()
	// 3. write series ids
	m.flusher.Prepare(tagKeyID)
	if err := m.flusher.WriteSeriesIDs(m.seriesIDs); err != nil {
		return err
	}

	// 4. write tag value ids by bitmap low container
	for idx, highKey := range highKeys {
		m.tagValueIDs = m.tagValueIDs[:0]
		container := m.seriesIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		for it.HasNext() {
			lowSeriesID := it.Next()
			// scan index data then merge tag value ids, sort by series id
			for _, scanner := range m.scanners {
				m.tagValueIDs = scanner.scan(highKey, lowSeriesID, m.tagValueIDs)
			}
		}
		// flush tag value ids by one container
		if err := m.flusher.WriteTagValueIDs(m.tagValueIDs); err != nil {
			return err
		}
	}
	return m.flusher.Commit()
}

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
	"github.com/lindb/lindb/pkg/encoding"
)

// for testing
var (
	newInvertedIndexFlusher = NewInvertedIndexFlusher
	bitmapUnmarshal         = encoding.BitmapUnmarshal
)
var InvertedIndexMerger kv.MergerType = "InvertedIndexMergerV1"

func init() {
	// register inverted index merger
	kv.RegisterMerger(InvertedIndexMerger, NewInvertedIndexMerger)
}

// invertedIndexMerger implements kv.Merger interface.
type invertedIndexMerger struct {
	flusher         InvertedIndexFlusher
	targetSeriesIDs *roaring.Bitmap // target merged series ids
	seriesIDs       *roaring.Bitmap
}

// NewInvertedIndexMerger creates an inverted index merger.
func NewInvertedIndexMerger(kvFlusher kv.Flusher) (kv.Merger, error) {
	flusher, err := newInvertedIndexFlusher(kvFlusher)
	if err != nil {
		return nil, err
	}
	return &invertedIndexMerger{
		flusher:         flusher,
		targetSeriesIDs: roaring.New(),
		seriesIDs:       roaring.New(),
	}, nil
}

func (m *invertedIndexMerger) Init(params map[string]interface{}) {}

// Merge merges series ids for key.
func (m *invertedIndexMerger) Merge(key uint32, values [][]byte) error {
	m.targetSeriesIDs.Clear()
	for _, val := range values {
		if _, err := bitmapUnmarshal(m.seriesIDs, val); err != nil {
			return err
		}
		m.targetSeriesIDs.Or(m.seriesIDs)
	}
	m.flusher.Prepare(key)
	if err := m.flusher.Write(m.targetSeriesIDs); err != nil {
		return err
	}
	return m.flusher.Commit()
}

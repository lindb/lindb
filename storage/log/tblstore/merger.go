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

package tblstore

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

var LogIndexMerger kv.MergerType = "LogIndexMerger"

// init registers metric data merger create function
func init() {
	kv.RegisterMerger(LogIndexMerger, NewMerger)
}

// merger implements kv.Merger for merging series data for each metric
type merger struct {
	flusher kv.Flusher

	result, temp *roaring.Bitmap
}

// NewMerger creates a metric data merger
func NewMerger(flusher kv.Flusher) (kv.Merger, error) {
	return &merger{
		flusher: flusher,
		result:  roaring.New(),
		temp:    roaring.New(),
	}, nil
}

// Init implements kv.Merger.
func (m *merger) Init(params map[string]any) {
}

// Merge implements kv.Merger.
func (m *merger) Merge(key uint32, values [][]byte) error {
	m.result.Clear()

	for _, op := range values {
		m.temp.Clear()
		encoding.BitmapUnmarshal(m.temp, op)

		// merge
		m.result.Or(m.temp)
	}

	d, err := m.result.ToBytes()
	if err != nil {
		return err
	}
	return m.flusher.Add(key, d)
}

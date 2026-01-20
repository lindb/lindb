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

package memdb

import (
	"fmt"
	"sort"

	"github.com/lindb/common/models"
	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/storage/metric/tblstore/metricsdata"
)

type ExemplarPage interface {
	encoding.TSDValueGetter

	write(slot uint16, traceID, spanID []byte, duration int64) error
	flush(flusher metricsdata.Flusher) error
}

type exemplar struct {
	traceID  []byte
	spanID   []byte
	duration int64
}

type exemplarPage struct {
	store map[uint16]exemplar
}

func (e *exemplarPage) GetExemplar(slot uint16) (*models.Exemplar, bool) {
	value, ok := e.store[slot]
	if !ok {
		return nil, false
	}
	return &models.Exemplar{
		TraceID:  string(value.traceID),
		SpanID:   string(value.spanID),
		Duration: value.duration,
	}, true
}

func (e *exemplarPage) GetValue(slot uint16) (float64, bool) {
	return 0, false
}

func newExemplarPage() ExemplarPage {
	return &exemplarPage{
		store: make(map[uint16]exemplar),
	}
}

func (e *exemplarPage) write(slot uint16, traceID, spanID []byte, duration int64) error {
	fmt.Println("write exemplar")
	e.store[slot] = exemplar{
		traceID:  traceID,
		spanID:   spanID,
		duration: duration,
	}
	return nil
}

func (e *exemplarPage) flush(flusher metricsdata.Flusher) error {
	slots := lo.Keys(e.store)
	sort.Slice(slots, func(i, j int) bool {
		return slots[i] < slots[j]
	})
	w := stream.NewBufferWriter(nil)
	w.PutUInt16(uint16(len(e.store)))
	for _, slot := range slots {
		ex := e.store[slot]
		w.PutUInt16(uint16(slot))
		w.PutUvarint32(uint32(len(ex.traceID)))
		w.PutBytes(ex.traceID)
		w.PutUvarint32(uint32(len(ex.spanID)))
		w.PutBytes(ex.spanID)
		w.PutVarint64(ex.duration)
	}

	data, err := w.Bytes()
	if err != nil {
		return err
	}
	fmt.Println("flush exemplar", len(data))
	return flusher.FlushField(data)
}

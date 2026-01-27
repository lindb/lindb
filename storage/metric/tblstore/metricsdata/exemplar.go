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

package metricsdata

import (
	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/strutil"
)

type exemplarTSDGetter struct {
	data map[uint16]*models.Exemplar
}

func newExemplarTSDGetter(data []byte) encoding.TSDValueGetter {
	r := stream.NewReader(data)
	size := r.ReadUint16() // read length
	var exemplars map[uint16]*models.Exemplar

	if size > 0 {
		exemplars = make(map[uint16]*models.Exemplar)
		for i := 0; i < int(size); i++ {
			slot := r.ReadUint16()
			traceIDLen := r.ReadUvarint32()
			traceID := r.ReadBytes(int(traceIDLen))
			spanIDLen := r.ReadUvarint32()
			spanID := r.ReadBytes(int(spanIDLen))
			duration := r.ReadVarint64()

			exemplars[slot] = &models.Exemplar{
				TraceID:  strutil.ByteSlice2String(traceID),
				SpanID:   strutil.ByteSlice2String(spanID),
				Duration: duration,
			}
		}
	}

	return &exemplarTSDGetter{
		data: exemplars,
	}
}

func (e *exemplarTSDGetter) GetExemplar(slot uint16) (*models.Exemplar, bool) {
	if e.data == nil {
		return nil, false
	}
	rs, ok := e.data[slot]
	return rs, ok
}

// GetValue implements [encoding.TSDValueGetter].
func (e *exemplarTSDGetter) GetValue(slot uint16) (float64, bool) {
	return 0, false
}

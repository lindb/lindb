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

package ingest

import (
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"
	"github.com/lindb/common/series"
	"github.com/lindb/lindb/pkg/strutil"
)

// Field represents field interface.
type Field interface {
	// writer field data into broker row builder.
	write(builder *series.RowBuilder) error
}

// simpleField represents simple field like(sum/min/max/last etc.)
type simpleField struct {
	name      string                        // field name
	v         float64                       // field value
	fieldType flatMetricsV1.SimpleFieldType // field type
}

// writer field data into broker row builder(flat protocol)
func (s *simpleField) write(builder *series.RowBuilder) error {
	return builder.AddSimpleField(strutil.String2ByteSlice(s.name), s.fieldType, s.v)
}

// Sum represents sum field, implements Field interface.
type Sum struct {
	simpleField
}

// NewSum creates a Sum filed.
func NewSum(name string, v float64) Field {
	return &Sum{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeDeltaSum,
			name:      name,
			v:         v,
		},
	}
}

// Min represents min field, implements Field interface.
type Min struct {
	simpleField
}

// NewMin creates a Min field.
func NewMin(name string, v float64) Field {
	return &Min{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeMin,
			name:      name,
			v:         v,
		},
	}
}

// Max represents max field, implements Field interface.
type Max struct {
	simpleField
}

// NewMax creates a Max field.
func NewMax(name string, v float64) Field {
	return &Max{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeMax,
			name:      name,
			v:         v,
		},
	}
}

// First represents first field, implements Field interface.
type First struct {
	simpleField
}

// NewFirst creates a First field.
func NewFirst(name string, v float64) Field {
	return &First{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeFirst,
			name:      name,
			v:         v,
		},
	}
}

// Last represents last field, implements Field interface.
type Last struct {
	simpleField
}

// NewLast creates a Last field.
func NewLast(name string, v float64) Field {
	return &Last{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeLast,
			name:      name,
			v:         v,
		},
	}
}

// Histogram represents histogram field(compound field).
type Histogram struct {
	min, max, sum, count float64
	values, bounds       []float64
}

// NewHistogram creates a Histogram field.
func NewHistogram(min, max, sum, count float64, values, bounds []float64) Field {
	return &Histogram{
		min:    min,
		max:    max,
		sum:    sum,
		count:  count,
		values: values,
		bounds: bounds,
	}
}

// writer histogram data into broker row builder.
func (h *Histogram) write(builder *series.RowBuilder) error {
	if err := builder.AddCompoundFieldMMSC(h.min, h.max, h.sum, h.count); err != nil {
		return err
	}
	return builder.AddCompoundFieldData(h.values, h.bounds)
}


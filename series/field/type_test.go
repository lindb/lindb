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

package field

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cfield "github.com/lindb/common/field"
)

func TestType_String(t *testing.T) {
	assert.Equal(t, "sum", cfield.Sum.String())
	assert.Equal(t, "max", cfield.Max.String())
	assert.Equal(t, "min", cfield.Min.String())
	assert.Equal(t, "last", cfield.Last.String())
	assert.Equal(t, "first", cfield.First.String())
	assert.Equal(t, "histogram", cfield.Histogram.String())
	assert.Equal(t, "unknown", cfield.Unknown.String())
	assert.Equal(t, "name", Name("name").String())
}

// TestFieldType_NumericAlignment verifies that field.Type constants align with
// the Arrow IPC wire encoding for zero-cost direct casting.
func TestFieldType_NumericAlignment(t *testing.T) {
	assert.Equal(t, cfield.Type(1), cfield.Sum, "Sum must be 1")
	assert.Equal(t, cfield.Type(2), cfield.Min, "Min must be 2")
	assert.Equal(t, cfield.Type(3), cfield.Max, "Max must be 3")
	assert.Equal(t, cfield.Type(4), cfield.Last, "Last must be 4")
	assert.Equal(t, cfield.Type(5), cfield.First, "First must be 5 (aligns with Arrow wire encoding)")
	assert.Equal(t, cfield.Type(6), cfield.Histogram, "Histogram must be 6")
}

func TestType_Aggregate(t *testing.T) {
	assert.Equal(t, 100.0, cfield.Sum.Aggregate(1, 99.0))

	assert.Equal(t, 1.0, cfield.Min.Aggregate(1, 99.0))
	assert.Equal(t, 1.0, cfield.Min.Aggregate(99.0, 1))

	assert.Equal(t, 99.0, cfield.Max.Aggregate(1, 99.0))
	assert.Equal(t, 99.0, cfield.Max.Aggregate(99.0, 1))

	assert.Equal(t, 99.0, cfield.Last.Aggregate(1, 99.0))

	assert.Equal(t, 1.0, cfield.First.Aggregate(1, 99.0))

	// Histogram uses sum aggregation
	assert.Equal(t, 100.0, cfield.Histogram.Aggregate(1, 99.0))

	assert.Panics(t, func() {
		cfield.Exemplar.Aggregate(1, 99.0)
	})
}

func TestPanicAggregate(t *testing.T) {
	assert.Panics(t, func() {
		cfield.Type(99).Aggregate(1, 99.0)
	})
}

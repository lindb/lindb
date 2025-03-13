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
)

func TestType_String(t *testing.T) {
	assert.Equal(t, "sum", SumField.String())
	assert.Equal(t, "max", MaxField.String())
	assert.Equal(t, "min", MinField.String())
	assert.Equal(t, "last", LastField.String())
	assert.Equal(t, "first", FirstField.String())
	assert.Equal(t, "histogram", HistogramField.String())
	assert.Equal(t, "unknown", Unknown.String())
	assert.Equal(t, "name", Name("name").String())
}

func TestAggType_Aggregate(t *testing.T) {
	assert.Equal(t, 100.0, SumField.AggType().Aggregate(1, 99.0))

	assert.Equal(t, 1.0, MinField.AggType().Aggregate(1, 99.0))
	assert.Equal(t, 1.0, MinField.AggType().Aggregate(99.0, 1))

	assert.Equal(t, 99.0, MaxField.AggType().Aggregate(1, 99.0))
	assert.Equal(t, 99.0, MaxField.AggType().Aggregate(99.0, 1))

	assert.Equal(t, 99.0, LastField.AggType().Aggregate(1, 99.0))

	assert.Equal(t, 1.0, FirstField.AggType().Aggregate(1, 99.0))

	assert.Panics(t, func() {
		AggType(22).Aggregate(1, 2)
	})
}

func TestPanicAgg(t *testing.T) {
	assert.Panics(t, func() {
		Type(99).AggType().Aggregate(1, 99.0)
	})
}

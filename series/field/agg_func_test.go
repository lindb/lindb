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

func TestGetAggFunc(t *testing.T) {
	assert.NotNil(t, Sum.AggFunc())
	assert.NotNil(t, Min.AggFunc())
	assert.NotNil(t, Max.AggFunc())
	assert.NotNil(t, Count.AggFunc())
	assert.NotNil(t, LastValue.AggFunc())
	assert.Nil(t, AggType(99).AggFunc())
}

func TestSumAgg(t *testing.T) {
	agg := Sum.AggFunc()
	assert.Equal(t, Sum, agg.AggType())
	assert.Equal(t, 100.0, agg.Aggregate(1, 99.0))
}

func TestMinAgg(t *testing.T) {
	agg := Min.AggFunc()
	assert.Equal(t, Min, agg.AggType())
	assert.Equal(t, 1.0, agg.Aggregate(1, 99.0))
	assert.Equal(t, 1.0, agg.Aggregate(99.0, 1))
}

func TestMaxAgg(t *testing.T) {
	agg := Max.AggFunc()
	assert.Equal(t, Max, agg.AggType())
	assert.Equal(t, 99.0, agg.Aggregate(1, 99.0))
	assert.Equal(t, 99.0, agg.Aggregate(99.0, 1))
}

func TestReplaceAgg(t *testing.T) {
	agg := LastValue.AggFunc()
	assert.Equal(t, LastValue, agg.AggType())
	assert.Equal(t, 99.0, agg.Aggregate(1, 99.0))
}

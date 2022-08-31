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

package aggregation

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
)

func TestTopNOrderBy(t *testing.T) {
	data := []float64{20, 1, 23, 40, 3, 50, 10, 43, 1000, 50, 20}
	topNAsc := NewTopNOrderBy([]*OrderByItem{{
		Desc: false,
	}}, 5)
	for _, d := range data {
		r := newRow(d)
		topNAsc.Push(r)
	}
	rows := topNAsc.ResultSet()
	var rs []float64
	for _, r := range rows {
		rs = append(rs, r.GetValue("", function.Count))
	}
	sort.Float64s(rs)
	assert.Equal(t, []float64{1, 3, 10, 20, 20}, rs)
}

func TestOrderByRow(t *testing.T) {
	t.Run("no data", func(t *testing.T) {
		mockFields := map[string]*collections.FloatArray{
			"f1": collections.NewFloatArray(10),
		}
		row := NewOrderByRow("tags", mockFields)
		tags, fields := row.ResultSet()
		assert.Equal(t, "tags", tags)
		assert.Equal(t, mockFields, fields)
		assert.Zero(t, row.GetValue("f1", function.Min))
		assert.Zero(t, row.GetValue("f1", function.Avg))
		assert.Zero(t, row.GetValue("f1", function.Stddev))
		assert.Zero(t, row.GetValue("f2", function.Min))

		mockFields = map[string]*collections.FloatArray{
			"f1": nil,
		}
		row = NewOrderByRow("tags", mockFields)
		assert.Zero(t, row.GetValue("f1", function.Min))
	})

	t.Run("has data", func(t *testing.T) {
		values := collections.NewFloatArray(3)
		values.SetValue(0, 2.0)
		values.SetValue(1, 3.0)
		values.SetValue(2, 1.0)
		mockFields := map[string]*collections.FloatArray{
			"f1": values,
		}
		row := NewOrderByRow("tags", mockFields)
		tags, fields := row.ResultSet()
		assert.Equal(t, "tags", tags)
		assert.Equal(t, mockFields, fields)
		assert.Equal(t, 1.0, row.GetValue("f1", function.Min))
		assert.Equal(t, 1.0, row.GetValue("f1", function.Min))
		assert.Equal(t, 3.0, row.GetValue("f1", function.Count))
		assert.Equal(t, 3.0, row.GetValue("f1", function.Max))
		assert.Equal(t, 6.0, row.GetValue("f1", function.Sum))
		assert.Equal(t, 6.0/3.0, row.GetValue("f1", function.Avg))
		assert.Equal(t, 2.0, row.GetValue("f1", function.First))
		assert.Equal(t, 1.0, row.GetValue("f1", function.Last))
		assert.NotZero(t, row.GetValue("f1", function.Stddev))
		assert.Zero(t, row.GetValue("f1", function.Unknown))
	})
}

func TestResultLimiter(t *testing.T) {
	limiter := NewResultLimiter(2)
	r1 := newRow(1)
	r2 := newRow(2)
	r3 := newRow(3)
	limiter.Push(r1)
	limiter.Push(r2)
	limiter.Push(r3)
	assert.Equal(t, []Row{r1, r2}, limiter.ResultSet())
}

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

type row struct {
	v float64
}

func newRow(v float64) Row {
	return &row{v: v}
}

func (*row) ResultSet() (tags string, fields map[string]*collections.FloatArray) {
	return "", nil
}

func (r *row) GetValue(_ string, _ function.FuncType) float64 {
	return r.v
}

func TestTopN(t *testing.T) {
	data := []float64{20, 1, 23, 40, 3, 50, 10, 43, 1000, 50, 20}
	topNAsc := newTopNHeap([]*OrderByItem{{
		Desc: false,
	}}, 5)
	topNDesc := newTopNHeap([]*OrderByItem{{Desc: true}}, 5)
	for _, d := range data {
		r := newRow(d)
		topNAsc.Add(r)
		topNDesc.Add(r)
	}
	rows := topNAsc.ResultSet()
	var rs []float64
	for _, r := range rows {
		rs = append(rs, r.GetValue("", function.Count))
	}
	sort.Float64s(rs)
	assert.Equal(t, []float64{1, 3, 10, 20, 20}, rs)

	rows = topNDesc.ResultSet()
	rs = []float64{}
	for _, r := range rows {
		rs = append(rs, r.GetValue("", function.Count))
	}
	sort.Float64s(rs)
	assert.Equal(t, []float64{40, 43, 50, 50, 1000}, rs)

	assert.Nil(t, topNAsc.Pop())
}

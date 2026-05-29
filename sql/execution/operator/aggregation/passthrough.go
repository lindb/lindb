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
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
)

// passThroughAgg copies the value at colIdx from the input batch as-is.
// Used for standard aggregations (sum, count, min, max, first, last) whose values
// are already computed by storage nodes; the broker only needs to relay them.
type passThroughAgg struct{ colIdx int }

func (p *passThroughAgg) fill(b array.Builder, batch arrow.RecordBatch, numRows int) {
	col := batch.Column(p.colIdx)
	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		appendColumnValue(b, col, rowIdx)
	}
}

// nullColumnAgg appends null for every row.  Used as a safe fallback when
// the expected input column is absent from the batch.
type nullColumnAgg struct{}

func (n *nullColumnAgg) fill(b array.Builder, _ arrow.RecordBatch, numRows int) {
	for i := 0; i < numRows; i++ {
		b.AppendNull()
	}
}

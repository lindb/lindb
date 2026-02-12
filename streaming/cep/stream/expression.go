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

package stream

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/tree"
)

var empty = roaring.New()

type column struct {
	name  string
	index int
}

type Expr interface {
	Eval(record arrow.RecordBatch) (*roaring.Bitmap, error)
}

type ComparisonExpr struct {
	column column
	value  string
	// TODO: add operator
}

func (e *ComparisonExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	col := record.Column(e.column.index)
	values, ok := col.(*array.String)
	if !ok {
		return nil, fmt.Errorf("column %s is not string type", e.column.name)
	}
	result := roaring.New()
	for i := 0; i < col.Len(); i++ {
		if values.Value(i) == e.value {
			result.Add(uint32(i))
		}
	}
	return result, nil
}

type InExpr struct {
	column column
	values []string
}

func (e *InExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	col := record.Column(e.column.index)
	values, ok := col.(*array.String)
	if !ok {
		return nil, fmt.Errorf("column %s is not string type", e.column.name)
	}
	result := roaring.New()
	for i := 0; i < col.Len(); i++ {
		if lo.Contains(e.values, values.Value(i)) {
			result.Add(uint32(i))
		}
	}
	return result, nil
}

type NotExpr struct {
	expr Expr
}

type LogicalExpr struct {
	op    tree.LogicalOperator
	exprs []Expr
}

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

package expression

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/tree"
)

// Logical evaluates AND / OR across two or more boolean sub-expressions
// element-wise, producing a single boolean Arrow array per RecordBatch.
type Logical struct {
	operator tree.LogicalOperator
	terms    []Expression
}

// NewLogical constructs a Logical expression.
func NewLogical(_ EvalContext, op tree.LogicalOperator, terms []Expression) Expression {
	return &Logical{operator: op, terms: terms}
}

func (l *Logical) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("logical: scalar evaluation not supported")
}

func (l *Logical) ResultType() ResultType {
	return Array
}

func (l *Logical) String() string {
	parts := make([]string, len(l.terms))
	for i, t := range l.terms {
		parts[i] = t.String()
	}
	sep := fmt.Sprintf(" %s ", l.operator)
	return "(" + strings.Join(parts, sep) + ")"
}

// Eval applies the logical operator across all term arrays element-wise.
// Each term must return a *array.Boolean (same length as the batch).
func (l *Logical) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	n := int(record.NumRows())
	result := make([]bool, n)

	// Short-circuit initialisation: AND starts all-true, OR starts all-false.
	init := l.operator == tree.LogicalAND
	for i := range n {
		result[i] = init
	}

	for _, term := range l.terms {
		arr, err := term.Eval(record)
		if err != nil {
			return nil, fmt.Errorf("logical %s: term eval: %w", l.operator, err)
		}
		boolArr, ok := arr.(*array.Boolean)
		if !ok {
			return nil, fmt.Errorf("logical %s: term returned %T, expected *array.Boolean", l.operator, arr)
		}
		for i := range n {
			v := boolArr.Value(i)
			if l.operator == tree.LogicalAND {
				result[i] = result[i] && v
			} else {
				result[i] = result[i] || v
			}
		}
	}

	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.Reserve(n)
	for _, v := range result {
		b.Append(v)
	}
	return b.NewArray(), nil
}

// ─────────────────────────────────────────────────────────────────────────────

// Not negates a single boolean sub-expression element-wise.
type Not struct {
	value Expression
}

// NewNot constructs a Not expression.
func NewNot(_ EvalContext, value Expression) Expression {
	return &Not{value: value}
}

func (n *Not) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("not: scalar evaluation not supported")
}

func (n *Not) ResultType() ResultType {
	return Array
}

func (n *Not) String() string {
	return fmt.Sprintf("NOT(%s)", n.value)
}

// Eval inverts every boolean entry produced by the inner expression.
func (n *Not) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	arr, err := n.value.Eval(record)
	if err != nil {
		return nil, fmt.Errorf("not: inner eval: %w", err)
	}
	boolArr, ok := arr.(*array.Boolean)
	if !ok {
		return nil, fmt.Errorf("not: inner expression returned %T, expected *array.Boolean", arr)
	}

	rows := int(record.NumRows())
	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.Reserve(rows)
	for i := range rows {
		b.Append(!boolArr.Value(i))
	}
	return b.NewArray(), nil
}

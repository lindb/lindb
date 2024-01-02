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

package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/tree"
)

func TestCalc_CalcExpr(t *testing.T) {
	add := &tree.BinaryExpr{
		Left:     &tree.FieldExpr{Name: "x"},
		Right:    &tree.FieldExpr{Name: "y"},
		Operator: tree.ADD,
	}
	sub := &tree.BinaryExpr{
		Left:     &tree.FieldExpr{Name: "x"},
		Right:    &tree.FieldExpr{Name: "y"},
		Operator: tree.SUB,
	}
	mul := &tree.BinaryExpr{
		Left:     &tree.FieldExpr{Name: "x"},
		Right:    &tree.FieldExpr{Name: "y"},
		Operator: tree.MUL,
	}
	div := &tree.BinaryExpr{
		Left:     &tree.FieldExpr{Name: "x"},
		Right:    &tree.FieldExpr{Name: "y"},
		Operator: tree.DIV,
	}

	x, y := 11.1, 22.2
	variables := map[string]float64{
		"x": x,
		"y": y,
	}
	for _, expr := range []*tree.BinaryExpr{add, sub, mul, div} {
		left, right := x, y
		var expected float64
		switch expr.Operator {
		case tree.ADD:
			expected = left + right
		case tree.SUB:
			expected = left - right
		case tree.MUL:
			expected = left * right
		case tree.DIV:
			expected = left / right
		}
		calc := NewCalc(expr)
		result, err := calc.CalcExpr(variables)
		assert.Nil(t, err)
		assert.IsType(t, float64(1), result)
		assert.Equal(t, expected, result.(float64))
	}

	errDiv := &tree.BinaryExpr{
		Left:     &tree.FieldExpr{Name: "x"},
		Right:    &tree.FieldExpr{Name: "y"},
		Operator: tree.DIV,
	}
	calc := NewCalc(errDiv)
	result, err := calc.CalcExpr(map[string]float64{"x": 11.1, "y": 0})
	assert.NotNil(t, err)
	assert.Nil(t, result)

	errAdd := &tree.BinaryExpr{
		Left:     &tree.FieldExpr{Name: "x"},
		Right:    &tree.FieldExpr{Name: "y"},
		Operator: tree.ADD,
	}
	calc = NewCalc(errAdd)
	result, err = calc.CalcExpr(map[string]float64{"y": 1})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCalc_calcBinary(t *testing.T) {
	less := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 1},
		Right:    &tree.NumberLiteral{Val: 2},
		Operator: tree.LESS,
	}
	greater := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 2},
		Right:    &tree.NumberLiteral{Val: 1},
		Operator: tree.GREATER,
	}
	lessEqual := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 1},
		Right:    &tree.NumberLiteral{Val: 2},
		Operator: tree.LESSEQUAL,
	}
	greaterEqual := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 2},
		Right:    &tree.NumberLiteral{Val: 1},
		Operator: tree.GREATEREQUAL,
	}
	equal := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 2},
		Right:    &tree.NumberLiteral{Val: 1},
		Operator: tree.EQUAL,
	}
	notEqual := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 2},
		Right:    &tree.NumberLiteral{Val: 1},
		Operator: tree.NOTEQUAL,
	}

	notLess := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 10},
		Right:    &tree.NumberLiteral{Val: 2},
		Operator: tree.LESS,
	}
	notGreater := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 2},
		Right:    &tree.NumberLiteral{Val: 10},
		Operator: tree.GREATER,
	}

	and := &tree.BinaryExpr{
		Left:     less,
		Right:    greater,
		Operator: tree.AND,
	}
	and2 := &tree.BinaryExpr{
		Left:     lessEqual,
		Right:    greaterEqual,
		Operator: tree.AND,
	}
	or := &tree.BinaryExpr{
		Left:     equal,
		Right:    notEqual,
		Operator: tree.OR,
	}
	notOr := &tree.BinaryExpr{
		Left:     notLess,
		Right:    notGreater,
		Operator: tree.OR,
	}

	for _, expr := range []*tree.BinaryExpr{and, and2, or} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.True(t, r.(bool))
	}

	for _, expr := range []*tree.BinaryExpr{notOr} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.False(t, r.(bool))
	}

	ands := &tree.BinaryExpr{
		Left:     and,
		Right:    or,
		Operator: tree.AND,
	}
	ors := &tree.BinaryExpr{
		Left:     and,
		Right:    notOr,
		Operator: tree.OR,
	}
	p := &tree.BinaryExpr{
		Left:     &tree.ParenExpr{Expr: or},
		Right:    and,
		Operator: tree.AND,
	}

	for _, expr := range []*tree.BinaryExpr{ands, ors, p} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.True(t, r.(bool))
	}

	notAnds := &tree.BinaryExpr{
		Left:     and,
		Right:    notOr,
		Operator: tree.AND,
	}
	for _, expr := range []*tree.BinaryExpr{notAnds} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.False(t, r.(bool))
	}
}

func TestCalc_calcEquation(t *testing.T) {
	add := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 11.1},
		Right:    &tree.NumberLiteral{Val: 22.2},
		Operator: tree.ADD,
	}
	sub := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 11.1},
		Right:    &tree.NumberLiteral{Val: 22.2},
		Operator: tree.SUB,
	}
	mul := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 11.1},
		Right:    &tree.NumberLiteral{Val: 22.2},
		Operator: tree.MUL,
	}
	div := &tree.BinaryExpr{
		Left:     &tree.NumberLiteral{Val: 11.1},
		Right:    &tree.NumberLiteral{Val: 22.2},
		Operator: tree.DIV,
	}
	for _, expr := range []*tree.BinaryExpr{add, sub, mul, div} {
		left, right := expr.Left.(*tree.NumberLiteral).Val, expr.Right.(*tree.NumberLiteral).Val
		var expected float64
		switch expr.Operator {
		case tree.ADD:
			expected = left + right
		case tree.SUB:
			expected = left - right
		case tree.MUL:
			expected = left * right
		case tree.DIV:
			expected = left / right
		}
		calc := NewCalc(expr)
		result, err := calc.CalcExpr(nil)
		assert.Nil(t, err)
		assert.IsType(t, float64(1), result)
		assert.Equal(t, expected, result.(float64))
	}
}

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

	"github.com/lindb/lindb/sql/stmt"
)

func TestCalc_CalcExpr(t *testing.T) {
	add := &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "x"},
		Right:    &stmt.FieldExpr{Name: "y"},
		Operator: stmt.ADD,
	}
	sub := &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "x"},
		Right:    &stmt.FieldExpr{Name: "y"},
		Operator: stmt.SUB,
	}
	mul := &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "x"},
		Right:    &stmt.FieldExpr{Name: "y"},
		Operator: stmt.MUL,
	}
	div := &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "x"},
		Right:    &stmt.FieldExpr{Name: "y"},
		Operator: stmt.DIV,
	}

	x, y := 11.1, 22.2
	variables := map[string]float64{
		"x": x,
		"y": y,
	}
	for _, expr := range []*stmt.BinaryExpr{add, sub, mul, div} {
		left, right := x, y
		var expected float64
		switch expr.Operator {
		case stmt.ADD:
			expected = left + right
		case stmt.SUB:
			expected = left - right
		case stmt.MUL:
			expected = left * right
		case stmt.DIV:
			expected = left / right
		}
		calc := NewCalc(expr)
		result, err := calc.CalcExpr(variables)
		assert.Nil(t, err)
		assert.IsType(t, float64(1), result)
		assert.Equal(t, expected, result.(float64))
	}

	errDiv := &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "x"},
		Right:    &stmt.FieldExpr{Name: "y"},
		Operator: stmt.DIV,
	}
	calc := NewCalc(errDiv)
	result, err := calc.CalcExpr(map[string]float64{"x": 11.1, "y": 0})
	assert.NotNil(t, err)
	assert.Nil(t, result)

	errAdd := &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "x"},
		Right:    &stmt.FieldExpr{Name: "y"},
		Operator: stmt.ADD,
	}
	calc = NewCalc(errAdd)
	result, err = calc.CalcExpr(map[string]float64{"y": 1})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCalc_calcBinary(t *testing.T) {
	less := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 1},
		Right:    &stmt.NumberLiteral{Val: 2},
		Operator: stmt.LESS,
	}
	greater := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 2},
		Right:    &stmt.NumberLiteral{Val: 1},
		Operator: stmt.GREATER,
	}
	lessEqual := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 1},
		Right:    &stmt.NumberLiteral{Val: 2},
		Operator: stmt.LESSEQUAL,
	}
	greaterEqual := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 2},
		Right:    &stmt.NumberLiteral{Val: 1},
		Operator: stmt.GREATEREQUAL,
	}
	equal := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 2},
		Right:    &stmt.NumberLiteral{Val: 1},
		Operator: stmt.EQUAL,
	}
	notEqual := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 2},
		Right:    &stmt.NumberLiteral{Val: 1},
		Operator: stmt.NOTEQUAL,
	}

	notLess := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 10},
		Right:    &stmt.NumberLiteral{Val: 2},
		Operator: stmt.LESS,
	}
	notGreater := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 2},
		Right:    &stmt.NumberLiteral{Val: 10},
		Operator: stmt.GREATER,
	}

	and := &stmt.BinaryExpr{
		Left:     less,
		Right:    greater,
		Operator: stmt.AND,
	}
	and2 := &stmt.BinaryExpr{
		Left:     lessEqual,
		Right:    greaterEqual,
		Operator: stmt.AND,
	}
	or := &stmt.BinaryExpr{
		Left:     equal,
		Right:    notEqual,
		Operator: stmt.OR,
	}
	notOr := &stmt.BinaryExpr{
		Left:     notLess,
		Right:    notGreater,
		Operator: stmt.OR,
	}

	for _, expr := range []*stmt.BinaryExpr{and, and2, or} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.True(t, r.(bool))
	}

	for _, expr := range []*stmt.BinaryExpr{notOr} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.False(t, r.(bool))
	}

	ands := &stmt.BinaryExpr{
		Left:     and,
		Right:    or,
		Operator: stmt.AND,
	}
	ors := &stmt.BinaryExpr{
		Left:     and,
		Right:    notOr,
		Operator: stmt.OR,
	}
	p := &stmt.BinaryExpr{
		Left:     &stmt.ParenExpr{Expr: or},
		Right:    and,
		Operator: stmt.AND,
	}

	for _, expr := range []*stmt.BinaryExpr{ands, ors, p} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.True(t, r.(bool))
	}

	notAnds := &stmt.BinaryExpr{
		Left:     and,
		Right:    notOr,
		Operator: stmt.AND,
	}
	for _, expr := range []*stmt.BinaryExpr{notAnds} {
		calc := NewCalc(nil)
		r, err := calc.calcBinary(expr, nil)
		assert.Nil(t, err)
		assert.IsType(t, true, r)
		assert.False(t, r.(bool))
	}
}

func TestCalc_calcEquation(t *testing.T) {
	add := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 11.1},
		Right:    &stmt.NumberLiteral{Val: 22.2},
		Operator: stmt.ADD,
	}
	sub := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 11.1},
		Right:    &stmt.NumberLiteral{Val: 22.2},
		Operator: stmt.SUB,
	}
	mul := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 11.1},
		Right:    &stmt.NumberLiteral{Val: 22.2},
		Operator: stmt.MUL,
	}
	div := &stmt.BinaryExpr{
		Left:     &stmt.NumberLiteral{Val: 11.1},
		Right:    &stmt.NumberLiteral{Val: 22.2},
		Operator: stmt.DIV,
	}
	for _, expr := range []*stmt.BinaryExpr{add, sub, mul, div} {
		left, right := expr.Left.(*stmt.NumberLiteral).Val, expr.Right.(*stmt.NumberLiteral).Val
		var expected float64
		switch expr.Operator {
		case stmt.ADD:
			expected = left + right
		case stmt.SUB:
			expected = left - right
		case stmt.MUL:
			expected = left * right
		case stmt.DIV:
			expected = left / right
		}
		calc := NewCalc(expr)
		result, err := calc.CalcExpr(nil)
		assert.Nil(t, err)
		assert.IsType(t, float64(1), result)
		assert.Equal(t, expected, result.(float64))
	}
}

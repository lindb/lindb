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
	"errors"
	"fmt"
	"reflect"

	"github.com/lindb/lindb/sql/tree"
)

// Calc represents Expr calculator
type Calc struct {
	expr tree.Expr
}

func NewCalc(expr tree.Expr) *Calc {
	return &Calc{expr: expr}
}

// calcBinary calculate a binary expr
func (q *Calc) calcBinary(expr *tree.BinaryExpr, variables map[string]float64) (result any, err error) {
	switch expr.Operator {
	case tree.AND, tree.OR:
		left, err0 := q.calcExpr(expr.Left, variables)
		if err0 != nil {
			return result, err0
		}
		right, err0 := q.calcExpr(expr.Right, variables)
		if err0 != nil {
			return result, err0
		}

		l, ok := left.(bool)
		if !ok {
			return result, errors.New("expected calcExpr returns bool type")
		}
		r, ok := right.(bool)
		if !ok {
			return result, errors.New("expected calcExpr returns bool type")
		}

		switch expr.Operator {
		case tree.AND:
			return l && r, nil
		case tree.OR:
			return l || r, nil
		}
	case tree.ADD, tree.SUB, tree.MUL, tree.DIV:
		left, err0 := q.calcEquation(expr.Left, variables)
		if err0 != nil {
			return result, err0
		}
		right, err0 := q.calcEquation(expr.Right, variables)
		if err0 != nil {
			return result, err0
		}

		var r float64
		switch expr.Operator {
		case tree.ADD:
			r = left + right
		case tree.SUB:
			r = left - right
		case tree.MUL:
			r = left * right
		case tree.DIV:
			if right == 0 {
				return result, errors.New("divisor cannot be zero")
			}
			r = left / right
		}
		return r, nil
	case tree.EQUAL, tree.NOTEQUAL, tree.GREATER, tree.GREATEREQUAL, tree.LESS, tree.LESSEQUAL, tree.LIKE:
		left, err0 := q.calcEquation(expr.Left, variables)
		if err0 != nil {
			return result, err0
		}
		right, err0 := q.calcEquation(expr.Right, variables)
		if err0 != nil {
			return result, err0
		}

		var r bool
		switch expr.Operator {
		case tree.EQUAL, tree.LIKE:
			r = left == right
		case tree.NOTEQUAL:
			r = left != right
		case tree.GREATER:
			r = left > right
		case tree.GREATEREQUAL:
			r = left >= right
		case tree.LESS:
			r = left < right
		case tree.LESSEQUAL:
			r = left <= right
		}
		return r, nil
	}

	return result, fmt.Errorf("calcBinary unknown operator %d", expr.Operator)
}

// calcEquation calculate an equation that may contain variables
func (q *Calc) calcEquation(expr tree.Expr, variables map[string]float64) (result float64, err error) {
	switch v := expr.(type) {
	case *tree.FieldExpr:
		if val, ok := variables[v.Name]; !ok {
			return result, fmt.Errorf("variable %s does not exist", v.Name)
		} else {
			return val, nil
		}
	case *tree.NumberLiteral:
		return v.Val, nil
	case *tree.ParenExpr:
		r, err0 := q.calcExpr(v.Expr, variables)
		if err0 != nil {
			return result, err0
		}
		if v, ok := r.(float64); !ok {
			return result, fmt.Errorf("expected float64 type got %v", reflect.TypeOf(r))
		} else {
			return v, nil
		}
	case *tree.BinaryExpr:
		r, err0 := q.calcBinary(v, variables)
		if err0 != nil {
			return result, err0
		}
		if v, ok := r.(float64); !ok {
			return result, fmt.Errorf("expected float64 type got %v", reflect.TypeOf(r))
		} else {
			return v, nil
		}
	default:
		return result, errors.New("calcEquation unknown type")
	}
}

// calcExpr calculate an expr that may be of type BinaryExpr or ParenExpr
func (q *Calc) calcExpr(expr tree.Expr, variables map[string]float64) (result any, err error) {
	switch v := expr.(type) {
	case *tree.BinaryExpr:
		return q.calcBinary(v, variables)
	case *tree.ParenExpr:
		return q.calcExpr(v.Expr, variables)
	default:
		return result, fmt.Errorf("unexpected type %v", reflect.TypeOf(expr))
	}
}

// CalcExpr calls calcExpr to compute an expr
func (q *Calc) CalcExpr(variables map[string]float64) (result any, err error) {
	r, err := q.calcExpr(q.expr, variables)
	if err != nil {
		return
	}
	return r, nil
}

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

package stmt

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/encoding"
)

//go:generate mockgen -source ./expr.go -destination=./expr_mock.go -package=stmt

// exprData represents inner wrapper of expr for json marshal
type exprData struct {
	Type string          `json:"type"`
	Expr json.RawMessage `json:"expr"`
}

// Expr represents a interface for all expression types
type Expr interface {
	// Rewrite rewrites the expr after parse
	Rewrite() string
}

// TagFilter represents tag filter for searching time series
type TagFilter interface {
	Expr
	// TagKey returns the filter's tag key
	TagKey() string
}

// SelectItem represents a select item from select statement
type SelectItem struct {
	Expr  Expr
	Alias string
}

// innerSelectItem represents inner wrapper of select item for json marshal
type innerSelectItem struct {
	exprData
	Alias string `json:"alias"`
}

// FieldExpr represents a field name for select list
type FieldExpr struct {
	Name string `json:"name"`
}

// NumberLiteral represents a number.
type NumberLiteral struct {
	Val float64 `json:"val"`
}

// CallExpr represents a function call expression
type CallExpr struct {
	FuncType function.FuncType
	Params   []Expr
}

// innerCallExpr represents inner wrapper of call expr for json marshal
type innerCallExpr struct {
	Type     string            `json:"type"`
	FuncType function.FuncType `json:"funcType"`
	Params   []json.RawMessage `json:"params"`
}

// ParenExpr represents a parenthesized expression
type ParenExpr struct {
	Expr Expr
}

// BinaryExpr represents an operations with two expressions
type BinaryExpr struct {
	Left, Right Expr
	Operator    BinaryOP
}

// innerBinaryExpr represents inner wrapper of binary expr for json marshal
type innerBinaryExpr struct {
	Type     string          `json:"type"`
	Left     json.RawMessage `json:"left"`
	Right    json.RawMessage `json:"right"`
	Operator BinaryOP        `json:"operator"`
}

// EqualsExpr represents an equals expression
type EqualsExpr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// InExpr represents an in expression
type InExpr struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// LikeExpr represents a like expression
type LikeExpr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RegexExpr represents a regular expression
type RegexExpr struct {
	Key    string `json:"key"`
	Regexp string `json:"regexp"`
}

// NotExpr represents a not expression
type NotExpr struct {
	Expr Expr
}

// Rewrite rewrites the select item expr after parse
func (e *SelectItem) Rewrite() string {
	if len(e.Alias) == 0 {
		return e.Expr.Rewrite()
	}
	return fmt.Sprintf("%s as %s", e.Expr.Rewrite(), e.Alias)
}

// Rewrite rewrites the field expr after parse
func (e *FieldExpr) Rewrite() string {
	return e.Name
}

// Rewrite rewrites the call expr after parse
func (e *CallExpr) Rewrite() string {
	var params []string
	for _, param := range e.Params {
		params = append(params, param.Rewrite())
	}
	return fmt.Sprintf("%s(%s)", e.FuncType, strings.Join(params, ","))
}

// Rewrite rewrites the paren expr after parse
func (e *ParenExpr) Rewrite() string {
	return fmt.Sprintf("(%s)", e.Expr.Rewrite())
}

// Rewrite rewrites the number literal after parse
func (e *NumberLiteral) Rewrite() string {
	return fmt.Sprintf("%.2f", e.Val)
}

// Rewrite rewrites the binary expr after parse
func (e *BinaryExpr) Rewrite() string {
	return fmt.Sprintf("%s%s%s", e.Left.Rewrite(), BinaryOPString(e.Operator), e.Right.Rewrite())
}

// Rewrite rewrites the not expr after parse
func (e *NotExpr) Rewrite() string {
	return fmt.Sprintf("not %s", e.Expr.Rewrite())
}

// Rewrite rewrites the equals expr after parse
func (e *EqualsExpr) Rewrite() string {
	return fmt.Sprintf("%s=%s", e.Key, e.Value)
}

// Rewrite rewrites the in expr after parse
func (e *InExpr) Rewrite() string {
	return fmt.Sprintf("%s in (%s)", e.Key, strings.Join(e.Values, ","))
}

// Rewrite rewrites the like expr after parse
func (e *LikeExpr) Rewrite() string {
	return fmt.Sprintf("%s like %s", e.Key, e.Value)
}

// Rewrite rewrites the regex expr after parse
func (e *RegexExpr) Rewrite() string {
	return fmt.Sprintf("%s=~%s", e.Key, e.Regexp)
}

// Marshal returns json of expr using custom json marshal
func Marshal(expr Expr) []byte {
	switch e := expr.(type) {
	case *RegexExpr:
		return encoding.JSONMarshal(&exprData{Type: "regex", Expr: encoding.JSONMarshal(expr)})
	case *LikeExpr:
		return encoding.JSONMarshal(&exprData{Type: "like", Expr: encoding.JSONMarshal(expr)})
	case *InExpr:
		return encoding.JSONMarshal(&exprData{Type: "in", Expr: encoding.JSONMarshal(expr)})
	case *EqualsExpr:
		return encoding.JSONMarshal(&exprData{Type: "equals", Expr: encoding.JSONMarshal(expr)})
	case *NumberLiteral:
		return encoding.JSONMarshal(&exprData{Type: "number", Expr: encoding.JSONMarshal(expr)})
	case *FieldExpr:
		return encoding.JSONMarshal(&exprData{Type: "field", Expr: encoding.JSONMarshal(expr)})
	case *NotExpr:
		return encoding.JSONMarshal(&exprData{Type: "not", Expr: Marshal(e.Expr)})
	case *ParenExpr:
		return encoding.JSONMarshal(&exprData{Type: "paren", Expr: Marshal(e.Expr)})
	case *SelectItem:
		inner := innerSelectItem{
			exprData: exprData{
				Type: "selectItem",
				Expr: Marshal(e.Expr),
			},
			Alias: e.Alias,
		}
		return encoding.JSONMarshal(&inner)
	case *CallExpr:
		inner := innerCallExpr{
			Type:     "call",
			FuncType: e.FuncType,
		}
		for _, param := range e.Params {
			inner.Params = append(inner.Params, Marshal(param))
		}
		return encoding.JSONMarshal(&inner)
	case *BinaryExpr:
		inner := innerBinaryExpr{
			Type:     "binary",
			Left:     Marshal(e.Left),
			Right:    Marshal(e.Right),
			Operator: e.Operator,
		}
		return encoding.JSONMarshal(&inner)
	default:
		return nil
	}
}

// Unmarshal parses value to expr
func Unmarshal(value []byte) (Expr, error) {
	var exprData exprData
	err := json.Unmarshal(value, &exprData)
	if err != nil {
		return nil, err
	}
	switch exprData.Type {
	case "regex":
		return unmarshal(&exprData, &RegexExpr{})
	case "like":
		return unmarshal(&exprData, &LikeExpr{})
	case "in":
		return unmarshal(&exprData, &InExpr{})
	case "equals":
		return unmarshal(&exprData, &EqualsExpr{})
	case "number":
		return unmarshal(&exprData, &NumberLiteral{})
	case field:
		return unmarshal(&exprData, &FieldExpr{})
	case "paren":
		e, err := Unmarshal(exprData.Expr)
		if err != nil {
			return nil, err
		}
		return &ParenExpr{Expr: e}, nil
	case "binary":
		return unmarshalBinary(value)
	case "selectItem":
		return unmarshalSelectItem(value)
	case "call":
		return unmarshalCall(value)
	case "not":
		e, err := Unmarshal(exprData.Expr)
		if err != nil {
			return nil, err
		}
		return &NotExpr{Expr: e}, nil
	default:
		return nil, fmt.Errorf("expr type not match:%s", exprData.Type)
	}
}

// unmarshalCall parses value to call expr
func unmarshalCall(value []byte) (Expr, error) {
	innerExpr := innerCallExpr{}
	err := encoding.JSONUnmarshal(value, &innerExpr)
	if err != nil {
		return nil, err
	}
	expr := &CallExpr{FuncType: innerExpr.FuncType}
	for _, param := range innerExpr.Params {
		e, err := Unmarshal(param)
		if err != nil {
			return nil, err
		}
		expr.Params = append(expr.Params, e)
	}
	return expr, nil
}

// unmarshalSelectItem parses value to select item expr
func unmarshalSelectItem(value []byte) (Expr, error) {
	innerExpr := innerSelectItem{}
	err := encoding.JSONUnmarshal(value, &innerExpr)
	if err != nil {
		return nil, err
	}
	e, err := Unmarshal(innerExpr.Expr)
	if err != nil {
		return nil, err
	}
	return &SelectItem{Alias: innerExpr.Alias, Expr: e}, nil
}

// unmarshalBinary parses value to binary expr
func unmarshalBinary(value []byte) (Expr, error) {
	innerExpr := innerBinaryExpr{}
	err := encoding.JSONUnmarshal(value, &innerExpr)
	if err != nil {
		return nil, err
	}
	left, err := Unmarshal(innerExpr.Left)
	if err != nil {
		return nil, err
	}
	right, err := Unmarshal(innerExpr.Right)
	if err != nil {
		return nil, err
	}
	expr := &BinaryExpr{
		Left:     left,
		Right:    right,
		Operator: innerExpr.Operator,
	}
	return expr, nil
}

// unmarshal parses expr data to expr
func unmarshal(exprData *exprData, expr Expr) (Expr, error) {
	if err := encoding.JSONUnmarshal(exprData.Expr, expr); err != nil {
		return nil, err
	}
	return expr, nil
}

// TagKey returns the equals filter's tag key
func (e *EqualsExpr) TagKey() string { return e.Key }

// TagKey returns the in filter's tag key
func (e *InExpr) TagKey() string { return e.Key }

// TagKey returns the like filter's tag key
func (e *LikeExpr) TagKey() string { return e.Key }

// TagKey returns the regex filter's tag key
func (e *RegexExpr) TagKey() string { return e.Key }

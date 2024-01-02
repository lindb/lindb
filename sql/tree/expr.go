package tree

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/aggregation/function"
)

//go:generate mockgen -source ./expr.go -destination=./expr_mock.go -package=stmt

// OrderByExpr represents order by expr item.
type OrderByExpr struct {
	Expr Expr // support field name/function for select field item.
	Desc bool
}

// Rewrite rewrites the order by expr after parse
func (e *OrderByExpr) Rewrite() string {
	sort := "asc"
	if e.Desc {
		sort = "desc"
	}
	return fmt.Sprintf("%s %s", e.Expr.Rewrite(), sort)
}

// innerOrderByExpr represents inner wrapper of order by for json marshal.
type innerOrderByExpr struct {
	exprData
	Desc bool `json:"desc"`
}

// exprData represents inner wrapper of expr for json marshal.
type exprData struct {
	Type string          `json:"type"`
	Expr json.RawMessage `json:"expr"`
}

// Expr represents a interface for all expression types.
type Expr interface {
	// Rewrite rewrites the expr after parse
	Rewrite() string
}

// SelectItem2 represents a select item from select statement.
type SelectItem2 struct {
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
	Name  string `json:"name"`
	Value string `json:"value"`
}

// InExpr represents an in expression
type InExpr struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// LikeExpr represents a like expression
type LikeExpr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RegexExpr represents a regular expression
type RegexExpr struct {
	Name   string `json:"name"`
	Regexp string `json:"regexp"`
}

// NotExpr represents a not expression
type NotExpr struct {
	Expr Expr
}

// Rewrite rewrites the select item expr after parse
func (e *SelectItem2) Rewrite() string {
	if e.Alias == "" {
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
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

// Rewrite rewrites the in expr after parse
func (e *InExpr) Rewrite() string {
	return fmt.Sprintf("%s in (%s)", e.Name, strings.Join(e.Values, ","))
}

// Rewrite rewrites the like expr after parse
func (e *LikeExpr) Rewrite() string {
	return fmt.Sprintf("%s like %s", e.Name, e.Value)
}

// Rewrite rewrites the regex expr after parse
func (e *RegexExpr) Rewrite() string {
	return fmt.Sprintf("%s=~%s", e.Name, e.Regexp)
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
	case *SelectItem2:
		inner := innerSelectItem{
			exprData: exprData{
				Type: "selectItem",
				Expr: Marshal(e.Expr),
			},
			Alias: e.Alias,
		}
		return encoding.JSONMarshal(&inner)
	case *OrderByExpr:
		inner := innerOrderByExpr{
			exprData: exprData{
				Type: "orderBy",
				Expr: Marshal(e.Expr),
			},
			Desc: e.Desc,
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
	var expr exprData
	err := encoding.JSONUnmarshal(value, &expr)
	if err != nil {
		return nil, err
	}
	switch expr.Type {
	case "regex":
		return unmarshal(&expr, &RegexExpr{})
	case "like":
		return unmarshal(&expr, &LikeExpr{})
	case "in":
		return unmarshal(&expr, &InExpr{})
	case "equals":
		return unmarshal(&expr, &EqualsExpr{})
	case "number":
		return unmarshal(&expr, &NumberLiteral{})
	case field:
		return unmarshal(&expr, &FieldExpr{})
	case "paren":
		e, err := Unmarshal(expr.Expr)
		if err != nil {
			return nil, err
		}
		return &ParenExpr{Expr: e}, nil
	case "binary":
		return unmarshalBinary(value)
	case "selectItem":
		return unmarshalSelectItem(value)
	case "orderBy":
		return unmarshalOrderByExpr(value)
	case "call":
		return unmarshalCall(value)
	case "not":
		e, err := Unmarshal(expr.Expr)
		if err != nil {
			return nil, err
		}
		return &NotExpr{Expr: e}, nil
	default:
		return nil, fmt.Errorf("expr type not match:%s", expr.Type)
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
	return &SelectItem2{Alias: innerExpr.Alias, Expr: e}, nil
}

// unmarshalOrderByExpr parses value to order by expr
func unmarshalOrderByExpr(value []byte) (Expr, error) {
	innerExpr := innerOrderByExpr{}
	err := encoding.JSONUnmarshal(value, &innerExpr)
	if err != nil {
		return nil, err
	}
	e, err := Unmarshal(innerExpr.Expr)
	if err != nil {
		return nil, err
	}
	return &OrderByExpr{Expr: e, Desc: innerExpr.Desc}, nil
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

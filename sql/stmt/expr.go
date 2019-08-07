package stmt

import (
	"fmt"
	"strings"

	"github.com/lindb/lindb/aggregation/function"
)

// Expr represents a interface for all expression types
type Expr interface {
	// Rewrite rewrites the expr after parse
	Rewrite() string
}

// TagFilter represents tag filter for searching time series
type TagFilter interface {
	// TagKey returns the filter's tag key
	TagKey() string
}

// SelectItem represents a select item from select statement
type SelectItem struct {
	Expr  Expr
	Alias string
}

// FieldExpr represents a field name for select list
type FieldExpr struct {
	Name string
}

// CallExpr represents a function call expression
type CallExpr struct {
	FuncType function.FuncType
	Params   []Expr
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

// EqualsExpr represents an equals expression
type EqualsExpr struct {
	Key   string
	Value string
}

// InExpr represents an in expression
type InExpr struct {
	Key    string
	Values []string
}

// LikeExpr represents a like expression
type LikeExpr struct {
	Key   string
	Value string
}

// RegexExpr represents a regular expression
type RegexExpr struct {
	Key    string
	Regexp string
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
	return fmt.Sprintf("%s(%s)", function.FuncTypeString(e.FuncType), strings.Join(params, ","))
}

// Rewrite rewrites the paren expr after parse
func (e *ParenExpr) Rewrite() string {
	return fmt.Sprintf("(%s)", e.Expr.Rewrite())
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

// TagKey returns the equals filter's tag key
func (e *EqualsExpr) TagKey() string { return e.Key }

// TagKey returns the in filter's tag key
func (e *InExpr) TagKey() string { return e.Key }

// TagKey returns the like filter's tag key
func (e *LikeExpr) TagKey() string { return e.Key }

// TagKey returns the regex filter's tag key
func (e *RegexExpr) TagKey() string { return e.Key }

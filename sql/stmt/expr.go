package stmt

import "github.com/lindb/lindb/pkg/function"

// Expr represents a interface for all expression types
type Expr interface {
	// expr ensures spec expression type need implement the interface
	expr()
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
	Type   function.Type
	Params []Expr
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

func (e *SelectItem) expr() {}
func (e *FieldExpr) expr()  {}
func (e *CallExpr) expr()   {}
func (e *ParenExpr) expr()  {}
func (e *BinaryExpr) expr() {}
func (e *NotExpr) expr()    {}
func (e *EqualsExpr) expr() {}
func (e *InExpr) expr()     {}
func (e *LikeExpr) expr()   {}
func (e *RegexExpr) expr()  {}

// TagKey returns the equals filter's tag key
func (e *EqualsExpr) TagKey() string { return e.Key }

// TagKey returns the in filter's tag key
func (e *InExpr) TagKey() string { return e.Key }

// TagKey returns the like filter's tag key
func (e *LikeExpr) TagKey() string { return e.Key }

// TagKey returns the regex filter's tag key
func (e *RegexExpr) TagKey() string { return e.Key }

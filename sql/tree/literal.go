package tree

import (
	"strconv"
	"strings"
)

type Literal interface {
	Node
}

type StringLiteral struct {
	BaseNode

	Value string `json:"value"`
}

func (n *StringLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type BooleanLiteral struct {
	BaseNode
	Value bool
}

func (n *BooleanLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func NewBooleanLiteral(location *NodeLocation, value string) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode: BaseNode{
			Location: location,
		},
		Value: strings.ToLower(value) == "true",
	}
}

type LongLiteral struct {
	BaseNode
	Value int64
}

func NewLongLiteral(location *NodeLocation, value string) *LongLiteral {
	// TODO: check error
	val, _ := strconv.ParseInt(value, 10, 64)
	return &LongLiteral{
		BaseNode: BaseNode{
			Location: location,
		},
		Value: val,
	}
}

func (n *LongLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type FloatLiteral struct {
	BaseNode
	Value float64
}

func NewFloatLiteral(location *NodeLocation, value string) *FloatLiteral {
	// TODO: parse value
	return &FloatLiteral{
		BaseNode: BaseNode{
			Location: location,
		},
		Value: 1.0,
	}
}

func (n *FloatLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

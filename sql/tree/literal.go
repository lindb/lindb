package tree

import (
	"strconv"
	"strings"

	"github.com/lindb/lindb/spi/types"
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

func NewBooleanLiteral(id NodeID, location *NodeLocation, value string) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode: BaseNode{
			ID:       id,
			Location: location,
		},
		Value: strings.ToLower(value) == "true",
	}
}

type LongLiteral struct {
	BaseNode
	Value int64
}

func NewLongLiteral(id NodeID, location *NodeLocation, value string) *LongLiteral {
	// TODO: check error
	val, _ := strconv.ParseInt(value, 10, 64)
	return &LongLiteral{
		BaseNode: BaseNode{
			ID:       id,
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

func NewFloatLiteral(id NodeID, location *NodeLocation, value string) *FloatLiteral {
	// TODO: parse value
	val, _ := strconv.ParseFloat(value, 64)
	return &FloatLiteral{
		BaseNode: BaseNode{
			ID:       id,
			Location: location,
		},
		Value: val,
	}
}

func (n *FloatLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type Constant struct {
	BaseNode
	Value any
	Type  types.DataType // TODO:
}

func (n *Constant) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

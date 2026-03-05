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

package tree

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
)

type (
	LogicalOperator    string
	ComparisonOperator string
	ArithmeticOperator string
)

var (
	LogicalAND LogicalOperator = "AND"
	LogicalOR  LogicalOperator = "OR"

	ComparisonEQ  ComparisonOperator = "="
	ComparisonNEQ ComparisonOperator = "!="
	ComparisonGT  ComparisonOperator = ">"
	ComparisonLT  ComparisonOperator = "<"

	Add      ArithmeticOperator = "+"
	Subtract ArithmeticOperator = "-"
	Multiply ArithmeticOperator = "*"
	Divide   ArithmeticOperator = "/"
	Modulus  ArithmeticOperator = "%"
)

func (op ArithmeticOperator) FunctionName() FuncName {
	switch op {
	case Add:
		return Plus
	case Subtract:
		return Minus
	case Multiply:
		return Mul
	case Divide:
		return Div
	case Modulus:
		return Mod
	default:
		panic(fmt.Sprintf("unknown arithmetic operator: %s", op))
	}
}

type Expression interface {
	Node
}

type ArrayExpression struct {
	BaseNode
	Elements []Expression
}

func (n *ArrayExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type Row struct {
	// FIXME: remove it
	BaseNode
	Items []Expression
}

func (n *Row) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type Cast struct {
	BaseNode
	Type       arrow.DataType `json:"type"`
	Expression Expression     `json:"expression"`
}

func (n *Cast) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type FieldReference struct {
	BaseNode

	FieldIndex FieldIndex
}

func (n *FieldReference) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type DereferenceExpression struct {
	BaseNode
	Base  Expression
	Field *Identifier
}

func (n *DereferenceExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

func (n *DereferenceExpression) ToQualifiedName() (name *QualifiedName) {
	if n.Field == nil {
		return
	}
	switch e := n.Base.(type) {
	case *Identifier:
		name = NewQualifiedName([]*Identifier{e, n.Field})
	case *DereferenceExpression:
		baseQualifiedName := e.ToQualifiedName()
		if baseQualifiedName != nil {
			parts := baseQualifiedName.OriginalParts
			parts = append(parts, n.Field)
			name = NewQualifiedName(parts)
		}
	}
	return
}

type ArithmeticBinaryExpression struct {
	BaseNode

	Left     Expression         `json:"left"`
	Right    Expression         `json:"right"`
	Operator ArithmeticOperator `json:"operator"` // TODO: add type
}

func (n *ArithmeticBinaryExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ComparisonExpression struct {
	BaseNode

	Left     Expression         `json:"left"`
	Right    Expression         `json:"right"`
	Operator ComparisonOperator `json:"operator"`
}

func (n *ComparisonExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type LogicalExpression struct {
	BaseNode
	Operator LogicalOperator
	Terms    []Expression
}

func (n *LogicalExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type InListExpression struct {
	BaseNode
	Values []Expression `json:"values"`
}

func (n *InListExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type NotExpression struct {
	BaseNode
	Value Expression `json:"value"`
}

func (n *NotExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

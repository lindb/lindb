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

type LikePredicate struct {
	BaseNode
	Value   Expression `json:"value"`
	Pattern Expression `json:"pattern"`
}

func (n *LikePredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type InPredicate struct {
	BaseNode
	Value     Expression `json:"value"`
	ValueList Expression `json:"valueList"`
}

func (n *InPredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type RegexPredicate struct {
	BaseNode
	Value   Expression `json:"value"`
	Pattern Expression `json:"pattern"`
}

func (n *RegexPredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type TimePredicate struct {
	BaseNode
	Operator ComparisonOperator
	Value    Expression
}

func (n *TimePredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type NullPredicate struct {
	BaseNode
	Value Expression `json:"value"`
	Not   bool       `json:"not"`
}

func (n *NullPredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

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

package types

var (
	Equal              OperatorType = NewOperatorType("=", 2)
	NotEqual           OperatorType = NewOperatorType("!=", 2)
	GreaterThan        OperatorType = NewOperatorType(">", 2)
	GreaterThanOrEqual OperatorType = NewOperatorType(">=", 2)
	LessThan           OperatorType = NewOperatorType("<", 2)
	LessThanOrEqual    OperatorType = NewOperatorType("<=", 2)
	Add                OperatorType = NewOperatorType("+", 2)
	Subtract           OperatorType = NewOperatorType("-", 2)
	Multiply           OperatorType = NewOperatorType("*", 2)
	Divide             OperatorType = NewOperatorType("/", 2)
	Modulus            OperatorType = NewOperatorType("%", 2)
)

type OperatorType struct {
	// FIXME: remove it
	Operator      string
	ArgumentCount int
}

func NewOperatorType(operator string, argumentCount int) OperatorType {
	return OperatorType{
		Operator:      operator,
		ArgumentCount: argumentCount,
	}
}

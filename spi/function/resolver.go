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

package function

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type ResolvedFunction struct {
	Signature *BoundSignature `json:"signature"`
}

type FunctionResolver struct{}

func NewFunctionResolver() *FunctionResolver {
	return &FunctionResolver{}
}

func (r *FunctionResolver) ResolveOperator(operatorType types.OperatorType, argumentTypes []types.Type) *ResolvedFunction {
	return &ResolvedFunction{
		// FIXME: function name/types
		Signature: NewBoundSignature(operatorType.Operator, types.DTFloat, []types.DataType{types.DTFloat}),
	}
}

func (r *FunctionResolver) ResolveFunction(name *tree.QualifiedName) *ResolvedFunction {
	return &ResolvedFunction{
		// FIXME: function name/types
		Signature: NewBoundSignature(name.Suffix, types.DTFloat, []types.DataType{types.DTFloat}),
	}
}

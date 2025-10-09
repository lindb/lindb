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
)

type DefaultTraversalVisitor struct {
	PreProcess  func(n Node)
	PostProcess func(n Node)
}

func (v *DefaultTraversalVisitor) Visit(context any, n Node) (r any) {
	fmt.Printf("express visit = %T value=%v\n", n, n)
	if v.PreProcess != nil {
		// do pre process if has pre func
		v.PreProcess(n)
	}
	switch node := n.(type) {
	case *ArithmeticBinaryExpression:
		_ = node.Left.Accept(context, v)
		_ = node.Right.Accept(context, v)
	case *FunctionCall:
		for _, arg := range node.Arguments {
			_ = arg.Accept(context, v)
		}
	case *LogicalExpression:
		for _, term := range node.Terms {
			_ = term.Accept(context, v)
		}
	case *ComparisonExpression:
		_ = node.Left.Accept(context, v)
		_ = node.Right.Accept(context, v)
	default:
		// TODO: remove
		fmt.Printf("default traversal visitor not support..................=%T\n", n)
	}
	if v.PostProcess != nil {
		// do post process if has post func
		v.PostProcess(n)
	}
	return
}

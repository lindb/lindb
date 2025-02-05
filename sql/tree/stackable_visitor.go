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

import "github.com/lindb/lindb/pkg/collections"

type StackableVisitorContext[C any] struct {
	stack   *collections.Stack
	context C
}

func NewStackableVisitorContext[C any](context C) *StackableVisitorContext[C] {
	return &StackableVisitorContext[C]{
		stack:   collections.NewStack(),
		context: context,
	}
}

func (c *StackableVisitorContext[C]) GetContext() (r C) {
	return c.context
}

func (c *StackableVisitorContext[C]) Push(node Node) {
	c.stack.Push(node)
}

func (c *StackableVisitorContext[C]) Pop() {
	_ = c.stack.Pop()
}

func (c *StackableVisitorContext[C]) GetPreviousNode() Node {
	if c.stack.Size() > 1 {
		return c.stack.Get(1).(Node)
	}
	return nil
}

type StackableAstVisitor[C any] struct{}

func (v *StackableAstVisitor[C]) Visit(context any, node Node) any {
	stackCxt := context.(*StackableVisitorContext[C])
	stackCxt.Push(node)
	defer func() {
		stackCxt.Pop()
	}()
	// FIXME: ??
	return nil
}

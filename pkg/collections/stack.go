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

package collections

// Stack represents a last-in-first-out(LIFO) stack of interface, using slice saving data.
// NOTICE: not safe for goroutine concurrent
type Stack struct {
	elements []interface{}
	len      int
}

// NewStack creates an empty stack
func NewStack() *Stack {
	return &Stack{}
}

// Push pushes an element onto the top of this stack
func (s *Stack) Push(element interface{}) {
	s.elements = append(s.elements, element)
	s.len++
}

// Pop removes the element at the top of this stack, if stack is empty return nil
func (s *Stack) Pop() interface{} {
	if s.len == 0 {
		return nil
	}
	element := s.Peek()

	s.elements = s.elements[0 : s.len-1]
	s.len--
	return element
}

// Peek looks at the element at the top of this stack without removing it from the stack
func (s *Stack) Peek() interface{} {
	if s.len == 0 {
		return nil
	}
	return s.elements[s.len-1]
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return s.len
}

// Empty tests if this stack is empty
func (s *Stack) Empty() bool {
	return s.len == 0
}

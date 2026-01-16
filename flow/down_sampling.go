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

package flow

type Stream[V any] interface {
	SetAtStep(step int, value V, fn func(a, b V) V)
	GetAtStep(step int) V
	Reset()
}

type stream[V any] struct {
	values []V
}

func NewStream[V any](size int) Stream[V] {
	return &stream[V]{
		values: make([]V, size),
	}
}

func (s *stream[V]) SetAtStep(step int, value V, fn func(a, b V) V) {
	s.values[step] = fn(s.values[step], value)
}

func (s *stream[V]) GetAtStep(step int) V {
	return s.values[step]
}

func (s *stream[V]) Reset() {
	// FIXME: implement reset logic
}

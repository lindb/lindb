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

package linmetric

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_DeltaCounterVec(t *testing.T) {
	scope := NewScope("vec")
	vec := scope.NewCounterVec("count", "1", "2")
	assert.Panics(t, func() {
		vec.WithTagValues("1", "2", "3")
	})
	assert.Panics(t, func() {
		scope.NewCounterVec("count2")
	})
	vec.WithTagValues("a", "b").Incr()
	vec.WithTagValues("a", "c").Incr()
	vec.WithTagValues("a", "b").Incr()
}

func Benchmark_DeltaCounterVec(b *testing.B) {
	scope := NewScope("vec_test")
	vec := scope.NewCounterVec("counter", "1", "2")

	for i := 0; i < b.N; i++ {
		vec.WithTagValues("3", "4").Incr()
	}
}

func Benchmark_DeltaCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewScope("counter_test", "1", "3", "2", "4").
			NewCounter("counter").
			Incr()
	}
}

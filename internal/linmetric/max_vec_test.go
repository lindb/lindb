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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MaxVec(t *testing.T) {
	scope := BrokerRegistry.NewScope("vecg")
	vec := scope.NewMaxVec("Max", "1", "2")
	assert.Panics(t, func() {
		vec.WithTagValues("1", "2", "3")
	})
	assert.Panics(t, func() {
		scope.NewMaxVec("count2")
	})
	vec.WithTagValues("a", "b").Update(1)
	vec.WithTagValues("a", "c").Update(2)
	vec.WithTagValues("a", "b").Update(1)
}

func Benchmark_MaxVec(b *testing.B) {
	scope := BrokerRegistry.NewScope("vec_test")
	vec := scope.NewMaxVec("Max", "1", "2")

	for i := 0; i < b.N; i++ {
		vec.WithTagValues("3", "4").Update(2323)
	}
}

func Benchmark_Max(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BrokerRegistry.NewScope("Max_test", "1", "3", "2", "4").
			NewMax("Max").
			Update(2)
	}
}

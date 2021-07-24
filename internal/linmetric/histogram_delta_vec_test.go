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
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_HistogramDeltaVec(t *testing.T) {
	scope := NewScope("44")
	vec := scope.NewDeltaHistogramVec("1", "2")

	assert.Panics(t, func() {
		vec.WithTagValues("1")
	})
	vec.WithLinearBuckets(time.Millisecond, time.Second*4, 10)
	vec.WithTagValues("1", "2").UpdateSeconds(1)
	vec.WithExponentBuckets(time.Millisecond, time.Second*4, 10)

	vec.WithTagValues("a", "b").UpdateSeconds(1)
	vec.WithTagValues("a", "c").UpdateSeconds(1)
	vec.WithTagValues("a", "b").UpdateSeconds(1)
}

func Benchmark_HistogramVec(b *testing.B) {
	scope := NewScope("vec_test")
	vec := scope.NewDeltaHistogramVec("1", "2").
		WithExponentBuckets(time.Millisecond, time.Second*4, 10)

	for i := 0; i < b.N; i++ {
		vec.WithTagValues("3", "4").UpdateSeconds(1)
	}
}

func Benchmark_histogram(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewScope("histogram_test", "1", "3", "2", "4").
			NewDeltaHistogram().
			WithExponentBuckets(time.Millisecond, time.Second*4, 10).
			UpdateSeconds(1)
	}
}

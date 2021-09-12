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

package linmetric_test

import (
	"math"
	"testing"
	"time"

	"github.com/lindb/lindb/internal/linmetric"

	"github.com/stretchr/testify/assert"
)

func Test_MetricScope(t *testing.T) {
	scope0 := linmetric.NewScope("0")
	scope0.Scope("x")
	scope0.Scope("x")

	scope1 := linmetric.NewScope("1",
		"k2", "v2", "k1", "v1", "k2", "v2")
	scope1.NewGauge("g1").Incr()
	scope1.NewCounter("c2").Incr()
	scope1.NewCounter("c2").Incr()
	scope1.NewMax("max3").Update(1)
	scope1.NewMax("max3").Update(2)
	scope1.NewMax("max3").Update(math.Inf(1))

	assert.Panics(t, func() {
		scope1.NewCounter("max3")
	})

	scope1.NewMin("min4").Update(2)
	scope1.NewMin("min4").Update(1)
	scope1.NewMin("min4").Update(math.Inf(-1))

	scope12 := scope1.Scope("2", "k1", "v1", "k3", "v3")
	scope12.NewGauge("g1").Update(1)
	scope12.NewGauge("g1").Update(2)
	scope12.NewHistogram().UpdateDuration(time.Second)
	scope12.NewHistogram().UpdateDuration(time.Second)
	time.Sleep(time.Second)
	gather := linmetric.NewGather(linmetric.WithReadRuntimeOption())
	_, _ = gather.Gather()
	_, _ = gather.Gather()
}

func Test_MetricScope_Scope(t *testing.T) {
	assert.Panics(t, func() {
		linmetric.NewScope("")
	})
	assert.Panics(t, func() {
		linmetric.NewScope("lindb", "1")
	})

	scope3 := linmetric.NewScope("3")
	scope3.NewCounter("c")
	assert.Panics(t, func() {
		scope3.NewGauge("c")
	})
	assert.Panics(t, func() {
		scope3.NewMax("c")
	})
	assert.Panics(t, func() {
		scope3.NewMin("c")
	})
	scope3.NewCounter("d")
	assert.Panics(t, func() {
		scope3.NewHistogramVec()
	})
	assert.Panics(t, func() {
		scope3.NewCounterVec("23")
	})
	assert.Panics(t, func() {
		scope3.NewGaugeVec("23")
	})
	assert.Panics(t, func() {
		scope3.NewMaxVec("23")
	})
	assert.Panics(t, func() {
		scope3.NewMinVec("23")
	})
	assert.Panics(t, func() {
		scope3.NewGauge("")
	})
}

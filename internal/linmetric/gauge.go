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
	"go.uber.org/atomic"

	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
)

// BoundGauge is a gauge which has Bound to a certain metric with field-name and tags
type BoundGauge struct {
	value     atomic.Float64
	fieldName string
}

func newGauge(fieldName string) *BoundGauge {
	return &BoundGauge{
		fieldName: fieldName,
		value:     *atomic.NewFloat64(0),
	}
}

// Update updates gauge with a new value
func (g *BoundGauge) Update(v float64) {
	g.value.Store(v)
}

// Add adds v to g.
func (g *BoundGauge) Add(v float64) {
	g.value.Add(v)
}

// Sub subs v to g.
func (g *BoundGauge) Sub(v float64) {
	g.value.Sub(v)
}

// Incr increments g.
func (g *BoundGauge) Incr() {
	g.value.Add(1)
}

// Decr decrements g.
func (g *BoundGauge) Decr() {
	g.value.Sub(1)
}

// Get returns the current gauge value
func (g *BoundGauge) Get() float64 {
	return g.value.Load()
}

func (g *BoundGauge) gather() float64 { return g.value.Load() }

func (g *BoundGauge) name() string { return g.fieldName }

func (g *BoundGauge) flatType() flatMetricsV1.SimpleFieldType {
	return flatMetricsV1.SimpleFieldTypeGauge
}

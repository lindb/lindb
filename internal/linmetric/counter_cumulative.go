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
)

// BoundCumulativeCounter is a counter which has been Bound to a certain metric
// with field-name and metrics, it supports update method
// make sure that data is cumulative
type BoundCumulativeCounter struct {
	value     atomic.Float64
	fieldName string
}

func newCumulativeCounter(fieldName string) *BoundCumulativeCounter {
	return &BoundCumulativeCounter{fieldName: fieldName}
}

// Incr increments c.
func (c *BoundCumulativeCounter) Incr() {
	c.value.Add(1)
}

// Decr decrements g.
func (c *BoundCumulativeCounter) Decr() {
	c.value.Sub(1)
}

// Add adds v to c.
func (c *BoundCumulativeCounter) Add(v float64) {
	c.value.Add(v)
}

// Sub subs v to c.
func (c *BoundCumulativeCounter) Sub(v float64) {
	c.value.Sub(v)
}

// Get returns the current cumulative counter value
func (c *BoundCumulativeCounter) Get() float64 {
	return c.value.Load()
}

// Update updates counter with a new value
func (c *BoundCumulativeCounter) Update(v float64) {
	c.value.Store(v)
}

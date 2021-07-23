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

import "go.uber.org/atomic"

// BoundDeltaCounter is a counter which has been Bound to a certain metric
// with field-name and metrics, it does not support update method.
// Get will resets the underlying delta value
type BoundDeltaCounter struct {
	delta     atomic.Float64
	fieldName string
}

func newDeltaCounter(fieldName string) *BoundDeltaCounter {
	return &BoundDeltaCounter{
		fieldName: fieldName,
	}
}

// Incr increments c.
func (c *BoundDeltaCounter) Incr() {
	c.delta.Add(1)
}

// Decr decrements g.
func (c *BoundDeltaCounter) Decr() {
	c.delta.Sub(1)
}

// Add adds v to c.
func (c *BoundDeltaCounter) Add(v float64) {
	c.delta.Add(v)
}

// Sub subs v to c.
func (c *BoundDeltaCounter) Sub(v float64) {
	c.delta.Sub(v)
}

// Get returns the current delta counter value
func (c *BoundDeltaCounter) Get() float64 {
	return c.delta.Load()
}

// getAndReset returns the current cumulative counter value
// and resets the delta value by spin lock.
func (c *BoundDeltaCounter) getAndReset() float64 {
	for {
		v := c.delta.Load()
		if c.delta.CAS(v, 0) {
			return v
		}
	}
}

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

// BoundCounter is a counter which has been Bound to a certain metric
// with field-name and metrics, it does not support update method.
// Get will resets the underlying delta value
type BoundCounter struct {
	delta     atomic.Float64
	fieldName string
}

func newCounter(fieldName string) *BoundCounter {
	return &BoundCounter{
		fieldName: fieldName,
	}
}

// Incr increments c.
func (c *BoundCounter) Incr() {
	c.delta.Add(1)
}

// Decr decrements g.
func (c *BoundCounter) Decr() {
	c.delta.Sub(1)
}

// Add adds v to c.
func (c *BoundCounter) Add(v float64) {
	c.delta.Add(v)
}

// Sub subs v to c.
func (c *BoundCounter) Sub(v float64) {
	c.delta.Sub(v)
}

// Get returns the current delta counter value
func (c *BoundCounter) Get() float64 {
	return c.delta.Load()
}

// gather returns the current cumulative counter value
// and resets the delta value by spin lock.
func (c *BoundCounter) gather() float64 {
	for {
		v := c.delta.Load()
		if c.delta.CAS(v, 0) {
			return v
		}
	}
}

func (c *BoundCounter) name() string { return c.fieldName }

func (c *BoundCounter) flatType() flatMetricsV1.SimpleFieldType {
	return flatMetricsV1.SimpleFieldTypeDeltaSum
}

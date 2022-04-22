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
	"math"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
)

type BoundMax struct {
	value     atomic.Float64
	fieldName string
}

func newMax(fieldName string) *BoundMax {
	return &BoundMax{
		fieldName: fieldName,
		value:     *atomic.NewFloat64(math.Inf(-1)),
	}
}

// Update updates Max with a new value
// Skip updating when newValue is smaller than v
func (m *BoundMax) Update(newValue float64) {
	if m.value.Load() >= newValue {
		return
	}
	for {
		v := m.value.Load()
		if newValue <= v {
			break
		}
		if m.value.CAS(v, newValue) {
			return
		}
	}
}

// Get returns the current max value
func (m *BoundMax) Get() float64 {
	return m.value.Load()
}

func (m *BoundMax) gather() float64 { return m.value.Load() }

func (m *BoundMax) name() string { return m.fieldName }

func (m *BoundMax) flatType() flatMetricsV1.SimpleFieldType {
	return flatMetricsV1.SimpleFieldTypeMax
}

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

func TestRegistry_FindMetricList(t *testing.T) {
	r := &Registry{
		series: make(map[uint64]*taggedSeries),
	}
	r.NewScope("test-1", "a", "a-1", "b", "b").NewCounter("f")
	r.NewScope("test-1", "a", "a-2", "b", "b").NewCounter("f")
	r.NewScope("test-2", "a", "a-2", "b", "b").NewCounter("f")

	rs := r.FindMetricList([]string{"test-1"}, nil)
	assert.Len(t, rs["test-1"], 2)

	rs = r.FindMetricList([]string{"test-1"}, map[string]string{"a": "a-1"})
	assert.Len(t, rs["test-1"], 1)
}

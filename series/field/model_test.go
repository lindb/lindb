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

package field

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fields(t *testing.T) {
	var fs Fields
	fs = append(fs,
		Field{Name: []byte("a"), Type: SumField, Value: float64(0)},
		Field{Name: []byte("c"), Type: HistogramField, Value: float64(0)},
		Field{Name: []byte("b"), Type: SummaryField, Value: float64(0)})
	sort.Sort(fs)

	fs = fs.Insert(Field{Name: []byte("b"), Type: MaxField, Value: float64(0)})
	assert.Equal(t, MaxField, fs[1].Type)

	fs = fs.Insert(Field{Name: []byte("d"), Type: MinField, Value: float64(0)})
	assert.Equal(t, MinField, fs[3].Type)
}

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

package aggregation

import (
	"github.com/lindb/lindb/sql/tree"
)

type countAggFactory struct{}

func (f *countAggFactory) NewAggregator(args []tree.Expression) Aggregator {
	return &countAgg{}
}

type countAgg struct {
	value int64
}

func NewCountAgg() Aggregator {
	return &countAgg{}
}

// Process implements Aggregator.
func (c *countAgg) Process(event any) {
	c.value++
	// fmt.Println("count....")
}

func (c *countAgg) GetValue() any {
	return c.value
}

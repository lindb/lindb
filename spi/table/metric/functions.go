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

package metric

import (
	"fmt"

	"github.com/lindb/common/models"

	"github.com/lindb/common/field"
	"github.com/lindb/lindb/sql/tree"
)

type aggregateFunc[V float64 | *models.Exemplar] func(a, b V) V

type getAggregateFunc[V float64 | *models.Exemplar] func(funcName tree.FuncName) aggregateFunc[V]

func getAggFunc(funcName tree.FuncName) aggregateFunc[float64] {
	switch funcName {
	case tree.Sum:
		return field.Sum.Aggregate
	case tree.Max:
		return field.Max.Aggregate
	case tree.Min:
		return field.Min.Aggregate
	case tree.Last:
		return field.Last.Aggregate
	case tree.First:
		return field.First.Aggregate
	default:
		panic(fmt.Sprintf("aggregation function not support: %s", funcName))
	}
}

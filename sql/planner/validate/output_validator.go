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

package validate

import (
	"errors"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type OutputValidator struct {
	Base[*plan.OutputNode]
}

func NewOutputValidator() Validator {
	v := &OutputValidator{}
	v.validate = func(ctx *context.PlannerContext, node *plan.OutputNode) error {
		_, ok := lo.Find(node.GetOutputSymbols(), func(item *plan.Symbol) bool {
			return item.DataType == types.DTTimestamp
		})
		if !ok {
			// output node has no timestamp column
			return nil
		}
		if _, ok = lo.Find(node.GetOutputSymbols(), func(item *plan.Symbol) bool {
			return item.DataType == types.DTTimeSeries
		}); !ok {
			return errors.New("push timestamp column failed, output node has no time series column")
		}
		return nil
	}
	return v
}

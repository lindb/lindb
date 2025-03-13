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

package printer

import (
	"fmt"
	"strings"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

func formatSymbols(symbols []*plan.Symbol) string {
	var columns []string
	for i := range symbols {
		columns = append(columns, symbols[i].String())
	}
	return "[" + strings.Join(columns, ", ") + "]"
}

func formatAggregation(aggregation *plan.Aggregation) string {
	var args []string
	for _, arg := range aggregation.Arguments {
		args = append(args, fmt.Sprintf("%q", tree.FormatExpression(arg)))
	}
	return fmt.Sprintf("%s(%s)", aggregation.Function, strings.Join(args, ", "))
}

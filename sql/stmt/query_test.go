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

package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestQuery_Marshal(t *testing.T) {
	query := Query{
		Namespace:  "ns",
		MetricName: "test",
		AllFields:  true,
		SelectItems: []Expr{
			&SelectItem{Expr: &FieldExpr{Name: "a"}},
			&SelectItem{Expr: &FieldExpr{Name: "b"}},
			&SelectItem{
				Expr: &CallExpr{
					FuncType: function.Stddev,
					Params: []Expr{
						&CallExpr{
							FuncType: function.Max,
							Params: []Expr{&CallExpr{
								FuncType: function.Sum,
								Params:   []Expr{&FieldExpr{Name: "c"}}},
							},
						},
					},
				},
			},
		},
		Condition: &BinaryExpr{
			Left: &ParenExpr{Expr: &BinaryExpr{
				Left:     &InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}},
				Operator: AND,
				Right:    &EqualsExpr{Key: "region", Value: "sh"},
			}},
			Operator: AND,
			Right: &ParenExpr{Expr: &BinaryExpr{
				Left:     &EqualsExpr{Key: "path", Value: "/data"},
				Operator: OR,
				Right:    &EqualsExpr{Key: "path", Value: "/home"},
			}},
		},
		TimeRange: timeutil.TimeRange{Start: 10, End: 30},
		Interval:  1000,
		GroupBy:   []string{"a", "b", "c"},
		OrderByItems: []Expr{
			&FieldExpr{Name: "b"},
			&CallExpr{
				FuncType: function.Max,
				Params:   []Expr{&FieldExpr{Name: "c"}},
			},
		},
		Limit: 100,
	}

	data := encoding.JSONMarshal(&query)
	query1 := Query{}
	err := encoding.JSONUnmarshal(data, &query1)
	assert.NoError(t, err)
	assert.Equal(t, query, query1)
	assert.True(t, query.HasGroupBy())
	assert.True(t, query.AllFields)
}

func TestQuery_Marshal_Fail(t *testing.T) {
	query := &Query{}
	err := query.UnmarshalJSON([]byte{1, 2, 3})
	assert.Error(t, err)
	err = query.UnmarshalJSON([]byte("{\"condition\":\"123\"}"))
	assert.Error(t, err)
	err = query.UnmarshalJSON([]byte("{\"selectItems\":[\"123\"]}"))
	assert.Error(t, err)
	err = query.UnmarshalJSON([]byte("{\"orderByItems\":[\"123\"]}"))
	assert.Error(t, err)
}

func TestQuery_StatementType(t *testing.T) {
	assert.Equal(t, QueryStatement, (&Query{}).StatementType())
}

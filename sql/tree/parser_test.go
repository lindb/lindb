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

package tree

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/grammar"
)

var (
	tableA = &AliasedRelation{
		Aliase: &Identifier{
			Value: "aa",
		},
		Relation: &Table{
			Name: &QualifiedName{
				Name:  "a",
				Parts: []string{"a"},
			},
		},
	}
	tableB = &Table{
		Name: &QualifiedName{
			Name:  "b",
			Parts: []string{"b"},
		},
	}
	table1 = &Table{
		Name: &QualifiedName{
			Name:  "table1",
			Parts: []string{"table1"},
		},
	}
	table2 = &Table{
		Name: &QualifiedName{
			Name:  "table2",
			Parts: []string{"table2"},
		},
	}

	expression1 = &ComparisonExpression{
		Operator: "=",
		Left: &DereferenceExpression{
			Base: &Identifier{
				Value: "aa",
			},
			Field: &Identifier{
				Value: "a",
			},
		},
		Right: &Identifier{
			Value: "b",
		},
	}
	expression2 = &ComparisonExpression{
		Operator: "=",
		Left: &Identifier{
			Value: "c",
		},
		Right: &Identifier{
			Value: "d",
		},
	}
)

func TestSQLParser_QueryStatement_Error(t *testing.T) {
	defer func() {
		newNodeLocation = NewNodeLocation
	}()
	newNodeLocation = func(line, column int) *NodeLocation {
		return nil
	}
	parser := GetParser()
	_, _ = parser.CreateStatement("s", NewNodeIDAllocator())
}

func TestSQLParser_QueryStatement(t *testing.T) {
	defer func() {
		newNodeLocation = NewNodeLocation
	}()
	newNodeLocation = func(line, column int) *NodeLocation {
		return nil
	}
	parser := GetParser()
	cases := []struct {
		sql  string
		stmt Statement
	}{
		{
			"select * from a aa group by a,b,c",
			&Query{
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{&AllColumns{}},
					},
					From: tableA,
					GroupBy: &GroupBy{
						GroupingElements: []GroupingElement{
							&SimpleGroupBy{
								Columns: []Expression{
									&Identifier{Value: "a"},
								},
							},
							&SimpleGroupBy{
								Columns: []Expression{
									&Identifier{Value: "b"},
								},
							},
							&SimpleGroupBy{
								Columns: []Expression{
									&Identifier{Value: "c"},
								},
							},
						},
					},
				},
			},
		},
		{
			"select a,b,c from b where aa.a=b and c=d",
			&Query{
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{
							&SingleColumn{
								Expression: &Identifier{Value: "a"},
							},
							&SingleColumn{
								Expression: &Identifier{Value: "b"},
							},
							&SingleColumn{
								Expression: &Identifier{Value: "c"},
							},
						},
					},
					From: tableB,
					Where: &LogicalExpression{
						Operator: LogicalAND,
						Terms: []Expression{
							expression1,
							expression2,
						},
					},
				},
			},
		},
		{
			"select * from a aa,b where aa.a=b and c=d and not(aa.a=b)",
			&Query{
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{&AllColumns{}},
					},
					From: &Join{
						Type:  IMPLICIT,
						Left:  tableA,
						Right: tableB,
					},
					Where: &LogicalExpression{
						Operator: LogicalAND,
						Terms: []Expression{
							expression1,
							expression2,
							&NotExpression{
								Value: expression1,
							},
						},
					},
				},
			},
		},
		{
			"select * from a aa,b where aa.a=b or c=d and not(aa.a=b)",
			&Query{
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{&AllColumns{}},
					},
					From: &Join{
						Type:  IMPLICIT,
						Left:  tableA,
						Right: tableB,
					},
					Where: &LogicalExpression{
						Operator: LogicalOR,
						Terms: []Expression{
							expression1,
							&LogicalExpression{
								Operator: LogicalAND,
								Terms: []Expression{
									expression2,
									&NotExpression{
										Value: expression1,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"select * from a aa,b where aa.a=b and c=d or not(aa.a=b)",
			&Query{
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{&AllColumns{}},
					},
					From: &Join{
						Type:  IMPLICIT,
						Left:  tableA,
						Right: tableB,
					},
					Where: &LogicalExpression{
						Operator: LogicalOR,
						Terms: []Expression{
							&LogicalExpression{
								Operator: LogicalAND,
								Terms: []Expression{
									expression1,
									expression2,
								},
							},
							&NotExpression{
								Value: expression1,
							},
						},
					},
				},
			},
		},
		{
			"select * from a aa where aa.a=b and (c=d or c=d or c=d)",
			&Query{
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{&AllColumns{}},
					},
					From: tableA,
					Where: &LogicalExpression{
						Operator: LogicalAND,
						Terms: []Expression{
							expression1,
							&LogicalExpression{
								Operator: LogicalOR,
								Terms: []Expression{
									expression2,
									expression2,
									expression2,
								},
							},
						},
					},
				},
			},
		},
		{
			`
			with 
			  table1 as select * from a aa where aa.a=b,
			  table2 as select * from b where c=d
			select * from table1 left join table2 using(a,b,c)
			`,
			&Query{
				With: &With{
					Queries: []*WithQuery{
						{
							Name: &Identifier{Value: "table1"},
							Query: &Query{
								QueryBody: &QuerySpecification{
									Select: &Select{
										SelectItems: []SelectItem{&AllColumns{}},
									},
									From:  tableA,
									Where: expression1,
								},
							},
						},
						{
							Name: &Identifier{Value: "table2"},
							Query: &Query{
								QueryBody: &QuerySpecification{
									Select: &Select{
										SelectItems: []SelectItem{&AllColumns{}},
									},
									From:  tableB,
									Where: expression2,
								},
							},
						},
					},
				},
				QueryBody: &QuerySpecification{
					Select: &Select{
						SelectItems: []SelectItem{&AllColumns{}},
					},
					From: &Join{
						Type:  LEFT,
						Left:  table1,
						Right: table2,
						Criteria: &JoinUsing{
							Columns: []*Identifier{
								{Value: "a"},
								{Value: "b"},
								{Value: "c"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := parser.CreateStatement(tt.sql, NewNodeIDAllocator())
			assert.NoError(t, err)
		})
	}
}

func TestSQLParse(t *testing.T) {
	parser := GetParser()
	_, err := parser.CreateStatement(`select 12*(idle*10+100)/10,node from
		"lindb.monitor.system.cpu_stat" group by node`,
		NewNodeIDAllocator())
	assert.NoError(t, err)
}

func TestSQLParse_CountStar(t *testing.T) {
	// count(*) must parse without error and produce a FunctionCall with no arguments.
	parser := GetParser()
	cases := []string{
		"select count(*) from logs",
		"SELECT COUNT(*) FROM logs",
		"select count(*) from logs where level = 'INFO'",
		"SELECT count(*) FROM logs GROUP BY level",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			stmt, err := parser.CreateStatement(sql, NewNodeIDAllocator())
			assert.NoError(t, err)
			assert.NotNil(t, stmt)

			// Drill into the AST to verify count(*) is a FunctionCall named "count"
			// with zero arguments — not with a star argument.
			q, ok := stmt.(*Query)
			assert.True(t, ok)
			spec, ok := q.QueryBody.(*QuerySpecification)
			assert.True(t, ok)
			// Find the count(*) SingleColumn
			found := false
			for _, item := range spec.Select.SelectItems {
				sc, ok := item.(*SingleColumn)
				if !ok {
					continue
				}
				fn, ok := sc.Expression.(*FunctionCall)
				if !ok {
					continue
				}
				if string(fn.Name) == "count" {
					assert.Empty(t, fn.Arguments, "count(*) should produce FunctionCall with no arguments")
					found = true
				}
			}
			assert.True(t, found, "count function not found in SELECT items")
		})
	}
}

func Test_KeyWorks(t *testing.T) {
	typ := reflect.TypeOf(&grammar.NonReservedContext{})
	var keyWords []string
	for i := range typ.NumMethod() {
		methodName := typ.Method(i).Name
		if strings.ToUpper(methodName) == methodName {
			fmt.Println(methodName)
			keyWords = append(keyWords, methodName)
		}
	}
	sort.Strings(keyWords)
	count := 0
	for _, name := range keyWords {
		fmt.Printf("%-15s", name)
		count++
		if count%6 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Println()
}

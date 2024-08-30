package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	for i := range cases {
		tt := cases[i]
		t.Run(tt.sql, func(t *testing.T) {
			stmt, err := parser.CreateStatement(tt.sql, NewNodeIDAllocator())
			assert.NoError(t, err)
			assert.Equal(t, tt.stmt, stmt)
		})
	}
}

func TestSQLParse(t *testing.T) {
	parser := GetParser()
	_, err := parser.CreateStatement(`select 12*(idle*10+100)/10,node from "lindb.monitor.system.cpu_stat" group by node`, NewNodeIDAllocator())
	assert.NoError(t, err)
}

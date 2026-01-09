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

package planner

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/optimization"
	"github.com/lindb/lindb/sql/planner/plan"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	printpkg "github.com/lindb/lindb/sql/planner/printer"
	"github.com/lindb/lindb/sql/tree"
)

type LogicalPlanner struct {
	context        *context.PlannerContext
	planOptimizers []optimization.PlanOptimizer
}

func NewLogicalPlanner(ctx *context.PlannerContext, planOptimizers []optimization.PlanOptimizer) *LogicalPlanner {
	return &LogicalPlanner{
		context:        ctx,
		planOptimizers: planOptimizers,
	}
}

func (p *LogicalPlanner) Plan() *planpkg.Plan {
	// plan
	root := p.planStatement()

	printer := printpkg.NewPlanPrinter(printpkg.NewTextRender(0))
	fmt.Printf("init plan:\n%s\n", printer.PrintLogicPlan(root))

	// TODO: check intermediate plan

	// optimizer
	for _, optimizer := range p.planOptimizers {
		root = p.runOptimizer(root, optimizer)
		printer = printpkg.NewPlanPrinter(printpkg.NewTextRender(0))
		fmt.Printf("after optimizer plan:%T\n%s\n", optimizer, printer.PrintLogicPlan(root))
	}
	printer = printpkg.NewPlanPrinter(printpkg.NewTextRender(0))
	fmt.Printf("after op plan:\n%s\n", printer.PrintLogicPlan(root))

	return &planpkg.Plan{
		Root: root,
	}
}

func (p *LogicalPlanner) planStatement() planpkg.PlanNode {
	relationPlan := p.planStatementWithoutOutput()
	return p.createOutputPlan(relationPlan)
}

func (p *LogicalPlanner) planStatementWithoutOutput() *RelationPlan {
	statement := p.context.AnalyzerContext.Analysis.GetStatement()
	fmt.Printf("statement type=%T\n", statement)
	switch stmt := statement.(type) {
	case *tree.Query:
		planner := NewRelationPlanner(p.context, nil, nil, nil)
		return stmt.Accept(nil, planner).(*RelationPlan)
	case *tree.Insert:
		return p.createInsertPlan(stmt)
	default:
		// TODO: plan other statement
		panic("not support statement type")
	}
}

func (p *LogicalPlanner) createInsertPlan(statement *tree.Insert) *RelationPlan {
	insert := p.context.AnalyzerContext.Analysis.GetInsert()
	planner := NewRelationPlanner(p.context, nil, nil, nil)
	queryPlan := planner.Visit(nil, statement.Query).(*RelationPlan)
	outputDescriptor := p.context.AnalyzerContext.Analysis.GetOutputDescriptor(statement.Query)
	var (
		columns []string
		outputs []*planpkg.Symbol
	)
	for i := range outputDescriptor.Fields {
		field := outputDescriptor.Fields[i]
		if field.Hidden {
			// ignore hidden column
			continue
		}
		name := field.Name
		if name == "" {
			name = fmt.Sprintf("_col%d", i)
		}
		fieldIdx := outputDescriptor.IndexOf(field)
		fmt.Printf("find field index=%v\n", fieldIdx)
		outputs = append(outputs, queryPlan.getSymbol(fieldIdx))
		columns = append(columns, name)
	}
	fmt.Printf("create insert plan output descriptor=%v,%v\n", outputDescriptor, outputs)

	project := &plan.ProjectionNode{
		BaseNode: plan.BaseNode{
			ID: p.context.PlanNodeIDAllocator.Next(),
		},
		Source: queryPlan.Root,
		Assignments: lo.Map(outputs, func(item *plan.Symbol, index int) *plan.Assignment {
			return &plan.Assignment{
				Symbol:     &plan.Symbol{Name: columns[index], DataType: item.DataType, AggType: item.AggType},
				Expression: item.ToSymbolReference(),
			}
		}),
	}
	fmt.Printf("create insert plan project=%v\n", project)

	return &RelationPlan{
		Scope: p.context.AnalyzerContext.Analysis.GetScope(statement),
		Root: &planpkg.InsertNode{
			Database: p.context.Database,
			Table:    insert.Table,
			Source:   project,
		},
	}
}

func (p *LogicalPlanner) createOutputPlan(plan *RelationPlan) planpkg.PlanNode {
	var (
		columns []string
		outputs []*planpkg.Symbol
	)
	analysis := p.context.AnalyzerContext.Analysis
	outputDescriptor := analysis.GetOutputDescriptor(analysis.GetRoot())
	for i := range outputDescriptor.Fields {
		field := outputDescriptor.Fields[i]
		if field.Hidden {
			// ignore hidden column
			continue
		}
		name := field.Name
		if name == "" {
			name = fmt.Sprintf("_col%d", i)
		}
		columns = append(columns, name)
		fieldIdx := outputDescriptor.IndexOf(field)
		fmt.Printf("find field index=%v\n", fieldIdx)
		outputs = append(outputs, plan.getSymbol(fieldIdx))
	}
	fmt.Printf("create output columns=%v,%v\n", columns, outputs)
	return &planpkg.OutputNode{
		BaseNode: planpkg.BaseNode{
			ID: p.context.PlanNodeIDAllocator.Next(),
		},
		Source:      plan.Root,
		ColumnNames: columns,
		Outputs:     outputs,
	}
}

func (p *LogicalPlanner) runOptimizer(root planpkg.PlanNode, optimizer optimization.PlanOptimizer) (result planpkg.PlanNode) {
	// FIXME:
	return optimizer.Optimize(p.context, root)
}

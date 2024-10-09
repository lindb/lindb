package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/optimization"
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
	fmt.Printf("after plan:\n%s\n", printer.PrintLogicPlan(root))

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
		planner := NewRelationPlanner(p.context, nil)
		return stmt.Accept(nil, planner).(*RelationPlan)
	default:
		// TODO: plan other statement
		panic("not support statement type")
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
		name := field.Name
		if name == "" {
			name = fmt.Sprintf("_col%d", i)
		}
		columns = append(columns, name)
		fieldIdx := outputDescriptor.IndexOf(field)
		outputs = append(outputs, plan.getSymbol(fieldIdx))
	}

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
	result = optimizer.Optimize(p.context, root)
	return
}

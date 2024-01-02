package planner

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type GroupingSetsPlan struct {
	subPlan      *PlanBuilder
	groupingSets [][]*plan.Symbol
}

type QueryPlanner struct {
	context *context.PlannerContext

	subQueryPlanner *SubQueryPlanner
}

func NewQueryPlanner(context *context.PlannerContext) *QueryPlanner {
	return &QueryPlanner{
		context:         context,
		subQueryPlanner: NewSubQueryPlanner(context),
	}
}

func (p *QueryPlanner) planQuery(node *tree.Query) *RelationPlan {
	builder := p.planQueryBody(node)

	selectExpressions := p.context.AnalyzerContext.Analysis.GetSelectExpressions(node)
	var outputs []tree.Expression
	for i := range selectExpressions {
		outputs = append(outputs, selectExpressions[i].Expression)
	}
	orderBy := p.context.AnalyzerContext.Analysis.GetOrderByExpressions(node)
	var orderByAndOutputs []tree.Expression
	orderByAndOutputs = append(orderByAndOutputs, orderBy...)
	orderByAndOutputs = append(orderByAndOutputs, outputs...)
	builder = builder.appendProjections(orderByAndOutputs)
	// FIXME:>>>>> ort/limit

	builder = builder.appendProjections(outputs)

	return &RelationPlan{
		Root:          builder.root,
		Scope:         p.context.AnalyzerContext.Analysis.GetScope(node),
		FieldMappings: p.computeOutputs(builder, outputs),
	}
}

func (p *QueryPlanner) planQueryBody(query *tree.Query) *PlanBuilder {
	planner := NewRelationPlanner(p.context)
	relationPlan := query.QueryBody.Accept(nil, planner).(*RelationPlan)
	return newPlanBuilder(p.context, relationPlan, nil)
}

func (p *QueryPlanner) planQuerySpecification(node *tree.QuerySpecification) *RelationPlan {
	// from clause
	builder := p.planFrom(node)
	builder = p.aggregate(builder, node)
	// TODO: agg
	// where clause
	builder = p.filter(builder, p.context.AnalyzerContext.Analysis.GetWhere(node), node)

	selectExpressions := p.context.AnalyzerContext.Analysis.GetSelectExpressions(node)
	outputs := p.outputExpressions(selectExpressions)
	return &RelationPlan{
		Root:          builder.root,
		Scope:         p.context.AnalyzerContext.Analysis.GetScope(node),
		FieldMappings: p.computeOutputs(builder, outputs),
	}
}

func (p *QueryPlanner) planFrom(node *tree.QuerySpecification) *PlanBuilder {
	planner := NewRelationPlanner(p.context)
	relationPlan := node.From.Accept(nil, planner).(*RelationPlan)

	return newPlanBuilder(p.context, relationPlan, nil)
}

func (p *QueryPlanner) aggregate(subPlan *PlanBuilder, node *tree.QuerySpecification) *PlanBuilder {
	if !p.context.AnalyzerContext.Analysis.IsAggregation(node) {
		return subPlan
	}
	// TODO: aggregates
	groupingSetAnalysis := p.context.AnalyzerContext.Analysis.GetGroupingSets(node)
	groupingSets := p.planGroupingSets(subPlan, node, groupingSetAnalysis)
	// TODO: group agg
	subPlan = p.planAggregation(groupingSets.subPlan)

	return p.planGroupingOperations(subPlan, node)
}

func (p *QueryPlanner) planGroupingSets(subPlan *PlanBuilder, node *tree.QuerySpecification, groupingSetAnalysis *analyzer.GroupingSetAnalysis) *GroupingSetsPlan {
	groupingSetMappings := make(map[*plan.Symbol]*plan.Symbol) // ouput -> input
	complexExpressions := make(map[tree.NodeID]*plan.Symbol)
	fields := make([]*plan.Symbol, len(subPlan.translations.fieldSymbols))
	for _, field := range groupingSetAnalysis.GetAllFields() {
		input := subPlan.translations.fieldSymbols[field.FieldIndex]
		output := p.context.SymbolAllocator.FromSymbol(input, "gid")
		fields[field.FieldIndex] = output
		groupingSetMappings[output] = input
	}

	for _, expression := range groupingSetAnalysis.GetComplexExpressions() {
		if _, ok := complexExpressions[expression.GetID()]; !ok {
			input := subPlan.translate(expression)
			output := p.context.SymbolAllocator.NewSymbol(expression, "gid")
			complexExpressions[expression.GetID()] = output
			groupingSetMappings[output] = input
		}
	}

	groupingSets := [][]*plan.Symbol{lo.Values(complexExpressions)}
	assignments := make(plan.Assignments)
	assignments.Add(subPlan.root.GetOutputSymbols())
	for k, v := range groupingSetMappings {
		assignments[k] = v.ToSymbolReference()
	}
	groupID := &plan.ProjectionNode{
		BaseNode: plan.BaseNode{
			ID: p.context.PlanNodeIDAllocator.Next(),
		},
		Source:      subPlan.root,
		Assignments: assignments,
	}
	subPlan = &PlanBuilder{
		root:         groupID,
		translations: subPlan.translations.withNewMappings(complexExpressions, fields),
	}
	return &GroupingSetsPlan{
		subPlan:      subPlan,
		groupingSets: groupingSets,
	}
}

func (p *QueryPlanner) planGroupingOperations(subPlan *PlanBuilder, node *tree.QuerySpecification) *PlanBuilder {
	return subPlan
}

func (p *QueryPlanner) planAggregation(subPlan *PlanBuilder) *PlanBuilder {
	return subPlan
}

func (p *QueryPlanner) filter(subPlan *PlanBuilder, predicate tree.Expression, node tree.Node) *PlanBuilder {
	if predicate == nil {
		return subPlan
	}
	subPlan = p.subQueryPlanner.handleSubQueries(subPlan, predicate, nil)

	return subPlan.withNewRoot(&plan.FilterNode{
		BaseNode: plan.BaseNode{
			ID: p.context.PlanNodeIDAllocator.Next(),
		},
		Source:    subPlan.root,
		Predicate: coerceIfNecessary(predicate, predicate), // FIXME:::
		// TODO:
	})
}

func coerceIfNecessary(original, rewritten tree.Expression) tree.Expression {
	// FIXME::
	return &tree.Cast{
		Expression: rewritten,
	}
}

func (p *QueryPlanner) computeOutputs(builder *PlanBuilder, outputs []tree.Expression) (outputSymbols []*plan.Symbol) {
	for _, expression := range outputs {
		outputSymbols = append(outputSymbols, builder.translate(expression))
	}
	return
}

func (p *QueryPlanner) outputExpressions(selectExpressions []*analyzer.SelectExpression) (outputs []tree.Expression) {
	// TODO: fixme unfolded express
	for i := range selectExpressions {
		outputs = append(outputs, selectExpressions[i].Expression)
	}
	return
}

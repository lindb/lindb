package planner

import (
	"fmt"

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
	context      *context.PlannerContext
	outerContext *TranslationMap

	subQueryPlanner *SubQueryPlanner
}

func NewQueryPlanner(context *context.PlannerContext, outerContext *TranslationMap) *QueryPlanner {
	return &QueryPlanner{
		context:         context,
		outerContext:    outerContext,
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
	// FIXME:>>>>> order/limit

	builder = builder.appendProjections(outputs)
	fmt.Printf("query plan outputs==%v\n", outputs)

	return &RelationPlan{
		Root:          builder.root,
		Scope:         p.context.AnalyzerContext.Analysis.GetScope(node),
		FieldMappings: p.computeOutputs(builder, outputs),
	}
}

func (p *QueryPlanner) planQueryBody(query *tree.Query) *PlanBuilder {
	planner := NewRelationPlanner(p.context, p.outerContext, nil)
	relationPlan := query.QueryBody.Accept(nil, planner).(*RelationPlan)
	return newPlanBuilder(p.context, relationPlan, nil)
}

func (p *QueryPlanner) planQuerySpecification(node *tree.QuerySpecification) *RelationPlan {
	// from clause
	builder := p.planFrom(node)
	// where clause
	builder = p.filter(builder, p.context.AnalyzerContext.Analysis.GetWhere(node), node)
	// agg/group by
	builder = p.aggregate(builder, node)
	// TODO: having
	// TODO: sub query

	selectExpressions := p.context.AnalyzerContext.Analysis.GetSelectExpressions(node)
	outputs := p.outputExpressions(selectExpressions)
	// TODO: sort/order by

	builder = builder.appendProjections(outputs)
	return &RelationPlan{
		Root:          builder.root,
		Scope:         p.context.AnalyzerContext.Analysis.GetScope(node),
		FieldMappings: p.computeOutputs(builder, outputs),
	}
}

func (p *QueryPlanner) planFrom(node *tree.QuerySpecification) *PlanBuilder {
	if node.From != nil {
		planner := NewRelationPlanner(p.context, p.outerContext, p.context.AnalyzerContext.Analysis.GetTimePredicates(node))
		relationPlan := node.From.Accept(nil, planner).(*RelationPlan)
		return newPlanBuilder(p.context, relationPlan, nil)
	}
	return &PlanBuilder{
		root: &plan.ValuesNode{
			BaseNode: plan.BaseNode{
				ID: p.context.PlanNodeIDAllocator.Next(),
			},
			RowCount: 1,
		},
		translations: &TranslationMap{
			scope:   p.context.AnalyzerContext.Analysis.GetImplicitFromScope(node),
			context: p.context,
		},
	}
}

func (p *QueryPlanner) aggregate(subPlan *PlanBuilder, node *tree.QuerySpecification) *PlanBuilder {
	if !p.context.AnalyzerContext.Analysis.IsGroupingSets(node) {
		// check if has group by sets
		return subPlan
	}
	// TODO: aggregates
	groupingSetAnalysis := p.context.AnalyzerContext.Analysis.GetGroupingSets(node)
	groupingSets := p.planGroupingSets(subPlan, node, groupingSetAnalysis)
	// TODO: group agg
	subPlan = p.planAggregation(groupingSets.subPlan, groupingSets.groupingSets, p.context.AnalyzerContext.Analysis.GetAggregates(node))

	return p.planGroupingOperations(subPlan, node)
}

func (p *QueryPlanner) planGroupingSets(subPlan *PlanBuilder, node *tree.QuerySpecification, groupingSetAnalysis *analyzer.GroupingSetAnalysis) *GroupingSetsPlan {
	groupingSetMappings := make(map[*plan.Symbol]*plan.Symbol) // ouput -> input
	complexExpressions := make(map[tree.NodeID]*plan.Symbol)
	fields := make([]*plan.Symbol, len(subPlan.translations.fieldSymbols))
	fmt.Printf("sub plan fields=%v\n", subPlan.translations.fieldSymbols)
	// TODO: remove it?
	copy(fields, subPlan.translations.fieldSymbols)
	fmt.Printf("plan grouping sets:%v\n", len(subPlan.translations.fieldSymbols))
	for _, field := range groupingSetAnalysis.GetAllFields() {
		input := subPlan.translations.fieldSymbols[field.FieldIndex]
		// add group field suffix
		// FIXME: add gid for symbol suffix
		output := p.context.SymbolAllocator.FromSymbol(input, "", input.DataType)
		fields[field.FieldIndex] = output
		groupingSetMappings[output] = input
	}

	for _, expression := range groupingSetAnalysis.GetComplexExpressions() {
		if _, ok := complexExpressions[expression.GetID()]; !ok {
			input := subPlan.translate(expression)
			// FIXME: add gid for symbol suffix
			output := p.context.SymbolAllocator.NewSymbol(expression, "", p.context.AnalyzerContext.Analysis.GetType(expression))
			complexExpressions[expression.GetID()] = output
			groupingSetMappings[output] = input
		}
	}
	columnOnlyGroupingSets := p.enumerateGroupingSets(groupingSetAnalysis)
	groupingSets := [][]*plan.Symbol{lo.Values(complexExpressions)}
	for _, gs := range columnOnlyGroupingSets {
		groupingSets = append(groupingSets, lo.Map(gs, func(item *analyzer.FieldID, index int) *plan.Symbol {
			return fields[item.FieldIndex]
		}))
	}

	var assignments plan.Assignments
	assignments = assignments.Add(subPlan.root.GetOutputSymbols())
	for k, v := range groupingSetMappings {
		assignments = assignments.Put(k, v.ToSymbolReference())
	}

	groupID := &plan.ProjectionNode{
		BaseNode: plan.BaseNode{
			ID: p.context.PlanNodeIDAllocator.Next(),
		},
		Source:      subPlan.root,
		Assignments: assignments.Unique(),
	}
	fmt.Printf("plan agg group... fields=%v\n", fields)
	subPlan = &PlanBuilder{
		root:         groupID,
		translations: subPlan.translations.withNewMappings(complexExpressions, fields),
	}
	return &GroupingSetsPlan{
		subPlan:      subPlan,
		groupingSets: groupingSets,
	}
}

func (p *QueryPlanner) enumerateGroupingSets(groupingSetAnalysis *analyzer.GroupingSetAnalysis) [][]*analyzer.FieldID {
	// TODO: add cube/rollup?
	var partialSet [][]*analyzer.FieldID

	partialSet = append(partialSet, groupingSetAnalysis.GetOrdinarySets()...)
	if len(partialSet) == 0 {
		return nil
	}
	// TODO: compute the cross product of the partial sets
	return partialSet
}

func (p *QueryPlanner) planGroupingOperations(subPlan *PlanBuilder, node *tree.QuerySpecification) *PlanBuilder {
	return subPlan
}

func (p *QueryPlanner) planAggregation(subPlan *PlanBuilder, groupingSets [][]*plan.Symbol, aggregates []*tree.FunctionCall) *PlanBuilder {
	fmt.Printf("planagg.....%v,func call=%v\n", groupingSets, aggregates)

	var aggregateMapping []*plan.AggregationAssignment
	additionalMapping := make(map[tree.NodeID]*plan.Symbol)
	// TODO: scopeAwareDistinct
	for _, function := range aggregates {
		symbol := p.context.SymbolAllocator.NewSymbol(function, "", p.context.AnalyzerContext.Analysis.GetType(function))
		aggregation := &plan.Aggregation{
			Function: p.context.AnalyzerContext.Analysis.GetResolvedFunction(function),
			Arguments: lo.Map(function.Arguments, func(arg tree.Expression, _ int) tree.Expression {
				if iden, ok := arg.(*tree.Identifier); ok {
					return p.context.SymbolAllocator.NewSymbol(iden, "", p.context.AnalyzerContext.Analysis.GetType(iden)).ToSymbolReference()
				}
				return arg
			}), // TODO: parse arg
		}
		aggregateMapping = append(aggregateMapping, &plan.AggregationAssignment{
			Symbol:        symbol,
			Aggregation:   aggregation,
			ASTExpression: function,
		})
		additionalMapping[function.GetID()] = symbol
	}
	groupingKeys := make(map[string]*plan.Symbol)
	for _, symbol := range lo.Flatten(groupingSets) {
		groupingKeys[symbol.Name] = symbol
	}
	aggregationNode := plan.NewAggregationNode(
		p.context.PlanNodeIDAllocator.Next(),
		subPlan.root,
		aggregateMapping,
		&plan.GroupingSetDescriptor{
			GroupingKeys: lo.Values(groupingKeys),
		},
		plan.SINGLE)
	return &PlanBuilder{
		root: aggregationNode,
		// add aggregation symbol(function call node=>symbol)
		translations: subPlan.translations.withAdditionalMapping(additionalMapping),
	}
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
		Predicate: coerceIfNecessary(p.context.AnalyzerContext.Analysis, predicate, subPlan.rewrite(predicate)),
		// TODO:
	})
}

func (p *QueryPlanner) computeOutputs(builder *PlanBuilder, outputs []tree.Expression) (outputSymbols []*plan.Symbol) {
	for _, expression := range outputs {
		fmt.Printf("output exp=%v,%T\n", expression, expression)
		outputSymbols = append(outputSymbols, builder.translate(expression))
	}
	fmt.Printf("output result==%v\n", outputSymbols)
	return
}

func (p *QueryPlanner) outputExpressions(selectExpressions []*analyzer.SelectExpression) (outputs []tree.Expression) {
	fmt.Printf("select output expression==%v\n", selectExpressions)
	// TODO: fixme unfolded express
	for i := range selectExpressions {
		fmt.Println(selectExpressions[i].Expression)
		outputs = append(outputs, selectExpressions[i].Expression)
	}
	return
}

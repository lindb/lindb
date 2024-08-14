package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PrintPlanVisitor struct {
	representation *PlanRepresentation
}

func NewVisitor(representation *PlanRepresentation) plan.Visitor {
	return &PrintPlanVisitor{
		representation: representation,
	}
}

func (v *PrintPlanVisitor) Visit(_ any, n plan.PlanNode) (r any) {
	switch node := n.(type) {
	case *plan.OutputNode:
		v.visitOutput(node)
	case *plan.JoinNode:
		v.visitJoin(node)
	case *plan.FilterNode:
		v.visitScanFilterAndProjection(node, node, nil)
	case *plan.ProjectionNode:
		v.visitProjection(node)
	case *plan.ExchangeNode:
		v.visitExchange(node)
	case *plan.RemoteSourceNode:
		v.visitRemoteSource(node)
	case *plan.TableScanNode:
		v.visitTableScan(node)
	case *plan.AggregationNode:
		v.visitAggregation(node)
	case *plan.GroupReference:
		v.visitGroupReference(node)
	default:
		panic(fmt.Sprintf("impl print %T", n))
	}
	return
}

func (v *PrintPlanVisitor) visitGroupReference(node *plan.GroupReference) {
	v.addNode(node, "GroupReference", map[string]string{"groupId": strconv.Itoa(int(node.GetNodeID()))}, nil)
}

func (v *PrintPlanVisitor) visitAggregation(node *plan.AggregationNode) {
	// TODO: add type
	descriptor := make(map[string]string)
	if node.Step != plan.SINGLE {
		descriptor["type"] = node.Step.String()
	}

	if node.GroupingSets != nil && len(node.GroupingSets.GroupingKeys) > 0 {
		descriptor["keys"] = formatSymbols(node.GroupingSets.GroupingKeys)
	}

	nodeOutput := v.addNode(node, "Aggregate", descriptor, node.GetSources())
	for _, agg := range node.Aggregations {
		nodeOutput.appendDetails(fmt.Sprintf("%s := %s", agg.Symbol.Name, formatAggregation(agg.Aggregation)))
	}
	// TODO: add agg func
	v.processChildren(node)
}

func (v *PrintPlanVisitor) visitJoin(node *plan.JoinNode) {
	var criteriaExpressions []tree.Expression
	for _, criteria := range node.Criteria {
		criteriaExpressions = append(criteriaExpressions, criteria.ToExpression()) // FIXME:
	}
	if node.IsCrossJoin() {
		panic("impl it")
	} else {
		descriptor := make(map[string]string)
		descriptor["criteria"] = strings.Join(v.anonymizeExpressions(criteriaExpressions), " AND ")
		v.addNode(node, string(node.Type), descriptor, node.GetSources())
	}

	node.Left.Accept(nil, v)
	node.Right.Accept(nil, v)
}

func (v *PrintPlanVisitor) visitTableScan(node *plan.TableScanNode) {
	descriptor := map[string]string{"table": node.Table.String()}
	outputNode := v.addNode(node, "TableScan", descriptor, node.GetSources())
	partitions := []string{}
	for node, shards := range node.Partitions {
		partitions = append(partitions, fmt.Sprintf("%s=%s", node.Address(), strings.Join(strings.Fields(fmt.Sprint(shards)), ",")))
	}
	if len(partitions) > 0 {
		outputNode.appendDetails("Partitions: [" + strings.Join(partitions, ", ") + "]")
	}
	v.printTableScanInfo(outputNode, node)
}

func (v *PrintPlanVisitor) visitOutput(node *plan.OutputNode) {
	outputNode := v.addNode(node, "Output",
		map[string]string{"columnNames": "[" + strings.Join(node.ColumnNames, ", ") + "]"},
		node.GetSources())
	for i := range node.ColumnNames {
		name := node.ColumnNames[i]
		symbol := node.Outputs[i]
		if name != symbol.Name {
			outputNode.appendDetails(fmt.Sprintf("%s := %s", name, symbol.Name))
		}
	}
	v.processChildren(node)
}

func (v *PrintPlanVisitor) visitProjection(node *plan.ProjectionNode) {
	if source, ok := node.Source.(*plan.FilterNode); ok {
		v.visitScanFilterAndProjection(node, source, node)
		return
	}
	v.visitScanFilterAndProjection(node, nil, node)
}

func (v *PrintPlanVisitor) visitExchange(node *plan.ExchangeNode) {
	descriptor := make(map[string]string)
	descriptor["type"] = string(node.Type)
	// TODO: add type check
	v.addNode(node, fmt.Sprintf("%sExchange", node.Scope), descriptor, node.GetSources())
	v.processChildren(node)
}

func (v *PrintPlanVisitor) visitRemoteSource(node *plan.RemoteSourceNode) {
	descriptor := make(map[string]string)
	strs := make([]string, len(node.SourceFragmentIDs))
	for i, v := range node.SourceFragmentIDs {
		strs[i] = fmt.Sprintf("%d", v)
	}
	descriptor["sourceFragmentIds"] = fmt.Sprintf("[%s]", strings.Join(strs, ", "))
	v.addNode(node, "Remote", descriptor, node.GetSources())
}

func (v *PrintPlanVisitor) processChildren(node plan.PlanNode) {
	sources := node.GetSources()
	for i := range sources {
		child := sources[i]
		_ = child.Accept(nil, v)
	}
}

func (v *PrintPlanVisitor) visitScanFilterAndProjection(node plan.PlanNode, filter *plan.FilterNode, projection *plan.ProjectionNode) {
	var sourceNode plan.PlanNode
	if filter != nil {
		sourceNode = filter.Source
	} else {
		sourceNode = projection.Source
	}

	var scanNode *plan.TableScanNode
	if tableScanNode, ok := sourceNode.(*plan.TableScanNode); ok {
		scanNode = tableScanNode
	}

	var operatorName string
	descriptor := make(map[string]string)
	if scanNode != nil {
		operatorName += "Scan"
		descriptor["table"] = scanNode.Table.String()
	}
	if filter != nil {
		operatorName += "Filter"
		if filter.Predicate != nil {
			// conjuncts := analyzer.ExtractConjuncts(filter.Predicate)
			descriptor["filterPredicate"] = tree.FormatExpression(filter.Predicate)
		}
		// FIXME:
	}
	if projection != nil {
		operatorName += "Projection"
	}

	outputNode := v.addNode(node, operatorName, descriptor, []plan.PlanNode{sourceNode})

	if projection != nil {
		// print assignments
		for _, assignment := range projection.Assignments {
			expression := assignment.Expression
			symbol := assignment.Symbol
			if symboleRef, ok := expression.(*tree.SymbolReference); ok && symboleRef.Name == symbol.Name {
				// skip identity assignments
				continue
			}
			outputNode.appendDetails(fmt.Sprintf("%s := %s", symbol.Name, tree.FormatExpression(expression)))
		}
	}

	if scanNode != nil {
		return
	}

	sourceNode.Accept(nil, v)
}

func (v *PrintPlanVisitor) addNode(node plan.PlanNode,
	name string, descriptor map[string]string,
	children []plan.PlanNode,
) *NodeRepresentation {
	var childrenIDs []plan.PlanNodeID
	for _, child := range children {
		childrenIDs = append(childrenIDs, child.GetNodeID())
	}

	outputNode := &NodeRepresentation{
		id:         node.GetNodeID(),
		name:       name,
		descriptor: descriptor,
		children:   childrenIDs,
		outputs:    node.GetOutputSymbols(),
	}

	v.representation.addNode(outputNode)
	return outputNode
}

func (v *PrintPlanVisitor) printTableScanInfo(outputNode *NodeRepresentation, node *plan.TableScanNode) {
}

func (v *PrintPlanVisitor) anonymizeExpressions(expressions []tree.Expression) (result []string) {
	for i := range expressions {
		result = append(result, tree.FormatExpression(expressions[i]))
	}
	return
}

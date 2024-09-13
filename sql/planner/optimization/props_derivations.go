package optimization

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

func deriveProps(node plan.PlanNode, inputProperties []*ActualProps) *ActualProps {
	output := node.Accept(inputProperties, &PropertyDerivationVisitor{})
	// TODO: add verify
	return output.(*ActualProps)
}

func filterIfMissing(columns []*plan.Symbol, column *plan.Symbol) *plan.Symbol {
	if lo.ContainsBy(columns, func(item *plan.Symbol) bool {
		return item.Name == column.Name
	}) {
		return column
	}
	return nil
}

func exchangeInputToOutput(node *plan.ExchangeNode, sourceIndex int) map[*plan.Symbol]*plan.Symbol {
	return nil
}

type PropertyDerivationVisitor struct{}

func (v *PropertyDerivationVisitor) Visit(context any, n plan.PlanNode) (r any) {
	inputProperties := context.([]*ActualProps)
	switch node := n.(type) {
	case *plan.OutputNode:
		return v.visitOutput(inputProperties, node)
	case *plan.AggregationNode:
		return v.visitAggregation(inputProperties, node)
	case *plan.ExchangeNode:
		return v.visitExchangeNode(inputProperties, node)
	case *plan.ProjectionNode:
		return v.visitProjection(inputProperties, node)
	case *plan.FilterNode:
		return v.visitFilter(inputProperties, node)
	case *plan.TableScanNode:
		return v.visitTableScan(inputProperties, node)
	default:
		panic(fmt.Sprintf("impl prop derivation visitor %T", n))
	}
}

func (v *PropertyDerivationVisitor) visitFilter(inputProperties []*ActualProps, node *plan.FilterNode) *ActualProps {
	props := inputProperties[0]
	// identities := computeIdentityTranslations(node.Assignments)
	// // TODO: expression rewrite
	// translated := props.translate(func(symbol *plan.Symbol) *plan.Symbol {
	// 	return identities[symbol.Name]
	// })
	// FIXME: add constant
	return BuilderFrom(props).Build()
}

func (v *PropertyDerivationVisitor) visitProjection(inputProperties []*ActualProps, node *plan.ProjectionNode) *ActualProps {
	props := inputProperties[0]
	identities := computeIdentityTranslations(node.Assignments)
	// TODO: expression rewrite
	translated := props.translate(func(symbol *plan.Symbol) *plan.Symbol {
		return identities[symbol.Name]
	})
	// FIXME: add constant
	return BuilderFrom(translated).Build()
}

func (v *PropertyDerivationVisitor) visitOutput(inputProperties []*ActualProps, node *plan.OutputNode) *ActualProps {
	return inputProperties[0].translate(func(column *plan.Symbol) *plan.Symbol {
		return filterIfMissing(node.GetOutputSymbols(), column)
	})
}

func (v *PropertyDerivationVisitor) visitAggregation(inputProperties []*ActualProps, node *plan.AggregationNode) *ActualProps {
	props := inputProperties[0]
	translated := props.translate(func(symbol *plan.Symbol) *plan.Symbol {
		if lo.ContainsBy(node.GetGroupingKeys(), func(item *plan.Symbol) bool {
			return item.Name == symbol.Name
		}) {
			return symbol
		}
		return nil
	})
	return BuilderFrom(translated).Build()
}

func (v *PropertyDerivationVisitor) visitTableScan(inputProperties []*ActualProps, node *plan.TableScanNode) *ActualProps {
	props := NewActualPropsBuilder(arbitraryPartition())
	return props.Build()
}

func (v *PropertyDerivationVisitor) visitExchangeNode(inputProperties []*ActualProps, node *plan.ExchangeNode) *ActualProps {
	// TODO: check
	switch node.Type {
	case plan.Gather:
		// TODO: check coord
		return NewActualPropsBuilder(singlePartition()).Build()
	case plan.Repartition:
		return NewActualPropsBuilder(partitionedOn(&plan.Partitioning{})).Build()
	default:
		panic("unknonw exchange type")
	}
}

func computeIdentityTranslations(assigments plan.Assignments) map[string]*plan.Symbol {
	inputToOutput := make(map[string]*plan.Symbol)
	for _, assignment := range assigments {
		if symbolRef, ok := assignment.Expression.(*tree.SymbolReference); ok {
			inputToOutput[plan.SymbolFrom(symbolRef).Name] = assignment.Symbol
		}
	}
	return inputToOutput
}

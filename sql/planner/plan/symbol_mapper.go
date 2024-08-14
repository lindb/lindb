package plan

type SymbolMap struct {
	From *Symbol
	To   *Symbol
}

type SymbolMapper struct {
	fn func(symbol *Symbol) *Symbol
}

func NewSymbolMapper(mapping map[string]*Symbol) *SymbolMapper {
	return &SymbolMapper{
		fn: func(symbol *Symbol) *Symbol {
			for {
				val, ok := mapping[symbol.Name]
				if ok && val.Name != symbol.Name {
					symbol = val
				} else {
					break
				}
			}
			return symbol
		},
	}
}

func (m *SymbolMapper) MapAggregation(node *AggregationNode, source PlanNode, newNodeID PlanNodeID) *AggregationNode {
	return NewAggregationNode(newNodeID, source, node.Aggregations, node.GroupingSets, node.Step)
}

func (m *SymbolMapper) MapSymbol(symobl *Symbol) *Symbol {
	return m.fn(symobl)
}

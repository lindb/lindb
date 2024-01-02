package tree

type GroupBy struct {
	BaseNode
	GroupingElements []GroupingElement
}

type GroupingElement interface {
}

type SimpleGroupBy struct {
	BaseNode
	Columns []Expression
}

type GroupByAllColumns struct {
}

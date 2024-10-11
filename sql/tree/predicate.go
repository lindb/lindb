package tree

type LikePredicate struct {
	BaseNode
	Value   Expression
	Pattern Expression
}

func (n *LikePredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type InPredicate struct {
	BaseNode
	Value     Expression
	ValueList Expression
}

func (n *InPredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type RegexPredicate struct {
	BaseNode
	Value   Expression
	Pattern Expression
}

func (n *RegexPredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type TimePredicate struct {
	BaseNode
	Operator ComparisonOperator
	Value    Expression
}

func (n *TimePredicate) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

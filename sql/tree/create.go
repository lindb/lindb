package tree

type CreateDatabase struct {
	BaseNode
	Name    string
	Options map[string]any
	Rollup  []RollupOption
}

func (n *CreateDatabase) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type CreateBroker struct {
	BaseNode
	Options map[string]any
	Name    string
}

func (n *CreateBroker) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type RollupOption struct {
	Options map[string]any
}

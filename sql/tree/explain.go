package tree

const (
	LogicalExplain     string = "LOGICAL"
	DistributedExplain string = "DISTRIBUTED"
)

type ExplainOption interface{}

type ExplainType struct {
	Type string // default is LOGICAL
}

type Explain struct {
	BaseNode
	Statement Statement
	Options   []ExplainOption
}

func (n *Explain) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

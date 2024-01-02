package tree

type Identifier struct {
	BaseNode

	Value     string `json:"value"`
	Delimited bool   `json:"delimited"`
}

func (n *Identifier) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

package tree

type StreamingApp struct {
	BaseNode

	Annotations []*Annotation
	Statements  []*StreamingStatement
}

func (n *StreamingApp) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type StreamingStatement struct {
	BaseNode

	Annotations []*Annotation
	Statement   Statement
}

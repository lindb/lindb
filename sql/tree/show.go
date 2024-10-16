package tree

type ShowNamespaces struct {
	BaseNode

	LikePattern string
}

func (n *ShowNamespaces) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowTableNames struct {
	BaseNode

	Namespace   *QualifiedName
	LikePattern string
}

func (n ShowTableNames) GetNamespace() string {
	if n.Namespace == nil {
		return ""
	}
	return n.Namespace.Name
}

func (n *ShowTableNames) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowColumns struct {
	BaseNode

	Table *Table
}

func (n *ShowColumns) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

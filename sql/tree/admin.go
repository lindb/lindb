package tree

// FlushDatabase represents the statement that flush database's memory database.
type FlushDatabase struct {
	BaseNode
	Database string
}

// Accept implements Statement
func (n *FlushDatabase) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

// CompactDatabase represents the statement that compact database's files.
type CompactDatabase struct {
	BaseNode
	Database string
}

func (n *CompactDatabase) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

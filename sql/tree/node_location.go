package tree

type NodeLocation struct {
	Line   int
	Column int
}

func NewNodeLocation(line, column int) *NodeLocation {
	return &NodeLocation{
		Line:   line,
		Column: column,
	}
}

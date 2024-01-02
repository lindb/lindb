package tree

// Use represents the statement that select search database.
type Use struct {
	BaseNode
	Database *Identifier
}

package execution

import (
	"context"

	"github.com/lindb/lindb/sql/tree"
)

// Session represents query session.
type Session struct {
	RequestID RequestID
	Context   context.Context
	Database  string

	Statement       *tree.PreparedStatement
	NodeIDAllocator *tree.NodeIDAllocator
}

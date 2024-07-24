package execution

import (
	"context"

	"github.com/lindb/lindb/execution/model"
	"github.com/lindb/lindb/sql/tree"
)

// Session represents query session.
type Session struct {
	RequestID model.RequestID
	Context   context.Context
	Database  string

	Statement       *tree.PreparedStatement
	NodeIDAllocator *tree.NodeIDAllocator
}

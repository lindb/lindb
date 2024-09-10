package execution

import (
	"context"

	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/tree"
)

// Session represents query language execution session.
type Session struct {
	Context  context.Context
	Database string

	Statement       *tree.PreparedStatement
	NodeIDAllocator *tree.NodeIDAllocator

	RequestID model.RequestID
}

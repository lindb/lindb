package rewrite

import "github.com/lindb/lindb/sql/tree"

type Rewrite interface {
	Rewrite(statement tree.Statement) tree.Statement
}

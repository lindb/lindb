package rewrite

import "github.com/lindb/lindb/sql/tree"

type StatementRewrite struct {
	rewrites []Rewrite
}

func NewStatementRewrite(rewrites []Rewrite) *StatementRewrite {
	return &StatementRewrite{
		rewrites: rewrites,
	}
}

func (rw *StatementRewrite) Rewrite(statement tree.Statement) tree.Statement {
	for _, rewrite := range rw.rewrites {
		statement = rewrite.Rewrite(statement)
	}
	return statement
}

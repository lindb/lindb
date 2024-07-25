package analyzer

import "github.com/lindb/lindb/sql/tree"

type AnalyzerContext struct {
	Database    string // default database name
	Analysis    *Analysis
	IDAllocator *tree.NodeIDAllocator
}

func NewAnalyzerContext(database string, stmt tree.Statement, idallocator *tree.NodeIDAllocator) *AnalyzerContext {
	return &AnalyzerContext{
		Database:    database,
		Analysis:    NewAnalysis(stmt),
		IDAllocator: idallocator,
	}
}

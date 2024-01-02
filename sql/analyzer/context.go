package analyzer

import "github.com/lindb/lindb/sql/tree"

type AnalyzerContext struct {
	Analysis    *Analysis
	IDAllocator *tree.NodeIDAllocator
}

func NewContext(stmt tree.Statement, idallocator *tree.NodeIDAllocator) *AnalyzerContext {
	return &AnalyzerContext{
		Analysis:    NewAnalysis(stmt),
		IDAllocator: idallocator,
	}
}

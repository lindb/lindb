package analyzer

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
)

type Analyzer struct {
	context     *AnalyzerContext
	metadataMgr spi.MetadataManager
}

func NewAnalyzer(context *AnalyzerContext, metadataMgr spi.MetadataManager) *Analyzer {
	return &Analyzer{
		context:     context,
		metadataMgr: metadataMgr,
	}
}

func (a *Analyzer) Analyze(statement tree.Statement) {
	// analyze statement
	analyzer := NewStatementAnalyzer(a.context, a.metadataMgr)
	analyzer.Analyze(statement)
}

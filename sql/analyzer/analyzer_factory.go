package analyzer

import (
	"github.com/lindb/lindb/spi"
)

type AnalyzerFactory struct {
	metadataMgr spi.MetadataManager
}

func NewAnalyzerFactory(metadataMgr spi.MetadataManager) *AnalyzerFactory {
	return &AnalyzerFactory{
		metadataMgr: metadataMgr,
	}
}

func (fct *AnalyzerFactory) CreateAnalyzer(ctx *AnalyzerContext) *Analyzer {
	return NewAnalyzer(ctx, fct.metadataMgr)
}

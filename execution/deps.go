package execution

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/rewrite"
)

type Deps struct {
	Repo             state.Repository
	CurrentNode      *models.InternalNode
	AnalyzerFct      *analyzer.AnalyzerFactory
	StatementRewrite *rewrite.StatementRewrite
}

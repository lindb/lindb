package execution

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/analyzer"
)

type Deps struct {
	Repo        state.Repository
	CurrentNode *models.InternalNode
	AnalyzerFct *analyzer.AnalyzerFactory
}

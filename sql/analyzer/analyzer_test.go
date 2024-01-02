package analyzer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/rewrite"
	"github.com/lindb/lindb/sql/tree"
)

func TestAnalyzer_Analyze(t *testing.T) {
	analyzer := NewAnalyzer(rewrite.NewStatementRewrite(nil))
	stmt, err := tree.CreateStatement("select a from b")
	assert.NoError(t, err)
	analysis := analyzer.Analyze(stmt)
	fmt.Println(analysis)
}

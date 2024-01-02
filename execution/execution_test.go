package execution

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/parser"
	"github.com/lindb/lindb/sql/rewrite"
)

func TestExecution_QueryExecution(t *testing.T) {
	parser := &parser.SQLParser{}
	stmt, err := parser.CreateStatement(`
select 
    t.*,n.*,sum(t.node)
from
    system.runtime.tasks t
left join 
    system.runtime.nodes n
on 
    t.node = n.node_id
where 
    t.state = 'FINISHED'
	`)
	//FIXME: test join
	assert.NoError(t, err)
	preparedStatment := &PreparedStatement{
		Statement: stmt,
	}
	ctx := context.WithValue(context.TODO(), constants.ContextKeySession, &models.Session{})
	exec := NewQueryExecution(ctx, &Deps{
		AnalyzerFct: analyzer.NewAnalyzerFactory(rewrite.NewStatementRewrite(nil)),
	}, preparedStatment)
	exec.Start(ctx)
}

func TestExecution_QueryExecution2(t *testing.T) {
	parser := &parser.SQLParser{}
	stmt, err := parser.CreateStatement(`
	select a.idle from "lindb.monitor.system.cpu_stat" a where role='Broker'
	`)
	//FIXME: test join
	assert.NoError(t, err)
	preparedStatment := &PreparedStatement{
		Statement: stmt,
	}
	ctx := context.WithValue(context.TODO(), constants.ContextKeySession, &models.Session{
		Databases: "_internal",
	})
	exec := NewQueryExecution(ctx, &Deps{
		AnalyzerFct: analyzer.NewAnalyzerFactory(rewrite.NewStatementRewrite(nil)),
	}, preparedStatment)
	exec.Start(ctx)
}

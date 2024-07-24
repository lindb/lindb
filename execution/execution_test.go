package execution

// func TestExecution_QueryExecution(t *testing.T) {
// 	stmt, err := tree.GetParser().CreateStatement(`
// select
//     t.*,n.*,sum(t.node)
// from
//     system.runtime.tasks t
// left join
//     system.runtime.nodes n
// on
//     t.node = n.node_id
// where
//     t.state = 'FINISHED'
// 	`, tree.NewNodeIDAllocator())
// 	// FIXME: test join
// 	assert.NoError(t, err)
// 	preparedStatment := &tree.PreparedStatement{
// 		Statement: stmt,
// 	}
// 	ctx := context.WithValue(context.TODO(), constants.ContextKeySession, &models.Session{})
// 	exec := NewQueryExecution(&Session{}, &Deps{
// 		AnalyzerFct: analyzer.NewAnalyzerFactory(rewrite.NewStatementRewrite(nil)),
// 	}, preparedStatment)
// 	exec.Start(ctx)
// }
//
// func TestExecution_QueryExecution2(t *testing.T) {
// 	stmt, err := tree.GetParser().CreateStatement(`
// 	select a.idle from "lindb.monitor.system.cpu_stat" a where role='Broker'
// 	`, tree.NewNodeIDAllocator())
// 	// FIXME: test join
// 	assert.NoError(t, err)
// 	preparedStatment := &PreparedStatement{
// 		Statement: stmt,
// 	}
// 	ctx := context.WithValue(context.TODO(), constants.ContextKeySession, &models.Session{
// 		Databases: "_internal",
// 	})
// 	exec := NewQueryExecution(ctx, &Deps{
// 		AnalyzerFct: analyzer.NewAnalyzerFactory(rewrite.NewStatementRewrite(nil)),
// 	}, preparedStatment)
// 	exec.Start(ctx)
// }

// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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

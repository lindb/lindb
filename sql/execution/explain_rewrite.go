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

import (
	"github.com/lindb/lindb/sql/interfaces"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/sql/utils"
)

type ExplainRewrite struct {
	session   *Session
	explainer *QueryExplainer
	builder   *utils.QueryBuilder
}

func NewExplainRewrite(session *Session, explainer *QueryExplainer) interfaces.Rewrite {
	return &ExplainRewrite{explainer: explainer, session: session, builder: utils.NewQueryBuilder(session.NodeIDAllocator)}
}

func (e *ExplainRewrite) Rewrite(statement tree.Statement) tree.Statement {
	if explain, ok := statement.(*tree.Explain); ok {
		return e.visitExplain(explain)
	}
	return statement
}

func (e *ExplainRewrite) visitExplain(node *tree.Explain) tree.Statement {
	explainType := tree.LogicalExplain
	for _, option := range node.Options {
		if eType, ok := option.(*tree.ExplainType); ok {
			explainType = eType.Type
		}
	}
	plan := e.explainer.ExplainPlan(e.session, node.Statement, explainType)
	return e.builder.SingleValueQuery("Query Plan", plan)
}

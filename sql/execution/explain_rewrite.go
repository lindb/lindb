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

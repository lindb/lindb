package rewrite

import (
	commonConstants "github.com/lindb/common/constants"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/sql/interfaces"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/sql/utils"
)

type ShowQueriesRewrite struct {
	db string
}

func NewShowQueriesRewrite(db string) interfaces.Rewrite {
	return &ShowQueriesRewrite{
		db: db,
	}
}

func (r *ShowQueriesRewrite) Rewrite(statement tree.Statement) tree.Statement {
	result := statement.Accept(nil, r)
	if rewritten, ok := result.(tree.Statement); ok {
		return rewritten
	}
	return statement
}

func (v *ShowQueriesRewrite) Visit(context any, n tree.Node) (r any) {
	switch node := n.(type) {
	case *tree.ShowColumns:
		return utils.SimpleQuery(
			utils.SelectItems("column_name", "data_type", "agg_type"),
			utils.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableColumns),
			utils.LogicalAnd(
				utils.StringEqual("table_schema", v.db),
				utils.StringEqual("namespace", node.Table.GetNamespace()),
				utils.StringEqual("table_name", node.Table.GetTableName()),
			),
		)
	}
	return nil
}

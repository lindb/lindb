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

package rewrite

import (
	commonConstants "github.com/lindb/common/constants"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi/table/infoschema"
	"github.com/lindb/lindb/sql/interfaces"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/sql/utils"
)

type ShowQueriesRewrite struct {
	builder *utils.QueryBuilder

	db string
}

func NewShowQueriesRewrite(db string, idAllocator *tree.NodeIDAllocator) interfaces.Rewrite {
	return &ShowQueriesRewrite{
		db:      db,
		builder: utils.NewQueryBuilder(idAllocator),
	}
}

func (v *ShowQueriesRewrite) Rewrite(statement tree.Statement) tree.Statement {
	result := statement.Accept(nil, v)
	if rewritten, ok := result.(tree.Statement); ok {
		return rewritten
	}
	return statement
}

func (v *ShowQueriesRewrite) Visit(context any, n tree.Node) (r any) {
	switch node := n.(type) {
	case *tree.Show:
		return node.Body.Accept(context, v)
	case *tree.ShowReplications:
		var terms []tree.Expression
		terms = append(terms, v.builder.StringEqual("table_schema", v.db)) // database
		return v.builder.SimpleQuery(
			v.builder.SelectItems(infoschema.GetShowSelectColumns(constants.TableReplications, 1)...),
			v.builder.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableReplications),
			v.builder.LogicalAnd(terms...),
		)
	case *tree.ShowMemoryDatabases:
		var terms []tree.Expression
		terms = append(terms, v.builder.StringEqual("table_schema", v.db)) // database
		return v.builder.SimpleQuery(
			v.builder.SelectItems(infoschema.GetShowSelectColumns(constants.TableMemoryDatabases, 1)...),
			v.builder.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableMemoryDatabases),
			v.builder.LogicalAnd(terms...),
		)
	case *tree.ShowNamespaces:
		var terms []tree.Expression
		terms = append(terms, v.builder.StringEqual("table_schema", v.db)) // database
		if node.LikePattern != "" {
			terms = append(terms, v.builder.Like("namespace", node.LikePattern)) // namespace like pattern
		}
		return v.builder.SimpleQuery(
			v.builder.SelectItems(infoschema.GetShowSelectColumns(constants.TableNamespaces, 1)...),
			v.builder.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableNamespaces),
			v.builder.LogicalAnd(terms...),
		)
	case *tree.ShowTableNames:
		var terms []tree.Expression
		terms = append(terms, v.builder.StringEqual("table_schema", v.db)) // database
		namespace := node.GetNamespace()
		if namespace == "" {
			namespace = commonConstants.DefaultNamespace
		}
		terms = append(terms, v.builder.StringEqual("namespace", namespace)) // namespace predicate
		if node.LikePattern != "" {
			terms = append(terms, v.builder.Like("table_name", node.LikePattern)) // table_name like pattern
		}
		return v.builder.SimpleQuery(
			v.builder.SelectItems(infoschema.GetShowSelectColumns(constants.TableTableNames, 1)...),
			v.builder.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableTableNames),
			v.builder.LogicalAnd(terms...),
		)
	case *tree.ShowColumns:
		return v.builder.SimpleQuery(
			v.builder.SelectItems(infoschema.GetShowSelectColumns(constants.TableColumns, 3)...),
			v.builder.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableColumns),
			v.builder.LogicalAnd(
				v.builder.StringEqual("table_schema", v.db),
				v.builder.StringEqual("namespace", node.Table.GetNamespace()),
				v.builder.StringEqual("table_name", node.Table.GetTableName()),
			),
		)
	case *tree.ShowDatabases:
		return v.builder.SimpleQuery(
			v.builder.AliasedSelectItem("schema_name", "Database"),
			v.builder.Table(constants.InformationSchema, commonConstants.DefaultNamespace, constants.TableSchemata),
			nil,
		)
	}
	return nil
}

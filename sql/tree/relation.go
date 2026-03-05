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

package tree

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/constants"
)

type Relation interface {
	Node
}

type Values struct {
	BaseNode

	Rows arrow.RecordBatch
}

func (n *Values) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type AliasedRelation struct {
	BaseNode

	Relation    Relation
	Aliase      *Identifier
	ColumnNames []*Identifier
}

func (n *AliasedRelation) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type Table struct {
	BaseNode
	Name *QualifiedName
}

func (n *Table) GetDatabase(defaultDB string) string {
	if len(n.Name.Parts) == 3 {
		return n.Name.Parts[0]
	}
	return defaultDB
}

func (n *Table) GetNamespace() string {
	switch len(n.Name.Parts) {
	case 3:
		return n.Name.Parts[1]
	case 2:
		return n.Name.Parts[0]
	default:
		return constants.DefaultNamespace
	}
}

func (n *Table) GetTableName() string {
	return n.Name.Parts[len(n.Name.Parts)-1]
}

func (n *Table) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type TableSubQuery struct {
	BaseNode
	Query *Query
}

func (n *TableSubQuery) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

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

type ShowBody interface {
	Node
}

type Show struct {
	BaseNode
	Body ShowBody
}

func (n *Show) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowNamespaces struct {
	BaseNode

	LikePattern string
}

func (n *ShowNamespaces) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowTableNames struct {
	BaseNode

	Namespace   *QualifiedName
	LikePattern string
}

func (n ShowTableNames) GetNamespace() string {
	if n.Namespace == nil {
		return ""
	}
	return n.Namespace.Name
}

func (n *ShowTableNames) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowColumns struct {
	BaseNode

	Table *Table
}

func (n *ShowColumns) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowReplications struct {
	BaseNode
}

func (n *ShowReplications) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowMemoryDatabases struct {
	BaseNode
}

func (n *ShowMemoryDatabases) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type ShowDatabases struct {
	BaseNode
}

func (n *ShowDatabases) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

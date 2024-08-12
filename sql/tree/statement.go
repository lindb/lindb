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

type StatementType int

const (
	UseStatement StatementType = iota + 1
	MetadataStatement
	SchemaStatement
	StorageStatement
	StateStatement
	MetricMetadataStatement
	QueryStatement
	RequestStatement
	BrokerStatement
	LimitStatement
)

// Statement represents LinDB query language statement
type Statement interface {
	Node
}

type BaseNode struct {
	ID       NodeID        `json:"id"`
	Location *NodeLocation `json:"-"`
}

func (n *BaseNode) GetID() NodeID {
	return n.ID
}

type PreparedStatement struct {
	PrepareSQL string
	Statement  Statement
}
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

const (
	LogicalExplain     string = "LOGICAL"
	DistributedExplain string = "DISTRIBUTED"
)

type ExplainOption interface{}

type ExplainType struct {
	Type string // default is LOGICAL
}

type Explain struct {
	BaseNode
	Statement Statement
	Options   []ExplainOption
}

func (n *Explain) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

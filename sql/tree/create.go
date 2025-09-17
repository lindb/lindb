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

import "github.com/lindb/lindb/pkg/option"

type CreateOption interface{}

type EngineOption struct {
	Type option.EngineType
}

type CreateDatabase struct {
	BaseNode
	Name          string
	CreateOptions []CreateOption
	Props         []*Property
	Rollup        []*RollupOption
}

func (n *CreateDatabase) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type CreateBroker struct {
	BaseNode
	Options map[string]any
	Name    string
}

func (n *CreateBroker) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type RollupOption struct {
	BaseNode
	Props []*Property
}

func (n *RollupOption) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

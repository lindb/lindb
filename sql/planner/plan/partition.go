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

package plan

import "github.com/lindb/lindb/sql/tree"

type PartitioningScheme struct {
	Partitioning *Partitioning `json:"partitioning"`
	OutputLayout []*Symbol     `json:"outputLayout"`
}

type Partitioning struct {
	Handle    *PartitioningHandle `json:"handle"`
	Arguments []*ArgumentBinding  `json:"arguments"`
}

func (p *Partitioning) Translate(translator func(symbol *Symbol) *Symbol) *Partitioning {
	// FIXME: imple translate
	return nil
}

type PartitioningHandle struct{}

func (h *PartitioningHandle) IsSingleNode() bool {
	return false
}

type ArgumentBinding struct {
	Expression tree.Expression `json:"expression"`
}

func (arg *ArgumentBinding) Translate() *ArgumentBinding {
	// FIXME: imple arg binding translate
	return nil
}

type PartitioningProps struct {
	PartitioningColumns []*Symbol
}

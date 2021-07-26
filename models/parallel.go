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

package models

// PhysicalPlan represents the distribution query's physical plan
type PhysicalPlan struct {
	Database      string         `json:"database"`      // database name
	Root          Root           `json:"root"`          // root node
	Intermediates []Intermediate `json:"intermediates"` // intermediate node if need
	Leafs         []Leaf         `json:"leafs"`         // leaf nodes(storage nodes of query database)
}

// NewPhysicalPlan creates the physical plan with root node
func NewPhysicalPlan(root Root) *PhysicalPlan {
	return &PhysicalPlan{Root: root}
}

// AddIntermediate adds an intermediate node into the intermediate node list
func (t *PhysicalPlan) AddIntermediate(intermediate Intermediate) {
	t.Intermediates = append(t.Intermediates, intermediate)
}

// AddLeaf adds a leaf node into the leaf node list
func (t *PhysicalPlan) AddLeaf(leaf Leaf) {
	t.Leafs = append(t.Leafs, leaf)
}

// Root represents the root node info
type Root struct {
	Indicator string `json:"indicator"`
	NumOfTask int32  `json:"numOfTask"`
}

type BaseNode struct {
	Parent    string `json:"parent"`    // parent node's indicator
	Indicator string `json:"indicator"` // current node's indicator
}

// Intermediate represents the intermediate node info
type Intermediate struct {
	BaseNode

	NumOfTask int32 `json:"numOfTask"`
}

// Leaf represents the leaf node info
type Leaf struct {
	BaseNode

	Receivers []Node  `json:"receivers"`
	ShardIDs  []int32 `json:"shardIDs"`
}

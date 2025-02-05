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

import "github.com/lindb/lindb/models"

type PlanFragment struct {
	Root          PlanNode                      `json:"root,omitempty"`
	ParentNode    *PlanNodeID                   `json:"parentNode,omitempty"`
	Partitions    map[models.InternalNode][]int `json:"-"`
	Receivers     []models.InternalNode         `json:"receivers"`
	RemoteSources []*RemoteSourceNode           `json:"remoteSources,omitempty"`
	ID            FragmentID                    `json:"id"`
}

func NewPlanFragment(id FragmentID, root PlanNode) *PlanFragment {
	fragment := &PlanFragment{
		ID:   id,
		Root: root,
	}

	fragment.findRemoteSources(root)
	return fragment
}

func (pf *PlanFragment) findRemoteSources(node PlanNode) {
	sources := node.GetSources()
	for i := range sources {
		pf.findRemoteSources(sources[i])
	}

	if remoteSource, ok := node.(*RemoteSourceNode); ok {
		pf.RemoteSources = append(pf.RemoteSources, remoteSource)
	}
}

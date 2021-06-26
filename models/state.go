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

import "encoding/json"

// StorageState represents storage cluster node state.
// NOTICE: it is not safe for concurrent use.
type StorageState struct {
	Name        string                 `json:"name"`
	ActiveNodes map[string]*ActiveNode `json:"activeNodes"`
}

// NewStorageState creates storage cluster state
func NewStorageState() *StorageState {
	return &StorageState{
		ActiveNodes: make(map[string]*ActiveNode),
	}
}

// AddActiveNode adds a node into active node list
func (s *StorageState) AddActiveNode(node *ActiveNode) {
	key := node.Node.Indicator()
	_, ok := s.ActiveNodes[key]
	if !ok {
		s.ActiveNodes[key] = node
	}
}

// RemoveActiveNode removes a node from active node list
func (s *StorageState) RemoveActiveNode(node string) {
	delete(s.ActiveNodes, node)
}

// GetActiveNodes returns all active nodes
func (s *StorageState) GetActiveNodes() []*ActiveNode {
	var nodes []*ActiveNode
	for _, node := range s.ActiveNodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// Stringer returns a human readable string
func (s *StorageState) String() string {
	content, _ := json.Marshal(s)
	return string(content)
}

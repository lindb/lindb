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

import (
	"fmt"

	"github.com/lindb/lindb/constants"
)

// PhysicalPlan represents the distribution query's physical plan
type PhysicalPlan struct {
	Database  string    `json:"database"` // database name
	Targets   []*Target `json:"targets"`
	Receivers []string  `json:"receivers"`
}

// AddReceiver adds a receiver.
func (t *PhysicalPlan) AddReceiver(receiver string) {
	t.Receivers = append(t.Receivers, receiver)
}

// AddTarget adds a target.
func (t *PhysicalPlan) AddTarget(target *Target) {
	t.Targets = append(t.Targets, target)
}

// Validate checks the plan if valid.
func (t *PhysicalPlan) Validate() error {
	if t.Database == "" {
		return constants.ErrDatabaseNameRequired
	}
	if len(t.Targets) == 0 {
		return constants.ErrTargetNodesNotFound
	}
	if len(t.Receivers) == 0 {
		return constants.ErrReceiveNodesNotFound
	}
	return nil
}

// Leaf represents the leaf node info
type Target struct {
	ReceiveOnly bool      `json:"receiverOnly"`
	Indicator   string    `json:"indicator"` // current node's indicator
	ShardIDs    []ShardID `json:"shardIDs"`
}

// Partition represents data partition info.
type Partition struct {
	ID   ShardID      `json:"id"`
	Node InternalNode `json:"node"`
}

type InternalNode struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

func (n InternalNode) Address() string {
	return fmt.Sprintf("%s:%d", n.IP, n.Port)
}

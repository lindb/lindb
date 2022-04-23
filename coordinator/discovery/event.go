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

package discovery

// EventType represents coordinator event type.
type EventType int

const (
	DatabaseConfigChanged EventType = iota + 1
	DatabaseConfigDeletion
	ShardAssignmentChanged
	ShardAssignmentDeletion
	NodeStartup
	NodeFailure
	StorageStateChanged
	StorageStateDeletion
	StorageConfigChanged
	StorageConfigDeletion
)

// String returns string value of EventType.
func (e EventType) String() string {
	switch e {
	case DatabaseConfigChanged:
		return "DatabaseConfigChanged"
	case DatabaseConfigDeletion:
		return "DatabaseConfigDeletion"
	case ShardAssignmentChanged:
		return "ShardAssignmentChanged"
	case ShardAssignmentDeletion:
		return "ShardAssignmentDeletion"
	case NodeStartup:
		return "NodeStartup"
	case NodeFailure:
		return "NodeFailure"
	case StorageStateChanged:
		return "StorageStateChanged"
	case StorageStateDeletion:
		return "StorageStateDeletion"
	case StorageConfigChanged:
		return "StorageConfigChanged"
	case StorageConfigDeletion:
		return "StorageConfigDeletion"
	default:
		return "unknown"
	}
}

// Event represents discovery state change event.
type Event struct {
	Type  EventType
	Key   string
	Value []byte

	Attributes map[string]string
}

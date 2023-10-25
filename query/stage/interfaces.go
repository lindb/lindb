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

package stage

import "github.com/lindb/common/models"

//go:generate mockgen -source=./interfaces.go -destination=./interfaces_mock.go -package=stage

// Type represents the type of stage.
type Type int

const (
	// Unknown represents unknown stage.
	Unknown Type = iota
	// MetadataLookup represents metadata lookup stage.
	MetadataLookup
	// ShardScan represents shard scan stage.
	ShardScan
	// Grouping represents grouping stage.
	Grouping
	// DataLoad represents data load stage.
	DataLoad
	// MetadataSuggest represents metadata suggest stage.
	MetadataSuggest
	// ShardLookup represents shard lookup stage.
	ShardLookup
	// PhysicalPlan represents physical plan stage.
	PhysicalPlan
	// TaskSend represents task send stage.
	TaskSend
)

// String returns string value of stage type.
func (t Type) String() string {
	switch t {
	case MetadataLookup:
		return "MetadataLookup"
	case ShardScan:
		return "ShardScan"
	case Grouping:
		return "Grouping"
	case DataLoad:
		return "DataLoad"
	case MetadataSuggest:
		return "MetadataSuggest"
	case ShardLookup:
		return "ShardLookup"
	case PhysicalPlan:
		return "PhysicalPlan"
	case TaskSend:
		return "TaskSend"
	default:
		return "Unknown"
	}
}

// Stage represents stage under execute pipeline.
type Stage interface {
	// Identifier returns identifier value of current stage.
	Identifier() string
	// Stats returns the execution stats of current stage.
	Stats() []*models.OperatorStats
	// Type returns the type of stage.
	Type() Type
	// Plan plans sub execute tree for this stage.
	Plan() PlanNode
	// NextStages returns the next stages after this stage completed.
	NextStages() []Stage
	// Execute executes the plan tree of this stage.
	Execute(node PlanNode, completeHandle func(), errHandle func(err error))
	// Complete completes this stage, does some resource release operate.
	Complete()
	// IsAsync returns stage if stage async execute.
	IsAsync() bool
}

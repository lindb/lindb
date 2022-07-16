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

package storagequery

import (
	"github.com/google/uuid"

	stagepkg "github.com/lindb/lindb/query/stage"
)

//go:generate mockgen -source=./pipeline.go -destination=./pipeline_mock.go -package=storagequery

// Pipeline represents the pipeline execution model, pipeline executes all query stage.
type Pipeline interface {
	// Execute executes the stage(sub plan tree).
	Execute(stage stagepkg.Stage)
}

// pipeline implements Pipeline interface.
type pipeline struct {
	sm *pipelineStateMachine
}

// NewExecutePipeline creates a Pipeline instance for executing query stage.
func NewExecutePipeline(needStats bool, completeCallback func(err error)) Pipeline {
	return &pipeline{
		sm: newPipelineStateMachine(needStats, completeCallback),
	}
}

// Execute executes the stage(sub plan tree).
func (p *pipeline) Execute(stage stagepkg.Stage) {
	p.executeStage(stage)
}

// executeStage executes current the plan tree of current stage,
// if it executes success, plan next stages and executes them.
//
//       +----------------+ 1.Plan&Execute +-----------+
//       | current stage  |--------------->| plan tree |
//       +----------------+                +-----------+
//               | 2. Plan Next Stages
//               v
//  +-------------------------+
//  |  +-------+   +-------+  |
//  |  |stage1 |   |stage2 |  |
//  |  +-------+   +-------+  |
//  +-------------------------+
func (p *pipeline) executeStage(stage stagepkg.Stage) {
	if stage == nil || p.sm.isCompleted() {
		return
	}

	stageID := uuid.New().String()
	p.sm.executeStage(stageID, stage)

	stage.Execute(stage.Plan(), func() {
		// after current stage execute completed, then plan next stages
		nextStages := stage.NextStages()
		for idx := range nextStages {
			p.executeStage(nextStages[idx])
		}

		// completed current stage, change stage state
		p.sm.completeStage(stageID, nil)
		stage.Complete()
	}, func(err error) {
		// complete stage with err
		p.sm.completeStage(stageID, err)
	})
}

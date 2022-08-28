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

package query

import (
	"github.com/google/uuid"

	"github.com/lindb/lindb/models"
	errorpkg "github.com/lindb/lindb/pkg/error"
	"github.com/lindb/lindb/pkg/logger"
	stagepkg "github.com/lindb/lindb/query/stage"
	trackerpkg "github.com/lindb/lindb/query/tracker"
)

//go:generate mockgen -source=./pipeline.go -destination=./pipeline_mock.go -package=query

// Pipeline represents the pipeline execution model, pipeline executes all query stage.
type Pipeline interface {
	// Execute executes the stage(sub plan tree).
	Execute(stage stagepkg.Stage)
	// Stats returns the stats of stages.
	Stats() []*models.StageStats
}

// pipeline implements Pipeline interface.
type pipeline struct {
	sm *pipelineStateMachine

	logger *logger.Logger
}

// NewExecutePipeline creates a Pipeline instance for executing query stage.
func NewExecutePipeline(tracker *trackerpkg.StageTracker, completeCallback func(err error)) Pipeline {
	return &pipeline{
		sm:     newPipelineStateMachine(tracker, completeCallback),
		logger: logger.GetLogger("Query", "Pipeline"),
	}
}

// Execute executes the stage(sub plan tree).
func (p *pipeline) Execute(stage stagepkg.Stage) {
	defer func() {
		if r := recover(); r != nil {
			err := errorpkg.Error(r)
			p.sm.complete(err)
			p.logger.Error("execute query pipeline panic", logger.Error(err), logger.Stack())
		}
	}()

	p.executeStage("", stage)
}

// Stats returns the stats of stages.
func (p *pipeline) Stats() []*models.StageStats {
	return p.sm.GetStats()
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
func (p *pipeline) executeStage(parentStageID string, stage stagepkg.Stage) {
	if stage == nil || p.sm.isCompleted() {
		return
	}

	stageID := uuid.New().String()
	p.sm.executeStage(parentStageID, stageID, stage)

	stage.Execute(stage.Plan(), func() {
		// after current stage execute completed, then plan next stages
		nextStages := stage.NextStages()
		for idx := range nextStages {
			p.executeStage(stageID, nextStages[idx])
		}

		// completed current stage, change stage state
		p.sm.completeStage(stageID, nil)
	}, func(err error) {
		// complete stage with err
		p.sm.completeStage(stageID, err)
	})
}

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

import (
	"context"
	"errors"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/models"
)

// baseStage represents common implements for Stage interface.
type baseStage struct {
	ctx context.Context

	stageType Type
	execPool  concurrent.Pool
	track     bool

	operators []*models.OperatorStats
}

// Track tracks the stage exeucte stats.
func (stage *baseStage) Track() {
	stage.track = true
}

// Stats returns the stats of current stage.
func (stage *baseStage) Stats() []*models.OperatorStats {
	return stage.operators
}

// Type returns the type of stage.
func (stage *baseStage) Type() Type {
	return stage.stageType
}

// IsAsync returns stage if stage async execute.
func (stage *baseStage) IsAsync() bool {
	return stage.execPool != nil && stage.ctx != nil
}

// Execute executes the plan node, if it executes success invoke completeHandle func else invoke errHande func.
func (stage *baseStage) Execute(node PlanNode, completeHandle func(), errHandle func(err error)) {
	execFn := func() {
		// execute sub plan tree for current stage
		if err := stage.execute(node); err != nil {
			errHandle(err)
		} else {
			completeHandle()
		}
	}
	if stage.IsAsync() {
		stage.execPool.Submit(stage.ctx, concurrent.NewTask(func() {
			execFn()
		}, errHandle))
	} else {
		execFn()
	}
}

// execute the plan node under current stage.
func (stage *baseStage) execute(node PlanNode) (err error) {
	if node == nil {
		return nil
	}

	var stats *models.OperatorStats
	// execute current plan node logic
	if stage.track {
		stats, err = node.ExecuteWithStats()
		if stats != nil {
			stage.operators = append(stage.operators, stats)
		}
	} else {
		err = node.Execute()
	}
	if err != nil {
		if node.IgnoreNotFound() && errors.Is(err, constants.ErrNotFound) {
			return nil
		}
		return err
	}

	// if it has child node, need execute child node logic
	children := node.Children()
	for idx := range children {
		if err := stage.execute(children[idx]); err != nil {
			return err
		}
	}
	return nil
}

// Complete completes current stage.
func (stage *baseStage) Complete() {
}

// NextStages returns the next stages in the pipeline.
func (stage *baseStage) NextStages() []Stage {
	return nil
}

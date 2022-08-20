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
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	stagepkg "github.com/lindb/lindb/query/stage"
)

// stageTracker represents track the stat of the stage execution.
type stageTracker struct {
	stageID            string
	stage              stagepkg.Stage // current execute stage
	state              stagepkg.State // stage execute stage
	startTime, endTime time.Time      // stage start/end time
	stats              *models.StageStats
}

// pipelineStateMachine represents pipeline stage machine which track all stage execution state under this pipeline.
type pipelineStateMachine struct {
	stages              map[string]*stageTracker                    // store schedule stage
	pending             atomic.Int32                                // how many stages are pending, not completed
	completedCallbackFn func(stats []*models.StageStats, err error) // pipeline execute completed will invoke
	mutex               sync.Mutex
	completed           atomic.Bool

	needStats     bool
	statsOfStages []*models.StageStats
}

// newPipelineStateMachine creates a pipelineStateMachine instance.
func newPipelineStateMachine(needStats bool, completeCallback func(stats []*models.StageStats, err error)) *pipelineStateMachine {
	return &pipelineStateMachine{
		stages:              make(map[string]*stageTracker),
		completedCallbackFn: completeCallback,
		needStats:           needStats,
	}
}

// executeStage tracks stage start execution state.
func (sm *pipelineStateMachine) executeStage(parentStageID, stageID string, stage stagepkg.Stage) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.pending.Inc()
	ts := &stageTracker{
		stageID: stageID,
		stage:   stage,
		state:   stagepkg.ExecutingState,
	}
	sm.stages[stageID] = ts
	if sm.needStats {
		ts.stage.Track()
		ts.stats = &models.StageStats{}
		if parentStageID == "" {
			sm.statsOfStages = append(sm.statsOfStages, ts.stats)
		}
		ts.startTime = time.Now()
	}
}

// completeStage tracks stage complete execution state.
func (sm *pipelineStateMachine) completeStage(parentStageID, stageID string, err error) {
	sm.mutex.Lock()
	if s, ok := sm.stages[stageID]; ok {
		var errMsg string
		if err != nil {
			s.state = stagepkg.ErrorState
			errMsg = err.Error()
		} else {
			s.state = stagepkg.CompleteState
		}

		if sm.needStats {
			s.stats.Operators = s.stage.Stats()
			s.endTime = time.Now()
			s.stats.Identifier = s.stage.Identifier()
			s.stats.Start = s.startTime.UnixNano()
			s.stats.End = s.endTime.UnixNano()
			s.stats.Cost = s.endTime.Sub(s.startTime).Nanoseconds()
			s.stats.State = s.state.String()
			s.stats.ErrMsg = errMsg
			s.stats.Async = s.stage.IsAsync()

			if parentStageID != "" {
				parent := sm.stages[parentStageID]
				parent.stats.Children = append(parent.stats.Children, s.stats)
			}
		}
	}
	sm.mutex.Unlock()

	if sm.pending.Dec() == 0 {
		// check if all stages execute completed
		sm.complete(err)
	}
}

// complete executes pipeline completed, invokes completed callback.
func (sm *pipelineStateMachine) complete(err error) {
	if sm.completed.CAS(false, true) && sm.completedCallbackFn != nil {
		// check if all stages execute completed
		sm.completedCallbackFn(sm.statsOfStages, err)
	}
}

// isCompleted checks if the pipeline is completed.
func (sm *pipelineStateMachine) isCompleted() bool {
	return sm.completed.Load()
}

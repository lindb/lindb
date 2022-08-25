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

package tracker

import (
	"sync"
	"time"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
)

// State represents the state of stage.
type State int

const (
	// UnknownState represents unknown stage state.
	UnknownState State = iota
	// InitState represetnsi stage init.
	InitState
	// ExecutingState represents stage is executing.
	ExecutingState
	// CompleteState represents stage completed successfully.
	CompleteState
	// ErrorState represents stage completed with error.
	ErrorState
)

// String returns string value of stage state.
func (s State) String() string {
	switch s {
	case InitState:
		return "Init"
	case ExecutingState:
		return "Executing"
	case CompleteState:
		return "Complete"
	case ErrorState:
		return "Error"
	default:
		return "Unknown"
	}
}

// StageTracker represents a tracker which track the state of stage execution.
type StageTracker struct {
	mutex   sync.Mutex
	taskCtx *flow.TaskContext

	stages               []*models.StageStats
	groupingCollectStage *models.StageStats

	stats *models.LeafNodeStats
}

// NewStageTracker creates a StageTracker instance.
func NewStageTracker(taskCtx *flow.TaskContext) *StageTracker {
	return &StageTracker{
		taskCtx: taskCtx,
	}
}

// AddStage adds a stage execution stats.
func (s *StageTracker) AddStage(stage *models.StageStats) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stages = append(s.stages, stage)
}

// GetStages returns all stages' execution stats with lock.
func (s *StageTracker) GetStages() (rs []*models.StageStats) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getStages()
}

// SetGroupingCollectStageValues sets grouping collect stage stats via callback.
func (s *StageTracker) SetGroupingCollectStageValues(fn func(stage *models.StageStats)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// if grouping stage stats is nil, new a instance
	if s.groupingCollectStage == nil {
		s.groupingCollectStage = &models.StageStats{}
	}
	// invoke func do set values logic
	fn(s.groupingCollectStage)
}

// Complete completes stage stats track, build result.
func (s *StageTracker) Complete() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	end := time.Now()
	s.stats = &models.LeafNodeStats{
		Start:     s.taskCtx.Start.UnixNano(),
		End:       end.UnixNano(),
		TotalCost: end.Sub(s.taskCtx.Start).Nanoseconds(),
		Stages:    s.getStages(),
	}
}

// GetStats returns the track stats result.
func (s *StageTracker) GetStats() *models.LeafNodeStats {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.stats
}

// getStages returns all stages' execution stats without lock.
func (s *StageTracker) getStages() (rs []*models.StageStats) {
	rs = append(rs, s.stages...)

	// if has grouping stage append state last
	if s.groupingCollectStage != nil {
		rs = append(rs, s.groupingCollectStage)
	}
	return rs
}

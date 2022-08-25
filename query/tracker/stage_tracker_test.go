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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
)

func TestState_String(t *testing.T) {
	assert.Equal(t, "Unknown", UnknownState.String())
	assert.Equal(t, "Init", InitState.String())
	assert.Equal(t, "Executing", ExecutingState.String())
	assert.Equal(t, "Complete", CompleteState.String())
	assert.Equal(t, "Error", ErrorState.String())
}

func TestStageTracker(t *testing.T) {
	taskCtx := flow.NewTaskContextWithTimeout(context.TODO(), time.Minute)
	tracker := NewStageTracker(taskCtx)

	assert.Empty(t, tracker.GetStages())

	tracker.AddStage(&models.StageStats{})
	assert.Len(t, tracker.GetStages(), 1)
	tracker.SetGroupingCollectStageValues(func(stage *models.StageStats) {
		stage.Identifier = "test"
	})
	tracker.Complete()
	assert.Len(t, tracker.GetStages(), 2)
}

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
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
)

// groupingExecuteTrack represents the track for grouping execute.
type groupingExecuteTrack struct {
	ctx       *executeContext
	queryFlow flow.StorageQueryFlow

	pendingTask atomic.Int32 // track how many tasks are pending
	completed   atomic.Bool
}

// newGroupingExecuteTrack creates a grouping execute track.
func newGroupingExecuteTrack(ctx *executeContext, queryFlow flow.StorageQueryFlow) *groupingExecuteTrack {
	return &groupingExecuteTrack{
		ctx:       ctx,
		queryFlow: queryFlow,
	}
}

// submitTask submits group by execute task.
func (t *groupingExecuteTrack) submitTask(stage flow.Stage, task func()) {
	t.pendingTask.Inc()
	t.queryFlow.Submit(stage, func() {
		defer func() {
			t.pendingTask.Dec()
			t.collectGroupByTagValues()
		}()
		task()
	})
}

// collectGroupByTagValues collects group tag values
func (t *groupingExecuteTrack) collectGroupByTagValues() {
	if t.pendingTask.Load() != 0 || !t.completed.CAS(false, true) {
		return
	}
	// all shard pending query tasks and grouping task completed, start collect tag values
	for idx, tagKeyID := range t.ctx.storageExecuteCtx.GroupByTags {
		tagKey := tagKeyID
		tagValueIDs := t.ctx.storageExecuteCtx.GroupingTagValueIDs[idx]
		tagIndex := idx
		if tagValueIDs == nil || tagValueIDs.IsEmpty() {
			t.queryFlow.ReduceTagValues(tagIndex, nil)
			continue
		}
		t.queryFlow.Submit(flow.ScannerStage, func() {
			tagValues := make(map[uint32]string) // tag value id => tag value
			task := newCollectTagValuesTaskFunc(t.ctx, t.ctx.getMetadata(), tagKey, tagValueIDs, tagValues)
			if err := task.Run(); err != nil {
				t.queryFlow.Complete(err)
				return
			}
			t.queryFlow.ReduceTagValues(tagIndex, tagValues)
		})
	}
}

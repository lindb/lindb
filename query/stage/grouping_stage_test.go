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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/context"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestGroupingStage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{}).AnyTimes()
	dataLoadCtx := &flow.DataLoadContext{}
	shard := tsdb.NewMockShard(ctrl)
	stage := NewGroupingStage(&context.LeafExecuteContext{
		TaskCtx:  &flow.TaskContext{},
		Database: db,
		GroupingCtx: context.NewLeafGroupingContext(&context.LeafExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{Query: &stmtpkg.Query{}},
			Database:          db,
		}),
	}, dataLoadCtx, shard)

	assert.NotNil(t, stage.Plan())
	stage.Complete()
	shard.EXPECT().ShardID().Return(models.ShardID(19))
	assert.Equal(t, "Grouping[Shard(19)]", stage.Identifier())

	t.Run("group not found", func(t *testing.T) {
		dataLoadCtx.IsGrouping = true
		assert.Empty(t, stage.NextStages())
	})

	t.Run("group found", func(t *testing.T) {
		dataLoadCtx.IsGrouping = false
		dataLoadCtx.ShardExecuteCtx = &flow.ShardExecuteContext{
			TimeSegmentContext: flow.NewTimeSegmentContext(),
		}
		dataLoadCtx.PendingDataLoadTasks = atomic.NewInt32(0)
		dataLoadCtx.ShardExecuteCtx.TimeSegmentContext.TimeSegments[10] = &flow.TimeSegmentResultSet{}
		assert.NotEmpty(t, stage.NextStages())
	})
}

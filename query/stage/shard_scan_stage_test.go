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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	contextpkg "github.com/lindb/lindb/query/context"
	trackerpkg "github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestShardScanStage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	storageCtx := &flow.StorageExecuteContext{
		Query: &stmt.Query{
			Condition: &stmt.EqualsExpr{},
			GroupBy:   []string{"key"},
		},
		ShardIDs: []models.ShardID{1, 2},
	}
	ctx := &contextpkg.LeafExecuteContext{
		TaskCtx:           flow.NewTaskContextWithTimeout(context.TODO(), time.Minute),
		Database:          db,
		StorageExecuteCtx: storageCtx,
	}
	ctx.Tracker = trackerpkg.NewStageTracker(ctx.TaskCtx)
	ctx.GroupingCtx = contextpkg.NewLeafGroupingContext(ctx)
	shard := tsdb.NewMockShard(ctrl)
	shardExecuteCtx := flow.NewShardExecuteContext(storageCtx)
	db.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{}).AnyTimes()
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	s := NewShardScanStage(ctx, shardExecuteCtx, shard)

	t.Run("no family", func(t *testing.T) {
		shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil)
		assert.Nil(t, s.Plan())
	})
	t.Run("all series", func(t *testing.T) {
		storageCtx.Query.Condition = nil
		shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).
			Return([]tsdb.DataFamily{tsdb.NewMockDataFamily(ctrl)})
		assert.NotNil(t, s.Plan())
	})
	t.Run("query condition", func(t *testing.T) {
		storageCtx.Query.Condition = &stmt.EqualsExpr{}
		shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).
			Return([]tsdb.DataFamily{tsdb.NewMockDataFamily(ctrl)})
		assert.NotNil(t, s.Plan())
	})

	shardExecuteCtx.SeriesIDsAfterFiltering = roaring.BitmapOf(1, 2, 3)
	assert.NotEmpty(t, s.NextStages())
	s.Complete()

	shard.EXPECT().ShardID().Return(models.ShardID(19))
	assert.Equal(t, "Shard Scan[Shard(19)]", s.Identifier())
}

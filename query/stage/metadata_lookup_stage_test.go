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
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestMetadataLookupStage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	meta.EXPECT().MetadataDatabase().Return(metaDB)
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	storageCtx := &flow.StorageExecuteContext{
		Query: &stmt.Query{
			Condition: &stmt.EqualsExpr{},
		},
		ShardIDs: []models.ShardID{1, 2},
	}
	ctx := &context.LeafExecuteContext{
		TaskCtx:           &flow.TaskContext{},
		Database:          db,
		StorageExecuteCtx: storageCtx,
	}
	ctx.GroupingCtx = context.NewLeafGroupingContext(ctx)
	s := NewMetadataLookupStage(ctx)
	assert.NotNil(t, s.Plan())

	shard := tsdb.NewMockShard(ctrl)
	db.EXPECT().GetShard(gomock.Any()).Return(shard, true).MaxTimes(2)
	db.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{}).MaxTimes(2)
	assert.NotEmpty(t, s.NextStages())

	assert.Equal(t, "Metadata Lookup", s.Identifier())
}

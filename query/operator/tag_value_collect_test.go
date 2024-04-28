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

package operator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/series/tag"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestTagValueCollect_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("db").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().ShardID().Return(models.ShardID(10)).AnyTimes()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	db.EXPECT().MetaDB().Return(metaDB).AnyTimes()
	indexSegment := index.NewMockMetricIndexSegment(ctrl)
	shard.EXPECT().IndexSegment().Return(indexSegment).AnyTimes()
	indexSegment.EXPECT().GetGroupingContext(gomock.Any()).Return(nil).AnyTimes()
	metaDB.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	ctx := &context.LeafMetadataContext{
		Database:          db,
		Request:           &stmtpkg.MetricMetadata{},
		StorageExecuteCtx: &flow.StorageExecuteContext{},
		TagKeyID:          tag.KeyID(10),
	}
	shardCtx := flow.NewShardExecuteContext(ctx.StorageExecuteCtx)
	shardCtx.SeriesIDsAfterFiltering = roaring.BitmapOf(1, 2, 3)
	shardCtx.GroupingContext = flow.NewGroupContext([]tag.KeyID{10}, map[tag.KeyID][]flow.GroupingScanner{})

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "get grouping context failure",
			prepare: func() {
				indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(nil, fmt.Errorf("err")).AnyTimes()
			},
		},
		{
			name: "collect tag value failure",
			prepare: func() {
				indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(nil, nil).AnyTimes()
				metaDB.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
			},
		},
		{
			name: "collect tag value successfully",
			prepare: func() {
				indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(nil, nil).AnyTimes()
				metaDB.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ tag.KeyID,
					_ *roaring.Bitmap,
					tagValues map[uint32]string) error {
					tagValues[10] = "value10"
					return nil
				}).AnyTimes()
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := NewTagValueCollect(ctx, shardCtx, shard)
			if tt.prepare != nil {
				tt.prepare()
			}
			assert.NoError(t, op.Execute())
		})
	}
}

func TestTagValueCollect_Identifier(t *testing.T) {
	assert.Equal(t, "Tag Value Collect", NewTagValueCollect(nil, nil, nil).Identifier())
}

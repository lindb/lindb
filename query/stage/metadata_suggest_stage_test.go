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
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestMetadataSuggestStage_Plan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	db.EXPECT().Metadata().Return(meta)

	ctx := context.NewLeafMetadataContext(&stmtpkg.MetricMetadata{}, db, nil)

	cases := []struct {
		name string
		in   *stmtpkg.MetricMetadata
	}{
		{
			name: "namespace suggest",
			in:   &stmtpkg.MetricMetadata{Type: stmtpkg.Namespace},
		},
		{
			name: "metric suggest",
			in:   &stmtpkg.MetricMetadata{Type: stmtpkg.Metric},
		},
		{
			name: "tag key suggest",
			in:   &stmtpkg.MetricMetadata{Type: stmtpkg.TagKey},
		},
		{
			name: "field suggest",
			in:   &stmtpkg.MetricMetadata{Type: stmtpkg.Field},
		},
		{
			name: "tag value suggest without condition",
			in:   &stmtpkg.MetricMetadata{Type: stmtpkg.TagValue},
		},
		{
			name: "tag value suggest with condition",
			in: &stmtpkg.MetricMetadata{
				Type:      stmtpkg.TagValue,
				Condition: &stmtpkg.EqualsExpr{},
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				ctx.Request = &stmtpkg.MetricMetadata{}
			}()
			ctx.Request = tt.in
			s := NewMetadataSuggestStage(ctx)
			assert.NotNil(t, s.Plan())
		})
	}

	t.Run("unknown type", func(t *testing.T) {
		ctx.Request = &stmtpkg.MetricMetadata{}
		s := NewMetadataSuggestStage(ctx)
		assert.Nil(t, s.Plan())
	})
}

func TestMetadataSuggestStage_NextStages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("not tag value suggest", func(t *testing.T) {
		ctx := &context.LeafMetadataContext{
			Request: &stmtpkg.MetricMetadata{},
		}
		assert.Empty(t, NewMetadataSuggestStage(ctx).NextStages())
	})
	t.Run("tag filter result not found", func(t *testing.T) {
		ctx := &context.LeafMetadataContext{
			Request:           &stmtpkg.MetricMetadata{Type: stmtpkg.TagValue},
			StorageExecuteCtx: &flow.StorageExecuteContext{},
		}
		assert.Empty(t, NewMetadataSuggestStage(ctx).NextStages())
	})
	t.Run("plan next stages", func(t *testing.T) {
		db := tsdb.NewMockDatabase(ctrl)
		ctx := &context.LeafMetadataContext{
			Request: &stmtpkg.MetricMetadata{Type: stmtpkg.TagValue},
			StorageExecuteCtx: &flow.StorageExecuteContext{
				TagFilterResult: map[string]*flow.TagFilterResult{"test": nil},
			},
			ShardIDs: []models.ShardID{1, 2},
			Database: db,
		}
		db.EXPECT().GetShard(models.ShardID(1)).Return(nil, false)
		db.EXPECT().GetShard(models.ShardID(2)).Return(nil, true)
		assert.Len(t, NewMetadataSuggestStage(ctx).NextStages(), 1)
	})
}

func TestMetadataSuggest_Identifier(t *testing.T) {
	assert.Equal(t, "Metadata Suggest", NewMetadataSuggestStage(nil).Identifier())
}

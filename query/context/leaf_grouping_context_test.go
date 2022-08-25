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

package context

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/series/tag"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestLeafGroupingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	meta.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	c, cancel := context.WithCancel(context.TODO())
	storageCtx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{
			GroupBy: []string{"key"},
		},
		TaskCtx: &flow.TaskContext{
			Cancel: cancel,
			Ctx:    c,
		},
	}
	ctx := NewLeafGroupingContext(&LeafExecuteContext{
		Tracker:           tracker.NewStageTracker(flow.NewTaskContextWithTimeout(context.TODO(), time.Minute)),
		StorageExecuteCtx: storageCtx,
		Database:          db,
		LeafNode:          &models.Leaf{},
	})

	cases := []struct {
		name    string
		prepare func()
		assert  func()
	}{
		{
			name: "no grouping",
			prepare: func() {
				storageCtx.Query.GroupBy = nil
			},
		},
		{
			name: "empty tag value",
			prepare: func() {
				storageCtx.GroupByTags = tag.Metas{{
					Key: "key",
					ID:  tag.KeyID(1),
				}}
				storageCtx.GroupingTagValueIDs = []*roaring.Bitmap{nil}
				storageCtx.Query.GroupBy = []string{"key"}
				ctx.tagValuesMap = make([]map[uint32]string, 1)
			},
		},
		{
			name: "collect tag value failure",
			prepare: func() {
				storageCtx.GroupByTags = tag.Metas{{
					Key: "key",
					ID:  tag.KeyID(1),
				}}
				storageCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf(1, 2, 3)}
				storageCtx.Query.GroupBy = []string{"key"}
				tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
		},
		{
			name: "collect tag value successfully",
			prepare: func() {
				storageCtx.GroupByTags = tag.Metas{{
					Key: "key",
					ID:  tag.KeyID(1),
				}}
				storageCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf(1, 2, 3)}
				storageCtx.Query.GroupBy = []string{"key"}
				tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				storageCtx.Query.GroupBy = nil
				storageCtx.GroupByTags = nil
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			ctx.ForkGroupingTask()
			ctx.CompleteGroupingTask()

			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestLeafGroupingContext_getTagValues(t *testing.T) {
	ctx := &LeafGroupingContext{
		tagsMap: map[string]string{
			"key": "value",
		},
		tagValuesMap: []map[uint32]string{{1: "value1"}},
		tagValues:    make([]string, 1),
	}
	t.Run("get value from cache", func(t *testing.T) {
		assert.Equal(t, "value", ctx.getTagValues("key"))
	})
	t.Run("get value", func(t *testing.T) {
		assert.Equal(t, "value1", ctx.getTagValues(string([]byte{1, 0, 0, 0})))
	})
	t.Run("tag value not found", func(t *testing.T) {
		assert.Equal(t, tagValueNotFound, ctx.getTagValues(string([]byte{2, 0, 0, 0})))
	})
}

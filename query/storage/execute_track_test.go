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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestGroupingExecuteTrack_submitTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newCollectTagValuesTaskFunc = newCollectTagValuesTask
		ctrl.Finish()
	}()

	db := tsdb.NewMockDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()

	qf := flow.NewMockStorageQueryFlow(ctrl)
	ctx := &executeContext{
		storageExecuteCtx: &flow.StorageExecuteContext{
			GroupByTags: tag.Metas{{Key: "ip"}},
		},
		database: db,
	}

	qf.EXPECT().Submit(gomock.Any(), gomock.Any()).DoAndReturn(func(_ flow.Stage, fn func()) {
		fn()
	}).AnyTimes()
	task := flow.NewMockQueryTask(ctrl)
	newCollectTagValuesTaskFunc = func(ctx *executeContext, metadata metadb.Metadata,
		tagKey tag.Meta, tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) flow.QueryTask {
		return task
	}

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "tag value ids not found",
			prepare: func() {
				ctx.storageExecuteCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf()}
				qf.EXPECT().ReduceTagValues(0, nil)
			},
		},
		{
			name: "find tag value id failure",
			prepare: func() {
				ctx.storageExecuteCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf(1)}
				task.EXPECT().Run().Return(fmt.Errorf("err"))
				qf.EXPECT().Complete(fmt.Errorf("err"))
			},
		},
		{
			name: "find tag value id successfully",
			prepare: func() {
				ctx.storageExecuteCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf(1)}
				task.EXPECT().Run().Return(nil)
				newCollectTagValuesTaskFunc = func(ctx *executeContext, metadata metadb.Metadata,
					tagKey tag.Meta, tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) flow.QueryTask {
					tagValues[1] = "1.1.1.1"
					return task
				}
				qf.EXPECT().ReduceTagValues(0, map[uint32]string{1: "1.1.1.1"})
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			track := newGroupingExecuteTrack(ctx, qf)
			track.submitTask(flow.GroupingStage, func() {})
			// collect tag value once
			track.collectGroupByTagValues()
		})
	}
}

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
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestLeafExecuteContext_SendResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := tsdb.NewMockDatabase(ctrl)
	taskServerFct := rpc.NewMockTaskServerFactory(ctrl)
	stream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	leaf := &models.Leaf{Receivers: []models.StatelessNode{{}}}
	cases := []struct {
		name    string
		in      error
		prepare func(ctx *LeafExecuteContext)
		assert  func()
	}{
		{
			name: "send response with err",
			in:   fmt.Errorf("err"),
			prepare: func(ctx *LeafExecuteContext) {
				leaf.Receivers = nil
			},
		},
		{
			name: "not found send stream",
			in:   fmt.Errorf("err"),
			prepare: func(ctx *LeafExecuteContext) {
				taskServerFct.EXPECT().GetStream(gomock.Any()).Return(nil)
			},
		},
		{
			name: "send response failure",
			in:   fmt.Errorf("err"),
			prepare: func(ctx *LeafExecuteContext) {
				taskServerFct.EXPECT().GetStream(gomock.Any()).Return(stream)
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
			},
		},
		{
			name: "send response with grouping",
			in:   nil,
			prepare: func(ctx *LeafExecuteContext) {
				ctx.StorageExecuteCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf(1, 2)}
				ctx.GroupingCtx.collectGroupingTagsCompleted = make(chan struct{})
				taskServerFct.EXPECT().GetStream(gomock.Any()).Return(stream)
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
				time.AfterFunc(time.Millisecond*100, func() {
					close(ctx.GroupingCtx.collectGroupingTagsCompleted)
				})
			},
		},
		{
			name: "time out",
			in:   nil,
			prepare: func(ctx *LeafExecuteContext) {
				ctx.StorageExecuteCtx.GroupingTagValueIDs = []*roaring.Bitmap{roaring.BitmapOf(1, 2)}
				ctx.GroupingCtx.collectGroupingTagsCompleted = make(chan struct{})
				taskServerFct.EXPECT().GetStream(gomock.Any()).Return(stream)
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
				time.AfterFunc(time.Millisecond*100, func() {
					ctx.TaskCtx.Release()
				})
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				leaf.Receivers = []models.StatelessNode{{}}
			}()
			c, cancel := context.WithCancel(context.TODO())
			taskCtx := &flow.TaskContext{
				Ctx:    c,
				Cancel: cancel,
			}
			ctx := NewLeafExecuteContext(taskCtx, &stmtpkg.Query{},
				&protoCommonV1.TaskRequest{}, taskServerFct, leaf, db)

			if tt.prepare != nil {
				tt.prepare(ctx)
			}

			ctx.SendResponse(tt.in)
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

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

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataContext_WaitResponse(t *testing.T) {
	t.Run("task complete", func(t *testing.T) {
		ctx := NewMetadataContext(&MetadataDeps{
			Ctx:       context.TODO(),
			Statement: &stmt.MetricMetadata{},
		})
		ctx.SetTracker(tracker.NewStageTracker(flow.NewTaskContextWithTimeout(context.TODO(), time.Minute)))
		go func() {
			ctx.Complete(fmt.Errorf("err"))
		}()
		rs, err := ctx.WaitResponse()
		assert.Error(t, err)
		assert.Nil(t, rs)
	})
	t.Run("task cancel", func(t *testing.T) {
		c, cancel := context.WithCancel(context.TODO())
		ctx := NewMetadataContext(&MetadataDeps{
			Ctx:       c,
			Statement: &stmt.MetricMetadata{},
		})
		go func() {
			cancel()
		}()
		rs, err := ctx.WaitResponse()
		assert.Equal(t, constants.ErrTimeout, err)
		assert.Nil(t, rs)
	})
}

func TestMetadataContext_MakePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chooseMgr := flow.NewMockNodeChoose(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "choose fail",
			prepare: func() {
				chooseMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "no target",
			prepare: func() {
				chooseMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "plan invalid",
			prepare: func() {
				chooseMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{}}, nil)
			},
			wantErr: true,
		},
		{
			name: "make plan successfully",
			prepare: func() {
				chooseMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).
					Return([]*models.PhysicalPlan{{Database: "test", Targets: []*models.Target{{}}}}, nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			ctx := NewMetadataContext(&MetadataDeps{
				Ctx:         context.TODO(),
				Request:     &models.Request{},
				Statement:   &stmt.MetricMetadata{},
				Choose:      chooseMgr,
				CurrentNode: models.StatelessNode{},
			})
			err := ctx.MakePlan()
			if tt.wantErr != (err != nil) {
				t.Fatalf("%s fail", tt.name)
			}
		})
	}
}

func TestMetadataContext_HandleResponse(t *testing.T) {
	ctx := NewMetadataContext(&MetadataDeps{
		Statement: &stmt.MetricMetadata{},
	})
	ctx.SetTracker(tracker.NewStageTracker(flow.NewTaskContextWithTimeout(context.TODO(), time.Minute)))
	ctx.HandleResponse(&protoCommonV1.TaskResponse{}, "leaf")
}

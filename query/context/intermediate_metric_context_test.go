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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	commonmodels "github.com/lindb/common/models"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestIntermediateMetricContext_WaitResponse(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		metricCtx := NewIntermediateMetricContext(ctx, nil, nil, nil,
			models.StatelessNode{}, &models.PhysicalPlan{}, &stmt.Query{}, []string{"root"})
		go func() {
			cancel()
		}()
		resp, err := metricCtx.WaitResponse()
		assert.Nil(t, resp)
		assert.Equal(t, constants.ErrTimeout, err)
	})
	t.Run("complete with result", func(t *testing.T) {
		metricCtx := NewIntermediateMetricContext(context.TODO(), nil, nil,
			&protoCommonV1.TaskRequest{}, models.StatelessNode{}, &models.PhysicalPlan{},
			&stmt.Query{}, []string{"root"})
		go func() {
			close(metricCtx.doneCh)
		}()
		resp, err := metricCtx.WaitResponse()
		assert.NotNil(t, resp)
		assert.NoError(t, err)
	})
}

func TestIntermediateMetricContext_MakePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	stateMgr := broker.NewMockStateManager(ctrl)
	cfg := models.Database{
		Option: &option.DatabaseOption{
			Intervals: option.Intervals{
				{Interval: timeutil.Interval(commontimeutil.OneSecond)},
				{Interval: timeutil.Interval(commontimeutil.OneMinute)},
			},
		},
	}
	metricCtx := NewIntermediateMetricContext(context.TODO(), nil, stateMgr,
		&protoCommonV1.TaskRequest{}, models.StatelessNode{}, &models.PhysicalPlan{},
		&stmt.Query{}, []string{"root"})
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "choose plan failure",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "no replica",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "database config not found",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{}}, nil)
				stateMgr.EXPECT().GetDatabaseCfg(gomock.Any()).Return(models.Database{}, false)
			},
			wantErr: true,
		},
		{
			name: "plan invalid",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{}}, nil)
				stateMgr.EXPECT().GetDatabaseCfg(gomock.Any()).Return(cfg, true)
			},
			wantErr: true,
		},
		{
			name: "make plan successfully",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{
					Database: "test",
					Targets:  []*models.Target{{}},
				}}, nil)
				stateMgr.EXPECT().GetDatabaseCfg(gomock.Any()).Return(cfg, true)
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			err := metricCtx.MakePlan()
			if (err != nil) != tt.wantErr {
				t.Fatalf("run test failure, case: %s", tt.name)
			}
		})
	}
}

func TestIntermediateMetricContext_makeTaskResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metricCtx := NewIntermediateMetricContext(context.TODO(), nil, nil,
		&protoCommonV1.TaskRequest{}, models.StatelessNode{}, &models.PhysicalPlan{},
		&stmt.Query{}, []string{"root"})
	metricCtx.stats = &commonmodels.NodeStats{}
	metricCtx.aggregatorSpecs = map[string]*protoCommonV1.AggregatorSpec{"f": {}}
	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
	groupIt := series.NewMockGroupedIterator(ctrl)
	it := series.NewMockIterator(ctrl)
	groupAgg.EXPECT().ResultSet().Return(series.GroupedIterators{groupIt})
	groupIt.EXPECT().HasNext().Return(true)
	groupIt.EXPECT().Next().Return(it)
	it.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
	groupIt.EXPECT().HasNext().Return(true)
	groupIt.EXPECT().Next().Return(it)
	it.EXPECT().MarshalBinary().Return([]byte{1, 2, 2}, nil)
	it.EXPECT().FieldName().Return(field.Name("f"))
	groupIt.EXPECT().Tags().Return("tags")
	groupIt.EXPECT().HasNext().Return(false)
	metricCtx.groupAgg = groupAgg
	resp := metricCtx.makeTaskResponse()
	assert.NotNil(t, resp)
}

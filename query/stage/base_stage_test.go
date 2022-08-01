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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/models"
)

type mockPool struct {
}

func (p *mockPool) Submit(_ context.Context, task *concurrent.Task) {
	task.Exec()
}
func (p *mockPool) Stopped() bool {
	return false
}
func (p *mockPool) Stop() {
}

func TestBaseStage_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := &mockPool{}
	s := &baseStage{
		ctx:       context.TODO(),
		stageType: Grouping,
		execPool:  pool,
	}
	s.Complete()
	assert.Equal(t, Grouping, s.Type())
	assert.Nil(t, s.NextStages())

	cases := []struct {
		name           string
		plan           PlanNode
		prepare        func(p *MockPlanNode)
		completeHandle func()
		errHandler     func(err error)
	}{
		{
			name: "empty plan node",
			plan: nil,
			completeHandle: func() {
				assert.True(t, true)
			},
		},
		{
			name: "execute plan failure",
			plan: NewMockPlanNode(ctrl),
			prepare: func(p *MockPlanNode) {
				p.EXPECT().IgnoreNotFound().Return(false)
				p.EXPECT().Execute().Return(fmt.Errorf("err"))
			},
			errHandler: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "execute plan failure, ignore not found err",
			plan: NewMockPlanNode(ctrl),
			prepare: func(p *MockPlanNode) {
				p.EXPECT().IgnoreNotFound().Return(true)
				p.EXPECT().Execute().Return(constants.ErrNotFound)
			},
			completeHandle: func() {
				assert.True(t, true)
			},
		},
		{
			name: "execute children plan failure",
			plan: NewMockPlanNode(ctrl),
			prepare: func(p *MockPlanNode) {
				p.EXPECT().Execute().Return(nil)
				p1 := NewMockPlanNode(ctrl)
				p1.EXPECT().IgnoreNotFound().Return(false)
				p1.EXPECT().Execute().Return(fmt.Errorf("err"))
				p.EXPECT().Children().Return([]PlanNode{p1})
			},
			errHandler: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "execute plan sync",
			plan: NewMockPlanNode(ctrl),
			prepare: func(p *MockPlanNode) {
				s.ctx = nil
				p.EXPECT().Execute().Return(nil)
				p.EXPECT().Children().Return(nil)
			},
			completeHandle: func() {
				assert.True(t, true)
			},
		},
		{
			name: "execute plan successfully",
			plan: NewMockPlanNode(ctrl),
			prepare: func(p *MockPlanNode) {
				p.EXPECT().Execute().Return(nil)
				p1 := NewMockPlanNode(ctrl)
				p1.EXPECT().Execute().Return(nil)
				p1.EXPECT().Children().Return(nil)
				p.EXPECT().Children().Return([]PlanNode{p1})
			},
			completeHandle: func() {
				assert.True(t, true)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				s.ctx = context.TODO()
			}()
			if tt.prepare != nil {
				tt.prepare(tt.plan.(*MockPlanNode))
			}
			s.Execute(tt.plan, tt.completeHandle, tt.errHandler)
		})
	}
}

func TestBaseStage_Track(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := &mockPool{}
	s := &baseStage{
		ctx:       context.TODO(),
		stageType: Grouping,
		execPool:  pool,
	}
	s.Track()

	p := NewMockPlanNode(ctrl)
	p.EXPECT().ExecuteWithStats().Return(&models.OperatorStats{}, nil)
	p.EXPECT().Children().Return(nil)
	s.Execute(p, func() {
	}, func(err error) {
	})
	assert.NotNil(t, s.Stats())
}

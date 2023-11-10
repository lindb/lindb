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

package query

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/stage"
	trackerpkg "github.com/lindb/lindb/query/tracker"
)

func TestPipeline_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tracker := trackerpkg.NewStageTracker(flow.NewTaskContextWithTimeout(context.TODO(), time.Minute))
	t.Run("execute nil stage", func(_ *testing.T) {
		p := NewExecutePipeline(tracker, nil)
		p.Execute(nil)
	})
	t.Run("execute stage after pipeline completed", func(_ *testing.T) {
		p := NewExecutePipeline(tracker, nil)
		p1 := p.(*pipeline)
		s := stage.NewMockStage(ctrl)
		p1.sm.completed.Store(true)
		p.Execute(s)
	})
	t.Run("stage execute successfully", func(t *testing.T) {
		p := NewExecutePipeline(tracker, func(err error) {
			assert.NoError(t, err)
		})
		s := stage.NewMockStage(ctrl)
		s2 := stage.NewMockStage(ctrl)
		s.EXPECT().Plan()
		s.EXPECT().NextStages().Return([]stage.Stage{s2})
		s.EXPECT().Complete()
		s.EXPECT().Stats()
		s.EXPECT().Identifier()
		s.EXPECT().IsAsync().Return(true)
		s.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Do(
			func(_ stage.PlanNode, completeFn func(), _ func(err error)) {
				completeFn()
			})
		s2.EXPECT().Plan()
		s2.EXPECT().NextStages().Return([]stage.Stage{nil})
		s2.EXPECT().Complete()
		s2.EXPECT().Stats()
		s2.EXPECT().Identifier()
		s2.EXPECT().IsAsync().Return(true)
		s2.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Do(
			func(_ stage.PlanNode, completeFn func(), _ func(err error)) {
				completeFn()
			})
		p.Execute(s)
		assert.NotNil(t, p.Stats())
	})
	t.Run("stage execute failure", func(t *testing.T) {
		p := NewExecutePipeline(tracker, func(err error) {
			assert.Error(t, err)
		})
		s := stage.NewMockStage(ctrl)
		s.EXPECT().Plan()
		s.EXPECT().Stats()
		s.EXPECT().Identifier()
		s.EXPECT().IsAsync().Return(true)
		s.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Do(
			func(_ stage.PlanNode, _ func(), errFn func(err error)) {
				errFn(fmt.Errorf("err"))
			})
		s.EXPECT().Complete()
		p.Execute(s)
	})
	t.Run("panic", func(t *testing.T) {
		p := NewExecutePipeline(tracker, func(err error) {
			assert.Error(t, err)
		})
		s := stage.NewMockStage(ctrl)
		s.EXPECT().Identifier().Do(func() string {
			panic("xx")
		})
		p.Execute(s)
	})
}

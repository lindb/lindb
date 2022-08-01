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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/stage"
)

func TestPipeline_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("execute nil stage", func(_ *testing.T) {
		p := NewExecutePipeline(false, nil)
		p.Execute(nil)
	})
	t.Run("execute stage after pipeline completed", func(_ *testing.T) {
		p := NewExecutePipeline(false, nil)
		p1 := p.(*pipeline)
		s := stage.NewMockStage(ctrl)
		p1.sm.completed.Store(true)
		p.Execute(s)
	})
	t.Run("stage execute successfully", func(t *testing.T) {
		p := NewExecutePipeline(true, func(_ []*models.StageStats, err error) {
			assert.NoError(t, err)
		})
		s := stage.NewMockStage(ctrl)
		s.EXPECT().Plan()
		s.EXPECT().NextStages().Return([]stage.Stage{nil})
		s.EXPECT().Complete()
		s.EXPECT().Track()
		s.EXPECT().Stats()
		s.EXPECT().Identifier()
		s.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Do(
			func(_ stage.PlanNode, completeFn func(), _ func(err error)) {
				completeFn()
			})
		p.Execute(s)
	})
	t.Run("stage execute failure", func(t *testing.T) {
		p := NewExecutePipeline(true, func(_ []*models.StageStats, err error) {
			assert.Error(t, err)
		})
		s := stage.NewMockStage(ctrl)
		s.EXPECT().Plan()
		s.EXPECT().Track()
		s.EXPECT().Stats()
		s.EXPECT().Identifier()
		s.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Do(
			func(_ stage.PlanNode, _ func(), errFn func(err error)) {
				errFn(fmt.Errorf("err"))
			})
		p.Execute(s)
	})
}

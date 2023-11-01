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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/query/context"
)

func TestPhysicalPlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskCtx := context.NewMockTaskContext(ctrl)
	taskCtx.EXPECT().MakePlan().Return(nil)
	op := NewPhysicalPlan(taskCtx)
	assert.NoError(t, op.Execute())
	assert.Equal(t, "Physical Plan", op.Identifier())
}

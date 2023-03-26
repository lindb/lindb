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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestDataLoadStage_Plan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{}).AnyTimes()
	rs := flow.NewMockFilterResultSet(ctrl)

	now := timeutil.Now()
	stage := NewDataLoadStage(
		&context.LeafExecuteContext{
			TaskCtx:  &flow.TaskContext{},
			Database: db,
		},
		&flow.DataLoadContext{
			ShardExecuteCtx: &flow.ShardExecuteContext{
				StorageExecuteCtx: &flow.StorageExecuteContext{
					Query: &stmt.Query{
						Interval:        1,
						IntervalRatio:   1.0,
						StorageInterval: 1,
					},
				},
			},
		},
		&flow.TimeSegmentResultSet{
			FilterRS:   []flow.FilterResultSet{rs},
			FamilyTime: now,
		})
	assert.NotEmpty(t, stage.Plan())
	id := fmt.Sprintf("Data Load[%s]", timeutil.FormatTimestamp(now, timeutil.DataTimeFormat2))
	assert.Equal(t, id, stage.Identifier())
}

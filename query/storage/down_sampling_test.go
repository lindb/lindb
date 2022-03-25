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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func Test_buildDownSamplingTimeRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	opt := &option.DatabaseOption{Intervals: option.Intervals{{Interval: timeutil.Interval(10 * timeutil.OneSecond)}}}
	db.EXPECT().GetOption().Return(opt)

	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{
				TimeRange: timeutil.TimeRange{
					Start: 35 * timeutil.OneSecond,
					End:   65 * timeutil.OneSecond},
				Interval: timeutil.Interval(30 * timeutil.OneSecond),
			},
		},
	}

	buildDownSamplingTimeRange(ctx)

	assert.Equal(t, 3, ctx.storageExecuteCtx.QueryIntervalRatio)
	assert.Equal(t, timeutil.Interval(30*timeutil.OneSecond), ctx.storageExecuteCtx.QueryInterval)
	assert.Equal(t, timeutil.TimeRange{
		Start: 30 * timeutil.OneSecond,
		End:   60 * timeutil.OneSecond,
	}, ctx.storageExecuteCtx.QueryTimeRange)
}

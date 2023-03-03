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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/tsdb"
)

func TestSeriesLimit_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).MaxTimes(3)
	ctx := flow.NewShardExecuteContext(nil)
	op := NewSeriesLimit(ctx, shard)
	assert.NoError(t, op.Execute())

	ctx.SeriesIDsAfterFiltering.Add(1)
	ctx.SeriesIDsAfterFiltering.Add(2)
	limit := models.NewDefaultLimits()
	db.EXPECT().GetLimits().Return(limit).MaxTimes(3)
	assert.NoError(t, op.Execute())

	limit.MaxSeriesPerQuery = 1
	assert.Equal(t, constants.ErrTooManySeriesFound, op.Execute())
	limit.MaxSeriesPerQuery = 0
	assert.NoError(t, op.Execute())
}

func TestSeriesLimit_Identifier(t *testing.T) {
	assert.Equal(t, "Series Limit", NewSeriesLimit(nil, nil).Identifier())
}

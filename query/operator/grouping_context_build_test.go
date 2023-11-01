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

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

func TestGroupingContextBuild_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB)
	indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(nil)
	ctx := flow.NewShardExecuteContext(nil)
	ctx.TimeSegmentContext.SeriesIDs.Add(1)
	op := NewGroupingContextBuild(ctx, shard)
	assert.NoError(t, op.Execute())

	// no series found
	ctx.TimeSegmentContext.SeriesIDs.Clear()
	op = NewGroupingContextBuild(ctx, shard)
	assert.NoError(t, op.Execute())
}

func TestGroupingContextBuild_Identifier(t *testing.T) {
	assert.Equal(t, "Grouping Context Build", NewGroupingContextBuild(nil, nil).Identifier())
}

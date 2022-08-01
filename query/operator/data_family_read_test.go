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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb"
)

func TestDataFamilyRead_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := tsdb.NewMockDataFamily(ctrl)
	shardCtx := &flow.ShardExecuteContext{
		TimeSegmentContext: flow.NewTimeSegmentContext(),
	}

	t.Run("filter data failure", func(t *testing.T) {
		op := NewDataFamilyRead(shardCtx, family)
		family.EXPECT().Filter(gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})

	t.Run("filter data success", func(t *testing.T) {
		rs := flow.NewMockFilterResultSet(ctrl)
		rs.EXPECT().FamilyTime().Return(int64(1010))
		rs.EXPECT().SlotRange().Return(timeutil.SlotRange{})
		rs.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2, 3))
		op := NewDataFamilyRead(shardCtx, family)
		family.EXPECT().Interval().Return(timeutil.Interval(10))
		family.EXPECT().Filter(gomock.Any()).Return([]flow.FilterResultSet{rs}, nil)
		assert.NoError(t, op.Execute())
	})

	op := NewDataFamilyRead(nil, nil)
	assert.NotEmpty(t, op.Identifier())
}

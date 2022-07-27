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
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb"
)

func TestDataFamilyReader_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := tsdb.NewMockDataFamily(ctrl)
	t.Run("filter data failure", func(t *testing.T) {
		op := NewDataFamilyRead(nil, family)
		family.EXPECT().Filter(gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("found data", func(t *testing.T) {
		op := NewDataFamilyRead(&flow.ShardExecuteContext{
			TimeSegmentContext: flow.NewTimeSegmentContext(),
		}, family)
		rs := flow.NewMockFilterResultSet(ctrl)
		family.EXPECT().Filter(gomock.Any()).Return([]flow.FilterResultSet{rs}, nil)
		family.EXPECT().Interval().Return(timeutil.Interval(10))
		rs.EXPECT().FamilyTime().Return(int64(10))
		rs.EXPECT().SeriesIDs().Return(roaring.New())
		rs.EXPECT().SlotRange().Return(timeutil.SlotRange{})
		assert.NoError(t, op.Execute())
	})
}

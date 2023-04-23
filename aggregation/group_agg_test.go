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

package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestGroupByAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	gIt := series.NewMockGroupedIterator(ctrl)
	sIt := series.NewMockIterator(ctrl)
	fIt := series.NewMockFieldIterator(ctrl)

	now, _ := timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "field agg not found",
			prepare: func() {
				gIt.EXPECT().Tags().Return("tags")
				gIt.EXPECT().HasNext().Return(true)
				gIt.EXPECT().Next().Return(sIt)
				sIt.EXPECT().FieldName().Return(field.Name("c"))
				gIt.EXPECT().HasNext().Return(false)
			},
		},
		{
			name: "field data not found",
			prepare: func() {
				gIt.EXPECT().Tags().Return("tags")
				gIt.EXPECT().HasNext().Return(true)
				gIt.EXPECT().Next().Return(sIt)
				sIt.EXPECT().FieldName().Return(field.Name("b"))
				sIt.EXPECT().HasNext().Return(true)
				sIt.EXPECT().Next().Return(int64(1), nil)
				sIt.EXPECT().HasNext().Return(false)
				gIt.EXPECT().HasNext().Return(false)
			},
		},
		{
			name: "merge field data",
			prepare: func() {
				gIt.EXPECT().Tags().Return("tags")
				gIt.EXPECT().HasNext().Return(true)
				gIt.EXPECT().Next().Return(sIt)
				sIt.EXPECT().FieldName().Return(field.Name("b"))
				sIt.EXPECT().HasNext().Return(true)
				sIt.EXPECT().Next().Return(familyTime, fIt)
				fIt.EXPECT().HasNext().Return(false)
				sIt.EXPECT().HasNext().Return(false)
				gIt.EXPECT().HasNext().Return(false)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			agg := NewGroupingAggregator(
				timeutil.Interval(timeutil.OneSecond),
				1,
				timeutil.TimeRange{
					Start: now,
					End:   now + 3*timeutil.OneHour,
				},
				AggregatorSpecs{
					NewAggregatorSpec("b", field.SumField),
					NewAggregatorSpec("a", field.SumField),
				})
			if tt.prepare != nil {
				tt.prepare()
			}

			agg.Aggregate(gIt)
			rs := agg.ResultSet()
			assert.NotNil(t, rs)
			assert.NotNil(t, agg.Fields())
		})
	}

	agg := NewGroupingAggregator(
		timeutil.Interval(timeutil.OneSecond),
		1,
		timeutil.TimeRange{
			Start: now,
			End:   now + 3*timeutil.OneHour,
		},
		AggregatorSpecs{})
	rs := agg.ResultSet()
	assert.Nil(t, rs)
	assert.Equal(t, timeutil.Interval(timeutil.OneSecond), agg.Interval())
	assert.Empty(t, agg.Fields())
	assert.Equal(t,
		timeutil.TimeRange{
			Start: now,
			End:   now + 3*timeutil.OneHour,
		}, agg.TimeRange())
}

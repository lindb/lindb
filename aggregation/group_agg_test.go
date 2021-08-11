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

//func TestGroupByAggregator_Aggregate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	gIt := series.NewMockGroupedIterator(ctrl)
//	sIt := series.NewMockIterator(ctrl)
//	fIt := series.NewMockFieldIterator(ctrl)
//
//	now, _ := timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
//	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
//	agg := NewGroupingAggregator(
//		timeutil.Interval(timeutil.OneSecond),
//		1,
//		timeutil.TimeRange{
//			Start: now,
//			End:   now + 3*timeutil.OneHour,
//		},
//		AggregatorSpecs{
//			NewAggregatorSpec("b", field.SumField),
//			NewAggregatorSpec("a", field.SumField),
//		})
//
//	gomock.InOrder(
//		gIt.EXPECT().Tags().Return("1.1.1.1"),
//		gIt.EXPECT().HasNext().Return(true),
//		gIt.EXPECT().Next().Return(sIt),
//		// series it
//		sIt.EXPECT().FieldName().Return(field.Name("a")),
//		sIt.EXPECT().HasNext().Return(true),
//		sIt.EXPECT().Next().Return(familyTime, fIt),
//		//fIt.EXPECT().HasNext().Return(false),
//		fIt.EXPECT().MarshalBinary().Return([]byte("abcdefg"), nil),
//		// series it
//		sIt.EXPECT().HasNext().Return(true),
//		sIt.EXPECT().Next().Return(familyTime, nil),
//		sIt.EXPECT().HasNext().Return(false),
//		// series it
//		gIt.EXPECT().HasNext().Return(true),
//		gIt.EXPECT().Next().Return(sIt),
//		sIt.EXPECT().FieldName().Return(field.Name("c")),
//
//		gIt.EXPECT().HasNext().Return(false),
//	)
//	agg.Aggregate(gIt)
//	rs := agg.ResultSet()
//	assert.Equal(t, 1, len(rs))

//
//gomock.InOrder(
//	gIt.EXPECT().Tags().Return("1.1.1.2"),
//	gIt.EXPECT().HasNext().Return(true),
//	gIt.EXPECT().Next().Return(sIt),
//	// series it
//	sIt.EXPECT().FieldName().Return(field.Name("a")),
//	sIt.EXPECT().FieldType().Return(field.SumField),
//	sIt.EXPECT().HasNext().Return(true),
//	sIt.EXPECT().Next().Return(familyTime, fIt),
//	fIt.EXPECT().HasNext().Return(false),
//	// series it
//	sIt.EXPECT().HasNext().Return(true),
//	sIt.EXPECT().Next().Return(familyTime, nil),
//	sIt.EXPECT().HasNext().Return(false),
//	// series it
//	gIt.EXPECT().HasNext().Return(true),
//	gIt.EXPECT().Next().Return(sIt),
//	sIt.EXPECT().FieldName().Return(field.Name("c")),
//	sIt.EXPECT().FieldType().Return(field.SumField),
//
//	gIt.EXPECT().HasNext().Return(false),
//)
//agg.Aggregate(gIt)
//
//rs = agg.ResultSet()
//assert.Equal(t, 2, len(rs))

//agg = NewGroupingAggregator(
//	timeutil.Interval(timeutil.OneSecond),
//	1,
//	timeutil.TimeRange{
//		Start: now,
//		End:   now + 3*timeutil.OneHour,
//	},
//	AggregatorSpecs{})
//rs = agg.ResultSet()
//assert.Nil(t, rs)
//}

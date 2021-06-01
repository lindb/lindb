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

//TODO need impl
//func TestFieldAggregator_Aggregate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")
//
//	aggSpec := NewDownSamplingSpec("f", field.SumField)
//	aggSpec.AddFunctionType(function.Sum)
//
//	selector1 := selector.NewIndexSlotSelector(15, 55, 1)
//	agg := NewFieldAggregator(baseTime, selector1)
//	block, ok := agg.GetBlock(1, func() series.Block {
//		return nil
//	})
//	assert.False(t, ok)
//	assert.Nil(t, block)
//	it := MockSumFieldIterator(ctrl, field.PrimitiveID(1), map[int]interface{}{
//		5:  5.5,
//		15: 5.6,
//		17: 5.7,
//		16: 5.8,
//		56: 5.9,
//	})
//	agg.Aggregate(it)
//
//	expect := map[int]float64{
//		15: 5.6,
//		16: 5.8,
//		17: 5.7,
//	}
//
//	start, fieldIt := agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	AssertPrimitiveIt(t, fieldIt.Next(), expect)
//	assert.False(t, fieldIt.HasNext())
//
//	it = MockSumFieldIterator(ctrl, field.PrimitiveID(1), map[int]interface{}{
//		5:  5.5,
//		15: 5.5,
//		19: 5.5,
//		16: 1.1,
//		56: 5.5,
//	})
//	agg.Aggregate(it)
//
//	expect = map[int]float64{
//		15: 11.1,
//		16: 6.9,
//		17: 5.7,
//		19: 5.5,
//	}
//
//	start, fieldIt = agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	AssertPrimitiveIt(t, fieldIt.Next(), expect)
//	assert.False(t, fieldIt.HasNext())
//
//	// not match query time range case 1
//	agg.reset()
//	it = MockSumFieldIterator(ctrl, field.PrimitiveID(1), map[int]interface{}{
//		4:  1.1,
//		56: 5.5,
//	})
//	agg.Aggregate(it)
//	start, fieldIt = agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	assert.False(t, fieldIt.Next().HasNext())
//	assert.False(t, fieldIt.HasNext())
//}

//func TestDownSamplingFieldAggregator(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	//baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")
//
//	aggSpec := NewDownSamplingSpec("f", field.SummaryField)
//	aggSpec.AddFunctionType(function.Sum)
//	aggSpec.AddFunctionType(function.Max)
//	aggSpec.AddFunctionType(function.Avg)
//
//	agg := NewDownSamplingFieldAggregator(aggSpec, 2)
//	agg.Aggregate(nil)
//	it := series.NewMockFieldIterator(ctrl)
//	agg.Aggregate(it)
//	_, rs := agg.ResultSet()
//	assert.Nil(t, rs)
//	//assert.Equal(t, baseTime, startTime)
//	agg.reset()
//}
//
//func TestDownSamplingFieldAggregator_GetBlock(t *testing.T) {
//	aggSpec := NewDownSamplingSpec("f", field.SummaryField)
//	agg := NewDownSamplingFieldAggregator(aggSpec, 2)
//	block, ok := agg.GetBlock(2, func() series.Block {
//		return nil
//	})
//	assert.False(t, ok)
//	assert.Nil(t, block)
//
//	block, ok = agg.GetBlock(1, func() series.Block {
//		return series.NewBlock(0, 10)
//	})
//	assert.True(t, ok)
//	assert.NotNil(t, block)
//}

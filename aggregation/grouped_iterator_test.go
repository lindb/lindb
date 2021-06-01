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

//func TestGroupedIterator_HasNext(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	sIt1 := NewMockSeriesAggregator(ctrl)
//	sIt2 := NewMockSeriesAggregator(ctrl)
//	fIt := series.NewMockIterator(ctrl)
//	tagValues := "1.1.1.1,disk"
//	it := newGroupedIterator(tagValues, FieldAggregates{sIt1, sIt2})
//	gomock.InOrder(
//		sIt1.EXPECT().ResultSet().Return(fIt),
//		sIt2.EXPECT().ResultSet().Return(fIt),
//	)
//	assert.Equal(t, tagValues, it.Tags())
//	assert.True(t, it.HasNext())
//	assert.Equal(t, fIt, it.Next())
//	assert.True(t, it.HasNext())
//	assert.Equal(t, fIt, it.Next())
//	assert.False(t, it.HasNext())
//}

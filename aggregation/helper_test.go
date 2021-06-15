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

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

///////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

// MockSumFieldIterator returns mock an iterator of sum field
//func MockSumFieldIterator(ctrl *gomock.Controller, fieldID field.PrimitiveID, points map[int]interface{}) *series.MockFieldIterator {
//	it := series.NewMockFieldIterator(ctrl)
//	it.EXPECT().HasNext().Return(true)
//
//	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
//	it.EXPECT().Next().Return(primitiveIt)
//
//	primitiveIt.EXPECT().FieldID().Return(fieldID)
//	primitiveIt.EXPECT().AggType().Return(field.Sum)
//
//	var keys []int
//	for timeSlot := range points {
//		keys = append(keys, timeSlot)
//	}
//	sort.Slice(keys, func(i, j int) bool {
//		return keys[i] < keys[j]
//	})
//
//	for _, timeSlot := range keys {
//		primitiveIt.EXPECT().HasNext().Return(true)
//		primitiveIt.EXPECT().Next().Return(timeSlot, points[timeSlot])
//	}
//	// mock nil primitive iterator
//	it.EXPECT().HasNext().Return(true)
//	it.EXPECT().Next().Return(nil)
//
//	// return hasNext=>false, finish primitive iterator
//	primitiveIt.EXPECT().HasNext().Return(false).AnyTimes()
//
//	// sum field only has one primitive field
//	it.EXPECT().HasNext().Return(false).AnyTimes()
//	return it
//}

func AssertFieldIt(t *testing.T, it series.FieldIterator, expect map[int]float64) {
	count := 0
	for it.HasNext() {
		pIt := it.Next()
		for pIt.HasNext() {
			timeSlot, value := pIt.Next()
			assert.Equal(t, expect[timeSlot], value)
			count++
		}
	}
	assert.Equal(t, count, len(expect))
}

func generateFloatArray(values []float64) collections.FloatArray {
	if values == nil {
		return nil
	}
	floatArray := collections.NewFloatArray(len(values))
	for idx, value := range values {
		floatArray.SetValue(idx, value)
	}
	return floatArray
}

// mockSingleIterator returns mock an iterator of single field
func mockSingleIterator(ctrl *gomock.Controller, aggType field.AggType) series.FieldIterator {
	it := series.NewMockFieldIterator(ctrl)
	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	primitiveIt.EXPECT().AggType().Return(aggType)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(4, 4.0)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(50, 50.0)
	primitiveIt.EXPECT().HasNext().Return(false)
	it.EXPECT().HasNext().Return(false)
	return it
}

func mockTimeSeries(ctrl *gomock.Controller, startTime int64,
	fieldName field.Name, fieldType field.Type,
	aggType field.AggType,
) series.Iterator {
	timeSeries := series.NewMockIterator(ctrl)
	timeSeries.EXPECT().FieldType().Return(fieldType)
	timeSeries.EXPECT().FieldName().Return(fieldName)
	it := mockSingleIterator(ctrl, aggType)
	timeSeries.EXPECT().HasNext().Return(true)
	timeSeries.EXPECT().Next().Return(startTime, it)
	timeSeries.EXPECT().HasNext().Return(false)
	return timeSeries
}

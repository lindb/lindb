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

package fields

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestNewDynamicField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := NewDynamicField(field.SumField, 10, 10, 10)
	f.SetValue(mockSingleIterator(ctrl))
	values := f.GetDefaultValues()
	assert.Equal(t, 1, len(values))
	assert.Equal(t, 1.1, values[0].GetValue(4))
	assert.Equal(t, 1, values[0].Size())

	f.Reset()
	fIt := series.NewMockIterator(ctrl)
	fIt.EXPECT().HasNext().Return(true)
	fIt.EXPECT().Next().Return(int64(10), nil)
	fIt.EXPECT().HasNext().Return(false)
	f.SetValue(fIt)
	values = f.GetDefaultValues()
	assert.Equal(t, 1, len(values))
	assert.Equal(t, 0, values[0].Size())

	f = NewDynamicField(field.SumField, 10, 10, 10)
	f.SetValue(nil)
	values = f.GetDefaultValues()
	assert.Nil(t, values)
}

func TestDynamicField_UnknownType(t *testing.T) {
	f := NewDynamicField(field.Unknown, 10, 10, 10)
	values := f.GetDefaultValues()
	assert.Nil(t, values)
	values = f.GetValues(function.Sum)
	assert.Nil(t, values)
}

// mockSingleIterator returns mock an iterator of single field
func mockSingleIterator(ctrl *gomock.Controller) series.Iterator {
	fIt := series.NewMockIterator(ctrl)
	it := series.NewMockFieldIterator(ctrl)
	fIt.EXPECT().HasNext().Return(true)
	fIt.EXPECT().Next().Return(int64(10), it)
	fIt.EXPECT().HasNext().Return(false)
	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	it.EXPECT().HasNext().Return(false)
	primitiveIt.EXPECT().AggType().Return(field.Sum)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(4, 1.1)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(112, 1.1)
	primitiveIt.EXPECT().HasNext().Return(false)
	return fIt
}

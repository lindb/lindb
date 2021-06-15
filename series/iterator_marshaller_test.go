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

package series

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

func Test_MarshalBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := stream.NewBufferWriter(nil)
	writer.PutByte(byte(field.SumField))
	writer.PutVarint64(10)
	writer.PutVarint32(int32(2))
	writer.PutBytes([]byte{1, 2})
	writer.PutVarint64(10)
	writer.PutVarint32(int32(0))
	data, err := writer.Bytes()
	assert.NoError(t, err)

	it := NewMockIterator(ctrl)
	fIt := NewMockFieldIterator(ctrl)
	gomock.InOrder(
		it.EXPECT().FieldType().Return(field.SumField),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().MarshalBinary().Return([]byte{1, 2}, nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().MarshalBinary().Return([]byte{}, nil),
		it.EXPECT().HasNext().Return(false),
	)
	data2, err := MarshalIterator(it)
	assert.NoError(t, err)
	assert.Equal(t, data, data2)

	data2, err = MarshalIterator(nil)
	assert.NoError(t, err)
	assert.Nil(t, data2)

	gomock.InOrder(
		it.EXPECT().FieldType().Return(field.SumField),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().MarshalBinary().Return([]byte{1, 2}, nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().MarshalBinary().Return([]byte{}, fmt.Errorf("err")),
	)
	_, err = MarshalIterator(it)
	assert.Error(t, err)
}

package series

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

func TestEncodeSeries(t *testing.T) {
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
		fIt.EXPECT().Bytes().Return([]byte{1, 2}, nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().Bytes().Return([]byte{}, nil),
		it.EXPECT().HasNext().Return(false),
	)
	data2, err := EncodeSeries(it)
	assert.NoError(t, err)
	assert.Equal(t, data, data2)

	data2, err = EncodeSeries(nil)
	assert.NoError(t, err)
	assert.Nil(t, data2)

	gomock.InOrder(
		it.EXPECT().FieldType().Return(field.SumField),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().Bytes().Return([]byte{1, 2}, nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), nil),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().Bytes().Return([]byte{}, fmt.Errorf("err")),
	)
	_, err = EncodeSeries(it)
	assert.Error(t, err)
}

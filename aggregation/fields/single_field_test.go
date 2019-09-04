package fields

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/series"
)

func TestNewSingleField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := series.NewMockFieldIterator(ctrl)
	it.EXPECT().HasNext().Return(false)
	it.EXPECT().FieldType().Return(field.SumField)
	f := NewSingleField(10, it)
	assert.Nil(t, f)

	it = series.NewMockFieldIterator(ctrl)
	it.EXPECT().FieldType().Return(field.Unknown)
	f = NewSingleField(10, it)
	assert.Nil(t, f)

	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	primitiveIt.EXPECT().HasNext().Return(false)
	it.EXPECT().FieldType().Return(field.SumField).AnyTimes()

	f = NewSingleField(10, it)
	assert.NotNil(t, f)

	f = NewSingleField(10, mockSingleIterator(ctrl, field.SumField))
	assert.NotNil(t, f)
}

func TestSingleField_Sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := mockSingleIterator(ctrl, field.SumField)

	f := NewSingleField(10, it)
	assert.NotNil(t, f)
	assert.Nil(t, f.GetValues(function.Avg))
	assertFieldValues(t, f.GetDefaultValues())
	assertFieldValues(t, f.GetValues(function.Sum))
}

func TestSingleField_Max(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := mockSingleIterator(ctrl, field.MaxField)

	f := NewSingleField(10, it)
	assert.NotNil(t, f)
	assert.Nil(t, f.GetValues(function.Avg))
	assert.Nil(t, f.GetValues(function.Sum))
	assertFieldValues(t, f.GetDefaultValues())
	assertFieldValues(t, f.GetValues(function.Max))
}

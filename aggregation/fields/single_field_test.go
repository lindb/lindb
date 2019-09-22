package fields

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestNewSingleField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := series.NewMockFieldIterator(ctrl)
	it.EXPECT().HasNext().Return(false)
	f := NewSingleField(10, field.SumField, it)
	assert.Nil(t, f)

	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	primitiveIt.EXPECT().HasNext().Return(false)

	f = NewSingleField(10, field.SumField, it)
	assert.NotNil(t, f)

	f = NewSingleField(10, field.SumField, mockSingleIterator(ctrl))
	assert.NotNil(t, f)
}

func TestSingleField_Sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := mockSingleIterator(ctrl)

	f := NewSingleField(10, field.SumField, it)
	assert.NotNil(t, f)
	assert.Nil(t, f.GetValues(function.Avg))
	assertFieldValues(t, f.GetDefaultValues())
	assertFieldValues(t, f.GetValues(function.Sum))
}

func TestSingleField_Max(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := mockSingleIterator(ctrl)

	f := NewSingleField(10, field.MaxField, it)
	assert.NotNil(t, f)
	assert.Nil(t, f.GetValues(function.Avg))
	assert.Nil(t, f.GetValues(function.Sum))
	assertFieldValues(t, f.GetDefaultValues())
	assertFieldValues(t, f.GetValues(function.Max))
}

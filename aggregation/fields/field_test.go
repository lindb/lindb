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
	it.EXPECT().AggType().Return(field.Sum)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(4, 1.1)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(112, 1.1)
	it.EXPECT().HasNext().Return(false)
	return fIt
}

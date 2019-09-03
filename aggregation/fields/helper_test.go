package fields

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

// mockSingleIterator returns mock an iterator of single field
func mockSingleIterator(ctrl *gomock.Controller, fieldType field.Type) series.FieldIterator {
	it := series.NewMockFieldIterator(ctrl)
	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().FieldType().Return(fieldType)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(4, 1.1)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(112, 1.1)
	primitiveIt.EXPECT().HasNext().Return(false)
	return it
}

func assertFieldValues(t *testing.T, values []collections.FloatArray) {
	assert.Equal(t, 1, len(values))

	list := values[0]
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, 1.1, list.GetValue(4))
}

package fields

import (
	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

// mockSingleIterator returns mock an iterator of single field
func mockSingleIterator(ctrl *gomock.Controller) series.Iterator {
	fIt := series.NewMockIterator(ctrl)
	it := series.NewMockFieldIterator(ctrl)
	fIt.EXPECT().HasNext().Return(true)
	fIt.EXPECT().Next().Return(int64(10), it)
	fIt.EXPECT().HasNext().Return(false)
	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	primitiveIt.EXPECT().FieldID().Return(field.PrimitiveID(1))
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	it.EXPECT().HasNext().Return(false)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(4, 1.1)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(112, 1.1)
	primitiveIt.EXPECT().HasNext().Return(false)
	return fIt
}

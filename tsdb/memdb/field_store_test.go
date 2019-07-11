package memdb

import (
	"testing"

	"github.com/eleme/lindb/pkg/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_getSegmentStore(t *testing.T) {
	fStore := newFieldStore("sum", field.SumField)
	sStore, _ := fStore.getSegmentStore(11)
	assert.Nil(t, sStore)
	assert.Equal(t, field.SumField, fStore.getFieldType())
}

func Test_mustGetFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fStore := newFieldStore("sum", field.SumField)
	mockGen := makeMockIDGenerator(ctrl)
	assert.NotZero(t, fStore.mustGetFieldID(22, mockGen))
	assert.NotZero(t, fStore.mustGetFieldID(22, mockGen))
}

func Test_flushFieldTo_write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tw := makeMockTableWriter(ctrl)
	gen := makeMockIDGenerator(ctrl)
	p := makeMockPoint(ctrl)

	fStore := newFieldStore("sum", field.SumField)
	assert.Equal(t, fStore.getFamiliesCount(), 0)
	fStore.segments[2] = newSimpleFieldStore(field.GetAggFunc(field.Sum))

	// not exist in fs.segments
	fStore.flushFieldTo(tw, 32, 1, gen)
	// exist in fs.segments
	assert.Equal(t, fStore.getFamiliesCount(), 1)
	assert.Equal(t, fStore.getFamiliesCount(), 1)
	fStore.flushFieldTo(tw, 32, 2, gen)
	assert.Equal(t, fStore.getFamiliesCount(), 0)

	for _, f := range p.Fields() {
		fStore.write(newBlockStore(10), 5, 3, f)
		fStore.flushFieldTo(tw, 32, 2, gen)
	}
}

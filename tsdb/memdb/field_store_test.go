package memdb

import (
	"testing"

	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/timeutil"

	pb "github.com/eleme/lindb/rpc/proto/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_getSegmentStore(t *testing.T) {
	fStore := newFieldStore(field.SumField)
	sStore, _ := fStore.getSegmentStore(11)
	assert.Nil(t, sStore)
	assert.Equal(t, field.SumField, fStore.getFieldType())
}

func Test_mustGetFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fStore := newFieldStore(field.SumField)
	mockGen := makeMockIDGenerator(ctrl)
	assert.NotZero(t, fStore.mustGetFieldID(mockGen, 22, "sum"))
	assert.NotZero(t, fStore.mustGetFieldID(mockGen, 22, "sum"))
}

func Test_flushFieldTo_write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tw := makeMockTableWriter(ctrl)
	gen := makeMockIDGenerator(ctrl)
	p := &pb.Metric{
		Name:      "cpu.load",
		Timestamp: timeutil.Now(),
		Tags:      "idle",
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: 1.0}},
		},
	}

	fStore := newFieldStore(field.SumField)
	assert.Equal(t, fStore.getFamiliesCount(), 0)
	fStore.segments[2] = newSimpleFieldStore(field.GetAggFunc(field.Sum))

	// not exist in fs.segments
	fStore.flushFieldTo(tw, 32, gen, 1, "sum")
	// exist in fs.segments
	assert.Equal(t, fStore.getFamiliesCount(), 1)
	assert.Equal(t, fStore.getFamiliesCount(), 1)
	fStore.flushFieldTo(tw, 2, gen, 32, "sum")
	assert.Equal(t, fStore.getFamiliesCount(), 0)

	for _, f := range p.Fields {
		fStore.write(newBlockStore(10), 5, 3, f)
		fStore.flushFieldTo(tw, 32, gen, 2, "sum")
	}
}

package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

func TestFieldType_getFieldType(t *testing.T) {
	assert.Equal(t, field.Unknown, getFieldType(&pb.Field{}))
	assert.Equal(t, field.SumField, getFieldType(&pb.Field{Type: pb.FieldType_Sum}))
	assert.Equal(t, field.MinField, getFieldType(&pb.Field{Type: pb.FieldType_Min}))
	assert.Equal(t, field.MaxField, getFieldType(&pb.Field{Type: pb.FieldType_Max}))
	assert.Equal(t, field.GaugeField, getFieldType(&pb.Field{Type: pb.FieldType_Gauge}))
}

func TestBuildFieldKey(t *testing.T) {
	assert.Equal(t, FieldKey(uint32(10)|uint32(1)<<16), buildFieldKey(1, 10))
}

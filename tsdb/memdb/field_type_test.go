package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

func TestFieldType_getFieldType(t *testing.T) {
	assert.Equal(t, field.Unknown, getFieldType(&pb.Field{}))
	assert.Equal(t, field.SumField, getFieldType(&pb.Field{Field: &pb.Field_Sum{}}))
	assert.Equal(t, field.MinField, getFieldType(&pb.Field{Field: &pb.Field_Min{}}))
	assert.Equal(t, field.MaxField, getFieldType(&pb.Field{Field: &pb.Field_Max{}}))
	assert.Equal(t, field.GaugeField, getFieldType(&pb.Field{Field: &pb.Field_Gauge{}}))
	assert.Equal(t, field.SummaryField, getFieldType(&pb.Field{Field: &pb.Field_Summary{}}))
}

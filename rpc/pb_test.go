package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

func TestPBModel(t *testing.T) {
	metric := &pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:   "f1",
			Type:   pb.FieldType_Sum,
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}

	data, _ := metric.Marshal()
	metric2 := &pb.Metric{}
	_ = metric2.Unmarshal(data)
	assert.Equal(t, *metric, *metric2)
}

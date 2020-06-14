package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

func TestPBModel(t *testing.T) {
	metric := &pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 1.0,
		}},
	}

	data, _ := metric.Marshal()
	metric2 := &pb.Metric{}
	_ = metric2.Unmarshal(data)
	assert.Equal(t, *metric, *metric2)
}

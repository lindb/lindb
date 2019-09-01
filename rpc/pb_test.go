package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc/proto/field"
)

func TestPBModel(t *testing.T) {
	metric := &field.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*field.Field{
			{Name: "f1", Field: &field.Field_Sum{Sum: &field.Sum{Value: 1.0}}},
		},
	}

	data, _ := metric.Marshal()
	metric2 := &field.Metric{}
	_ = metric2.Unmarshal(data)
	assert.Equal(t, *metric, *metric2)
}

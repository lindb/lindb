package replication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

type mockIOWriter struct {
}

func (mw *mockIOWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("err")
}

func TestChunk_Append(t *testing.T) {
	chunk := newChunk(2)
	assert.False(t, chunk.IsFull())
	assert.True(t, chunk.IsEmpty())
	assert.Equal(t, 0, chunk.Size())
	chunk.Append(&pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 1.0,
		}},
		Tags: map[string]string{"host": "1.1.1.1"},
	})
	assert.False(t, chunk.IsEmpty())
	assert.False(t, chunk.IsFull())
	assert.Equal(t, 1, chunk.Size())
	chunk.Append(&pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 1.0,
		}},
		Tags: map[string]string{"host": "1.1.1.1"},
	})
	assert.False(t, chunk.IsEmpty())
	assert.True(t, chunk.IsFull())
	assert.Equal(t, 2, chunk.Size())
}

func TestChunk_MarshalBinary(t *testing.T) {
	c1 := newChunk(2)
	data, err := c1.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)
	testMarshal(c1, 2, t)
	testMarshal(c1, 1, t)

	c2 := c1.(*chunk)
	c2.writer = snappy.NewBufferedWriter(&mockIOWriter{})

	c2.Append(&pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 1.0,
		}},
		Tags: map[string]string{"host": "1.1.1.1"},
	})
	data, err = c2.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, data)

	// mock write err
	c2.writer = snappy.NewBufferedWriter(&mockIOWriter{})
	_, err = c2.writer.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	err = c2.writer.Flush()
	assert.Error(t, err)
	c2.Append(&pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 1.0,
		}},
		Tags: map[string]string{"host": "1.1.1.1"},
	})
	data, err = c2.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func testMarshal(chunk Chunk, size int, t *testing.T) {
	rs := pb.MetricList{}
	for i := 0; i < size; i++ {
		metric := &pb.Metric{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			Fields: []*pb.Field{{
				Name:  "f1",
				Type:  pb.FieldType_Sum,
				Value: 1.0,
			}},
			Tags: map[string]string{"host": "1.1.1.1"},
		}
		chunk.Append(metric)
		rs.Metrics = append(rs.Metrics, metric)
	}
	data, err := chunk.MarshalBinary()
	assert.NoError(t, err)
	assert.NotNil(t, data)
	reader := snappy.NewReader(bytes.NewReader(data))
	data, err = ioutil.ReadAll(reader)
	assert.NoError(t, err)
	var metricList pb.MetricList
	err = metricList.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, rs, metricList)
}

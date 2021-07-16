// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package replication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
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
	chunk.Append(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
		Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
	})
	assert.False(t, chunk.IsEmpty())
	assert.False(t, chunk.IsFull())
	assert.Equal(t, 1, chunk.Size())
	chunk.Append(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
		Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
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

	c2.Append(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
		Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
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
	c2.Append(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
		Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
	})
	data, err = c2.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func testMarshal(chunk Chunk, size int, t *testing.T) {
	rs := protoMetricsV1.MetricList{}
	for i := 0; i < size; i++ {
		metric := &protoMetricsV1.Metric{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
			Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
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
	var metricList protoMetricsV1.MetricList
	err = metricList.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, rs, metricList)
}

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

package replica

import (
	"testing"

	"github.com/lindb/common/pkg/ltoml"
	"github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/compress"
	"github.com/lindb/lindb/series/metric"
)

func makeTestBrokerRows() metric.BrokerRow {
	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	var brokerRow metric.BrokerRow
	_ = converter.ConvertTo(&protoMetricsV1.Metric{
		Name:      "cpu",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
		},
	}, &brokerRow)
	return brokerRow
}

func TestChunk_Append(t *testing.T) {
	chunk := newChunk(ltoml.Size(1024))
	assert.False(t, chunk.IsFull())
	assert.True(t, chunk.IsEmpty())
	assert.Equal(t, ltoml.Size(0), chunk.Size())

	row := makeTestBrokerRows()
	_, _ = row.WriteTo(chunk)
	assert.False(t, chunk.IsEmpty())
	assert.False(t, chunk.IsFull())
	assert.NotZero(t, chunk.Size())
	_, _ = row.WriteTo(chunk)

	assert.False(t, chunk.IsEmpty())
	for i := 0; i < 10; i++ {
		_, _ = row.WriteTo(chunk)
	}
	assert.True(t, chunk.IsFull())
}

func TestChunk_MarshalBinary(t *testing.T) {
	c1 := newChunk(ltoml.Size(2))
	compressed, err := c1.Compress()
	assert.NoError(t, err)
	assert.Nil(t, compressed)
	testMarshal(c1, 2, t)
	testMarshal(c1, 1, t)

	c2 := c1.(*chunk)

	row := makeTestBrokerRows()
	_, _ = row.WriteTo(c2)

	compressed, err = c2.Compress()
	assert.NoError(t, err)
	assert.NotNil(t, compressed)
}

func testMarshal(chunk Chunk, count int, t *testing.T) {
	compressed, err := chunk.Compress()
	assert.Nil(t, compressed)
	assert.Nil(t, err)

	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	for i := 0; i < count; i++ {
		var row metric.BrokerRow
		_ = converter.ConvertTo(&protoMetricsV1.Metric{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
			},
			Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
		}, &row)
		_, _ = row.WriteTo(chunk)
	}
	compressed, err = chunk.Compress()
	assert.NoError(t, err)
	assert.NotNil(t, compressed)

	reader := compress.NewSnappyReader()
	dst, err := reader.Uncompress(compressed)
	assert.NoError(t, err)
	var batch metric.StorageBatchRows
	assert.NotPanics(t, func() {
		batch.UnmarshalRows(dst)
	})
}

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

package native

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"
)

var testMetricList = &protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{
	{
		Namespace: "",
		Name:      "a",
		Timestamp: 0,
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "counter", Type: protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM, Value: 23},
		}},
}}

func makeGzipData(testMetricList *protoMetricsV1.MetricList) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	data, _ := testMetricList.Marshal()
	_, _ = writer.Write(data)
	_ = writer.Close()
	return buf.Bytes()
}

func Test_Parse(t *testing.T) {
	r := bytes.NewReader(makeGzipData(testMetricList))

	req, err := http.NewRequest(http.MethodPut, "", r)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	req.Header.Set("Content-Encoding", "gzip")

	enrichedTags := []tag.Tag{
		tag.NewTag([]byte("ip"), []byte("1.1.1.1")),
		tag.NewTag([]byte("region"), []byte("nj")),
	}
	metrics, err := Parse(req, enrichedTags, "ns")
	assert.Nil(t, err)
	assert.NotNil(t, metrics)
	assert.Len(t, metrics.Metrics, 1)
	assert.Len(t, metrics.Metrics[0].Tags, 2)
}

func Test_Parse_badGzipData(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "", strings.NewReader("bad-data"))
	assert.Nil(t, err)
	assert.NotNil(t, req)
	req.Header.Set("Content-Encoding", "gzip")
	_, err = Parse(req, nil, "ns")
	assert.NotNil(t, err)
}

func Test_Parse_error(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, "", strings.NewReader("bad-data"))
	_, err := Parse(req, nil, "ns")
	assert.NotNil(t, err)
}

func Test_Parser_empty(t *testing.T) {
	var m = &protoMetricsV1.MetricList{}
	data, _ := m.Marshal()
	req, _ := http.NewRequest(http.MethodPut, "", bytes.NewReader(data))
	_, err := Parse(req, nil, "ns")
	assert.NotNil(t, err)
}

func Test_parseProtoMetric(t *testing.T) {
	data, _ := testMetricList.Marshal()
	ms, err := parseProtoMetric(data, nil, "ns")
	assert.Nil(t, err)
	assert.Equal(t, "ns", ms.Metrics[0].Namespace)
	assert.Len(t, ms.Metrics[0].Tags, 0)
}

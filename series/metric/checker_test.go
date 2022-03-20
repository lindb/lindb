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

package metric

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/golang/snappy"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
)

func Test_Sanitize(t *testing.T) {
	assert.Equal(t, "aaaa", SanitizeNamespace("aaaa"))
	assert.Equal(t, "aa_aa", SanitizeNamespace("aa|aa"))

	assert.Equal(t, "aaaa", SanitizeMetricName("aaaa"))
	assert.Equal(t, "aa_aa", SanitizeMetricName("aa|aa"))

	assert.Equal(t, "aa|bb", JoinNamespaceMetric("aa", "bb"))
}

func Benchmark_SerializeFlatMetric(b *testing.B) {
	builder := flatbuffers.NewBuilder(1024)

	for n := 0; n < b.N; n++ {
		buildFlatMetric(builder)
	}
}

func Benchmark_UnmarshalFlatMetric_10Fields(b *testing.B) {
	builder := flatbuffers.NewBuilder(1024)
	buildFlatMetric(builder)
	data := builder.FinishedBytes()
	var simpleField flatMetricsV1.SimpleField
	var m flatMetricsV1.Metric

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		m.Init(data, flatbuffers.GetUOffsetT(data))
		for i := 0; i < m.SimpleFieldsLength(); i++ {
			m.SimpleFields(&simpleField, i)
			_ = simpleField.Value()
			_ = simpleField.Type()
			_ = simpleField.Name()
		}
	}
}

func Test_SanitizeFieldName(t *testing.T) {
	assert.Equal(t, []byte("_HistogramTest"), SanitizeFieldName([]byte("HistogramTest")))
	assert.Equal(t, []byte("_bucket_1"), SanitizeFieldName([]byte("__bucket_1")))
	assert.Equal(t, []byte("bucket_1"), SanitizeFieldName([]byte("bucket_1")))
}

func Benchmark_FlatMetric_Unmarshal10KeyValues(b *testing.B) {
	builder := flatbuffers.NewBuilder(1024)
	buildFlatMetric(builder)
	data := builder.FinishedBytes()
	var kv flatMetricsV1.KeyValue
	var m flatMetricsV1.Metric

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		m.Init(data, flatbuffers.GetUOffsetT(data))
		for i := 0; i < m.KeyValuesLength(); i++ {
			m.KeyValues(&kv, i)
			_ = kv.Key()
			_ = kv.Value()
		}
	}
}

func Benchmark_Marshal_Proto(b *testing.B) {
	m := protoMetricsV1.Metric{Name: "hello", Namespace: "default-ns", Timestamp: fasttime.UnixMilliseconds()}

	for i := 0; i < 10; i++ {
		m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
			Name: "counter" + strconv.Itoa(i), Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: float64(i)})
		m.Tags = append(m.Tags, &protoMetricsV1.KeyValue{Key: "key" + strconv.Itoa(i), Value: "value" + strconv.Itoa(i)})
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Marshal()
	}
}

func Benchmark_Unmarshal_Proto_10Fields(b *testing.B) {
	m := protoMetricsV1.Metric{Name: "hello", Namespace: "default-ns", Timestamp: fasttime.UnixMilliseconds()}

	for i := 0; i < 10; i++ {
		m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
			Name: "counter" + strconv.Itoa(i), Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: float64(i)})
		m.Tags = append(m.Tags, &protoMetricsV1.KeyValue{Key: "key" + strconv.Itoa(i), Value: "value" + strconv.Itoa(i)})
	}
	data, _ := m.Marshal()

	b.ResetTimer()
	b.ReportAllocs()
	var metric protoMetricsV1.Metric
	for i := 0; i < b.N; i++ {
		_ = metric.Unmarshal(data)
		for x := 0; x < 10; x++ {
			f := metric.SimpleFields[x]
			_ = f.Name
			_ = f.Value
			_ = f.Type
		}
	}
}

func Test_FlatMetric(t *testing.T) {
	builder := flatbuffers.NewBuilder(1024)
	// serialize metric
	metricName := builder.CreateString("hello")
	namespace := builder.CreateString("default-ns")
	flatMetricsV1.MetricStart(builder)
	flatMetricsV1.MetricAddNamespace(builder, namespace)
	flatMetricsV1.MetricAddName(builder, metricName)
	end := flatMetricsV1.MetricEnd(builder)
	builder.Finish(end)

	data := builder.FinishedBytes()
	m := flatMetricsV1.GetRootAsMetric(data, 0)

	assert.Equal(t, "hello", string(m.Name()))
	assert.Equal(t, "default-ns", string(m.Namespace()))
	assert.Zero(t, m.Timestamp())
	assert.Zero(t, m.Hash())
}

func Test_FlatMetric_Size(t *testing.T) {
	builder := flatbuffers.NewBuilder(256)
	buildFlatMetric(builder)
	data := builder.FinishedBytes()

	t.Log("flat raw size", len(data))
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	_, _ = w.Write(data)
	_ = w.Flush()
	t.Log("flat gzip compressed size", len(buf.Bytes()))

	w2 := snappy.NewBufferedWriter(&buf)
	buf.Reset()
	_, _ = w2.Write(data)
	_ = w2.Flush()
	t.Log("flat snappy compressed size", len(buf.Bytes()))
}

func Benchmark_GzipCompress(b *testing.B) {
	builder := flatbuffers.NewBuilder(1024)
	buildFlatMetric(builder)
	data := builder.FinishedBytes()

	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = w.Write(data)
		_ = w.Flush()
	}
}

func Benchmark_SnappyStreamCompress(b *testing.B) {
	builder := flatbuffers.NewBuilder(1024)
	buildFlatMetric(builder)
	data := builder.FinishedBytes()

	var buf bytes.Buffer
	w := snappy.NewBufferedWriter(&buf)

	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = w.Write(data)
		_ = w.Flush()
	}
}

func Benchmark_SnappyBlockCompress(b *testing.B) {
	builder := flatbuffers.NewBuilder(1024)
	buildFlatMetric(builder)
	data := builder.FinishedBytes()

	var block []byte
	for i := 0; i < b.N; i++ {
		block = block[:0]
		block = snappy.Encode(block, data)
	}
}

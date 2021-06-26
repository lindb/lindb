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

package prometheus

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/lindb/lindb/series/tag"

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"
)

const goodText = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# 	TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{method="post",code="400",quantile="0"} 4.9351e-05
go_gc_duration_seconds{quantile="0.25",} 7.424100000000001e-05
go_gc_duration_seconds{quantile="0.5"} 8.383
go_gc_duration_seconds{quantile="0.75"} 8.3835
go_gc_duration_seconds{ quantile="0.9",method="post",code="400",} 8.3835
go_gc_duration_seconds {method="post",code="400", quantile = "0.95" } 8.38
go_gc_duration_seconds {method="post",code="400", quantile = "0.99" } 8.383
go_gc_duration_seconds {method="post",code="400", quantile = "0.999" } 8.383
go_gc_duration_seconds { quantile = "0.9999" } 8.38
go_gc_duration_seconds_count 9
go_gc_duration_seconds_sum 90
go_gc_duration_seconds_count{method="post",code="400"} 9
go_gc_duration_seconds_sum{method="post",code="400"} 90
# HELP prometheus_tsdb_compaction_chunk_range Final time range of chunks on their first compaction
# TYPE prometheus_tsdb_compaction_chunk_range histogram
prometheus_tsdb_compaction_chunk_range_bucket{le="100"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="400"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="1600"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="6400"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="25600"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="102400"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="409600"} 0
prometheus_tsdb_compaction_chunk_range_bucket{le="1.6384e+06"} 260
prometheus_tsdb_compaction_chunk_range_bucket{le="6.5536e+06"} 780
prometheus_tsdb_compaction_chunk_range_bucket{le="2.62144e+07"} 780
prometheus_tsdb_compaction_chunk_range_bucket{le="+Inf"} 780
prometheus_tsdb_compaction_chunk_range_sum 1.1540798e+09
prometheus_tsdb_compaction_chunk_range_count 780
# Hrandom comment starting with prefix of HELP
#
wind_speed{A="2",c="3"} 12345
some:aggregate:rate5m{a_b="c"}	1
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 33  	123123
# A normal comment.
#
# TYPE name counter
name{labelname="val1",basename="basevalue"} NaN
name {labelname="val2",basename="base\"v\\al\nue"} 0.23 1234567890
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE _metric_starting_with_underscore counter 
_metric_starting_with_underscore 1
testmetric{_label_starting_with_underscore="foo"} 1
testmetric{label="\"bar\""} 1`

func makeGzipData(body []byte) []byte {
	var w bytes.Buffer
	gw := gzip.NewWriter(&w)
	_, _ = gw.Write(body)
	_ = gw.Close()
	return w.Bytes()
}

func Test_Parse(t *testing.T) {
	r1 := bytes.NewReader(makeGzipData([]byte(goodText)))

	req, _ := http.NewRequest(http.MethodPut, "", r1)
	req.Header.Set("Content-Encoding", "gzip")

	metricList, err := Parse(req, []tag.Tag{
		tag.NewTag([]byte("zone"), []byte("bj")),
		tag.NewTag([]byte("host"), []byte("abcd")),
	}, "ns")
	assert.Nil(t, err)
	assert.NotNil(t, metricList)

	// bad gzip data
	r2 := bytes.NewReader([]byte(goodText))
	req, _ = http.NewRequest(http.MethodPut, "", r2)
	req.Header.Set("Content-Encoding", "gzip")
	metricList, err = Parse(req, []tag.Tag{}, "ns")
	assert.NotNil(t, err)
	assert.Nil(t, metricList)
}

func TestPromParse(t *testing.T) {
	input := goodText
	input += "\n# HELP metric foo\x00bar"
	input += "\nnull_byte_metric{a=\"abc\x00\"} 1\n"

	metrics, err := promParse(strings.NewReader(input), tag.Tags{}, "ns")
	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)

	metrics, err = promParse(strings.NewReader("empty"), tag.Tags{}, "ns")
	assert.Error(t, err)
	assert.Empty(t, metrics)
	input = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# 	TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{method="post",code="400",quantile="0"} 4.9351e-05
go_gc_duration_seconds { quantile = "0.9999" } 8.38`
	metrics, err = promParse(strings.NewReader(input), tag.Tags{}, "ns")
	assert.NoError(t, err)
	assert.Empty(t, metrics)
	input = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# 	TYPE go_gc_duration_seconds summary
go_gc_duration_seconds { quantile = "0.9999" } NaN
go_gc_duration_seconds_count 9
go_gc_duration_seconds_sum 90
`
	metrics, err = promParse(strings.NewReader(input), tag.Tags{}, "ns")
	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

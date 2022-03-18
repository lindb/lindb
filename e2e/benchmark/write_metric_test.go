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

//go:build benchmark
// +build benchmark

package benchmark

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/series/metric"
)

func TestWriteSumMetric(b *testing.T) {
	timestamp := timeutil.Now()
	cli := resty.New()
	count := 0
	for i := 0; i < 40; i++ {
		var buf bytes.Buffer
		for j := 0; j < 20; j++ {
			for k := 0; k < 4; k++ {
				var brokerRow metric.BrokerRow
				converter := metric.NewProtoConverter()
				err := converter.ConvertTo(&protoMetricsV1.Metric{
					Name:      "host_disk0_200",
					Timestamp: timestamp,
					Tags: []*protoMetricsV1.KeyValue{
						{Key: "host", Value: "host" + strconv.Itoa(i)},
						{Key: "disk", Value: "disk" + strconv.Itoa(k)},
						{Key: "partition", Value: "partition" + strconv.Itoa(j)},
					},
					SimpleFields: []*protoMetricsV1.SimpleField{
						{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
				}, &brokerRow)
				_, _ = brokerRow.WriteTo(&buf)
				if err != nil {
					panic(err)
				}
			}
		}
		body := buf.Bytes()
		_, err := cli.R().SetBody(body).Put("http://127.0.0.1:9000/api/flat/write?db=_internal")
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(count)
}

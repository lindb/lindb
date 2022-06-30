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
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/series/metric"
)

func TestWriteSumMetric(b *testing.T) {
	timestamp := timeutil.Now()
	cli := resty.New()
	count := 0
	for i := 0; i < 4; i++ {
		var buf bytes.Buffer
		for j := 0; j < 20; j++ {
			for k := 0; k < 4; k++ {
				var brokerRow metric.BrokerRow
				converter := metric.NewProtoConverter()
				err := converter.ConvertTo(&protoMetricsV1.Metric{
					Name:      "host_disk_700",
					Timestamp: timestamp,
					Tags: []*protoMetricsV1.KeyValue{
						{Key: "host", Value: "host" + strconv.Itoa(i)},
						{Key: "disk", Value: "disk" + strconv.Itoa(i)},
						{Key: "partition", Value: "partition" + strconv.Itoa(i)},
					},
					SimpleFields: []*protoMetricsV1.SimpleField{
						{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: float64(k)},
						{Name: "f2", Type: protoMetricsV1.SimpleFieldType_LAST, Value: float64(k)},
						{Name: "f3", Type: protoMetricsV1.SimpleFieldType_FIRST, Value: float64(k)},
					},
				}, &brokerRow)
				_, _ = brokerRow.WriteTo(&buf)
				if err != nil {
					panic(err)
				}
			}
		}
		body := buf.Bytes()
		_, err := cli.R().SetBody(body).Put("http://127.0.0.1:9000/api/v1/flat/write?db=_internal")
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(count)
}

func TestWriteSumMetric_OneDay(b *testing.T) {
	timestamp, _ := timeutil.ParseTimestamp("2022-03-27 00:00:00")
	cli := resty.New()
	oneDay := timeutil.OneHour * 2 / 2000
	for n := int64(0); n < oneDay; n++ {
		var buf bytes.Buffer
		for i := 0; i < 1; i++ {
			for j := 0; j < 1; j++ {
				for k := 0; k < 1; k++ {
					var brokerRow metric.BrokerRow
					converter := metric.NewProtoConverter()
					err := converter.ConvertTo(&protoMetricsV1.Metric{
						Name:      "host_disk_3400",
						Timestamp: timestamp + n*2000,
						Tags: []*protoMetricsV1.KeyValue{
							{Key: "host", Value: "host" + strconv.Itoa(i)},
							{Key: "disk", Value: "disk" + strconv.Itoa(k)},
							{Key: "partition", Value: "partition" + strconv.Itoa(j)},
						},
						SimpleFields: []*protoMetricsV1.SimpleField{
							{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
						},
					}, &brokerRow)
					_, _ = brokerRow.WriteTo(&buf)
					if err != nil {
						panic(err)
					}
				}
			}
		}
		body := buf.Bytes()
		_, err := cli.R().SetBody(body).Put("http://127.0.0.1:9000/api/v1/flat/write?db=_internal")
		if err != nil {
			panic(err)
		}
		fmt.Println(n)
	}
}

func TestWriteSumMetric_7Days(b *testing.T) {
	timestamp, _ := timeutil.ParseTimestamp("2022-06-05 00:00:00")
	cli := resty.New()
	for d := int64(0); d < 24; d++ {
		var buf bytes.Buffer
		for n := int64(0); n < 60; n++ {
			for i := 0; i < 1; i++ {
				var brokerRow metric.BrokerRow
				converter := metric.NewProtoConverter()
				fmt.Println(timeutil.FormatTimestamp(timestamp+n*timeutil.OneMinute+d*timeutil.OneHour, timeutil.DataTimeFormat2))
				err := converter.ConvertTo(&protoMetricsV1.Metric{
					Name:      "host_disk_170",
					Timestamp: timestamp + n*timeutil.OneMinute + d*timeutil.OneHour,
					Tags: []*protoMetricsV1.KeyValue{
						{Key: "host", Value: "host" + strconv.Itoa(i)},
						{Key: "disk", Value: "disk" + strconv.Itoa(i)},
						{Key: "partition", Value: "partition" + strconv.Itoa(i)},
					},
					SimpleFields: []*protoMetricsV1.SimpleField{
						{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
					},
				}, &brokerRow)
				_, _ = brokerRow.WriteTo(&buf)
				if err != nil {
					panic(err)
				}
			}
		}
		body := buf.Bytes()
		_, err := cli.R().SetBody(body).Put("http://127.0.0.1:9000/api/v1/flat/write?db=test")
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}

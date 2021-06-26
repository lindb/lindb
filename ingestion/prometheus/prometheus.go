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
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cespare/xxhash"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/tag"
)

// todo: line-based-parser

// Parse parses prometheus text
func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*pb.MetricList, error) {
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}

	return promParse(reader, enrichedTags, namespace)
}

// promParse parses prometheus text prometheus to LinDB pb prometheus.
func promParse(reader io.Reader, enrichedTags tag.Tags, namespace string) (*pb.MetricList, error) {
	parser := &expfmt.TextParser{}
	out, err := parser.TextToMetricFamilies(reader)
	if err != nil && len(out) == 0 {
		return nil, err
	}
	metricList := &pb.MetricList{}
	for name, pm := range out {
		metricType := *pm.Type
		if metricType == dto.MetricType_UNTYPED {
			// not support untyped metric type
			continue
		}
		for _, m := range pm.Metric {
			f := getFieldType(metricType, m)
			if f == nil {
				continue
			}

			metric := &pb.Metric{
				Name:      name,
				Namespace: namespace,
				Fields:    []*pb.Field{f},
			}
			if m.TimestampMs != nil {
				metric.Timestamp = *m.TimestampMs
			} else {
				metric.Timestamp = timeutil.Now()
			}
			tagCount := len(m.Label)
			if tagCount > 0 {
				tags := make(map[string]string, tagCount)
				for _, label := range m.Label {
					tags[*label.Name] = *label.Value
				}
				metric.Tags = tags
			}
			if enrichedTags.Size() > 0 {
				if metric.Tags == nil {
					metric.Tags = make(map[string]string)
				}
				for _, enrichedTag := range enrichedTags {
					metric.Tags[string(enrichedTag.Key)] = string(enrichedTag.Value)
				}
			}

			if metric.Tags != nil && len(metric.Tags) > 0 {
				metric.TagsHash = xxhash.Sum64String(tag.Concat(metric.Tags))
			} else {
				metric.TagsHash = xxhash.Sum64String(metric.Name)
			}

			metricList.Metrics = append(metricList.Metrics, metric)
		}
	}
	return metricList, nil
}

func getFieldType(metricType dto.MetricType, metric *dto.Metric) *pb.Field {
	switch metricType {
	case dto.MetricType_COUNTER:
		if metric.Counter != nil && metric.Counter.Value != nil {
			return &pb.Field{
				Name:  "counter",
				Type:  pb.FieldType_Sum,
				Value: *metric.Counter.Value,
			}
		}
	case dto.MetricType_GAUGE:
		if metric.Gauge != nil && metric.Gauge.Value != nil {
			return &pb.Field{
				Name:  "gauge",
				Type:  pb.FieldType_Gauge,
				Value: *metric.Gauge.Value,
			}
		}
	case dto.MetricType_SUMMARY:
		if metric.Summary == nil || metric.Summary.SampleCount == nil || metric.Summary.SampleSum == nil {
			return nil
		}
		f := &pb.Field{
			Name: "summary",
			//Type: pb.FieldType_Summary,
		}
		//FIXME stone1100
		//count := float64(*metric.Summary.SampleCount)
		//f.Fields = append(f.Fields, &pb.PrimitiveField{
		//	PrimitiveID: int32(1),
		//	Value:       *metric.Summary.SampleSum,
		//}, &pb.PrimitiveField{
		//	PrimitiveID: int32(2),
		//	Value:       count,
		//})
		//quantile := metric.Summary.Quantile
		//for _, q := range quantile {
		//	switch *q.Quantile {
		//	case 0.5:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(50),
		//			Value:       *q.Value * count,
		//		})
		//	case 0.75:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(75),
		//			Value:       *q.Value * count,
		//		})
		//	case 0.90:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(90),
		//			Value:       *q.Value * count,
		//		})
		//	case 0.95:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(95),
		//			Value:       *q.Value * count,
		//		})
		//	case 0.99:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(99),
		//			Value:       *q.Value * count,
		//		})
		//	case 0.999:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(39),
		//			Value:       *q.Value * count,
		//		})
		//	case 0.9999:
		//		f.Fields = append(f.Fields, &pb.PrimitiveField{
		//			PrimitiveID: int32(49),
		//			Value:       *q.Value * count,
		//		})
		//	}
		//}
		return f
	}
	return nil
}

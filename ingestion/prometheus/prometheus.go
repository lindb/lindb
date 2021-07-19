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
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

// todo: line-based-parser

// Parse parses prometheus text
func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*protoMetricsV1.MetricList, error) {
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
func promParse(reader io.Reader, enrichedTags tag.Tags, namespace string) (*protoMetricsV1.MetricList, error) {
	parser := &expfmt.TextParser{}
	out, err := parser.TextToMetricFamilies(reader)
	if err != nil && len(out) == 0 {
		return nil, err
	}
	metricList := &protoMetricsV1.MetricList{}
	for name, pm := range out {
		metricType := *pm.Type
		if metricType == dto.MetricType_UNTYPED {
			// not support untyped metric type
			continue
		}
		for _, m := range pm.Metric {
			metric := &protoMetricsV1.Metric{
				Name:      name,
				Namespace: namespace,
			}
			set := setField(metric, metricType, m)
			if !set {
				continue
			}
			if m.TimestampMs != nil {
				metric.Timestamp = *m.TimestampMs
			} else {
				metric.Timestamp = timeutil.Now()
			}
			tagCount := len(m.Label)
			if tagCount > 0 {
				var tags = make([]*protoMetricsV1.KeyValue, len(m.Label))
				for idx := range m.Label {
					tags[idx] = &protoMetricsV1.KeyValue{
						Key: *m.Label[idx].Name, Value: *m.Label[idx].Value}
				}
				metric.Tags = tags
			}
			if enrichedTags.Size() > 0 {
				if metric.Tags == nil {
					metric.Tags = nil
				}
				for _, enrichedTag := range enrichedTags {
					for idx := range metric.Tags {
						if metric.Tags[idx].Key == string(enrichedTag.Key) {
							continue
						} else {
							metric.Tags = append(metric.Tags, &protoMetricsV1.KeyValue{
								Key:   string(enrichedTag.Key),
								Value: string(enrichedTag.Value)})
						}
					}
				}
			}

			if metric.Tags != nil && len(metric.Tags) > 0 {
				metric.TagsHash = xxhash.Sum64String(tag.ConcatKeyValues(metric.Tags))
			} else {
				metric.TagsHash = xxhash.Sum64String(metric.Name)
			}

			metricList.Metrics = append(metricList.Metrics, metric)
		}
	}
	return metricList, nil
}

func setField(metric *protoMetricsV1.Metric, metricType dto.MetricType, dtoMetric *dto.Metric) bool {
	switch metricType {
	case dto.MetricType_COUNTER:
		if dtoMetric.Counter == nil || dtoMetric.Counter.Value == nil {
			return false
		}
		metric.SimpleFields = []*protoMetricsV1.SimpleField{{
			Name:  "counter",
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			Value: *dtoMetric.Counter.Value,
		}}
	case dto.MetricType_GAUGE:
		if dtoMetric.Gauge == nil && dtoMetric.Gauge.Value == nil {
			return false
		}
		metric.SimpleFields = []*protoMetricsV1.SimpleField{{
			Name:  "gauge",
			Type:  protoMetricsV1.SimpleFieldType_GAUGE,
			Value: *dtoMetric.Gauge.Value,
		}}
	case dto.MetricType_HISTOGRAM:
		if dtoMetric.Histogram == nil || len(dtoMetric.Histogram.Bucket) == 0 {
			return false
		}
		var (
			explicitBounds = make([]float64, len(dtoMetric.Histogram.GetBucket()))
			values         = make([]float64, len(dtoMetric.Histogram.GetBucket()))
		)
		for idx := range dtoMetric.Histogram.GetBucket() {
			bkt := dtoMetric.Histogram.GetBucket()[idx]
			if bkt == nil {
				return false
			}
			explicitBounds[idx] = bkt.GetUpperBound()
			values[idx] = float64(bkt.GetCumulativeCount())
		}
		metric.CompoundField = &protoMetricsV1.CompoundField{
			Type:           protoMetricsV1.CompoundFieldType_DELTA_HISTOGRAM,
			Sum:            dtoMetric.Histogram.GetSampleSum(),
			Count:          float64(dtoMetric.Histogram.GetSampleCount()),
			ExplicitBounds: explicitBounds,
			Values:         values,
		}

	case dto.MetricType_SUMMARY:
		// todo: record not-support data
		return false
	}
	return true
}

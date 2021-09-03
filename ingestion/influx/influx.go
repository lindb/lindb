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

package influx

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cespare/xxhash/v2"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/strutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

var (
	influxLogger = logger.GetLogger("ingestion", "InfluxDB")
)

// Parse parses influxdb line protocol data to LinDB pb prometheus.
// https://docs.influxdata.com/influxdb/v2.0/write-data/developer-tools/api/#example-api-write-request
func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*protoMetricsV1.MetricList, error) {
	qry := req.URL.Query()
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			influxCorruptedDataCounter.Incr()
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}
	// precision
	multiplier := getPrecisionMultiplier(qry.Get("precision"))

	cr := ingestCommon.GetChunkReader(reader)
	defer ingestCommon.PutChunkReader(cr)

	metricList := &protoMetricsV1.MetricList{}
	for cr.HasNext() {
		nextLine := cr.Next()
		influxReadBytesCounter.Add(float64(len(nextLine)))
		metric, err := parseInfluxLine(nextLine, namespace, multiplier)
		if err != nil {
			influxLogger.Warn("ingest error",
				logger.String("line", string(nextLine)),
				logger.Error(err))
			droppedMetricsCounter.Incr()
			continue
		}
		if metric == nil || len(metric.SimpleFields) == 0 {
			droppedMetricsCounter.Incr()
			continue
		}
		ingestedMetricsCounter.Incr()
		ingestedFieldsCounter.Add(float64(len(metric.SimpleFields)))
		// enrich tags
		for _, enrichedTag := range enrichedTags {
			tagKey := strutil.ByteSlice2String(enrichedTag.Key)
			for idx := range metric.Tags {
				if metric.Tags[idx].Key == tagKey {
					continue
				}
				metric.Tags = append(metric.Tags, &protoMetricsV1.KeyValue{
					Key:   tagKey,
					Value: strutil.ByteSlice2String(enrichedTag.Value),
				})
			}
			metric.TagsHash = xxhash.Sum64String(tag.ConcatKeyValues(metric.Tags))
		}
		metricList.Metrics = append(metricList.Metrics, metric)
	}
	if cr.Error() == nil || cr.Error() == io.EOF {
		return metricList, nil
	}
	return metricList, cr.Error()
}

// getPrecisionMultiplier returns a multiplier for the precision specified.
// https://docs.influxdata.com/influxdb/v2.0/api/#operation/PostWrite
// timestamp in lindb is milliseconds
// when multiplier > 0, real_timestamp = timestamp * multiplier
// when multiplier < 0, real_timestamp = timestamp / (-1 * multiplier)
func getPrecisionMultiplier(precision string) int64 {
	switch strings.ToLower(precision) {
	case "ns":
		return -1e6
	case "us":
		return -1e3
	case "ms":
		return 1
	case "s":
		return 1000
	case "m":
		return 1000 * 60
	case "h":
		return 1000 * 3600
	default:
		return 0
	}
}

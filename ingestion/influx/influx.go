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
	"github.com/cespare/xxhash"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/pkg/strutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/tag"

	"fmt"
	"io"
	"net/http"
	"strings"
)

// Parse parses influxdb line protocol data to LinDB pb prometheus.
// https://docs.influxdata.com/influxdb/v2.0/write-data/developer-tools/api/#example-api-write-request
func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*pb.MetricList, error) {
	qry := req.URL.Query()
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}
	// precision
	multiplier := getPrecisionMultiplier(qry.Get("precision"))

	cr := ingestCommon.GetChunkReader(reader)
	defer ingestCommon.PutChunkReader(cr)

	metricList := &pb.MetricList{}
	for cr.HasNext() {
		metric, err := parseInfluxLine(cr.Next(), namespace, multiplier)
		if err != nil {
			return nil, err
		}
		if metric == nil {
			continue
		}
		// enrich tags
		for _, enrichedTag := range enrichedTags {
			tagKey := strutil.ByteSlice2String(enrichedTag.Key)
			if _, ok := metric.Tags[tagKey]; !ok {
				metric.Tags[tagKey] = strutil.ByteSlice2String(enrichedTag.Value)
			}
			metric.TagsHash = xxhash.Sum64String(tag.Concat(metric.Tags))
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
		return 1
	}
}

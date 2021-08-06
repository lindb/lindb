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
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/cespare/xxhash"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/strutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

var (
	nativeIngestionScope         = linmetric.NewScope("lindb.ingestion").Scope("native")
	nativeCorruptedDataCounter   = nativeIngestionScope.NewDeltaCounter("data_corrupted_count")
	nativeUnmarshalMetricCounter = nativeIngestionScope.NewDeltaCounter("ingested_metrics")
	nativeReadBytesCounter       = nativeIngestionScope.NewDeltaCounter("read_bytes_count")
)

func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*protoMetricsV1.MetricList, error) {
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			nativeCorruptedDataCounter.Incr()
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	nativeReadBytesCounter.Add(float64(len(data)))
	ms, err := parseProtoMetric(data, enrichedTags, namespace)
	if err != nil {
		nativeCorruptedDataCounter.Incr()
		return nil, err
	}
	if len(ms.Metrics) == 0 {
		return nil, fmt.Errorf("empty metrics")
	}
	nativeUnmarshalMetricCounter.Add(float64(len(ms.Metrics)))
	return ms, nil
}

func parseProtoMetric(data []byte, enrichedTags tag.Tags, namespace string) (*protoMetricsV1.MetricList, error) {
	var ms protoMetricsV1.MetricList
	if err := ms.Unmarshal(data); err != nil {
		return nil, err
	}
	for _, m := range ms.Metrics {
		m.Namespace = namespace
		if len(enrichedTags) > 0 {
			var newKeyValues tag.KeyValues
			for _, t := range enrichedTags {
				newKeyValues = append(newKeyValues, &protoMetricsV1.KeyValue{
					Key:   strutil.ByteSlice2String(t.Key),
					Value: strutil.ByteSlice2String(t.Value),
				})
			}
			newKeyValues = append(newKeyValues, m.Tags...)
			sort.Sort(newKeyValues)
			m.Tags = newKeyValues
		} else {
			var kvs tag.KeyValues = m.Tags
			sort.Sort(kvs)
			m.Tags = kvs
		}
		m.TagsHash = xxhash.Sum64String(tag.ConcatKeyValues(m.Tags))
	}
	return &ms, nil
}

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

package proto

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

var (
	protoIngestionStatistics = metrics.NewNativeIngestionStatistics()
)

func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*metric.BrokerBatchRows, error) {
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			protoIngestionStatistics.CorruptedData.Incr()
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	protoIngestionStatistics.ReadBytes.Add(float64(len(data)))
	batch, err := parseProtoMetric(data, enrichedTags, namespace)
	if err != nil {
		protoIngestionStatistics.CorruptedData.Incr()
		return nil, err
	}
	if batch.Len() == 0 {
		return nil, fmt.Errorf("empty metrics")
	}
	protoIngestionStatistics.IngestedMetrics.Add(float64(batch.Len()))
	return batch, nil
}

func parseProtoMetric(
	data []byte,
	enrichedTags tag.Tags,
	namespace string,
) (
	batch *metric.BrokerBatchRows, err error,
) {
	batch = metric.NewBrokerBatchRows()

	converter, releaseFunc := metric.NewBrokerRowProtoConverter(strutil.String2ByteSlice(namespace), enrichedTags)
	defer releaseFunc(converter)

	var ms protoMetricsV1.MetricList
	if err := ms.Unmarshal(data); err != nil {
		return nil, err
	}
	for _, m := range ms.Metrics {
		m := m
		if err := batch.TryAppend(func(row *metric.BrokerRow) error {
			return converter.ConvertTo(m, row)
		}); err != nil {
			protoIngestionStatistics.DroppedMetrics.Incr()
		}
	}
	return batch, nil
}

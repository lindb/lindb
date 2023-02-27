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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	commonseries "github.com/lindb/common/series"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

var (
	influxLogger = logger.GetLogger("Ingestion", "InfluxDB")
)

// Parse parses influxdb line protocol data to LinDB pb prometheus.
// https://docs.influxdata.com/influxdb/v2.0/write-data/developer-tools/api/#example-api-write-request
func Parse(req *http.Request, enrichedTags tag.Tags, namespace string, limits *models.Limits) (*metric.BrokerBatchRows, error) {
	qry := req.URL.Query()
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			influxIngestionStatistics.CorruptedData.Incr()
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}
	// precision
	multiplier := getPrecisionMultiplier(qry.Get("precision"))

	cr := GetChunkReader(reader)
	defer PutChunkReader(cr)

	rowBuilder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(rowBuilder)

	batch := metric.NewBrokerBatchRows()

	for cr.HasNext() {
		nextLine := cr.Next()
		// reset for constructing next row
		rowBuilder.Reset()

		influxIngestionStatistics.ReadBytes.Add(float64(len(nextLine)))
		// skip comment line
		if bytes.HasPrefix(nextLine, []byte{'#'}) {
			continue
		}
		if err := parseInfluxLine(rowBuilder, nextLine, namespace, multiplier, limits); err != nil {
			influxLogger.Warn("ingest error",
				logger.String("line", string(nextLine)),
				logger.Error(err))
			influxIngestionStatistics.DroppedMetrics.Incr()
			continue
		}

		for _, enrichedTag := range enrichedTags {
			if err := rowBuilder.AddTag(enrichedTag.Key, enrichedTag.Value); err != nil {
				return nil, err
			}
		}
		if err := batch.TryAppend(func(row *metric.BrokerRow) error {
			data, err := rowBuilder.Build()
			if err != nil {
				return err
			}
			row.FromBlock(data)
			return nil
		}); err != nil {
			influxIngestionStatistics.DroppedMetrics.Incr()
			continue
		}

		influxIngestionStatistics.IngestedMetrics.Incr()
		influxIngestionStatistics.IngestedFields.Add(float64(rowBuilder.SimpleFieldsLen()))
	}
	if cr.Error() == nil || cr.Error() == io.EOF {
		return batch, nil
	}
	return batch, cr.Error()
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

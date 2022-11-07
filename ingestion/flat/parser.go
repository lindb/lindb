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

package flat

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

var (
	flatIngestionStatistics = metrics.NewFlatIngestionStatistics()
)

var flatLogger = logger.GetLogger("Ingestion", "Flat")

func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*metric.BrokerBatchRows, error) {
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			flatIngestionStatistics.CorruptedData.Incr()
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}
	bufioReader, releaseBufioReaderFunc := ingestCommon.NewBufioReader(reader)
	defer releaseBufioReaderFunc(bufioReader)

	batch, err := parseFlatMetric(reader, enrichedTags, namespace)
	if err != nil {
		flatIngestionStatistics.CorruptedData.Incr()
		return nil, err
	}
	if batch.Len() == 0 {
		return nil, fmt.Errorf("empty metrics")
	}
	flatIngestionStatistics.IngestedMetrics.Add(float64(batch.Len()))
	return batch, nil
}

func parseFlatMetric(
	reader io.Reader,
	enrichedTags tag.Tags,
	namespace string,
) (
	batch *metric.BrokerBatchRows, err error,
) {
	batch = metric.NewBrokerBatchRows()

	decoder, releaseFunc := metric.NewBrokerRowFlatDecoder(
		reader,
		strutil.String2ByteSlice(namespace),
		enrichedTags,
	)
	defer releaseFunc(decoder)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("bad flat metrics binary")
			flatLogger.Error("decode panic", logger.Any("error", r), logger.Stack())
		}
	}()
	for decoder.HasNext() {
		if err := batch.TryAppend(decoder.DecodeTo); err != nil {
			flatLogger.Warn("failed ingesting flat metric", logger.Error(err))
			flatIngestionStatistics.DroppedMetric.Incr()
		}
	}

	switch {
	case decoder.ReadLen() < 10*1024:
		flatIngestionStatistics.LT10KiBCounter.Incr()
	case decoder.ReadLen() < 100*1024:
		flatIngestionStatistics.LT100KiBCounter.Incr()
	case decoder.ReadLen() < 1024*1024:
		flatIngestionStatistics.LT1MiBCounter.Incr()
	case decoder.ReadLen() < 10*1024*1024:
		flatIngestionStatistics.LT10MiBCounter.Incr()
	default:
		flatIngestionStatistics.GT10MiBCounter.Incr()
	}
	flatIngestionStatistics.ReadBytes.Add(float64(decoder.ReadLen()))

	return batch, nil
}

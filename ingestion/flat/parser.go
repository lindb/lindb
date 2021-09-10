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
	"io/ioutil"
	"net/http"
	"strings"

	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

var (
	flatIngestionScope         = linmetric.NewScope("lindb.ingestion.flat")
	flatCorruptedDataCounter   = flatIngestionScope.NewCounter("data_corrupted_count")
	flatDroppedMetricCounter   = flatIngestionScope.NewCounter("dropped_metrics")
	flatUnmarshalMetricCounter = flatIngestionScope.NewCounter("ingested_metrics")
	flatReadBytesCounter       = flatIngestionScope.NewCounter("read_bytes")
)

var flatLogger = logger.GetLogger("ingestion", "Flat")

func Parse(req *http.Request, enrichedTags tag.Tags, namespace string) (*metric.BrokerBatchRows, error) {
	var reader = req.Body
	if strings.EqualFold(req.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := ingestCommon.GetGzipReader(req.Body)
		if err != nil {
			flatCorruptedDataCounter.Incr()
			return nil, fmt.Errorf("ingestion corrupted gzip data: %w", err)
		}
		defer ingestCommon.PutGzipReader(gzipReader)
		reader = gzipReader
	}

	// todo: @codingcrush, use stream chunk reader
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	flatReadBytesCounter.Add(float64(len(data)))
	batch, err := parseFlatMetric(data, enrichedTags, namespace)
	if err != nil {
		flatCorruptedDataCounter.Incr()
		return nil, err
	}
	if batch.Len() == 0 {
		return nil, fmt.Errorf("empty metrics")
	}
	flatUnmarshalMetricCounter.Add(float64(batch.Len()))
	return batch, nil
}

func parseFlatMetric(
	data []byte,
	enrichedTags tag.Tags,
	namespace string,
) (
	batch *metric.BrokerBatchRows, err error,
) {
	batch = metric.NewBrokerBatchRows()

	decoder, releaseFunc := metric.NewBrokerRowFlatDecoder(data, strutil.String2ByteSlice(namespace), enrichedTags)
	defer releaseFunc(decoder)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("bad flat metrics binary")
			flatLogger.Error("decode panic", logger.Stack())
		}
	}()
	for decoder.HasNext() {
		if err := batch.TryAppend(decoder.DecodeTo); err != nil {
			flatLogger.Warn("failed ingesting flat metric", logger.Error(err))
			flatDroppedMetricCounter.Incr()
		}
	}
	return batch, nil
}

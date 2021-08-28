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

package metricsdata

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

// metricLoader implements flow.DataLoader interface that loads metric data from file storage.
type metricLoader struct {
	reader             MetricReader
	lowContainer       roaring.Container
	lowKeyOffsets      *encoding.FixedOffsetDecoder
	seriesEntriesBlock []byte
}

// newMetricLoader creates a file storage metric loader.
func newMetricLoader(
	reader MetricReader,
	seriesEntriesBlock []byte,
	lowContainer roaring.Container,
	lowKeyOffsets *encoding.FixedOffsetDecoder,
) flow.DataLoader {
	return &metricLoader{
		seriesEntriesBlock: seriesEntriesBlock,
		reader:             reader,
		lowContainer:       lowContainer,
		lowKeyOffsets:      lowKeyOffsets,
	}
}

// Load load the metric data by given series id from file storage.
func (s *metricLoader) Load(lowSeriesID uint16) (timeutil.SlotRange, [][]byte) {
	// check low series id if exist
	if !s.lowContainer.Contains(lowSeriesID) {
		return s.reader.GetTimeRange(), nil
	}
	// get the index of low series id in container
	idx := s.lowContainer.Rank(lowSeriesID)
	// scan the data and aggregate the values

	seriesEntry, _ := s.lowKeyOffsets.GetBlock(idx-1, s.seriesEntriesBlock)
	// read series data of fields
	return s.reader.GetTimeRange(), s.reader.readSeriesData(seriesEntry)
}

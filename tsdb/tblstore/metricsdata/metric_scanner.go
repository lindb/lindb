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
)

// metricScanner implements flow.Scanner interface that scans metric data from file storage.
type metricScanner struct {
	reader        Reader
	lowContainer  roaring.Container
	seriesOffsets *encoding.FixedOffsetDecoder
}

// newMetricScanner creates a file storage metric scanner.
func newMetricScanner(reader Reader,
	lowContainer roaring.Container,
	seriesOffsets *encoding.FixedOffsetDecoder,
) flow.Scanner {
	return &metricScanner{
		reader:        reader,
		lowContainer:  lowContainer,
		seriesOffsets: seriesOffsets,
	}
}

// Scan scans the metric data by given series id from file storage.
func (s *metricScanner) Scan(lowSeriesID uint16) [][]byte {
	// check low series id if exist
	if !s.lowContainer.Contains(lowSeriesID) {
		return nil
	}
	// get the index of low series id in container
	idx := s.lowContainer.Rank(lowSeriesID)
	// scan the data and aggregate the values
	seriesPos, _ := s.seriesOffsets.Get(idx - 1)
	// read series data of fields
	return s.reader.readSeriesData(seriesPos)
}

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
	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
)

//go:generate mockgen -source ./series_merger.go -destination=./series_merger_mock.go -package metricsdata

// SeriesMerger represents series data merger which merge multi fields under same series id
type SeriesMerger interface {
	// merge the multi-fields' data with same series id
	merge(mergeCtx *mergerContext,
		decodeStreams []*encoding.TSDDecoder,
		fieldReaders []FieldReader,
	) error
}

// seriesMerger implements SeriesMerger interface
type seriesMerger struct {
	flusher Flusher
}

// newSeriesMerger creates a series merger
func newSeriesMerger(flusher Flusher) SeriesMerger {
	return &seriesMerger{
		flusher: flusher,
	}
}

// merge the multi-fields' data with same series id
func (sm *seriesMerger) merge(
	mergeCtx *mergerContext,
	streams []*encoding.TSDDecoder,
	fieldReaders []FieldReader,
) error {
	for idx, f := range mergeCtx.targetFields {
		fieldID := f.ID
		encodeStream := sm.flusher.GetEncoder(idx)
		encodeStream.RestWithStartTime(mergeCtx.targetRange.Start)

		for idx, reader := range fieldReaders {
			if reader == nil {
				// if series id not exist, metricReader is nil
				continue
			}
			fieldData := reader.GetFieldData(fieldID)
			if len(fieldData) > 0 {
				if streams[idx] == nil {
					// new tsd decoder
					streams[idx] = encoding.GetTSDDecoder()
				}
				oldSlotRange := reader.SlotRange()
				// reset tsd data
				streams[idx].ResetWithTimeRange(fieldData, oldSlotRange.Start, oldSlotRange.End)
			}
		}
		// merges field data from source time range => target time range,
		// compact merge: source range = target range and ratio = 1
		// rollup merge: source range[5,182]=>target range[0,6], ratio:30, source interval:10s, target interval:5min
		aggregation.DownSamplingMultiSeriesInto(
			mergeCtx.targetRange, mergeCtx.ratio, mergeCtx.baseSlot,
			f.Type, streams,
			encodeStream.EmitDownSamplingValue,
		)

		data, err := encodeStream.BytesWithoutTime()
		if err != nil {
			return err
		}

		// flush field data
		if err := sm.flusher.FlushField(data); err != nil {
			return err
		}
		encodeStream.Reset() // reset tsd compress stream for next loop
	}

	// need mark metricReader completed, because next series id maybe haven't field data in metricReader,
	// if it doesn't mark metricReader completed, some data will read duplicate.
	for _, reader := range fieldReaders {
		if reader != nil {
			reader.Close()
		}
	}
	return nil
}

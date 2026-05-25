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
	"math"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./series_merger.go -destination=./series_merger_mock.go -package metricsdata

// SeriesMerger represents series data merger which merge multi fields under same series id
type SeriesMerger interface {
	// merge the multi-fields' data with same series id
	merge(mergeCtx *mergerContext,
		highKey, lowSeriesID uint16,
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
	highKey, lowSeriesID uint16,
) (err error) {
	for idx, f := range mergeCtx.allFields {
		fieldID := f.ID
		isExemplar := f.Type.IsExemplar()
		var encodeStream *encoding.TSDEncoder
		if isExemplar {
		} else {
			encodeStream = sm.flusher.GetEncoder(idx)
			encodeStream.RestWithStartTime(mergeCtx.targetRange.Start)
		}

		// merger field data from multi scanners(same series)
		for _, scanner := range mergeCtx.scanners {
			seriesEntry := scanner.scan(highKey, lowSeriesID)
			if len(seriesEntry) == 0 {
				// maybe series id not exist in some values block
				continue
			}
			// initial field reader
			timeRange := scanner.slotRange()
			if mergeCtx.fieldReader == nil {
				mergeCtx.fieldReader, err = newFieldReader(scanner.fieldIndexes(), seriesEntry)
			} else {
				err = mergeCtx.fieldReader.Reset(scanner.fieldIndexes(), seriesEntry)
			}
			if err != nil {
				return err
			}
			fieldData := mergeCtx.fieldReader.GetFieldData(fieldID)
			if len(fieldData) == 0 {
				// maybe field not exist in some values block
				continue
			}

			if isExemplar {
				getter := newExemplarTSDGetter(fieldData)
				downsampling(mergeCtx, timeRange, mergeCtx.exemplars, getter.GetExemplar, field.ExemplarAggregate)
			} else {
				// reset tsd data
				mergeCtx.decoder.ResetWithTimeRange(fieldData, timeRange.Start, timeRange.End)
				downsampling(mergeCtx, timeRange, mergeCtx.values, mergeCtx.decoder.GetValue, f.Type.Aggregate)
			}
		}

		if isExemplar {
			w := stream.NewBufferWriter(nil)

			w.PutUInt16(uint16(mergeCtx.exemplars.Size()))
			for pos := mergeCtx.targetRange.Start; pos <= mergeCtx.targetRange.End; pos++ {
				targetPos := int(pos)

				if mergeCtx.exemplars.HasValue(targetPos) {
					ex := mergeCtx.exemplars.GetValue(targetPos)
					w.PutUInt16(pos)
					w.PutUvarint32(uint32(len(ex.TraceID)))
					w.PutBytes(strutil.String2ByteSlice(ex.TraceID))
					w.PutUvarint32(uint32(len(ex.SpanID)))
					w.PutBytes(strutil.String2ByteSlice(ex.SpanID))
					w.PutVarint64(ex.Duration)
				}
			}

			data, err := w.Bytes()
			if err != nil {
				return err
			}
			// flush field data
			if err := sm.flusher.FlushField(data); err != nil {
				return err
			}

			mergeCtx.exemplars.Reset() // reset exemplars for next field
		} else {
			for pos := mergeCtx.targetRange.Start; pos <= mergeCtx.targetRange.End; pos++ {
				targetPos := int(pos)

				if mergeCtx.values.HasValue(targetPos) {
					encodeStream.AppendTime(bit.One)
					encodeStream.AppendValue(math.Float64bits(mergeCtx.values.GetValue(targetPos)))
				} else {
					encodeStream.AppendTime(bit.Zero)
				}
			}

			data, err := encodeStream.BytesWithoutTime()
			if err != nil {
				return err
			}

			// flush field data
			if err := sm.flusher.FlushField(data); err != nil {
				return err
			}
			mergeCtx.values.Reset() // reset values for next field
			encodeStream.Reset()    // reset tsd compress stream for next loop
		}
	}

	return nil
}

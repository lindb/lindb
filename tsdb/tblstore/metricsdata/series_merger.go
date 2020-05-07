package metricsdata

import (
	"encoding/binary"
	"math"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./series_merger.go -destination=./series_merger_mock.go -package metricsdata

// SeriesMerger represents series data merger which merge multi fields under same series id
type SeriesMerger interface {
	// merge merges the multi fields data with same series id
	merge(mergeCtx *mergerContext,
		decodeStreams []*encoding.TSDDecoder, encodeStream encoding.TSDEncoder,
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

// merge merges the multi fields data with same series id
func (sm *seriesMerger) merge(mergeCtx *mergerContext,
	streams []*encoding.TSDDecoder, encodeStream encoding.TSDEncoder,
	fieldReaders []FieldReader,
) error {

	for _, f := range mergeCtx.targetFields {
		schema := f.Type.GetSchema()
		fieldID := f.ID
		primitiveFields := schema.GetAllPrimitiveFields()
		for _, primitiveID := range primitiveFields {
			aggFunc := schema.GetAggFunc(primitiveID)
			for idx, reader := range fieldReaders {
				if reader == nil {
					// if series id not exist, reader is nil
					continue
				}
				fieldData := reader.getPrimitiveData(fieldID, primitiveID)
				if len(fieldData) > 0 {
					if streams[idx] == nil {
						// new tsd decoder
						streams[idx] = encoding.GetTSDDecoder()
					}
					oldStart, oldEnd := reader.slotRange()
					// reset tsd data
					streams[idx].ResetWithTimeRange(fieldData, oldStart, oldEnd)
				}
			}
			// merge field data
			sm.mergeField(mergeCtx, aggFunc, encodeStream, streams)
			data, err := encodeStream.BytesWithoutTime()
			if err != nil {
				return err
			}

			// flush field data
			sm.flusher.FlushField(field.Key(binary.LittleEndian.Uint16([]byte{byte(fieldID), byte(primitiveID)})), data)
			encodeStream.Reset() // reset tsd compress stream for next loop
		}
	}

	// need mark reader completed, because next series id maybe haven't field data in reader,
	// if don't mark reader completed, some data will read duplicate.
	for _, reader := range fieldReaders {
		if reader != nil {
			reader.close()
		}
	}
	return nil
}

// mergeField merges field data from source time range => target time range,
// compact merge: source range = target range and ratio = 1
// rollup merge: source range[5,182]=>target range[0,6], ratio:30, source interval:10s, target interval:5min
func (sm *seriesMerger) mergeField(mergeCtx *mergerContext, aggFunc field.AggFunc,
	stream encoding.TSDEncoder, values []*encoding.TSDDecoder,
) {
	hasValue := false
	pos := mergeCtx.sourceStart
	result := 0.0
	// first loop: target slot range
	for j := mergeCtx.targetStart; j <= mergeCtx.targetEnd; j++ {
		// second loop: source slot range and ratio(target interval/source interval)
		intervalEnd := mergeCtx.ratio * (j + 1)
		for pos <= mergeCtx.sourceEnd && pos < intervalEnd {
			// 1. merge data by time slot
			for _, value := range values {
				if value == nil {
					// if series id not exist, value maybe nil
					continue
				}
				if value.HasValueWithSlot(pos) {
					if !hasValue {
						// if target value not exist, set it
						result = math.Float64frombits(value.Value())
						hasValue = true
					} else {
						// if target value exist, do aggregate
						result = aggFunc.Aggregate(result, math.Float64frombits(value.Value()))
					}
				}
			}
			pos++
		}
		// 2. add data into tsd stream
		if hasValue {
			stream.AppendTime(bit.One)
			stream.AppendValue(math.Float64bits(result))
			// reset has value for next loop
			hasValue = false
			result = 0.0
		} else {
			stream.AppendTime(bit.Zero)
		}
	}
}

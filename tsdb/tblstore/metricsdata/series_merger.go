package metricsdata

import (
	"encoding/binary"
	"math"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./series_merger.go -destination=./series_merger_mock.go -package metricsdata

type SeriesMerger interface {
	merge(fields field.Metas,
		decodeStreams []*encoding.TSDDecoder, encodeStream encoding.TSDEncoder,
		fieldReaders []FieldReader,
		start, end uint16,
	) error
}

type seriesMerger struct {
	flusher Flusher
}

func newSeriesMerger(flusher Flusher) SeriesMerger {
	return &seriesMerger{
		flusher: flusher,
	}
}

func (sm *seriesMerger) merge(fields field.Metas,
	streams []*encoding.TSDDecoder, encodeStream encoding.TSDEncoder,
	fieldReaders []FieldReader,
	start, end uint16,
) error {
	for _, f := range fields {
		schema := f.Type.GetSchema()
		fieldID := f.ID
		primitiveFields := schema.GetAllPrimitiveFields()
		for _, primitiveID := range primitiveFields {
			aggFunc := schema.GetAggFunc(primitiveID)
			for idx, reader := range fieldReaders {
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
			sm.mergeField(aggFunc, encodeStream, streams, start, end)
			data, err := encodeStream.BytesWithoutTime()
			if err != nil {
				return err
			}

			// flush field data
			sm.flusher.FlushField(field.Key(binary.LittleEndian.Uint16([]byte{byte(fieldID), byte(primitiveID)})), data)
			encodeStream.Reset() // reset tsd compress stream for next loop
		}
	}
	return nil
}

func (sm *seriesMerger) mergeField(aggFunc field.AggFunc,
	stream encoding.TSDEncoder, values []*encoding.TSDDecoder,
	start, end uint16,
) {
	hasValue := false
	target := 0.0
	for i := start; i <= end; i++ {
		// 1. merge data by time slot
		for _, value := range values {
			if value.HasValueWithSlot(i) {
				if !hasValue {
					// if target value not exist, set it
					target = math.Float64frombits(value.Value())
					hasValue = true
				} else {
					// if target value exist, do aggregate
					target = aggFunc.Aggregate(target, math.Float64frombits(value.Value()))
				}
			}
		}
		// 2. add data into tsd stream
		if hasValue {
			stream.AppendTime(bit.One)
			stream.AppendValue(math.Float64bits(target))
			// reset has value for next loop
			hasValue = false
		} else {
			stream.AppendTime(bit.Zero)
		}
	}
}

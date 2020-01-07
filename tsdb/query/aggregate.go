package query

import (
	"math"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
)

var aggLogger = logger.GetLogger("tsdb", "fieldAgg")

// Aggregate aggregates the field data for query aggregator
func Aggregate(fieldType field.Type, agg aggregation.FieldAggregator, tsd *encoding.TSDDecoder, data []byte) {
	switch fieldType {
	case field.SumField, field.GaugeField, field.MinField, field.MaxField:
		tsd.Reset(data)
		aggSimpleField(agg, tsd)
	case field.SummaryField:
		aggComplexField(agg, data, tsd)
	default:
		aggLogger.Error("unknown field type when does query aggregate")
	}
}

// aggSimpleField aggregates the simple field data by simple field store layout
func aggSimpleField(agg aggregation.FieldAggregator, tsd *encoding.TSDDecoder) {
	aggregators := agg.GetAllAggregators()
	for tsd.Next() {
		if tsd.HasValue() {
			timeSlot := tsd.Slot()
			val := tsd.Value()
			for _, a := range aggregators {
				a.Aggregate(timeSlot, math.Float64frombits(val))
			}
		}
	}
}

func aggComplexField(agg aggregation.FieldAggregator, data []byte, tsd *encoding.TSDDecoder) {
	//FIXME stone1100, need read primitive field ids
	//fieldIDs := collections.NewBitArray(data[0:1])
	offsets := encoding.NewFixedOffsetDecoder(data[1:])
	dataOffset := 1 + offsets.ValueWidth()*2 + 1
	aggregators := agg.GetAllAggregators()
	for i := 0; i < 2; i++ {
		offset := offsets.Get(i)
		tsd.Reset(data[offset+dataOffset:])
		for tsd.Next() {
			if tsd.HasValue() {
				timeSlot := tsd.Slot()
				val := tsd.Value()
				if i == 0 {
					aggregators[i].Aggregate(timeSlot, math.Float64frombits(val))
				} else {
					aggregators[i].Aggregate(timeSlot, float64(encoding.ZigZagDecode(val)))
				}
			}
		}
	}
}

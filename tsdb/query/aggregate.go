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

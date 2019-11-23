package memdb

import (
	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
)

// memScanContext represents the memory metric store scan context
type memScanContext struct {
	fieldIDs    []uint16
	aggregators aggregation.FieldAggregates
	tsd         *encoding.TSDDecoder

	fieldCount int
}

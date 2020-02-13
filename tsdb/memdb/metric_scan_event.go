package memdb

import (
	"github.com/lindb/lindb/pkg/encoding"
)

// memScanContext represents the memory metric store scan context
type memScanContext struct {
	fieldAggs []*fieldAggregator
	tsd       *encoding.TSDDecoder
}

package series

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanContext_GetAggregator(t *testing.T) {
	sCtx := &ScanContext{
		FieldIDs: []uint16{3, 4, 5},
	}
	sCtx.Aggregates = sync.Pool{
		New: func() interface{} {
			return "mock_agg"
		},
	}
	agg := sCtx.GetAggregator()
	assert.Equal(t, "mock_agg", agg)
	sCtx.Release(agg)
}

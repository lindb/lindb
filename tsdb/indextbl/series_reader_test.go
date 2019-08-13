package indextbl

import (
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/stretchr/testify/assert"
)

func Test_SeriesIndexReader(t *testing.T) {
	reader := NewSeriesIndexReader(nil)
	assert.NotNil(t, reader)

	// GetTagValues
	tagValues, err := reader.GetTagValues(1, nil, 0)
	assert.Nil(t, tagValues)
	assert.Nil(t, err)

	// FindSeriesIDsByExpr
	set, err := reader.FindSeriesIDsByExpr(1, nil, timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.Nil(t, err)

	// GetSeriesIDsForTag
	set, err = reader.GetSeriesIDsForTag(1, "", timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.Nil(t, err)
}

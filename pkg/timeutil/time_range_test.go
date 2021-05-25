package timeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeRange_IsEmpty(t *testing.T) {
	assert.True(t, (&TimeRange{Start: 10, End: 10}).IsEmpty())
	assert.True(t, (&TimeRange{Start: 100, End: 10}).IsEmpty())
	assert.False(t, (&TimeRange{Start: 10, End: 100}).IsEmpty())
}

func TestTimeRange_Contains(t *testing.T) {
	timeRange := &TimeRange{Start: 10, End: 100}
	assert.True(t, timeRange.Contains(10))
	assert.True(t, timeRange.Contains(59))
	assert.True(t, timeRange.Contains(100))

	assert.False(t, timeRange.Contains(5))
	assert.False(t, timeRange.Contains(101))
}

func TestTimeRange_Overlap(t *testing.T) {
	timeRange := &TimeRange{Start: 10, End: 100}
	assert.True(t, timeRange.Overlap(&TimeRange{Start: 10, End: 1000}))
	assert.True(t, timeRange.Overlap(&TimeRange{Start: 6, End: 100}))
	assert.True(t, timeRange.Overlap(&TimeRange{Start: 60, End: 70}))

	assert.False(t, timeRange.Overlap(&TimeRange{Start: 6, End: 9}))
	assert.False(t, timeRange.Overlap(&TimeRange{Start: 600, End: 900}))
}

func TestTimeRange_Intersect(t *testing.T) {
	timeRange := &TimeRange{Start: 10, End: 100}
	assert.Equal(t, TimeRange{Start: 10, End: 100}, *timeRange.Intersect(&TimeRange{Start: 10, End: 100}))
	assert.Equal(t, TimeRange{Start: 50, End: 60}, *timeRange.Intersect(&TimeRange{Start: 50, End: 60}))
	assert.Equal(t, TimeRange{Start: 50, End: 100}, *timeRange.Intersect(&TimeRange{Start: 50, End: 1000}))
	assert.Equal(t, TimeRange{Start: 10, End: 100}, *timeRange.Intersect(&TimeRange{Start: 5, End: 1000}))
	assert.Equal(t, TimeRange{Start: 10, End: 60}, *timeRange.Intersect(&TimeRange{Start: 5, End: 60}))

	assert.True(t, timeRange.Intersect(&TimeRange{Start: 7, End: 5}).IsEmpty())
	assert.True(t, timeRange.Intersect(&TimeRange{Start: 5, End: 7}).IsEmpty())
	assert.True(t, timeRange.Intersect(&TimeRange{Start: 500, End: 7000}).IsEmpty())
	assert.True(t, timeRange.Intersect(&TimeRange{Start: 5000, End: 700}).IsEmpty())
}

func TestSlotRange(t *testing.T) {
	sr := NewSlotRange(10, 20)
	sr.SetSlot(15)
	start, end := sr.GetRange()
	assert.Equal(t, uint16(10), start)
	assert.Equal(t, uint16(20), end)
	sr.SetSlot(5)
	start, end = sr.GetRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(20), end)
	sr.SetSlot(27)
	start, end = sr.GetRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(27), end)
}

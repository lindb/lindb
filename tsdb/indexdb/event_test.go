package indexdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappingEvent(t *testing.T) {
	e := newMappingEvent()
	assert.True(t, e.isEmpty())
	assert.False(t, e.isFull())
	e.addSeriesID(1, 10, 100)
	e.addSeriesID(1, 20, 120)
	e.addSeriesID(2, 30, 100)
	e.addSeriesID(2, 40, 200)
	assert.Len(t, e.events, 2)
	assert.Equal(t, []seriesEvent{{seriesID: 100, tagsHash: 10}, {seriesID: 120, tagsHash: 20}}, e.events[1].events)
	assert.Equal(t, uint32(120), e.events[1].metricIDSeq)
	assert.Equal(t, []seriesEvent{{seriesID: 100, tagsHash: 30}, {seriesID: 200, tagsHash: 40}}, e.events[2].events)
	assert.Equal(t, uint32(200), e.events[2].metricIDSeq)
	assert.False(t, e.isEmpty())
	for i := 0; i < full; i++ {
		e.addSeriesID(2, uint64(i), uint32(200+i))
	}
	assert.True(t, e.isFull())
}

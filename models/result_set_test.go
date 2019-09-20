package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultSet(t *testing.T) {
	rs := NewResultSet()
	series := NewSeries(map[string]string{"key": "value"})
	rs.AddSeries(series)
	points := NewPoints()
	points.AddPoint(int64(10), 10.0)
	series.AddField("f1", points)

	assert.Equal(t, 1, len(rs.Series))
	s := rs.Series[0]
	assert.Equal(t, map[string]string{"key": "value"}, s.Tags)
	assert.Equal(t, map[int64]float64{int64(10): 10.0}, s.Fields["f1"])
}

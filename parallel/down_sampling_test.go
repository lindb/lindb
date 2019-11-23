package parallel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
)

func Test_downSamplingTimeRange(t *testing.T) {
	timeRange, intervalRatio, interval := downSamplingTimeRange(
		timeutil.Interval(30*timeutil.OneSecond),
		timeutil.Interval(10*timeutil.OneSecond),
		timeutil.TimeRange{
			Start: 35 * timeutil.OneSecond,
			End:   65 * timeutil.OneSecond,
		})
	assert.Equal(t, 3, intervalRatio)
	assert.Equal(t, 30*timeutil.OneSecond, interval.Int64())
	assert.Equal(t, timeutil.TimeRange{
		Start: 30 * timeutil.OneSecond,
		End:   60 * timeutil.OneSecond,
	}, timeRange)
}

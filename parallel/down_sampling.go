package parallel

import (
	"github.com/lindb/lindb/pkg/timeutil"
)

// downSamplingTimeRange returns down sampling time range and interval ratio
func downSamplingTimeRange(queryInterval,
	storageInterval timeutil.Interval,
	queryTimeRange timeutil.TimeRange,
) (
	timeRange timeutil.TimeRange,
	intervalRatio int,
	interval timeutil.Interval,
) {
	// 1. calc interval, default use storage interval's interval if user not input
	interval = storageInterval
	intervalRatio = 1
	if queryInterval > 0 {
		intervalRatio = timeutil.CalIntervalRatio(queryInterval.Int64(), interval.Int64())
		interval = queryInterval
	}
	// 2. truncate time range
	timeRange = timeutil.TimeRange{
		Start: timeutil.Truncate(queryTimeRange.Start, interval.Int64()),
		End:   timeutil.Truncate(queryTimeRange.End, interval.Int64()),
	}
	return
}

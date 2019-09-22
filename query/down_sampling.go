package query

import (
	"github.com/lindb/lindb/pkg/timeutil"
)

// downSamplingTimeRange returns down sampling time range and interval ratio
func downSamplingTimeRange(queryInterval,
	storageInterval int64,
	queryTimeRange timeutil.TimeRange,
) (timeRange timeutil.TimeRange, intervalRatio int, interval int64) {
	// 1. calc interval, default use storage interval's interval if user not input
	interval = storageInterval
	intervalRatio = 1
	if queryInterval > 0 {
		intervalRatio = timeutil.CalIntervalRatio(queryInterval, interval)
		interval = queryInterval
	}
	// 2. truncate time range
	timeRange = timeutil.TimeRange{
		Start: timeutil.Truncate(queryTimeRange.Start, interval),
		End:   timeutil.Truncate(queryTimeRange.End, interval),
	}
	return
}

package timeutil

import "time"

const (
	// OneHour is the number of millisecond for an hour
	OneHour = 60 * 60 * 1000
	// OneDay is the number of millisecond for an day
	OneDay = 24 * 60 * 60 * 1000
)

// FormatTimestamp returns timestamp format based on layout
func FormatTimestamp(timestamp int64, layout string) string {
	t := time.Unix(timestamp/1000, 0)
	return t.Format(layout)
}

// ParseTimestamp parses timestamp str value based on layout using local zone
func ParseTimestamp(timestampStr, layout string) (int64, error) {
	tm, err := time.ParseInLocation(layout, timestampStr, time.Local)
	if err != nil {
		return 0, err
	}
	return tm.UnixNano() / 1000000, nil
}

// Now returns t as a Unix time, the number of millisecond elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func Now() int64 {
	return time.Now().UnixNano() / 1000000
}

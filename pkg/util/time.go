package util

import "time"

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

package models

import "time"

// Query represents search condition
type Query interface {
	// MetricName returns metric name for search, like table name
	MetricName() string
	// TimeRange returns query time range
	TimeRange() TimeRange
	// Interval returns query time interval
	Interval() time.Duration
}

// TimeRange represents time range with start/end timestamp
type TimeRange struct {
	Start, End int64
}

// TagFilter is a filter of metric-tag.
type TagFilter struct {
	TagName  string
	TagValue string
}

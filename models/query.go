package models

import (
	"time"

	"github.com/eleme/lindb/pkg/timeutil"
)

// Query represents search condition
type Query interface {
	// MetricName returns metric name for search, like table name
	MetricName() string
	// TimeRange returns query time range
	TimeRange() timeutil.TimeRange
	// Interval returns query time interval
	Interval() time.Duration
}

// TagFilter is a filter of metric-tag.
type TagFilter struct {
	Key    string
	Values []string
}

type Condition struct {
	TagFilters []*TagFilter
}

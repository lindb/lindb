package tree

import (
	"encoding/json"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/pkg/timeutil"
)

// Query1 represents search statement
type Query1 struct {
	Explain     bool   // need explain query execute stat
	Namespace   string // namespace
	MetricName  string // like table name
	SelectItems []Expr // select list, such as field, function call, math expression etc.
	AllFields   bool   // select all fields under metric
	Condition   Expr   // tag filter condition expression

	// broker plan maybe reset
	TimeRange       timeutil.TimeRange // query time range
	Interval        timeutil.Interval  // down sampling storage interval
	StorageInterval timeutil.Interval  // down sampling storage interval, data find
	IntervalRatio   int                // down sampling interval ratio(query interval/storage Interval)
	AutoGroupByTime bool               // auto fix group by interval based on query time range

	GroupBy      []string // group by tag keys
	Having       Expr     // having clause
	OrderByItems []Expr   // order by field expr list
	Limit        int      // num. of time series list for result
}

// StatementType returns metric query type.
func (q *Query1) StatementType() StatementType {
	return QueryStatement
}

// HasGroupBy returns whether query has grouping tag keys
func (q *Query1) HasGroupBy() bool {
	return len(q.GroupBy) > 0
}

// innerQuery represents a wrapper of query for json encoding
type innerQuery struct {
	Explain     bool              `json:"explain,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	MetricName  string            `json:"metricName,omitempty"`
	SelectItems []json.RawMessage `json:"selectItems,omitempty"`
	AllFields   bool              `json:"allFields,omitempty"`
	Condition   json.RawMessage   `json:"condition,omitempty"`

	TimeRange       timeutil.TimeRange `json:"timeRange,omitempty"`
	Interval        timeutil.Interval  `json:"interval,omitempty"`
	StorageInterval timeutil.Interval  `json:"storageInterval,omitempty"`
	IntervalRatio   int                `json:"intervalRatio,omitempty"`
	AutoGroupByTime bool               `json:"autoGroupByTime,omitempty"`

	GroupBy      []string          `json:"groupBy,omitempty"`
	Having       json.RawMessage   `json:"having,omitempty"`
	OrderByItems []json.RawMessage `json:"orderByItems,omitempty"`
	Limit        int               `json:"limit,omitempty"`
}

// MarshalJSON returns json data of query
func (q *Query1) MarshalJSON() ([]byte, error) {
	inner := innerQuery{
		Explain:         q.Explain,
		MetricName:      q.MetricName,
		AllFields:       q.AllFields,
		Namespace:       q.Namespace,
		Condition:       Marshal(q.Condition),
		TimeRange:       q.TimeRange,
		Interval:        q.Interval,
		IntervalRatio:   q.IntervalRatio,
		AutoGroupByTime: q.AutoGroupByTime,
		StorageInterval: q.StorageInterval,
		GroupBy:         q.GroupBy,
		Having:          Marshal(q.Having),
		Limit:           q.Limit,
	}
	for _, item := range q.SelectItems {
		inner.SelectItems = append(inner.SelectItems, Marshal(item))
	}
	for _, item := range q.OrderByItems {
		inner.OrderByItems = append(inner.OrderByItems, Marshal(item))
	}
	return encoding.JSONMarshal(&inner), nil
}

// UnmarshalJSON parses json data to query
func (q *Query1) UnmarshalJSON(value []byte) error {
	inner := innerQuery{}
	if err := encoding.JSONUnmarshal(value, &inner); err != nil {
		return err
	}

	if inner.Condition != nil {
		condition, err := Unmarshal(inner.Condition)
		if err != nil {
			return err
		}
		q.Condition = condition
	}

	if inner.Having != nil {
		having, err := Unmarshal(inner.Having)
		if err != nil {
			return err
		}
		q.Having = having
	}

	// select list
	var selectItems []Expr
	for _, item := range inner.SelectItems {
		selectItem, err := Unmarshal(item)
		if err != nil {
			return err
		}
		selectItems = append(selectItems, selectItem)
	}
	// order by list
	var orderByItems []Expr
	for _, item := range inner.OrderByItems {
		orderByItem, err := Unmarshal(item)
		if err != nil {
			return err
		}
		orderByItems = append(orderByItems, orderByItem)
	}

	q.Explain = inner.Explain
	q.MetricName = inner.MetricName
	q.Namespace = inner.Namespace
	q.SelectItems = selectItems
	q.AllFields = inner.AllFields
	q.TimeRange = inner.TimeRange
	q.Interval = inner.Interval
	q.IntervalRatio = inner.IntervalRatio
	q.AutoGroupByTime = inner.AutoGroupByTime
	q.StorageInterval = inner.StorageInterval
	q.GroupBy = inner.GroupBy
	q.OrderByItems = orderByItems
	q.Limit = inner.Limit
	return nil
}

// MetricMetadataType represents metric metadata suggest type
type MetricMetadataType uint8

// Defines all types of metric metadata suggest
const (
	Namespace MetricMetadataType = iota + 1
	Metric
	TagKey
	TagValue
	Field
)

// String returns string value of metadata type
func (m MetricMetadataType) String() string {
	switch m {
	case Namespace:
		return "namespace"
	case Metric:
		return "metric"
	case Field:
		return field
	case TagKey:
		return "tagKey"
	case TagValue:
		return "tagValue"
	default:
		return unknown
	}
}

// MetricMetadata represents search metric metadata statement
type MetricMetadata struct {
	Namespace  string             // namespace
	MetricName string             // like table name
	Type       MetricMetadataType // metadata suggest type
	TagKey     string
	Prefix     string
	Condition  Expr // tag filter condition expression
	Limit      int  // result set limit
}

// StatementType returns metadata query type.
func (q *MetricMetadata) StatementType() StatementType {
	return MetricMetadataStatement
}

// innerMetadata represents a wrapper of metadata for json encoding
type innerMetadata struct {
	Namespace  string             `json:"namespace,omitempty"`
	MetricName string             `json:"metricName,omitempty"`
	Type       MetricMetadataType `json:"type,omitempty"`
	TagKey     string             `json:"tagKey,omitempty"`
	Condition  json.RawMessage    `json:"condition,omitempty"`
	Prefix     string             `json:"prefix,omitempty"`
	Limit      int                `json:"limit,omitempty"`
}

// MarshalJSON returns json data of query
func (q *MetricMetadata) MarshalJSON() ([]byte, error) {
	inner := innerMetadata{
		MetricName: q.MetricName,
		Namespace:  q.Namespace,
		Condition:  Marshal(q.Condition),
		TagKey:     q.TagKey,
		Type:       q.Type,
		Prefix:     q.Prefix,
		Limit:      q.Limit,
	}
	return encoding.JSONMarshal(&inner), nil
}

// UnmarshalJSON parses json data to metadata
func (q *MetricMetadata) UnmarshalJSON(value []byte) error {
	inner := innerMetadata{}
	if err := encoding.JSONUnmarshal(value, &inner); err != nil {
		return err
	}
	if inner.Condition != nil {
		condition, err := Unmarshal(inner.Condition)
		if err != nil {
			return err
		}
		q.Condition = condition
	}
	q.Namespace = inner.Namespace
	q.MetricName = inner.MetricName
	q.Type = inner.Type
	q.TagKey = inner.TagKey
	q.Prefix = inner.Prefix
	q.Limit = inner.Limit
	return nil
}

const (
	// StateRepoSource represents from state persist repo.
	StateRepoSource = iota + 1
	// StateMachineSource represents from state machine in current memory.
	StateMachineSource
)

// MetadataType represents metadata type.
type MetadataType int

const (
	// MetadataTypes represents all metadata types.
	MetadataTypes MetadataType = iota + 1
	// BrokerMetadata represent broker metadata.
	BrokerMetadata
	// MasterMetadata represent master metadata.
	MasterMetadata
	// StorageMetadata represent storage metadata.
	StorageMetadata
	// RootMetadata represent root metadata.
	RootMetadata
)

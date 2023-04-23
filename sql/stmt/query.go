// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package stmt

import (
	"encoding/json"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

// Query represents search statement
type Query struct {
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
	OrderByItems []Expr   // order by field expr list
	Limit        int      // num. of time series list for result
}

// StatementType returns metric query type.
func (q *Query) StatementType() StatementType {
	return QueryStatement
}

// HasGroupBy returns whether query has grouping tag keys
func (q *Query) HasGroupBy() bool {
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
	OrderByItems []json.RawMessage `json:"orderByItems,omitempty"`
	Limit        int               `json:"limit,omitempty"`
}

// MarshalJSON returns json data of query
func (q *Query) MarshalJSON() ([]byte, error) {
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
func (q *Query) UnmarshalJSON(value []byte) error {
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

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

package models

import (
	"fmt"
	"path"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lindb/lindb/pkg/timeutil"
)

// SuggestResult represents the suggest result set
type SuggestResult struct {
	Values []string `json:"values"`
}

// ResultSet represents the query result set
type ResultSet struct {
	MetricName string      `json:"metricName,omitempty"`
	GroupBy    []string    `json:"groupBy,omitempty"`
	Fields     []string    `json:"fields,omitempty"`
	StartTime  int64       `json:"startTime,omitempty"`
	EndTime    int64       `json:"endTime,omitempty"`
	Interval   int64       `json:"interval,omitempty"`
	Series     []*Series   `json:"series,omitempty"`
	Stats      *QueryStats `json:"stats,omitempty"`
}

// NewResultSet creates a new result set
func NewResultSet() *ResultSet {
	return &ResultSet{}
}

// AddSeries adds a new series
func (rs *ResultSet) AddSeries(series *Series) {
	rs.Series = append(rs.Series, series)
}

// row represents a record in table.
type row struct {
	timestamp int64
	values    map[string]float64
	tags      map[string]string
}

// ToTable returns the result of query as table if it has value, else return empty string.
func (rs *ResultSet) ToTable() (rows int, tableStr string) {
	// if explain query return query plan
	if rs.Stats != nil {
		return rs.Stats.ToTable()
	}
	if len(rs.Series) == 0 {
		return 0, ""
	}
	// 1. set headers
	headers := table.Row{}
	for _, k := range rs.GroupBy {
		headers = append(headers, k)
	}
	headers = append(headers, "timestamp")
	for _, f := range rs.Fields {
		headers = append(headers, f)
	}
	// 2. build table rows
	tableRows := make(map[string]*row)
	var pks []string
	for _, s := range rs.Series {
		var values []string
		for _, tagKey := range rs.GroupBy {
			values = append(values, s.Tags[tagKey])
		}
		key := path.Join(values...)
		for n, f := range s.Fields {
			for timestamp, v := range f {
				k := fmt.Sprintf("%s_%d", key, timestamp)
				i, ok := tableRows[k]
				if !ok {
					i = &row{values: make(map[string]float64), tags: s.Tags, timestamp: timestamp}
					pks = append(pks, k)
					tableRows[k] = i
				}
				i.values[n] = v
			}
		}
	}
	// 3. format as table
	result := NewTableFormatter()
	result.AppendHeader(headers)
	sort.Strings(pks)
	for _, pk := range pks {
		r := tableRows[pk]
		row := table.Row{}
		for _, tagKey := range rs.GroupBy {
			row = append(row, r.tags[tagKey])
		}
		row = append(row, timeutil.FormatTimestamp(r.timestamp, timeutil.DataTimeFormat2))
		for _, f := range rs.Fields {
			row = append(row, r.values[f])
		}
		result.AppendRow(row)
	}
	return len(rs.Series), result.Render()
}

// Series represents one time series for metric.
type Series struct {
	Tags   map[string]string            `json:"tags,omitempty"`
	Fields map[string]map[int64]float64 `json:"fields,omitempty"`

	TagValues string `json:"-"` // return series in order by tag values
}

// NewSeries creates a new series.
func NewSeries(tags map[string]string, tagValues string) *Series {
	return &Series{
		Tags:      tags,
		Fields:    make(map[string]map[int64]float64),
		TagValues: tagValues,
	}
}

// AddField adds a field
func (s *Series) AddField(fieldName string, points *Points) {
	dataPoints, ok := s.Fields[fieldName]
	if !ok {
		s.Fields[fieldName] = points.Points
		return
	}
	for t, v := range points.Points {
		dataPoints[t] = v
	}
}

// Points represents the data points of the field
type Points struct {
	Points map[int64]float64 `json:"points,omitempty"`
}

// NewPoints creates the data point
func NewPoints() *Points {
	return &Points{Points: make(map[int64]float64)}
}

// AddPoint adds point
func (p *Points) AddPoint(timestamp int64, value float64) {
	p.Points[timestamp] = value
}

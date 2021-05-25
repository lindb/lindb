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

// SuggestResult represents the suggest result set
type SuggestResult struct {
	Values []string `json:"values"`
}

// ResultSet represents the query result set
type ResultSet struct {
	MetricName string      `json:"metricName,omitempty"`
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

// Series represents one time series for metric
type Series struct {
	Tags   map[string]string            `json:"tags,omitempty"`
	Fields map[string]map[int64]float64 `json:"fields,omitempty"`
}

// NewSeries creates a new series
func NewSeries(tags map[string]string) *Series {
	return &Series{Tags: tags, Fields: make(map[string]map[int64]float64)}
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

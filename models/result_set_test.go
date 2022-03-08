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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
)

func TestResultSet(t *testing.T) {
	rs := NewResultSet()
	series := NewSeries(map[string]string{"key": "value"})
	rs.AddSeries(series)
	points := NewPoints()
	points.AddPoint(int64(10), 10.0)
	series.AddField("f1", points)
	points = NewPoints()
	points.AddPoint(int64(20), 10.0)
	series.AddField("f1", points)

	assert.Equal(t, 1, len(rs.Series))
	s := rs.Series[0]
	assert.Equal(t, map[string]string{"key": "value"}, s.Tags)
	assert.Equal(t, map[int64]float64{
		int64(10): 10.0,
		int64(20): 10.0},
		s.Fields["f1"])
}

func TestResultSet_ToTable(t *testing.T) {
	rows, rs := NewResultSet().ToTable()
	assert.Zero(t, rows)
	assert.Empty(t, rs)

	rows, rs = (&ResultSet{
		MetricName: "cpu",
		GroupBy:    []string{"host", "ip"},
		Fields:     []string{"usage", "load"},
		Series: []*Series{{
			Tags:   map[string]string{"host": "host1", "ip": "1.1.1.1"},
			Fields: map[string]map[int64]float64{"usage": {timeutil.Now(): 1.1}, "load": {timeutil.Now(): 1.1}},
		}, {
			Tags:   map[string]string{"host": "host2", "ip": "1.1.1.1"},
			Fields: map[string]map[int64]float64{"usage": {timeutil.Now(): 1.1}, "load": {timeutil.Now(): 1.1}},
		}},
	}).ToTable()
	assert.Equal(t, rows, 2)
	assert.NotEmpty(t, rs)
}

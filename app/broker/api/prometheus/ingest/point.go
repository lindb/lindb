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

package ingest

import (
	"strings"
	"time"
)

// Point represents time series data points with name/tags/fields etc.
type Point struct {
	namespace  string
	metricName string
	tags       map[string]string
	timestamp  time.Time
	fields     []Field
}

// NewPoint creates a data point with name.
func NewPoint(metricName string) *Point {
	return &Point{
		metricName: strings.Trim(metricName, " "),
		timestamp:  time.Now(), // default timestamp
	}
}

// SetNamespace sets namespace.
func (p *Point) SetNamespace(namespace string) *Point {
	p.namespace = namespace
	return p
}

// Namespace returns namespace.
func (p *Point) Namespace() string {
	return p.namespace
}

// MetricName returns metric name.
func (p *Point) MetricName() string {
	return p.metricName
}

// SetMetricName sets metricName.
func (p *Point) SetMetricName(metricName string) {
	p.metricName = strings.Trim(metricName, " ")
}

// SetTimestamp sets timestamp.
func (p *Point) SetTimestamp(timestamp time.Time) *Point {
	p.timestamp = timestamp
	return p
}

// Timestamp returns timestamp.
func (p *Point) Timestamp() time.Time {
	return p.timestamp
}

// AddTag adds tag(key,value).
func (p *Point) AddTag(key, value string) *Point {
	if p.tags == nil {
		p.tags = make(map[string]string)
	}
	p.tags[key] = value
	return p
}

// Tags returns tags.
func (p *Point) Tags() map[string]string {
	return p.tags
}

// AddField adds field.
func (p *Point) AddField(field Field) *Point {
	p.fields = append(p.fields, field)
	return p
}

// Fields returns fields.
func (p *Point) Fields() []Field {
	return p.fields
}

// Valid returns if point is valid.
func (p *Point) Valid() bool {
	if p.metricName == "" {
		return false
	}
	if len(p.fields) == 0 {
		return false
	}
	return true
}

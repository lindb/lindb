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

	commonconstants "github.com/lindb/common/constants"
	commonseries "github.com/lindb/common/series"
)

// Limits represents all the limit for database level; can be used to describe global
// default limits, or per-database limits vis toml config.
type Limits struct {
	// Write limits
	MaxNamespaces       int    `toml:"max-namespaces"`
	MaxNamespaceLength  int    `toml:"max-namespace-length"`
	MaxMetrics          int    `toml:"max-metrics"`
	MaxMetricNameLength int    `toml:"max-metric-name-length"`
	MaxFieldNameLength  int    `toml:"max-field-name-length"`
	MaxFieldsPerMetric  int    `toml:"max-fields-per-metric"`
	MaxTagNameLength    int    `toml:"max-tag-name-length"`
	MaxTagValueLength   int    `toml:"max-tag-value-length"`
	MaxTagsPerMetric    int    `toml:"max-tags-per-metric"`
	MaxSeriesPerMetric  uint32 `toml:"max-series-per-metric"`
	// max series limit for metric
	Metrics map[string]uint32 `toml:"metrics"`

	// Read Limits
	MaxSeriesPerQuery int `toml:"max-series-per-query"`
}

// NewDefaultLimits creates a default limits.
func NewDefaultLimits() *Limits {
	return &Limits{
		// Write limits
		MaxNamespaces:       0,
		MaxNamespaceLength:  256,
		MaxMetrics:          0,
		MaxMetricNameLength: 256,
		MaxFieldNameLength:  128,
		MaxFieldsPerMetric:  256,
		MaxTagNameLength:    128,
		MaxTagValueLength:   1024,
		MaxTagsPerMetric:    32,
		MaxSeriesPerMetric:  200000,
		Metrics:             make(map[string]uint32),
		// Read limits
		MaxSeriesPerQuery: 200000,
	}
}

// TOML returns limits' configuration string as toml format.
func (l *Limits) TOML() string {
	return fmt.Sprintf(`
## 0 to disable the limit.
## It is a per-instance limit which no special describes.

## Maximum number of active namespaces.
## Default: %d
max-namespaces = %d
## Maximum number of active metrics per namespace. 
## Default: %d
max-metrics = %d
## Maximum length accepted for namespace. 
## Default: %d
max-namespace-length = %d
## Maximum length accepted for metric name.
## Default: %d
max-metric-name-length = %d
## Maximum number of active fields per metric.
## Default: %d
max-fields-per-metric = %d
## Maximum number of active tags per metric.
## Default: %d
max-tags-per-metric = %d
## Maximum number of active series per metric.
## Default: %d
max-series-per-metric = %d
## Maximum length accepted for field name.
## Default: %d
max-field-name-length = %d
## Maximum length accepted for tag name.
## Default: %d
max-tag-name-length = %d
## Maximum length accepted for tag value.
## Default: %d
max-tag-value-length = %d

## Maximum number of series for which a query can fetch.
## Default: %d
max-series-per-query = %d

## Maximum number of active series for special metric.
## Must be the last limit configure item.
## Example: "system.cpu" = 100000
## Example: "namespace|system.cpu" = 100000
[metrics]
%s
		`,
		l.MaxNamespaces,
		l.MaxNamespaces,
		l.MaxMetrics,
		l.MaxMetrics,
		l.MaxNamespaceLength,
		l.MaxNamespaceLength,
		l.MaxMetricNameLength,
		l.MaxMetricNameLength,
		l.MaxFieldsPerMetric,
		l.MaxFieldsPerMetric,
		l.MaxTagsPerMetric,
		l.MaxTagsPerMetric,
		l.MaxSeriesPerMetric,
		l.MaxSeriesPerMetric,
		l.MaxFieldNameLength,
		l.MaxFieldNameLength,
		l.MaxTagNameLength,
		l.MaxTagNameLength,
		l.MaxTagValueLength,
		l.MaxTagValueLength,
		l.MaxSeriesPerQuery,
		l.MaxSeriesPerQuery,
		l.metricsTOML(),
	)
}

// metricsTOML returns limits' configuration for metric level.
func (l *Limits) metricsTOML() string {
	rs := ""
	for k, v := range l.Metrics {
		rs += fmt.Sprintf("%q = %d\n", k, v)
	}
	return rs
}

// GetSeriesLimit returns the limit by given namespace/metric name.
func (l *Limits) GetSeriesLimit(namespace, metricName string) uint32 {
	if len(l.Metrics) == 0 {
		return l.MaxSeriesPerMetric
	}
	key := metricName
	if namespace != commonconstants.DefaultNamespace {
		key = commonseries.JoinNamespaceMetric(namespace, metricName)
	}
	limit, ok := l.Metrics[key]
	if ok {
		return limit
	}
	return l.MaxSeriesPerMetric
}

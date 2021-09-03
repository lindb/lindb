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

package constants

import (
	"math"
)

const (
	// DefaultMaxSeriesIDsCount represents series count limit, uses this limit of metric-level when maxSeriesIDsLimit is not set
	DefaultMaxSeriesIDsCount = 10000000
	// DefaultMaxTagKeysCount represents tag key count limit, uses this limit of max tag keys of a metric
	DefaultMaxTagKeysCount = 32
	// DefaultMaxFieldsCount represents field count limit, uses this limit of max fields of a metric
	DefaultMaxFieldsCount = math.MaxUint8
	// MaxSuggestions represents the max number of suggestions count
	MaxSuggestions = 10000

	// MetricMaxAheadDuration controls the global max write ahead duration.
	// If current timestamp is 2021-08-19 23:00:00, metric after 2021-08-20 23:00:00 will be dropped.
	MetricMaxAheadDuration = 24 * 60 * 60 * 1000
	// MetricMaxBehindDuration controls the global max write behind duration.
	// If current timestamp is 2021-08-19 23:00:00, metric before 2021-08-18 23:00:00 will be dropped.
	MetricMaxBehindDuration = 24 * 60 * 60 * 1000

	// TagValueIDForTag represents tag value id placeholder for store all series ids under tag.
	TagValueIDForTag = uint32(0)
	// DefaultNamespace represents default namespace if not set
	DefaultNamespace = "default-ns"
	// SeriesIDWithoutTags represents the series ids under spec metric, but without nothing tags.
	SeriesIDWithoutTags = uint32(0)

	// EmptyValue represents the empty value.
	EmptyValue = 0.0
)

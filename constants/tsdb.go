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

import "math"

const (
	// DefaultMaxSeriesIDsCount represents series count limit, uses this limit of metric-level when maxSeriesIDsLimit is not set
	DefaultMaxSeriesIDsCount = 10000000
	// DefaultMaxTagKeysCount represents tag key count limit, uses this limit of max tag keys of a metric
	DefaultMaxTagKeysCount = 32
	// DefaultMaxFieldsCount represents field count limit, uses this limit of max fields of a metric
	DefaultMaxFieldsCount = math.MaxUint8
	// MaxSuggestions represents the max number of suggestions count
	MaxSuggestions = 10000

	// MemoryHighWaterMark checks if the global memory usage is greater than the limit,
	// If so, engine will flush the biggest shard's memory database until we are down to the lower mark.
	MemoryHighWaterMark = 80
	// MemoryLowWaterMark checks if the global memory usage is low water mark.
	MemoryLowWaterMark = 60
	// ShardMemoryUsedThreshold checks if shard's memory usage is greater than this limit,
	// If so, engine will flush this shard to disk.
	ShardMemoryUsedThreshold = 500 * 1024 * 1024
	// FlushConcurrency controls the concurrent number of flush jobs.
	FlushConcurrency = 4

	// TagValueIDForTag represents tag value id placeholder for store all series ids under tag.
	TagValueIDForTag = uint32(0)
	// DefaultNamespace represents default namespace if not set
	DefaultNamespace = "default-ns"
	// SeriesIDWithoutTags represents the series ids under spec metric, but without nothing tags.
	SeriesIDWithoutTags = uint32(0)

	// EmptyValue represents the empty value.
	EmptyValue = 0.0
)

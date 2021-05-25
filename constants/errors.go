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

import "errors"

var (
	ErrDatabaseNotFound = errors.New("database not found")
	ErrShardNotFound    = errors.New("shard not found")

	// ErrNotFound represents the data not found
	ErrNotFound = errors.New("not found")

	// ErrNilMetric represents write nil metric error
	ErrNilMetric = errors.New("metric is nil")
	// ErrEmptyMetricName represents metric name is empty when write data
	ErrEmptyMetricName = errors.New("metric name is empty")
	// ErrEmptyField represents field is empty when write data
	ErrEmptyField = errors.New("field is empty")

	// ErrDataFileCorruption represents data in tsdb's file is corrupted
	ErrDataFileCorruption = errors.New("data corruption")
)

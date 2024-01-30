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

// WriteOptions represents writer configuration.
type WriteOptions struct {
	// database config
	databaseConfig DatabaseConfig
	// Number of series sent in single writer request, default 1000.
	batchSize int
	// Flush interval(ms) which is buffer flushed if it has not been already written, default 1000.
	flushInterval int64
	// Default tags are added to each written series.
	defaultTags map[string]string
}

// SetBatchSize sets batch size in single writer request.
func (opt *WriteOptions) SetBatchSize(batchSize int) *WriteOptions {
	opt.batchSize = batchSize
	return opt
}

// BatchSize returns the number of batch size in single writer request.
func (opt *WriteOptions) BatchSize() int {
	return opt.batchSize
}

// SetFlushInterval sets flush interval(ms).
func (opt *WriteOptions) SetFlushInterval(interval int64) *WriteOptions {
	opt.flushInterval = interval
	return opt
}

// FlushInterval returns the flush interval(ms).
func (opt *WriteOptions) FlushInterval() int64 {
	return opt.flushInterval
}

// AddDefaultTag adds default tag.
func (opt *WriteOptions) AddDefaultTag(key, value string) *WriteOptions {
	if opt.defaultTags == nil {
		opt.defaultTags = make(map[string]string)
	}
	opt.defaultTags[key] = value
	return opt
}

// DefaultTags returns the default tags for all metrics.
func (opt *WriteOptions) DefaultTags() map[string]string {
	return opt.defaultTags
}

// DefaultWriteOptions creates a WriteOptions with default.
func DefaultWriteOptions(databaseConfig DatabaseConfig) *WriteOptions {
	return &WriteOptions{
		databaseConfig: databaseConfig,
		batchSize:      1_000,
		flushInterval:  1_000, // 1s
	}
}

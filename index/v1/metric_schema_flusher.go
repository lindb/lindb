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

package v1

import (
	"io"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source ./metric_schema_flusher.go -destination=./metric_schema_flusher_mock.go -package v1

// Flusher is a wrapper of kv.Builder, provides ability to flush metric schema file to disk.
type MetricSchemaFlusher interface {
	io.Closer

	// Prepare prepares to write a new metric schema block.
	Prepare(metricID uint32)
	// Write writes metric schema to kv writer.
	Write(schema *metric.Schema) error
	// Commit ends writing metric schema block.
	Commit() error
}

// metricSchemaFlusher implements MetricSchemaFlusher interface.
type metricSchemaFlusher struct {
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter
}

// NewMetricSchemaFlusher creates a metric schema flusher.
func NewMetricSchemaFlusher(kvFlusher kv.Flusher) (MetricSchemaFlusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &metricSchemaFlusher{
		kvFlusher: kvFlusher,
		kvWriter:  kvWriter,
	}, nil
}

// Prepare prepares to write a new metric schema block.
func (f *metricSchemaFlusher) Prepare(metricID uint32) {
	f.kvWriter.Prepare(metricID)
}

// Write writes metric schema to kv writer.
func (f *metricSchemaFlusher) Write(schema *metric.Schema) error {
	return schema.Write(f.kvWriter)
}

// Commit ends writing metric schema block.
func (f *metricSchemaFlusher) Commit() error {
	return f.kvWriter.Commit()
}

// Close closes kv flusher.
func (f *metricSchemaFlusher) Close() error {
	return f.kvFlusher.Commit()
}

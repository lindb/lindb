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
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/metric"
)

var MetricSchemaMerger kv.MergerType = "MetricSchemaMergerV1"

func init() {
	// register metric schema merger
	kv.RegisterMerger(MetricSchemaMerger, NewMetricScheamMerger)
}

// metricSchemaMerger implements kv.Merger interface for merging metric schema.
type metricSchemaMerger struct {
	flusher MetricSchemaFlusher
}

// NewMetricScheamMerger creates a MetricSchemaMerger instance.
func NewMetricScheamMerger(kvFlusher kv.Flusher) (kv.Merger, error) {
	flusher, err := NewMetricSchemaFlusher(kvFlusher)
	if err != nil {
		return nil, err
	}
	return &metricSchemaMerger{
		flusher: flusher,
	}, nil
}

func (m *metricSchemaMerger) Init(params map[string]interface{}) {}

// Merge merges metric schema.
func (m *metricSchemaMerger) Merge(metricID uint32, values [][]byte) error {
	schema := &metric.Schema{}
	for _, val := range values {
		schema.Unmarshal(val)
	}
	m.flusher.Prepare(metricID)
	// write schema
	if err := m.flusher.Write(schema); err != nil {
		return err
	}
	// commit metric schema
	return m.flusher.Commit()
}

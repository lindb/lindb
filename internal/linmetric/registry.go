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

package linmetric

import (
	"sync"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
)

// registry is a set of metrics
// Metrics in this registry will be exported to lindb's native ingestion.
type registry struct {
	mu     sync.RWMutex
	buffer []*taggedSeries // store metrics in buffer to prevent long waiting during flushing
	series map[uint64]*taggedSeries
}

var defaultRegistry = &registry{
	series: make(map[uint64]*taggedSeries),
}

// Register registers a namedmetric
func (r *registry) Register(seriesID uint64, series *taggedSeries) *taggedSeries {
	r.mu.Lock()
	defer r.mu.Unlock()

	old, exist := r.series[seriesID]
	if exist {
		return old
	}
	r.series[seriesID] = series
	return series
}

// gatherMetricList transforms event-metrics to native lindb dto-proto format
func (r *registry) gatherMetricList() ([]*protoMetricsV1.Metric, int) {
	r.mu.Lock()
	r.buffer = r.buffer[:0]
	for _, nm := range r.series {
		r.buffer = append(r.buffer, nm)
	}
	r.mu.Unlock()

	var (
		ml    []*protoMetricsV1.Metric
		count int
	)
	for _, s := range r.buffer {
		gatheredMetric := s.gatherMetric()
		if gatheredMetric == nil {
			continue
		}
		count++
		ml = append(ml, gatheredMetric)
	}
	return ml, count
}

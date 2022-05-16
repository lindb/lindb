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
	"io"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/metric"
)

var (
	// BrokerRegistry represents broker level metric Registry.
	BrokerRegistry = &Registry{
		series: make(map[uint64]*taggedSeries),
	}
	// StorageRegistry represents storage level metric Registry.
	StorageRegistry = &Registry{
		series: make(map[uint64]*taggedSeries),
	}
)

// Registry is a set of metrics
// Metrics in this Registry will be exported to LinDB's native ingestion.
type Registry struct {
	mu     sync.RWMutex
	series map[uint64]*taggedSeries
}

func (r *Registry) NewScope(metricName string, tagList ...string) Scope {
	assertMetricName(metricName)

	m := tagList2Tags(tagList...)
	ms := newTaggedSeries(r, metricName, m)
	return ms
}

// NewGather returns a gather to gather metrics from sdk and runtime.
func (r *Registry) NewGather(options ...GatherOption) Gather {
	g := &gather{r: r}
	for _, o := range options {
		o.ApplyConfig(g)
	}
	return g
}

// FindMetricList returns metric list by given names/tags.
func (r *Registry) FindMetricList(names []string, includeTags map[string]string) map[string][]*models.StateMetric {
	nameMap := make(map[string]struct{})
	for _, name := range names {
		nameMap[name] = struct{}{}
	}
	var rs []*taggedSeries
	r.mu.RLock()
	for _, nm := range r.series {
		_, ok := nameMap[nm.metricName]
		if ok {
			rs = append(rs, nm)
		}
	}
	r.mu.RUnlock()

	result := make(map[string][]*models.StateMetric)
	for _, s := range rs {
		stateMetric := s.toStateMetric(includeTags)
		if stateMetric == nil {
			continue
		}
		metrics, ok := result[s.metricName]
		if ok {
			metrics = append(metrics, stateMetric)
			result[s.metricName] = metrics
		} else {
			result[s.metricName] = []*models.StateMetric{stateMetric}
		}
	}
	return result
}

// register registers a named metric
func (r *Registry) register(seriesID uint64, series *taggedSeries) *taggedSeries {
	r.mu.Lock()
	defer r.mu.Unlock()

	old, exist := r.series[seriesID]
	if exist {
		return old
	}
	r.series[seriesID] = series
	return series
}

// gatherMetricList transforms event-metrics to native LinDB dto-proto format
func (r *Registry) gatherMetricList(
	writer io.Writer, merger func(builder *metric.RowBuilder),
) (count int) {
	// store metrics in buffer to prevent long waiting during flushing
	var buffer []*taggedSeries
	r.mu.RLock()
	for _, nm := range r.series {
		buffer = append(buffer, nm)
	}
	r.mu.RUnlock()

	builder, releaseFunc := metric.NewRowBuilder()
	defer releaseFunc(builder)

	for _, s := range buffer {
		builder.Reset()

		if !s.buildFlatMetric(builder) {
			// no metric
			continue
		}
		merger(builder)

		data, err := builder.Build()
		if err != nil {
			continue
		}
		_, _ = writer.Write(data)
		count++
	}
	return count
}

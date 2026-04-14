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
	"github.com/lindb/arrow/pkg/model"

	"github.com/lindb/lindb/series/tag"
)

// ArrowGather gathers metrics from the registry as Arrow model.Metric slices,
// ready to be appended to a MetricBuilder for Arrow IPC serialization.
type ArrowGather interface {
	// GatherMetrics collects all registered series and returns them as []*model.Metric.
	GatherMetrics() []*model.Metric
}

// arrowGather implements ArrowGather. It holds a reference to the Registry and
// optional runtime observer / global tags / namespace injected via GatherOption.
type arrowGather struct {
	r               *Registry
	namespace       string
	runtimeObserver Observer
	tags            tag.Tags // global tags merged into every metric's Attributes
}

// GatherMetrics triggers the optional runtime observer, then converts every
// registered taggedSeries into a *model.Metric, merging in namespace and
// global tags. Nil entries (series with no payload yet) are omitted.
func (ag *arrowGather) GatherMetrics() []*model.Metric {
	if ag.runtimeObserver != nil {
		ag.runtimeObserver.Observe()
	}

	// snapshot the series map under read lock to avoid long hold during conversion
	var buffer []*taggedSeries
	ag.r.mu.RLock()
	for _, s := range ag.r.series {
		buffer = append(buffer, s)
	}
	ag.r.mu.RUnlock()

	metrics := make([]*model.Metric, 0, len(buffer))
	for _, s := range buffer {
		m := s.buildArrowMetric(ag.namespace, ag.tags)
		if m == nil {
			continue
		}
		metrics = append(metrics, m)
	}
	return metrics
}

// applyArrowGatherOption applies a GatherOption to an arrowGather by delegating
// through a temporary gather struct (for option compatibility).
func applyArrowGatherOption(ag *arrowGather, opt GatherOption) {
	// Reuse the existing gather struct as a bridge so GatherOptions remain compatible.
	bridge := &gather{
		r:         ag.r,
		namespace: ag.namespace,
		tags:      ag.tags,
	}
	opt.ApplyConfig(bridge)
	ag.namespace = bridge.namespace
	ag.tags = bridge.tags
	ag.runtimeObserver = bridge.runtimeObserver
}

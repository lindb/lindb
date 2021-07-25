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
	"sort"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

// Gather gathers native lindb dto metrics
type Gather interface {
	// Gather gathers and returns the gathered metrics
	Gather() []*protoMetricsV1.Metric
}

// NewGather returns a gather to gather metrics from sdk and runtime.
func NewGather(options ...GatherOption) Gather {
	g := &gather{}
	for _, o := range options {
		o.ApplyConfig(g)
	}
	return g
}

type GatherOption interface {
	ApplyConfig(g *gather)
}

type gather struct {
	namespace       string
	runtimeObserver *runtimeObserver
	keyValues       tag.KeyValues
}

func (g *gather) mergeTags(m *protoMetricsV1.Metric) {
	if len(g.keyValues) == 0 {
		return
	}
	var tags = make(map[string]string)
	for _, globalKV := range g.keyValues {
		tags[globalKV.Key] = globalKV.Value
	}
	for _, t := range m.Tags {
		tags[t.Key] = t.Value
	}
	newKVs := tag.KeyValuesFromMap(tags)
	sort.Sort(newKVs)
	m.Tags = newKVs
}

func (g *gather) Gather() []*protoMetricsV1.Metric {
	if g.runtimeObserver != nil {
		g.runtimeObserver.Observe()
	}

	metrics, _ := defaultRegistry.gatherMetricList()
	// enrich global tagKeyValues
	for _, m := range metrics {
		g.mergeTags(m)
		m.Namespace = g.namespace
	}

	return metrics
}

type readRuntimeOption struct{}

func (o *readRuntimeOption) ApplyConfig(g *gather) { g.runtimeObserver = newRuntimeObserver() }

func WithReadRuntimeOption() GatherOption { return &readRuntimeOption{} }

type globalKeyValuesOption struct {
	keyValues tag.KeyValues
}

func (o *globalKeyValuesOption) ApplyConfig(g *gather) {
	g.keyValues = o.keyValues.DeDup()
}

func WithGlobalKeyValueOption(kvs tag.KeyValues) GatherOption {
	return &globalKeyValuesOption{keyValues: kvs}
}

type namespaceOption struct{ namespace string }

func (o *namespaceOption) ApplyConfig(g *gather) {
	g.namespace = o.namespace
}

func WithNamespaceOption(namespace string) GatherOption {
	return &namespaceOption{namespace: namespace}
}

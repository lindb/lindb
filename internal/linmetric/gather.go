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
	"bytes"

	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

// Gather gathers native lindb dto metrics
type Gather interface {
	// Gather gathers and returns the gathered metrics
	Gather() ([]byte, int)
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
	tags            tag.Tags
	buf             bytes.Buffer
}

func (g *gather) enrichTagsNameSpace(builder *metric.RowBuilder) {
	if len(g.tags) == 0 {
		return
	}
	for _, kv := range g.tags {
		_ = builder.AddTag(kv.Key, kv.Value)
	}
	builder.AddNameSpace([]byte(g.namespace))
}

func (g *gather) Gather() ([]byte, int) {
	if g.runtimeObserver != nil {
		g.runtimeObserver.Observe()
	}

	g.buf.Reset()

	n := defaultRegistry.gatherMetricList(&g.buf, g.enrichTagsNameSpace)
	return g.buf.Bytes(), n
}

type readRuntimeOption struct{}

func (o *readRuntimeOption) ApplyConfig(g *gather) { g.runtimeObserver = newRuntimeObserver() }

func WithReadRuntimeOption() GatherOption { return &readRuntimeOption{} }

type globalKeyValuesOption struct {
	keyValues tag.Tags
}

func (o *globalKeyValuesOption) ApplyConfig(g *gather) {
	g.tags = o.keyValues
}

func WithGlobalKeyValueOption(kvs tag.Tags) GatherOption {
	return &globalKeyValuesOption{keyValues: kvs}
}

type namespaceOption struct{ namespace string }

func (o *namespaceOption) ApplyConfig(g *gather) {
	g.namespace = o.namespace
}

func WithNamespaceOption(namespace string) GatherOption {
	return &namespaceOption{namespace: namespace}
}

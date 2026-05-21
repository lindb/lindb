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
	"strings"
	"sync"
	"time"

	"github.com/lindb/lindb/series/tag"
)

type DeltaHistogramVec struct {
	r               *Registry
	tags            tag.Tags // unique tags
	tagKeys         []string
	metricName      string // concated metric name
	fieldName       string // logical histogram field name; empty = unnamed single histogram
	mu              sync.RWMutex
	deltaHistograms map[string]*BoundHistogram
	setBucketsFunc  func(h *BoundHistogram)
}

// NewNamedHistogramVec creates a DeltaHistogramVec where each BoundHistogram carries
// fieldName so the storage layer can decompose it as a named compound field
// (e.g. sent_duration_sum, sent_duration.__bucket_*) within the parent metric.
func NewNamedHistogramVec(r *Registry, metricName, fieldName string, tags tag.Tags, tagKey ...string) *DeltaHistogramVec {
	return newHistogramVec(r, metricName, fieldName, tags, tagKey...)
}

func newHistogramVec(r *Registry, metricName, fieldName string, tags tag.Tags, tagKey ...string) *DeltaHistogramVec {
	return &DeltaHistogramVec{
		r:               r,
		metricName:      metricName,
		fieldName:       fieldName,
		tags:            tags,
		tagKeys:         tagKey,
		deltaHistograms: make(map[string]*BoundHistogram),
	}
}

func (hv *DeltaHistogramVec) WithExponentBuckets(lower, upper time.Duration, count int) *DeltaHistogramVec {
	hv.mu.Lock()
	defer hv.mu.Unlock()

	hv.setBucketsFunc = func(h *BoundHistogram) {
		h.WithExponentBuckets(lower, upper, count)
	}
	return hv
}

func (hv *DeltaHistogramVec) WithLinearBuckets(lower, upper time.Duration, count int) *DeltaHistogramVec {
	hv.mu.Lock()
	defer hv.mu.Unlock()

	hv.setBucketsFunc = func(h *BoundHistogram) {
		h.WithLinearBuckets(lower, upper, count)
	}
	return hv
}

func (hv *DeltaHistogramVec) WithTagValues(tagValues ...string) *BoundHistogram {
	if len(tagValues) != len(hv.tagKeys) {
		panic("count of tagKey and tagValue not match")
	}
	id := strings.Join(tagValues, ",")
	hv.mu.RLock()
	h, ok := hv.deltaHistograms[id]
	hv.mu.RUnlock()
	if ok {
		return h
	}

	hv.mu.Lock()
	defer hv.mu.Unlock()

	h, ok = hv.deltaHistograms[id]
	if ok {
		return h
	}
	var tagsMap = hv.tags.Map()
	for i := range hv.tagKeys {
		tagsMap[hv.tagKeys[i]] = tagValues[i]
	}
	series := newTaggedSeries(hv.r, hv.metricName, tag.TagsFromMap(tagsMap))
	h = series.newNamedHistogram(hv.fieldName)
	if hv.setBucketsFunc != nil {
		hv.setBucketsFunc(h)
	}
	hv.deltaHistograms[id] = h
	return h
}

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
	tags            tag.KeyValues // unique tags
	tagKeys         []string
	metricName      string // concated metric name
	mu              sync.RWMutex
	deltaHistograms map[string]*BoundDeltaHistogram
	setBucketsFunc  func(h *BoundDeltaHistogram)
}

func NewHistogramVec(metricName string, tags tag.KeyValues, tagKey ...string) *DeltaHistogramVec {
	return &DeltaHistogramVec{
		metricName:      metricName,
		tags:            tags,
		tagKeys:         tagKey,
		deltaHistograms: make(map[string]*BoundDeltaHistogram),
	}
}

func (hv *DeltaHistogramVec) WithExponentBuckets(lower, upper time.Duration, count int) *DeltaHistogramVec {
	hv.mu.Lock()
	defer hv.mu.Unlock()

	hv.setBucketsFunc = func(h *BoundDeltaHistogram) {
		h.WithExponentBuckets(lower, upper, count)
	}
	return hv
}

func (hv *DeltaHistogramVec) WithLinearBuckets(lower, upper time.Duration, count int) *DeltaHistogramVec {
	hv.mu.Lock()
	defer hv.mu.Unlock()

	hv.setBucketsFunc = func(h *BoundDeltaHistogram) {
		h.WithLinearBuckets(lower, upper, count)
	}
	return hv
}

func (hv *DeltaHistogramVec) WithTagValues(tagValues ...string) *BoundDeltaHistogram {
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
	series := newTaggedSeries(hv.metricName, tag.KeyValuesFromMap(tagsMap))
	h = series.NewHistogram()
	if hv.setBucketsFunc != nil {
		hv.setBucketsFunc(h)
	}
	hv.deltaHistograms[id] = h
	return h
}

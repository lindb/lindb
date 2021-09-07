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

package linmetric //nolint:dupl

import (
	"strings"
	"sync"

	"github.com/lindb/lindb/series/tag"
)

type DeltaCounterVec struct {
	tags          tag.KeyValues // unique tags
	tagKeys       []string
	metricName    string // concated metric name
	fieldName     string
	mu            sync.RWMutex
	deltaCounters map[string]*BoundDeltaCounter
}

func NewCounterVec(metricName string, fieldName string, tags tag.KeyValues, tagKey ...string) *DeltaCounterVec {
	return &DeltaCounterVec{
		metricName:    metricName,
		fieldName:     fieldName,
		tags:          tags,
		tagKeys:       tagKey,
		deltaCounters: make(map[string]*BoundDeltaCounter),
	}
}

func (dcv *DeltaCounterVec) WithTagValues(tagValues ...string) *BoundDeltaCounter {
	if len(tagValues) != len(dcv.tagKeys) {
		panic("count of tagKey and tagValue not match")
	}
	id := strings.Join(tagValues, ",")
	dcv.mu.RLock()
	c, ok := dcv.deltaCounters[id]
	dcv.mu.RUnlock()
	if ok {
		return c
	}

	dcv.mu.Lock()
	defer dcv.mu.Unlock()

	c, ok = dcv.deltaCounters[id]
	if ok {
		return c
	}
	var tagsMap = dcv.tags.Map()
	for i := range dcv.tagKeys {
		tagsMap[dcv.tagKeys[i]] = tagValues[i]
	}
	series := newTaggedSeries(dcv.metricName, tag.KeyValuesFromMap(tagsMap))
	c = series.NewCounter(dcv.fieldName)

	dcv.deltaCounters[id] = c
	return c
}

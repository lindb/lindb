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

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
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

func newDeltaCounterVec(metricName string, fieldName string, tags tag.KeyValues, tagKey ...string) *DeltaCounterVec {
	if len(tagKey) == 0 {
		panic("tagKey length is zero")
	}
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
	var tags = dcv.tags.Clone()
	for i := range dcv.tagKeys {
		tags = append(tags, &protoMetricsV1.KeyValue{
			Key:   dcv.tagKeys[i],
			Value: tagValues[i],
		})
	}
	series := newTaggedSeries(dcv.metricName, tags)
	c = series.NewDeltaCounter(dcv.fieldName)

	dcv.deltaCounters[id] = c
	return c
}

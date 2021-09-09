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
	"fmt"
	"math"
	"sync"

	"github.com/cespare/xxhash/v2"

	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

// Scope is a namespace wrapper for linmetric.
// ensure that all emitted metric have a given prefix and tags.
// NewsScope("lindb").Scope("runtime").Scope("mem") make a point-concated metric-name: lindb.runtime.mem
type Scope interface {
	// Scope returns a child scope
	Scope(name string, tagList ...string) Scope
	// NewGauge returns a fast gauge which bounded to the scope
	NewGauge(fieldName string) *BoundGauge
	// NewCounter returns a fast counter which bounded to the scope
	NewCounter(fieldName string) *BoundCounter
	// NewMax returns a fast max which bounded to the scope
	NewMax(fieldName string) *BoundMax
	// NewMin returns a fast min which bounded to the scope
	NewMin(fieldName string) *BoundMin
	// NewHistogram returns a histogram which bounded to the scope
	NewHistogram() *BoundHistogram
	// NewHistogramVec initializes a vec by tagKeys
	NewHistogramVec(tagKey ...string) *DeltaHistogramVec
	// NewCounterVec initializes a vec by tagKeys and fieldName
	NewCounterVec(fieldName string, tagKey ...string) *DeltaCounterVec
	// NewGaugeVec initializes a vec by tagKeys and fieldName
	NewGaugeVec(fieldName string, tagKey ...string) *GaugeVec
	// NewMaxVec initializes a vec by tagKeys and fieldName
	NewMaxVec(fieldName string, tagKey ...string) *MaxVec
	// NewMinVec initializes a vec by tagKeys and fieldName
	NewMinVec(fieldName string, tagKey ...string) *MinVec
}

type taggedSeries struct {
	mu         sync.Mutex    // lock for modifying fields
	tagsID     uint64        // metric-name + tags
	metricName string        // concated metric name
	tags       tag.KeyValues // unique tags
	payload    *fieldPayload
}

type fieldPayload struct {
	gauges         []*BoundGauge   // BoundGauge list
	countersDelta  []*BoundCounter // BoundCounter list
	maxes          []*BoundMax     // BoundMax list
	mines          []*BoundMin     // BoundMin list
	histogramDelta *BoundHistogram
}

func NewScope(metricName string, tagList ...string) Scope {
	assertMetricName(metricName)

	m := tagList2Map(tagList...)
	ms := newTaggedSeries(metricName, tag.KeyValuesFromMap(m))
	return ms
}

func newTaggedSeries(metricName string, tags tag.KeyValues) *taggedSeries {
	ts := &taggedSeries{
		metricName: metricName,
		tags:       tags,
	}
	ts.tagsID = xxhash.Sum64String(ts.metricName + tag.ConcatKeyValues(ts.tags))
	// registered or replaced
	ts = defaultRegistry.Register(ts.tagsID, ts)
	return ts
}

func (s *taggedSeries) ensurePayload() {
	if s.payload == nil {
		s.payload = &fieldPayload{}
	}
}

func (s *taggedSeries) containsFieldName(fieldName string) bool {
	for _, g := range s.payload.gauges {
		if g.fieldName == fieldName {
			return true
		}
	}
	for _, dc := range s.payload.countersDelta {
		if dc.fieldName == fieldName {
			return true
		}
	}
	for _, dc := range s.payload.maxes {
		if dc.fieldName == fieldName {
			return true
		}
	}
	for _, dc := range s.payload.mines {
		if dc.fieldName == fieldName {
			return true
		}
	}
	return false
}

func assertMetricName(metricName string) {
	if len(metricName) == 0 {
		panic("metric-name cannot be empty string")
	}
}

func assertTagKeyList(tagKeyList ...string) {
	if len(tagKeyList) == 0 {
		panic("tag-key list cannot be empty")
	}
}

func assertFieldName(fieldName string) {
	if len(fieldName) == 0 {
		panic("field-name cannot be empty")
	}
}

func nextScopeKeyValues(oldTags tag.KeyValues, newTagList ...string) tag.KeyValues {
	if len(newTagList) == 0 {
		return oldTags.Clone()
	}
	nextScopeTags := tagList2Map(newTagList...)
	// newer add tags are higher priority
	for _, oldTag := range oldTags {
		if _, exist := nextScopeTags[oldTag.Key]; !exist {
			nextScopeTags[oldTag.Key] = nextScopeTags[oldTag.Value]
		}
	}
	return tag.KeyValuesFromMap(nextScopeTags)
}

func (s *taggedSeries) Scope(metricName string, tagList ...string) Scope {
	assertMetricName(metricName)

	nextMetricName := s.metricName + "." + metricName
	return newTaggedSeries(nextMetricName, nextScopeKeyValues(s.tags, tagList...))
}

func tagList2Map(tagList ...string) map[string]string {
	if len(tagList)%2 != 0 {
		panic("bad tags length ")
	}

	var m = make(map[string]string)
	for i := 0; i < len(tagList); i += 2 {
		m[tagList[i]] = tagList[i+1]
	}
	return m
}

func (s *taggedSeries) NewGauge(fieldName string) *BoundGauge {
	assertFieldName(fieldName)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	if !s.containsFieldName(fieldName) {
		bg := newGauge(fieldName)
		s.payload.gauges = append(s.payload.gauges, bg)
		return bg
	}
	for _, g := range s.payload.gauges {
		if g.fieldName == fieldName {
			return g
		}
	}
	panic(fmt.Sprintf("gauge field: %s has registered another type before", fieldName))
}

func (s *taggedSeries) NewCounter(fieldName string) *BoundCounter {
	assertFieldName(fieldName)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	if !s.containsFieldName(fieldName) {
		dc := newCounter(fieldName)
		s.payload.countersDelta = append(s.payload.countersDelta, dc)
		return dc
	}
	for _, dc := range s.payload.countersDelta {
		if dc.fieldName == fieldName {
			return dc
		}
	}
	panic(fmt.Sprintf("delta-counter field: %s has registered another type before", fieldName))
}
func (s *taggedSeries) NewMax(fieldName string) *BoundMax {
	assertFieldName(fieldName)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	if !s.containsFieldName(fieldName) {
		m := newMax(fieldName)
		s.payload.maxes = append(s.payload.maxes, m)
		return m
	}
	for _, m := range s.payload.maxes {
		if m.fieldName == fieldName {
			return m
		}
	}
	panic(fmt.Sprintf("max field: %s has registered another type before", fieldName))
}

func (s *taggedSeries) NewMin(fieldName string) *BoundMin {
	assertFieldName(fieldName)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	if !s.containsFieldName(fieldName) {
		m := newMin(fieldName)
		s.payload.mines = append(s.payload.mines, m)
		return m
	}
	for _, m := range s.payload.mines {
		if m.fieldName == fieldName {
			return m
		}
	}
	panic(fmt.Sprintf("min field: %s has registered another type before", fieldName))
}

func (s *taggedSeries) NewHistogram() *BoundHistogram {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	if s.payload.histogramDelta != nil {
		return s.payload.histogramDelta
	}
	s.payload.histogramDelta = NewHistogram()
	return s.payload.histogramDelta
}

func (s *taggedSeries) NewHistogramVec(tagKey ...string) *DeltaHistogramVec {
	assertTagKeyList(tagKey...)
	return NewHistogramVec(s.metricName, s.tags, tagKey...)
}

func (s *taggedSeries) NewCounterVec(fieldName string, tagKey ...string) *DeltaCounterVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return NewCounterVec(s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewGaugeVec(fieldName string, tagKey ...string) *GaugeVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return newGaugeVec(s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewMaxVec(fieldName string, tagKey ...string) *MaxVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return newMaxVec(s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewMinVec(fieldName string, tagKey ...string) *MinVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return newMinVec(s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) gatherMetric() *protoMetricsV1.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.payload == nil {
		return nil
	}
	var m = protoMetricsV1.Metric{
		Name:      s.metricName,
		Timestamp: timeutil.Now(),
		Tags:      s.tags,
		TagsHash:  s.tagsID,
	}
	// pick gauges
	for _, g := range s.payload.gauges {
		m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
			Name:  g.fieldName,
			Type:  protoMetricsV1.SimpleFieldType_GAUGE,
			Value: g.Get(),
		})
	}
	// pick delta counter
	for _, dc := range s.payload.countersDelta {
		m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
			Name:  dc.fieldName,
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			Value: dc.getAndReset(),
		})
	}
	// pick max
	for _, mf := range s.payload.maxes {
		v := mf.Get()
		if math.IsInf(v, 1) || math.IsInf(v, -1) {
			continue
		}
		m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
			Name:  mf.fieldName,
			Type:  protoMetricsV1.SimpleFieldType_Max,
			Value: v,
		})
	}
	// pick min
	for _, mf := range s.payload.mines {
		v := mf.Get()
		if math.IsInf(v, 1) || math.IsInf(v, -1) {
			continue
		}
		m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
			Name:  mf.fieldName,
			Type:  protoMetricsV1.SimpleFieldType_Min,
			Value: v,
		})
	}

	if s.payload.histogramDelta != nil {
		m.CompoundField = s.payload.histogramDelta.marshalToCompoundField()
	}
	return &m
}

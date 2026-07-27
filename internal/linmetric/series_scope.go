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
	"sync"

	xxhash "github.com/cespare/xxhash/v2"
	"github.com/lindb/arrow/pkg/model"
	"github.com/lindb/common/field"
	commonMetric "github.com/lindb/common/metric"
	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/strutil"
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
	// NewHistogramVec returns a DeltaHistogramVec for the given named histogram field.
	// fieldName must be non-empty; tagKey lists the dynamic tag dimensions.
	NewHistogramVec(fieldName string, tagKey ...string) *DeltaHistogramVec
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
	r          *Registry
	mu         sync.Mutex // lock for modifying fields
	seriesID   uint64     // metric-name + tags
	metricName string     // concated metric name
	tags       tag.Tags   // unique tags
	payload    *fieldPayload
}

type fieldPayload struct {
	simpleFields []simpleField     // Bound SimpleField list
	histograms   []*BoundHistogram // named or unnamed histograms; name="" means unnamed single histogram
}

func newTaggedSeries(r *Registry, metricName string, tags tag.Tags) *taggedSeries {
	ts := &taggedSeries{
		r:          r,
		metricName: metricName,
		tags:       tags,
	}
	ts.seriesID = xxhash.Sum64String(ts.metricName + string(ts.tags.AppendHashKey(nil)))
	// registered or replaced
	ts = r.register(ts.seriesID, ts)
	return ts
}

func (s *taggedSeries) ensurePayload() {
	if s.payload == nil {
		s.payload = &fieldPayload{}
	}
}

func assertMetricName(metricName string) {
	if metricName == "" {
		panic("metric-name cannot be empty string")
	}
}

func assertTagKeyList(tagKeyList ...string) {
	if len(tagKeyList) == 0 {
		panic("tag-key list cannot be empty")
	}
}

func assertFieldName(fieldName string) {
	if fieldName == "" {
		panic("field-name cannot be empty")
	}
}

func nextScopeKeyValues(oldTags tag.Tags, newTagList ...string) tag.Tags {
	if len(newTagList) == 0 {
		return oldTags.Clone()
	}
	if len(newTagList)%2 != 0 {
		panic("bad tags length ")
	}
	m := oldTags.Map()
	for i := 0; i < len(newTagList); i += 2 {
		m[newTagList[i]] = newTagList[i+1]
	}
	return tag.TagsFromMap(m)
}

func (s *taggedSeries) Scope(metricName string, tagList ...string) Scope {
	assertMetricName(metricName)

	nextMetricName := s.metricName + "." + metricName
	return newTaggedSeries(s.r, nextMetricName, nextScopeKeyValues(s.tags, tagList...))
}

func tagList2Tags(tagList ...string) tag.Tags {
	if len(tagList)%2 != 0 {
		panic("bad tags length ")
	}

	var ts tag.Tags
	for i := 0; i < len(tagList); i += 2 {
		ts = append(ts, tag.Tag{
			Key:   []byte(tagList[i]),
			Value: []byte(tagList[i+1]),
		})
	}
	return ts
}

func (s *taggedSeries) NewGauge(fieldName string) *BoundGauge {
	return s.findSimpleField(fieldName, field.Last, func() simpleField {
		return newGauge(fieldName)
	}).(*BoundGauge)
}

func (s *taggedSeries) NewCounter(fieldName string) *BoundCounter {
	return s.findSimpleField(fieldName, field.Sum, func() simpleField {
		return newCounter(fieldName)
	}).(*BoundCounter)
}

func (s *taggedSeries) NewMax(fieldName string) *BoundMax {
	return s.findSimpleField(fieldName, field.Max, func() simpleField {
		return newMax(fieldName)
	}).(*BoundMax)
}

func (s *taggedSeries) NewMin(fieldName string) *BoundMin {
	return s.findSimpleField(fieldName, field.Min, func() simpleField {
		return newMin(fieldName)
	}).(*BoundMin)
}

func (s *taggedSeries) findSimpleField(
	fieldName string,
	fieldType field.Type,
	createFunc func() simpleField,
) simpleField {
	assertFieldName(fieldName)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	for _, sf := range s.payload.simpleFields {
		if sf.name() == fieldName {
			if sf.fieldType() != fieldType {
				panic(fmt.Sprintf("field: %s has registered another type before", fieldName))
			}
			return sf
		}
	}
	sf := createFunc()
	s.payload.simpleFields = append(s.payload.simpleFields, sf)
	return sf
}

// newNamedHistogram returns a BoundHistogram with the given fieldName.
// If a histogram with the same name already exists it is returned unchanged.
func (s *taggedSeries) newNamedHistogram(fieldName string) *BoundHistogram {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensurePayload()
	for _, h := range s.payload.histograms {
		if h.name == fieldName {
			return h
		}
	}
	h := newBoundHistogram(fieldName)
	s.payload.histograms = append(s.payload.histograms, h)
	return h
}

// NewHistogramVec returns a DeltaHistogramVec for the given named histogram field.
// fieldName must be non-empty; tagKey lists the dynamic tag dimensions.
func (s *taggedSeries) NewHistogramVec(fieldName string, tagKey ...string) *DeltaHistogramVec {
	assertFieldName(fieldName)
	return NewNamedHistogramVec(s.r, s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewCounterVec(fieldName string, tagKey ...string) *DeltaCounterVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return NewCounterVec(s.r, s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewGaugeVec(fieldName string, tagKey ...string) *GaugeVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return newGaugeVec(s.r, s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewMaxVec(fieldName string, tagKey ...string) *MaxVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return newMaxVec(s.r, s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) NewMinVec(fieldName string, tagKey ...string) *MinVec {
	assertFieldName(fieldName)
	assertTagKeyList(tagKey...)
	return newMinVec(s.r, s.metricName, fieldName, s.tags, tagKey...)
}

func (s *taggedSeries) buildFlatMetric(builder *commonMetric.RowBuilder) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.payload == nil {
		return false
	}

	builder.AddMetricName(strutil.String2ByteSlice(s.metricName))
	builder.AddTimestamp(fasttime.UnixMilliseconds())
	for _, kv := range s.tags {
		_ = builder.AddTag(kv.Key, kv.Value)
	}

	// pick simple fields
	for _, sf := range s.payload.simpleFields {
		_ = builder.AddSimpleField(
			strutil.String2ByteSlice(sf.name()),
			toFlatSimpleFieldType(sf.fieldType()),
			sf.gather(),
		)
	}

	// FlatBuffer supports only one CompoundField per metric.
	// Emit the first unnamed histogram (backward compat); named histograms go via Arrow path.
	for _, h := range s.payload.histograms {
		if h.name == "" {
			h.marshalToCompoundField(builder)
			break
		}
	}
	return true
}

// buildArrowMetric converts the taggedSeries into a *model.Metric for Arrow IPC serialization.
// namespace and globalTags are merged into the resulting Attributes.
// Histogram fields are decomposed into atomic stat and bucket fields via AppendFields.
func (s *taggedSeries) buildArrowMetric(namespace string, globalTags tag.Tags) *model.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.payload == nil {
		return nil
	}

	m := &model.Metric{
		Namespace: namespace,
		Name:      s.metricName,
		Timestamp: fasttime.UnixNano(),
	}

	// merge series tags and global tags into Attributes
	totalLen := len(s.tags) + len(globalTags)
	if totalLen > 0 {
		m.Attributes = &model.Attributes{}
		for _, kv := range s.tags {
			m.Attributes.Append(string(kv.Key), string(kv.Value))
		}
		for _, kv := range globalTags {
			m.Attributes.Append(string(kv.Key), string(kv.Value))
		}
	}

	// map simple fields to Arrow model.Field, skipping unknown types
	for _, sf := range s.payload.simpleFields {
		m.Fields = append(m.Fields, model.Field{
			Name:  sf.name(),
			Kind:  sf.fieldType(),
			Value: sf.gather(),
		})
	}

	// append histogram fields as atomic named fields (e.g. <name>.__sum, <name>.__bucket_*)
	for _, h := range s.payload.histograms {
		h.AppendFields(m)
	}

	return m
}

func (s *taggedSeries) toStateMetric(includeTags map[string]string) *models.StateMetric {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.payload == nil {
		return nil
	}
	tags := s.tags.Map()
	if !isMapSubset(tags, includeTags) {
		return nil
	}
	rs := &models.StateMetric{}
	rs.Tags = tags

	// convert simple fields
	for _, sf := range s.payload.simpleFields {
		rs.Fields = append(rs.Fields, models.StateField{
			Name:  sf.name(),
			Type:  sf.fieldType().String(),
			Value: sf.Get(),
		})
	}

	return rs
}

// isMapSubset checks sub map if include given map.
func isMapSubset(m, sub map[string]string) bool {
	if len(m) < len(sub) {
		return false
	}
	for k, v := range sub {
		if vm, found := m[k]; !found || vm != v {
			return false
		}
	}
	return true
}

// toFlatSimpleFieldType maps a field.Type to the corresponding FlatBuffer SimpleFieldType
// used by the legacy flat-metrics wire format (buildFlatMetric path).
func toFlatSimpleFieldType(t field.Type) flatMetricsV1.SimpleFieldType {
	switch t {
	case field.Last:
		return flatMetricsV1.SimpleFieldTypeLast
	case field.Sum:
		return flatMetricsV1.SimpleFieldTypeDeltaSum
	case field.Max:
		return flatMetricsV1.SimpleFieldTypeMax
	case field.Min:
		return flatMetricsV1.SimpleFieldTypeMin
	case field.First:
		return flatMetricsV1.SimpleFieldTypeFirst
	default:
		return flatMetricsV1.SimpleFieldTypeDeltaSum
	}
}

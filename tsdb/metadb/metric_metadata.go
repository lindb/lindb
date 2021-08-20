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

package metadb

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./metric_metadata.go -destination=./metric_metadata_mock.go -package=metadb

// MetricMetadata represents the metric metadata for memory cache,
// it will cache all tag keys and fields in backend db
type MetricMetadata interface {
	// initialize initializes the metric metadata with tags/fields
	initialize(fields []field.Meta, tagKeys []tag.Meta)

	// getMetricID gets the metric id
	getMetricID() uint32
	// getField gets the field meta by field name, if not exist return false
	getField(fieldName field.Name) (field.Meta, bool)
	// getAllFields returns the all fields of the metric
	getAllFields() (fields []field.Meta)
	// getAllHistogramFields returns histogram bucket fields
	getAllHistogramFields() (fields field.Metas)
	// getTagKeyID gets the tag key id by tag key, if not exist return false
	getTagKeyID(tagKey string) (uint32, bool)
	// getAllTags returns the tag keys of the metric
	getAllTagKeys() (tagKeys []tag.Meta)

	// checkTagKeyCount checks the tag keys if limit, if limit return series.ErrTooManyTagKeys
	checkTagKeyCount() error

	// createField creates the field meta, if success return field id, else return series.ErrTooManyFields
	createField(fieldName field.Name, fieldType field.Type) (field.ID, error)
	// rollbackFieldID rollbacks field id
	rollbackFieldID(fieldID field.ID)
	// addField adds field meta
	addField(f field.Meta)
	// createTagKey creates the tag key
	createTagKey(tagKey string, tagKeyID uint32)
}

// metricMetadata implements MetricMetadata interface
type metricMetadata struct {
	metricID   uint32
	fieldIDSeq atomic.Int32
	fields     []field.Meta
	tagKeys    []tag.Meta
}

// newMetricMetadata creates the metric metadata with metric id and field id assign sequence
func newMetricMetadata(metricID uint32, fieldIDSeq int32) MetricMetadata {
	mm := &metricMetadata{
		metricID: metricID,
	}
	mm.fieldIDSeq.Store(fieldIDSeq)
	return mm
}

// initialize initializes the metric metadata with tags/fields
func (mm *metricMetadata) initialize(fields []field.Meta, tagKeys []tag.Meta) {
	mm.fields = fields
	mm.tagKeys = tagKeys
}

// getMetricID gets the metric id
func (mm *metricMetadata) getMetricID() uint32 {
	return mm.metricID
}

// getField gets the field meta by field name, if not exist return false
func (mm *metricMetadata) getField(fieldName field.Name) (field.Meta, bool) {
	for _, f := range mm.fields {
		if f.Name == fieldName {
			return f, true
		}
	}
	return field.Meta{}, false
}

// getAllFields returns the all fields of the metric
func (mm *metricMetadata) getAllFields() (fields []field.Meta) {
	return mm.fields
}

// getAllHistogramFields returns histogram buckets fields,
// with format like __bucket_${boundary}
func (mm *metricMetadata) getAllHistogramFields() (fields field.Metas) {
	for idx := range mm.fields {
		if mm.fields[idx].Type == field.HistogramField {
			fields = append(fields, mm.fields[idx])
		}
	}
	return
}

// getTagKeyID gets the tag key id by tag key, if not exist return false
func (mm *metricMetadata) getTagKeyID(tagKey string) (uint32, bool) {
	for _, t := range mm.tagKeys {
		if t.Key == tagKey {
			return t.ID, true
		}
	}
	return 0, false
}

// getAllTags returns the tag keys of the metric
func (mm *metricMetadata) getAllTagKeys() (tagKeys []tag.Meta) {
	return mm.tagKeys
}

// createField creates the field meta, if success return field id, else return series.ErrTooManyFields
func (mm *metricMetadata) createField(fieldName field.Name, fieldType field.Type) (field.ID, error) {
	// check fields count
	if mm.fieldIDSeq.Load() >= constants.DefaultMaxFieldsCount {
		return 0, series.ErrTooManyFields
	}
	// create new field
	fieldID := field.ID(mm.fieldIDSeq.Inc())
	return fieldID, nil
}

// rollbackFieldID rollbacks field id
func (mm *metricMetadata) rollbackFieldID(fID field.ID) {
	if mm.fieldIDSeq.Load() == int32(fID) {
		mm.fieldIDSeq.Dec()
	}
}

// addField adds field meta
func (mm *metricMetadata) addField(f field.Meta) {
	mm.fields = append(mm.fields, f)
}

// checkTagKeyCount checks the tag keys if limit, if limit return series.ErrTooManyTagKeys
func (mm *metricMetadata) checkTagKeyCount() error {
	// check tag keys count
	if len(mm.tagKeys) >= constants.DefaultMaxTagKeysCount {
		return series.ErrTooManyTagKeys
	}
	return nil
}

// createTagKey creates the tag key
func (mm *metricMetadata) createTagKey(tagKey string, tagKeyID uint32) {
	mm.tagKeys = append(mm.tagKeys, tag.Meta{ID: tagKeyID, Key: tagKey})
}

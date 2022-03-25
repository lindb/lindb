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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./metric_metadata.go -destination=./metric_metadata_mock.go -package=metadb

// MetricMetadata represents the metric metadata for memory cache,
// it will cache all tag keys and fields in backend db
type MetricMetadata interface {
	// initialize the metric metadata with tags/fields
	initialize(fields field.Metas, fieldMaxID int32, tagKeys tag.Metas)

	// getMetricID gets the metric id
	getMetricID() metric.ID

	// createField creates the field meta, if success return field id, else return series.ErrTooManyFields
	createField(fieldName field.Name, fieldType field.Type) (field.Meta, error)
	// getField gets the field meta by field name, if not exist return false
	getField(fieldName field.Name) (field.Meta, bool)
	// getAllFields returns the all fields of the metric
	getAllFields() (fields field.Metas)

	// createTagKey creates the tag key
	createTagKey(tagKey string, tagKeyID tag.KeyID)
	checkTagKey(tagKey string) error
	// getTagKeyID gets the tag key id by tag key, if not exist return false
	getTagKeyID(tagKey string) (tag.KeyID, bool)
	// getAllTags returns the tag keys of the metric
	getAllTagKeys() (tagKeys tag.Metas)
}

// metricMetadata implements MetricMetadata interface
type metricMetadata struct {
	metricID metric.ID
	fields   field.Metas
	tagKeys  tag.Metas

	fieldIDSeq atomic.Int32
}

// newMetricMetadata creates the metric metadata with metric id and field id assign sequence
func newMetricMetadata(metricID metric.ID) MetricMetadata {
	mm := &metricMetadata{
		metricID: metricID,
	}
	return mm
}

// initialize the metric metadata with tags/fields
func (mm *metricMetadata) initialize(fields field.Metas, fieldMaxID int32, tagKeys tag.Metas) {
	mm.fields = fields
	mm.tagKeys = tagKeys
	mm.fieldIDSeq.Store(fieldMaxID)
}

// getMetricID gets the metric id
func (mm *metricMetadata) getMetricID() metric.ID {
	return mm.metricID
}

// createField creates the field meta, if success return field id, else return series.ErrTooManyFields
func (mm *metricMetadata) createField(fieldName field.Name, fieldType field.Type) (field.Meta, error) {
	// check fields count
	// TODO add config????
	if mm.fieldIDSeq.Load() >= constants.DefaultMaxFieldsCount {
		return field.Meta{}, series.ErrTooManyFields
	}
	// create new field
	fieldID := field.ID(mm.fieldIDSeq.Inc())
	fieldMeta := field.Meta{
		ID:   fieldID,
		Name: fieldName,
		Type: fieldType,
	}
	mm.fields = append(mm.fields, fieldMeta)
	return fieldMeta, nil
}

// getField gets the field meta by field name, if not exist return false
func (mm *metricMetadata) getField(fieldName field.Name) (field.Meta, bool) {
	return mm.fields.Find(fieldName)
}

// getAllFields returns the all fields of the metric
func (mm *metricMetadata) getAllFields() (fields field.Metas) {
	length := len(mm.fields)
	if length == 0 {
		return
	}
	// need copy result
	fields = make(field.Metas, length)
	copy(fields, mm.fields)
	return
}

// createTagKey creates the tag key
// 1. checks the tag keys if limited, if limit return series.ErrTooManyTagKeys
func (mm *metricMetadata) createTagKey(tagKey string, tagKeyID tag.KeyID) {
	mm.tagKeys = append(mm.tagKeys, tag.Meta{Key: tagKey, ID: tagKeyID})
}

func (mm *metricMetadata) checkTagKey(_ string) error {
	// check tag keys count
	// TODO add config
	if len(mm.tagKeys) >= config.GlobalStorageConfig().TSDB.MaxTagKeysNumber {
		return series.ErrTooManyTagKeys
	}
	return nil
}

// getTagKeyID gets the tag key id by tag key, if not exist return false
func (mm *metricMetadata) getTagKeyID(tagKey string) (tag.KeyID, bool) {
	t, ok := mm.tagKeys.Find(tagKey)
	if ok {
		return t.ID, true
	}
	return tag.EmptyTagKeyID, false
}

// getAllTags returns the tag keys of the metric
func (mm *metricMetadata) getAllTagKeys() (tagKeys tag.Metas) {
	length := len(mm.tagKeys)
	if length == 0 {
		return
	}
	// need copy result
	tagKeys = make(tag.Metas, length)
	copy(tagKeys, mm.tagKeys)
	return
}

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestMetricMetadata_createField(t *testing.T) {
	mm := newMetricMetadata(1, 0)
	// test: create field id
	fieldID, err := mm.createField("f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(1), fieldID)

	// test: re-open metric metadata
	mm = newMetricMetadata(1, 1)
	fieldID, err = mm.createField("f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(2), fieldID)
	mm.addField(field.Meta{
		ID:   fieldID,
		Type: field.SumField,
		Name: "f",
	})

	// test: too many fields
	mm = newMetricMetadata(1, 0)
	for i := 1; i <= constants.DefaultMaxFieldsCount; i++ {
		fieldID, err = mm.createField(field.Name(fmt.Sprintf("f-%d", i)), field.SumField)
		assert.NoError(t, err)
		assert.Equal(t, field.ID(i), fieldID)
		mm.addField(field.Meta{
			ID:   fieldID,
			Type: field.SumField,
			Name: field.Name(fmt.Sprintf("f-%d", i)),
		})
	}
	fieldID, err = mm.createField("max-f", field.SumField)
	assert.Equal(t, series.ErrTooManyFields, err)
	assert.Equal(t, field.ID(0), fieldID)

	assert.Len(t, mm.getAllFields(), constants.DefaultMaxFieldsCount)

	for i := 1; i <= constants.DefaultMaxFieldsCount; i++ {
		f, ok := mm.getField(field.Name(fmt.Sprintf("f-%d", i)))
		assert.True(t, ok)
		assert.Equal(t, field.ID(i), f.ID)
	}
	_, ok := mm.getField("no-f")
	assert.False(t, ok)

	mm2 := newMetricMetadata(1, 0)
	mm2.initialize(mm.getAllFields(), mm.getAllTagKeys())
	assert.Len(t, mm2.getAllFields(), constants.DefaultMaxFieldsCount)

	_, ok = mm2.getField("max-f")
	assert.False(t, ok)
}

func TestMetricMetadata_createTag(t *testing.T) {
	mm := newMetricMetadata(1, 0)
	assert.Equal(t, uint32(1), mm.getMetricID())
	for i := 1; i <= constants.DefaultMaxTagKeysCount; i++ {
		err := mm.checkTagKeyCount()
		assert.NoError(t, err)
		mm.createTagKey(fmt.Sprintf("tag-%d", i), uint32(i))
	}

	err := mm.checkTagKeyCount()
	assert.Equal(t, series.ErrTooManyTagKeys, err)

	for i := 1; i <= constants.DefaultMaxTagKeysCount; i++ {
		tagKeyID, ok := mm.getTagKeyID(fmt.Sprintf("tag-%d", i))
		assert.True(t, ok)
		assert.Equal(t, uint32(i), tagKeyID)
	}
	assert.Len(t, mm.getAllTagKeys(), constants.DefaultMaxTagKeysCount)
	tagKeyID, ok := mm.getTagKeyID("no-tag")
	assert.False(t, ok)
	assert.Equal(t, uint32(0), tagKeyID)

	mm2 := newMetricMetadata(1, 0)
	mm2.initialize(mm.getAllFields(), mm.getAllTagKeys())
	assert.Len(t, mm2.getAllTagKeys(), constants.DefaultMaxTagKeysCount)
	err = mm.checkTagKeyCount()
	assert.Equal(t, series.ErrTooManyTagKeys, err)
}

func TestMetricMetadata_rollback(t *testing.T) {
	mm := newMetricMetadata(1, 0)
	// test: create field id
	fieldID, err := mm.createField("f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(1), fieldID)
	mm.rollbackFieldID(fieldID)
	fieldID, err = mm.createField("f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(1), fieldID)
	mm.rollbackFieldID(field.ID(0))
	fieldID, err = mm.createField("f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(2), fieldID)
}

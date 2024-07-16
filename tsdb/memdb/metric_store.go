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

package memdb

import (
	"sync"

	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock.go -package memdb

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// GenField generates field meta under memory database.
	GenField(fieldName field.Name, fieldType field.Type) (f field.Meta, created bool)
	// GetFields returns all field metas.
	GetFields() field.Metas
	// UpdateFieldMeta updates field meta after metric meta updated.
	UpdateFieldMeta(fieldID field.ID, fm field.Meta)
	// FindFields returns fields from store based on current written fields.
	FindFields(fields field.Metas) (found field.Metas)
}

// metricStore represents metric level storage, stores all series data, and fields/family times metadata
type metricStore struct {
	fields sync.Map // field metadata(field.Metas)

	fieldCount int
	lock       sync.RWMutex
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	var ms metricStore
	return &ms
}

// GetFields returns all field metas.
func (ms *metricStore) GetFields() (fields field.Metas) {
	ms.fields.Range(func(key, value any) bool {
		fields = append(fields, value.(field.Meta))
		return true
	})
	return fields
}

// GenField generates field meta under memory database.
func (ms *metricStore) GenField(name field.Name, fType field.Type) (f field.Meta, created bool) {
	fm, ok := ms.fields.Load(name)
	if ok {
		return fm.(field.Meta), false
	}

	ms.lock.Lock()
	defer ms.lock.Unlock()

	return ms.genField(name, fType)
}

func (ms *metricStore) genField(name field.Name, fType field.Type) (f field.Meta, created bool) {
	fm, ok := ms.fields.Load(name)
	if ok {
		return fm.(field.Meta), false
	}

	index := uint8(ms.fieldCount)
	f = field.Meta{
		Type:  fType,
		Name:  name, // TODO: check name
		Index: index,
	}
	ms.fieldCount++
	ms.fields.Store(name, f)
	return f, true
}

// UpdateFieldMeta updates field meta after metric meta updated.
func (ms *metricStore) UpdateFieldMeta(fieldID field.ID, fm field.Meta) {
	f, ok := ms.fields.Load(fm.Name)
	if ok {
		ms.lock.Lock()
		defer ms.lock.Unlock()

		fm := f.(field.Meta)
		fm.ID = fieldID
		fm.Persisted = true
		ms.fields.Store(fm.Name, fm)
	}
}

// FindFields returns fields from store based on current written fields.
func (ms *metricStore) FindFields(fields field.Metas) (found field.Metas) {
	for _, f := range fields {
		fm, ok := ms.fields.Load(f.Name)
		if ok {
			found = append(found, fm.(field.Meta))
		}
	}
	return
}

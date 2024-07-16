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
	"sort"
	"sync"

	"go.uber.org/atomic"

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
	fields atomic.Value // field metadata(field.Metas)

	lock sync.RWMutex
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	var ms metricStore
	// init field metas
	ms.fields.Store(field.Metas{})
	return &ms
}

// GetFields returns all field metas.
func (ms *metricStore) GetFields() field.Metas {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	mFields := ms.fields.Load().(field.Metas)
	return mFields.Clone()
}

// GenField generates field meta under memory database.
func (ms *metricStore) GenField(name field.Name, fType field.Type) (f field.Meta, created bool) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	// TODO: use sync.Map?
	fields := ms.fields.Load().(field.Metas)
	fm, ok := fields.GetFromName(name)
	if ok {
		return fm, false
	}

	index := uint8(len(fields))
	fm = field.Meta{
		Type:  fType,
		Name:  name, // TODO: check name
		Index: index,
	}
	fields = append(fields, fm)
	// sort by field name
	sort.Sort(fields)
	ms.fields.Store(fields)
	return fm, true
}

// UpdateFieldMeta updates field meta after metric meta updated.
func (ms *metricStore) UpdateFieldMeta(fieldID field.ID, fm field.Meta) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	fields := ms.fields.Load().(field.Metas)

	idx, ok := fields.FindIndexByName(fm.Name)
	if ok {
		fields[idx].ID = fieldID
		fields[idx].Persisted = true
	}

	ms.fields.Store(fields)
}

// FindFields returns fields from store based on current written fields.
func (ms *metricStore) FindFields(fields field.Metas) (found field.Metas) {
	mFields := ms.fields.Load().(field.Metas)
	for _, f := range fields {
		fm, ok := mFields.Find(f.Name)
		if ok {
			found = append(found, fm)
		}
	}
	return
}

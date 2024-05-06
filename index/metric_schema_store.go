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

package index

import (
	"math"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/lindb/lindb/constants"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

// for testing
var (
	newMetricSchemaFlusher = v1.NewMetricSchemaFlusher
)

// metricSchemaStore implements MetricSchemaStore interface.
type metricSchemaStore struct {
	family    kv.Family
	mutable   *imap.IntMap[*metric.Schema] // metric id -> schema
	immutable *imap.IntMap[*metric.Schema]

	cache *expirable.LRU[metric.ID, *metric.Schema]

	lock sync.RWMutex
}

// NewMetricSchemaStore creates a MetricSchemaStore instance.
func NewMetricSchemaStore(family kv.Family) MetricSchemaStore {
	// FIXME: add cfg when create database
	cache := expirable.NewLRU[metric.ID, *metric.Schema](1000, nil, time.Minute*10)
	return &metricSchemaStore{
		family:  family,
		mutable: imap.NewIntMap[*metric.Schema](),
		cache:   cache,
	}
}

// GetSchema returns metric schema by metric id, return nil if not exist.
func (s *metricSchemaStore) GetSchema(id metric.ID) (schema *metric.Schema, err error) {
	schema = s.getSchemaFromMem(id)
	if schema != nil {
		return schema, nil
	}
	schema, ok := s.cache.Get(id)
	if ok {
		return schema, nil
	}
	schema, err = s.getSchemaFromKV(id)
	if err != nil {
		return nil, err
	}
	if schema != nil {
		s.cache.Add(id, schema)
	}
	return
}

// genFieldID generates field id if field not exist.
func (s *metricSchemaStore) genFieldID(id metric.ID, f field.Meta) (fID field.ID, err error) {
	schema, err := s.GetSchema(id)
	if err != nil {
		return 0, err
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if schema == nil {
		// create new schema
		schema = &metric.Schema{}
	}
	// put into schema if schema not exist under mutable store
	s.mutable.PutIfNotExist(uint32(id), schema)

	fm, ok := schema.Fields.Find(f.Name)
	if ok {
		return fm.ID, nil
	}

	if len(schema.Fields) >= math.MaxUint8 {
		return 0, constants.ErrTooManyFields
	}
	fID = field.ID(len(schema.Fields))
	f.ID = fID
	schema.Fields = append(schema.Fields, f)
	return fID, nil
}

// genTagKeyID generates tag key id if tag key not exist.
func (s *metricSchemaStore) genTagKeyID(id metric.ID, tagKey []byte, createFn func() uint32) (tagKeyID tag.KeyID, err error) {
	schema, err := s.GetSchema(id)
	if err != nil {
		return 0, err
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if schema == nil {
		// create new schema
		schema = &metric.Schema{}
	}
	// put into schema if schema not exist under mutable store
	s.mutable.PutIfNotExist(uint32(id), schema)

	tm, ok := schema.TagKeys.Find(strutil.ByteSlice2String(tagKey))
	if ok {
		return tm.ID, nil
	}

	if len(schema.TagKeys) >= math.MaxUint8 {
		return 0, constants.ErrTooManyTagKeys
	}
	tm = tag.Meta{
		ID:  tag.KeyID(createFn()),
		Key: string(tagKey),
	}
	schema.TagKeys = append(schema.TagKeys, tm)
	return tm.ID, nil
}

// getSchemaFromKV gets schema from kv store.
func (s *metricSchemaStore) getSchemaFromKV(id metric.ID) (schema *metric.Schema, err error) {
	snapshot := s.family.GetSnapshot()
	defer snapshot.Close()

	key := uint32(id)
	if err := snapshot.Load(key, func(value []byte) error {
		if schema == nil {
			schema = &metric.Schema{}
		}
		schema.UnmarshalFromPersist(value)
		return nil
	}); err != nil {
		return nil, err
	}
	return schema, nil
}

// getSchemaFromMem gets schema from mem store.
func (s *metricSchemaStore) getSchemaFromMem(id metric.ID) *metric.Schema {
	key := uint32(id)
	getValue := func(mem *imap.IntMap[*metric.Schema]) *metric.Schema {
		if mem == nil {
			return nil
		}
		schema, _ := mem.Get(key)
		return schema
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	schema := getValue(s.mutable)
	if schema != nil {
		return schema
	}
	return getValue(s.immutable)
}

// PrepareFlush switches mutable/immutable mem store for flusing schema data.
func (s *metricSchemaStore) PrepareFlush() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.immutable == nil {
		s.immutable = s.mutable
		s.mutable = imap.NewIntMap[*metric.Schema]()
	}
}

func (s *metricSchemaStore) Flush() error {
	kvFlusher := s.family.NewFlusher()
	defer kvFlusher.Release()

	flusher, err := newMetricSchemaFlusher(kvFlusher)
	if err != nil {
		return err
	}
	err = s.immutable.WalkEntry(func(key uint32, value *metric.Schema) error {
		if !value.NeedWrite() {
			return nil
		}
		flusher.Prepare(key)
		if err0 := flusher.Write(value); err0 != nil {
			return err0
		}
		return flusher.Commit()
	})
	if err != nil {
		return err
	}
	err = flusher.Close()
	if err != nil {
		return err
	}

	s.lock.Lock()
	// mark schema persisted
	_ = s.immutable.WalkEntry(func(_ uint32, value *metric.Schema) error {
		value.MarkPersisted()
		return nil
	})
	s.immutable = nil
	s.cache.Purge()
	s.lock.Unlock()
	return nil
}

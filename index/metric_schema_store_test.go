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
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

func TestMetricSchemaStore_GenFieldID(t *testing.T) {
	name := "./gen_field"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	kvStore, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	family, err := kvStore.CreateFamily("field", kv.FamilyOption{Merger: string(v1.MetricSchemaMerger)})
	assert.NoError(t, err)
	s := NewMetricSchemaStore(family)
	fID, err := s.genFieldID(10, field.Meta{Name: "test", Type: field.SumField})
	assert.NoError(t, err)
	assert.Equal(t, field.ID(0), fID)
	schema, err := s.GetSchema(10)
	assert.NoError(t, err)
	fm, ok := schema.Fields.Find(field.Name("test"))
	assert.True(t, ok)
	assert.Equal(t, field.ID(0), fm.ID)

	// just get field
	fID, err = s.genFieldID(10, field.Meta{Name: "test", Type: field.SumField})
	assert.NoError(t, err)
	assert.Equal(t, field.ID(0), fID)
}

func TestMetricSchemaStore_GenTagKeyID(t *testing.T) {
	name := "./gen_tag_key"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	kvStore, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	family, err := kvStore.CreateFamily("tag_key", kv.FamilyOption{Merger: string(v1.MetricSchemaMerger)})
	assert.NoError(t, err)
	s := NewMetricSchemaStore(family)
	tagKeyID, err := s.genTagKeyID(10, []byte("key1"), func() uint32 { return 10 })
	assert.NoError(t, err)
	assert.Equal(t, tag.KeyID(10), tagKeyID)
	schema, err := s.GetSchema(10)
	assert.NoError(t, err)
	tagKey, ok := schema.TagKeys.Find("key1")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(10), tagKey.ID)
	s.PrepareFlush()
	assert.NoError(t, s.Flush())
	tagKeyID, err = s.genTagKeyID(10, []byte("key2"), func() uint32 { return 11 })
	assert.NoError(t, err)
	assert.Equal(t, tag.KeyID(11), tagKeyID)
	s.PrepareFlush()
	assert.NoError(t, s.Flush())
	family.Compact()
	time.Sleep(time.Second)

	schema, err = s.GetSchema(10)
	assert.NoError(t, err)
	tagKey, ok = schema.TagKeys.Find("key1")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(10), tagKey.ID)
	tagKey, ok = schema.TagKeys.Find("key2")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(11), tagKey.ID)

	// from cache
	schema, err = s.GetSchema(10)
	assert.NoError(t, err)
	assert.NotNil(t, schema)
	// just get
	tagKeyID, err = s.genTagKeyID(10, []byte("key2"), func() uint32 { return 133 })
	assert.NoError(t, err)
	assert.Equal(t, tag.KeyID(11), tagKeyID)
}

func TestMetricSchemaStore_Gen_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	family := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	ms := NewMetricSchemaStore(family)

	t.Run("get schema from kv error when gen field", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		fID, err := ms.genFieldID(10, field.Meta{Name: "test", Type: field.SumField})
		assert.Error(t, err)
		assert.Equal(t, field.ID(0), fID)
	})

	t.Run("get schema from kv error when gen tag", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		tID, err := ms.genTagKeyID(10, []byte("key"), func() uint32 { return 1 })
		assert.Error(t, err)
		assert.Equal(t, tag.KeyID(0), tID)
	})

	snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	t.Run("field limit", func(t *testing.T) {
		for i := 0; i < math.MaxUint8; i++ {
			fID, err := ms.genFieldID(10, field.Meta{Name: field.Name(fmt.Sprintf("test-%d", i)), Type: field.SumField})
			assert.NoError(t, err)
			assert.Equal(t, field.ID(i), fID)
		}
		fID, err := ms.genFieldID(10, field.Meta{Name: "limit", Type: field.SumField})
		assert.Equal(t, constants.ErrTooManyFields, err)
		assert.Equal(t, field.ID(0), fID)
	})

	t.Run("tag limit", func(t *testing.T) {
		for i := 0; i < math.MaxUint8; i++ {
			tID, err := ms.genTagKeyID(10, []byte(fmt.Sprintf("key-%d", i)), func() uint32 { return uint32(i) })
			assert.NoError(t, err)
			assert.Equal(t, tag.KeyID(i), tID)
		}
		tID, err := ms.genTagKeyID(10, []byte("limit"), func() uint32 { return 1 })
		assert.Equal(t, constants.ErrTooManyTagKeys, err)
		assert.Equal(t, tag.KeyID(0), tID)
	})
}

func TestMetricSchemaStore_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newMetricSchemaFlusher = v1.NewMetricSchemaFlusher
		ctrl.Finish()
	}()

	family := kv.NewMockFamily(ctrl)
	kvFlusher := kv.NewMockFlusher(ctrl)
	family.EXPECT().NewFlusher().Return(kvFlusher).AnyTimes()
	kvFlusher.EXPECT().Release().AnyTimes()
	flusher := v1.NewMockMetricSchemaFlusher(ctrl)
	flusher.EXPECT().Prepare(gomock.Any()).AnyTimes()
	sm := NewMetricSchemaStore(family)
	sm1 := sm.(*metricSchemaStore)

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "create flusher error",
			prepare: func() {
				newMetricSchemaFlusher = func(_ kv.Flusher) (v1.MetricSchemaFlusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "close flusher error",
			prepare: func() {
				newMetricSchemaFlusher = func(_ kv.Flusher) (v1.MetricSchemaFlusher, error) {
					return flusher, nil
				}
				sm1.immutable = imap.NewIntMap[*metric.Schema]()
				flusher.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "write schema error",
			prepare: func() {
				newMetricSchemaFlusher = func(_ kv.Flusher) (v1.MetricSchemaFlusher, error) {
					return flusher, nil
				}
				sm1.immutable = imap.NewIntMap[*metric.Schema]()
				sm1.immutable.Put(10, &metric.Schema{
					TagKeys: tag.Metas{{Key: "key"}},
				})
				flusher.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "commit schema error",
			prepare: func() {
				newMetricSchemaFlusher = func(_ kv.Flusher) (v1.MetricSchemaFlusher, error) {
					return flusher, nil
				}
				sm1.immutable = imap.NewIntMap[*metric.Schema]()
				sm1.immutable.Put(10, &metric.Schema{
					TagKeys: tag.Metas{{Key: "key"}},
				})
				flusher.EXPECT().Write(gomock.Any()).Return(nil)
				flusher.EXPECT().Commit().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "no schema flush",
			prepare: func() {
				newMetricSchemaFlusher = func(_ kv.Flusher) (v1.MetricSchemaFlusher, error) {
					return flusher, nil
				}
				sm1.immutable = imap.NewIntMap[*metric.Schema]()
				sm1.immutable.Put(10, &metric.Schema{
					TagKeys: tag.Metas{{Key: "key", Persisted: true}},
				})
				flusher.EXPECT().Close().Return(nil)
			},
			wantErr: false,
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newMetricSchemaFlusher = v1.NewMetricSchemaFlusher
			}()
			tt.prepare()
			err := sm.Flush()
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

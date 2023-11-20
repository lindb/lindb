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

package v1

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

func TestMetricSchemaFlusher_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kvFlusher := kv.NewMockFlusher(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "create kv write err",
			prepare: func() {
				kvFlusher.EXPECT().StreamWriter().Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create metric schema flusher successfully",
			prepare: func() {
				kvFlusher.EXPECT().StreamWriter().Return(nil, nil)
			},
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			f, err := NewMetricSchemaFlusher(kvFlusher)
			if ((err != nil) != tt.wantErr && f == nil) || (!tt.wantErr && f == nil) {
				t.Errorf("NewMetricSchemaFlusher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricSchema(t *testing.T) {
	name := "./GetMetricSchema"
	family := createMetricSchemaFamily(t, name)
	kvFlusher := family.NewFlusher()
	defer func() {
		_ = os.RemoveAll(name)
		kvFlusher.Release()
	}()

	flusher, err := NewMetricSchemaFlusher(kvFlusher)
	assert.NoError(t, err)
	assert.NotNil(t, flusher)
	flusher.Prepare(100)
	schema := &metric.Schema{TagKeys: tag.Metas{{Key: "key"}}}
	assert.NoError(t, flusher.Write(schema))
	assert.NoError(t, flusher.Commit())
	assert.NoError(t, flusher.Close())

	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(1000)
	assert.NoError(t, err)
	assert.Empty(t, readers)

	schema1 := &metric.Schema{}
	assert.NoError(t, snapshot.Load(100, func(value []byte) error {
		schema1.Unmarshal(value)
		return nil
	}))
	assert.Equal(t, schema, schema1)
}

func TestMetricSchema_Merge(t *testing.T) {
	name := "./MetricSchema_Merge"
	family := createMetricSchemaFamily(t, name)
	defer func() {
		_ = os.RemoveAll(name)
	}()

	write := func(i int) {
		kvFlusher := family.NewFlusher()
		defer kvFlusher.Release()
		flusher, err := NewMetricSchemaFlusher(kvFlusher)
		assert.NoError(t, err)
		assert.NotNil(t, flusher)
		flusher.Prepare(100)
		schema := &metric.Schema{TagKeys: tag.Metas{{Key: fmt.Sprintf("key-%d", i)}}}
		assert.NoError(t, flusher.Write(schema))
		assert.NoError(t, flusher.Commit())
		assert.NoError(t, flusher.Close())
	}
	for i := 0; i < 4; i++ {
		write(i)
	}

	family.Compact()
	time.Sleep(100 * time.Millisecond)

	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	schema1 := &metric.Schema{}
	assert.NoError(t, snapshot.Load(100, func(value []byte) error {
		schema1.Unmarshal(value)
		return nil
	}))

	schema := &metric.Schema{
		TagKeys: tag.Metas{
			{Key: "key-0"},
			{Key: "key-1"},
			{Key: "key-2"},
			{Key: "key-3"},
		},
	}
	sort.Sort(schema1.TagKeys)
	assert.Equal(t, schema, schema1)
}

func TestMetricSchema_Merge_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kvFlusher := kv.NewMockFlusher(ctrl)
	t.Run("create merger error", func(t *testing.T) {
		kvFlusher.EXPECT().StreamWriter().Return(nil, fmt.Errorf("err"))
		m, err := NewMetricScheamMerger(kvFlusher)
		assert.Nil(t, m)
		assert.Error(t, err)
	})

	t.Run("merge error", func(t *testing.T) {
		sw := table.NewMockStreamWriter(ctrl)
		kvFlusher.EXPECT().StreamWriter().Return(sw, nil)
		m, err := NewMetricScheamMerger(kvFlusher)
		assert.NoError(t, err)
		sw.EXPECT().Prepare(gomock.Any())
		sw.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
		m.Init(nil)
		assert.Error(t, m.Merge(12, nil))
	})
}

func createMetricSchemaFamily(t *testing.T, name string) kv.Family {
	store, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	assert.NotNil(t, store)

	family, err := store.CreateFamily(name, kv.FamilyOption{
		Merger: string(MetricSchemaMerger),
	})
	assert.NoError(t, err)
	assert.NotNil(t, family)
	return family
}

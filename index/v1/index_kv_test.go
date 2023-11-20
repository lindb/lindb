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
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/index/model"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

func TestIndexKV_NewFlusher_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := kv.NewMockFlusher(ctrl)
	flusher.EXPECT().StreamWriter().Return(nil, fmt.Errorf("err"))
	f, err := NewIndexKVFlusher(100, flusher)
	assert.Error(t, err)
	assert.Nil(t, f)
}

func TestIndexKV_GetValue(t *testing.T) {
	name := "./GetValue"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	family := createKVIndexFamily(t, name)

	flusher, err := NewIndexKVFlusher(1, family.NewFlusher())
	assert.NoError(t, err)
	assert.NotNil(t, flusher)
	flusher.PrepareBucket(100)
	err = flusher.WriteKVs([][]byte{[]byte("key1"), []byte("key2")}, []uint32{1, 2})
	assert.NoError(t, err)
	err = flusher.CommitBucket()
	assert.NoError(t, err)

	flusher.PrepareBucket(200)
	err = flusher.WriteKVs([][]byte{[]byte("key11"), []byte("key22")}, []uint32{11, 22})
	assert.NoError(t, err)
	err = flusher.CommitBucket()
	assert.NoError(t, err)

	err = flusher.Close()
	assert.NoError(t, err)

	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	reader := NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(100)
	assert.NoError(t, err)
	id, ok := bucket.GetValue([]byte("key2"))
	assert.True(t, ok)
	assert.Equal(t, uint32(2), id)

	id, ok = bucket.GetValue([]byte("key22"))
	assert.False(t, ok)
	assert.Equal(t, uint32(0), id)

	bucket, err = reader.GetBucket(200)
	assert.NoError(t, err)
	id, ok = bucket.GetValue([]byte("key22"))
	assert.True(t, ok)
	assert.Equal(t, uint32(22), id)
}

func TestIndexKV_Reader_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	reader := NewIndexKVReader(snapshot)
	t.Run("get bucket error", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		bucket, err := reader.GetBucket(10)
		assert.Error(t, err)
		assert.Nil(t, bucket)
	})
	t.Run("trie bucket unmarshal error", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Do(func(_ uint32, f func(value []byte) error) {
			_ = f([]byte{1, 0, 0, 0, 0})
		}).Return(fmt.Errorf("err"))
		bucket, err := reader.GetBucket(10)
		assert.Error(t, err)
		assert.Nil(t, bucket)
	})
}

func TestIndexKV_Merge(t *testing.T) {
	name := "./IndexKV_Merge"
	family := createKVIndexFamily(t, name)
	defer func() {
		_ = os.RemoveAll(name)
	}()

	write := func(i int) {
		kvFlusher := family.NewFlusher()
		defer kvFlusher.Release()
		flusher, err := NewIndexKVFlusher(1, family.NewFlusher())
		assert.NoError(t, err)
		flusher.PrepareBucket(100)
		err = flusher.WriteKVs([][]byte{[]byte(fmt.Sprintf("key-%d", i))}, []uint32{uint32(i)})
		assert.NoError(t, err)
		assert.NoError(t, flusher.CommitBucket())
		assert.NoError(t, flusher.Close())
	}
	for i := 0; i < 4; i++ {
		write(i)
	}

	family.Compact()
	time.Sleep(100 * time.Millisecond)

	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	reader := NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(100)
	assert.NoError(t, err)
	for i := 0; i < 4; i++ {
		id, ok := bucket.GetValue([]byte(fmt.Sprintf("key-%d", i)))
		assert.True(t, ok)
		assert.Equal(t, uint32(i), id)
	}
}

func TestIndexKV_Merge_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kvFlusher := kv.NewMockFlusher(ctrl)
	t.Run("create index kv merger error", func(t *testing.T) {
		kvFlusher.EXPECT().StreamWriter().Return(nil, fmt.Errorf("err"))
		f, err := NewIndexKVMerger(kvFlusher)
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("unmarshal bucket error", func(t *testing.T) {
		sm := table.NewMockStreamWriter(ctrl)
		kvFlusher.EXPECT().StreamWriter().Return(sm, nil)
		m, err := NewIndexKVMerger(kvFlusher)
		assert.NoError(t, err)
		assert.Error(t, m.Merge(10, [][]byte{{1, 0, 0, 0, 0}}))
	})

	t.Run("write bucket error", func(t *testing.T) {
		sm := table.NewMockStreamWriter(ctrl)
		kvFlusher.EXPECT().StreamWriter().Return(sm, nil)
		m, err := NewIndexKVMerger(kvFlusher)
		assert.NoError(t, err)
		m.Init(nil)
		sm.EXPECT().Prepare(gomock.Any())
		sm.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
		w := bytes.NewBuffer([]byte{})
		b := model.NewTrieBucketBuilder(10, w)
		assert.NoError(t, b.Write([][]byte{{1, 2, 3}}, []uint32{1}))
		data := w.Bytes()
		assert.Error(t, m.Merge(10, [][]byte{data}))
	})
}

func createKVIndexFamily(t *testing.T, name string) kv.Family {
	store, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	assert.NotNil(t, store)

	family, err := store.CreateFamily(name, kv.FamilyOption{
		Merger: string(IndexKVMerger),
	})
	assert.NoError(t, err)
	assert.NotNil(t, family)
	return family
}

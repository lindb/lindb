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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestNewInvertedIndexFlusher_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kvFlusher := kv.NewMockFlusher(ctrl)
	kvFlusher.EXPECT().StreamWriter().Return(nil, fmt.Errorf("err"))
	f, err := NewInvertedIndexFlusher(kvFlusher)
	assert.Nil(t, f)
	assert.Error(t, err)
}

func TestInvertedIndex(t *testing.T) {
	name := "./TestInvertedIndex"
	family := createInvertedFamily(t, name)
	kvFlusher := family.NewFlusher()
	defer func() {
		_ = os.RemoveAll(name)
		kvFlusher.Release()
	}()

	flusher, err := NewInvertedIndexFlusher(kvFlusher)
	assert.NoError(t, err)
	assert.NotNil(t, flusher)
	flusher.Prepare(100)
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	assert.NoError(t, flusher.Write(seriesIDs))
	assert.NoError(t, flusher.Commit())
	assert.NoError(t, flusher.Close())

	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(1000)
	assert.NoError(t, err)
	assert.Empty(t, readers)

	gotSeriesIDs := roaring.New()
	assert.NoError(t, snapshot.Load(100, func(value []byte) error {
		_, _ = gotSeriesIDs.FromBuffer(value)
		return nil
	}))
	assert.Equal(t, seriesIDs.ToArray(), gotSeriesIDs.ToArray())
}

func TestInvertedIndex_Merge(t *testing.T) {
	name := "./InvertedIndex_Merge"
	family := createInvertedFamily(t, name)
	defer func() {
		_ = os.RemoveAll(name)
	}()
	write := func(i int) {
		kvFlusher := family.NewFlusher()
		defer kvFlusher.Release()
		flusher, err := NewInvertedIndexFlusher(kvFlusher)
		assert.NoError(t, err)
		assert.NotNil(t, flusher)
		flusher.Prepare(100)
		assert.NoError(t, flusher.Write(roaring.BitmapOf(uint32(i))))
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

	gotSeriesIDs := roaring.New()
	assert.NoError(t, snapshot.Load(100, func(value []byte) error {
		_, _ = gotSeriesIDs.FromBuffer(value)
		return nil
	}))
	assert.Equal(t, roaring.BitmapOf(0, 1, 2, 3).ToArray(), gotSeriesIDs.ToArray())
}

func TestInvertedIndex_Merge_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newInvertedIndexFlusher = NewInvertedIndexFlusher
		bitmapUnmarshal = encoding.BitmapUnmarshal
		ctrl.Finish()
	}()

	t.Run("create merge error", func(t *testing.T) {
		newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (InvertedIndexFlusher, error) {
			return nil, fmt.Errorf("err")
		}
		f, err := NewInvertedIndexMerger(nil)
		assert.Error(t, err)
		assert.Nil(t, f)
	})
	f := NewMockInvertedIndexFlusher(ctrl)
	newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (InvertedIndexFlusher, error) {
		return f, nil
	}

	t.Run("unmarshal bitmap error", func(t *testing.T) {
		bitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) (int64, error) {
			return 0, fmt.Errorf("err")
		}
		m, err := NewInvertedIndexMerger(nil)
		assert.NoError(t, err)
		m.Init(nil)
		assert.Error(t, m.Merge(10, [][]byte{{1, 2}}))
	})
	t.Run("write series ids error", func(t *testing.T) {
		m, err := NewInvertedIndexMerger(nil)
		assert.NoError(t, err)
		f.EXPECT().Prepare(gomock.Any())
		f.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err"))
		assert.Error(t, m.Merge(10, nil))
	})
}

func createInvertedFamily(t *testing.T, name string) kv.Family {
	store, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	assert.NotNil(t, store)

	family, err := store.CreateFamily(name, kv.FamilyOption{
		Merger: string(InvertedIndexMerger),
	})
	assert.NoError(t, err)
	assert.NotNil(t, family)
	return family
}

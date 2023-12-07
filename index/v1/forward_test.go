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
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestForwardIndex(t *testing.T) {
	name := "./TestForwardIndex"
	family := createForwardFamily(t, name)
	kvFlusher := family.NewFlusher()
	defer func() {
		_ = os.RemoveAll(name)
		kvFlusher.Release()
	}()

	flusher, err := NewForwardIndexFlusher(kvFlusher)
	assert.NoError(t, err)
	assert.NotNil(t, flusher)
	flusher.Prepare(100)
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	assert.NoError(t, flusher.WriteSeriesIDs(seriesIDs))
	assert.NoError(t, flusher.WriteTagValueIDs([]uint32{1, 2, 3}))
	assert.NoError(t, flusher.Commit())
	assert.NoError(t, flusher.Close())

	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(1000)
	assert.NoError(t, err)
	assert.Empty(t, readers)

	readers, err = snapshot.FindReaders(100)
	assert.NoError(t, err)
	reader := NewForwardReader(readers)
	gotSeriesIDs, err := reader.GetSeriesIDsForTagKeyID(100)
	assert.NoError(t, err)
	assert.Equal(t, seriesIDs.ToArray(), gotSeriesIDs.ToArray())

	scanners, err := reader.GetGroupingScanner(100, roaring.BitmapOf(1, 2))
	assert.NoError(t, err)
	assert.Len(t, scanners, 1)
	lowSeriesIDs, tagValueIDs := scanners[0].GetSeriesAndTagValue(0)
	assert.Equal(t, []uint16{1, 2, 3}, lowSeriesIDs.ToArray())
	assert.Equal(t, []uint32{1, 2, 3}, tagValueIDs)
	lowSeriesIDs, tagValueIDs = scanners[0].GetSeriesAndTagValue(2)
	assert.Nil(t, lowSeriesIDs)
	assert.Nil(t, tagValueIDs)

	scanners, err = reader.GetGroupingScanner(100, roaring.BitmapOf(100))
	assert.NoError(t, err)
	assert.Len(t, scanners, 0)
}

func TestForwardIndex_Merge(t *testing.T) {
	name := "./ForwardIndex_Merge"
	family := createForwardFamily(t, name)
	defer func() {
		_ = os.RemoveAll(name)
	}()

	seriesIDs := roaring.New()
	write := func(i int) {
		kvFlusher := family.NewFlusher()
		defer kvFlusher.Release()
		flusher, err := NewForwardIndexFlusher(kvFlusher)
		assert.NoError(t, err)
		assert.NotNil(t, flusher)
		flusher.Prepare(100)
		seriesID := uint32(i * math.MaxUint16)
		seriesIDs.Add(seriesID)
		assert.NoError(t, flusher.WriteSeriesIDs(roaring.BitmapOf(seriesID)))
		assert.NoError(t, flusher.WriteTagValueIDs([]uint32{uint32(i)}))
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

	readers, err := snapshot.FindReaders(100)
	assert.NoError(t, err)
	reader := NewForwardReader(readers)
	gotSeriesIDs, err := reader.GetSeriesIDsForTagKeyID(100)
	assert.NoError(t, err)
	assert.Equal(t, seriesIDs.ToArray(), gotSeriesIDs.ToArray())
}

func TestForwardIndex_Merge_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newForwardIndexFlusher = NewForwardIndexFlusher
		newTagForwardReader = NewTagForwardReader
		ctrl.Finish()
	}()

	t.Run("create merge error", func(t *testing.T) {
		newForwardIndexFlusher = func(kvFlusher kv.Flusher) (ForwardIndexFlusher, error) {
			return nil, fmt.Errorf("Err")
		}
		f, err := NewForwardIndexMerger(nil)
		assert.Error(t, err)
		assert.Nil(t, f)
	})
	f := NewMockForwardIndexFlusher(ctrl)
	newForwardIndexFlusher = func(kvFlusher kv.Flusher) (ForwardIndexFlusher, error) {
		return f, nil
	}

	t.Run("new tag forward index reader error", func(t *testing.T) {
		newTagForwardReader = func(buf []byte) (TagForwardReader, error) {
			return nil, fmt.Errorf("err")
		}
		m, err := NewForwardIndexMerger(nil)
		assert.NoError(t, err)
		m.Init(nil)
		assert.Error(t, m.Merge(10, [][]byte{{1, 2}}))
	})
	reader := NewMockTagForwardReader(ctrl)
	newTagForwardReader = func(buf []byte) (TagForwardReader, error) {
		return reader, nil
	}
	f.EXPECT().Prepare(gomock.Any()).AnyTimes()
	reader.EXPECT().GetSeriesIDs().Return(roaring.BitmapOf(1, 2, 3)).AnyTimes()
	reader.EXPECT().GetSeriesAndTagValue(gomock.Any()).Return(nil, nil).AnyTimes()
	t.Run("write series ids error", func(t *testing.T) {
		m, err := NewForwardIndexMerger(nil)
		assert.NoError(t, err)
		f.EXPECT().WriteSeriesIDs(gomock.Any()).Return(fmt.Errorf("err"))
		assert.Error(t, m.Merge(10, [][]byte{{1, 2}}))
	})
	t.Run("write tag value ids error", func(t *testing.T) {
		m, err := NewForwardIndexMerger(nil)
		assert.NoError(t, err)
		f.EXPECT().WriteSeriesIDs(gomock.Any()).Return(nil)
		f.EXPECT().WriteTagValueIDs(gomock.Any()).Return(fmt.Errorf("err"))
		assert.Error(t, m.Merge(10, [][]byte{{1, 2}}))
	})
}

func TestForwardIndex_Flusher_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kvFlusher := kv.NewMockFlusher(ctrl)
	t.Run("create flusher error", func(t *testing.T) {
		kvFlusher.EXPECT().StreamWriter().Return(nil, fmt.Errorf("err"))
		f, err := NewForwardIndexFlusher(kvFlusher)
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	sw := table.NewMockStreamWriter(ctrl)
	kvFlusher.EXPECT().StreamWriter().Return(sw, nil)
	f, err := NewForwardIndexFlusher(kvFlusher)
	assert.NoError(t, err)
	t.Run("write seried ids error", func(t *testing.T) {
		sw.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
		assert.Error(t, f.WriteSeriesIDs(roaring.New()))
	})
	t.Run("write tag value ids error", func(t *testing.T) {
		sw.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
		assert.Error(t, f.WriteTagValueIDs([]uint32{1}))
	})
}

func TestForwardIndex_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		newTagForwardReader = NewTagForwardReader
	}()

	kvReader := table.NewMockReader(ctrl)
	fReader := NewForwardReader([]table.Reader{kvReader})
	t.Run("new tag forward reader error", func(t *testing.T) {
		kvReader.EXPECT().Get(gomock.Any()).Return(nil, nil)
		newTagForwardReader = func(buf []byte) (TagForwardReader, error) {
			return nil, fmt.Errorf("err")
		}
		seriesIDs, err := fReader.GetSeriesIDsForTagKeyID(10)
		assert.Nil(t, seriesIDs)
		assert.Error(t, err)
	})
	t.Run("read group scanner error", func(t *testing.T) {
		kvReader.EXPECT().Get(gomock.Any()).Return(nil, nil)
		newTagForwardReader = func(buf []byte) (TagForwardReader, error) {
			return nil, fmt.Errorf("err")
		}
		scanners, err := fReader.GetGroupingScanner(10, roaring.BitmapOf(1, 2))
		assert.Nil(t, scanners)
		assert.Error(t, err)
	})
	newTagForwardReader = NewTagForwardReader
	t.Run("get value from kv store error", func(t *testing.T) {
		kvReader.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
		seriesIDs, err := fReader.GetSeriesIDsForTagKeyID(10)
		assert.Empty(t, seriesIDs)
		assert.NoError(t, err)
	})
}

func TestForwardIndex_NewTagForwardReader_Error(t *testing.T) {
	defer func() {
		bitmapUnmarshal = encoding.BitmapUnmarshal
	}()
	bitmapUnmarshal = func(_ *roaring.Bitmap, _ []byte) (int64, error) {
		return 0, fmt.Errorf("err")
	}
	r, err := NewTagForwardReader(nil)
	assert.Nil(t, r)
	assert.Error(t, err)
}

func createForwardFamily(t *testing.T, name string) kv.Family {
	store, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	assert.NotNil(t, store)

	family, err := store.CreateFamily(name, kv.FamilyOption{
		Merger: string(ForwardIndexMerger),
	})
	assert.NoError(t, err)
	assert.NotNil(t, family)
	return family
}

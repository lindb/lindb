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

package indexdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

func TestInvertedIndex_buildInvertIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newInvertedReaderFunc = invertedindex.NewInvertedReader
		ctrl.Finish()
	}()
	reader := invertedindex.NewMockInvertedReader(ctrl)
	newInvertedReaderFunc = func(readers []table.Reader) invertedindex.InvertedReader {
		return reader
	}

	index := prepareInvertedIndex(ctrl)
	family := kv.NewMockFamily(ctrl)
	idx := index.(*invertedIndex)
	idx.invertedFamily = family
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()

	// case 1: get series ids by tag value ids
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err := index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(2, roaring.BitmapOf(2))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(2), seriesIDs)

	// case 2: tag key is not exist
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(4, roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)

	// case 3: tag value ids is not exist
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(10, 20))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)
	// case 4: tag key not exist
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(4, roaring.BitmapOf(10, 20))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)
	// case 5: get series ids, get empty reader
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(10, 20))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)
	// case 6: get kv readers err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(10, 20))
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 6: reader get data err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).AnyTimes()
	reader.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(10, 20))
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 6: reader get data success
	reader.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), seriesIDs)

	//case 7: get immutable data
	tagIndex := NewMockTagIndex(ctrl)
	idx.immutable = NewTagIndexStore()
	idx.immutable.Put(50, tagIndex)
	reader.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(), nil)
	tagIndex.EXPECT().getSeriesIDsByTagValueIDs(gomock.Any()).Return(roaring.BitmapOf(10, 200, 3000))
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(50, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(10, 200, 3000), seriesIDs)
}

func TestInvertedIndex_GetSeriesIDsForTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newForwardReaderFunc = invertedindex.NewForwardReader
		ctrl.Finish()
	}()
	reader := invertedindex.NewMockForwardReader(ctrl)
	newForwardReaderFunc = func(readers []table.Reader) invertedindex.ForwardReader {
		return reader
	}

	index := prepareInvertedIndex(ctrl)
	family := kv.NewMockFamily(ctrl)
	idx := index.(*invertedIndex)
	idx.forwardFamily = family
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()

	// case 1: get reader err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err := index.GetSeriesIDsForTags([]uint32{1, 2, 3})
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 2: reader get data success
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).AnyTimes()
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	seriesIDs, err = index.GetSeriesIDsForTags([]uint32{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), seriesIDs)
	// case 3: reader get series ids err
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = index.GetSeriesIDsForTags([]uint32{1, 2, 3})
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
}

func TestInvertedIndex_GetSeriesIDsForTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newForwardReaderFunc = invertedindex.NewForwardReader
		ctrl.Finish()
	}()
	reader := invertedindex.NewMockForwardReader(ctrl)
	newForwardReaderFunc = func(readers []table.Reader) invertedindex.ForwardReader {
		return reader
	}

	index := prepareInvertedIndex(ctrl)
	family := kv.NewMockFamily(ctrl)
	idx := index.(*invertedIndex)
	idx.forwardFamily = family
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()

	// case 1: reader get data success
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).AnyTimes()
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil)
	seriesIDs, err := index.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), seriesIDs)
}

func TestInvertedIndex_GetGroupingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newForwardReaderFunc = invertedindex.NewForwardReader
		ctrl.Finish()
	}()

	index := prepareInvertedIndex(ctrl)
	idx := index.(*invertedIndex)
	family := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	idx.forwardFamily = family

	// case 1: get sst file reader err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	ctx, err := index.GetGroupingContext([]uint32{3, 4}, nil)
	assert.Error(t, err)
	assert.Nil(t, ctx)
	// case 2: get empty reader
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil).Times(2)
	ctx, err = index.GetGroupingContext([]uint32{1, 2}, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
	// case 3: get scanner from file err
	reader := invertedindex.NewMockForwardReader(ctrl)
	newForwardReaderFunc = func(readers []table.Reader) invertedindex.ForwardReader {
		return reader
	}
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil)
	reader.EXPECT().GetGroupingScanner(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	ctx, err = index.GetGroupingContext([]uint32{1, 2}, roaring.BitmapOf(1, 2, 3))
	assert.Error(t, err)
	assert.Nil(t, ctx)
	// case 4: get scanner from file err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).Times(2)
	reader.EXPECT().GetGroupingScanner(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
	ctx, err = index.GetGroupingContext([]uint32{1, 2}, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
}

func TestInvertedIndex_FlushInvertedIndexTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newInvertedFlusherFunc = invertedindex.NewInvertedFlusher
		newForwardFlusherFunc = invertedindex.NewForwardFlusher
		ctrl.Finish()
	}()
	invertedFamily := kv.NewMockFamily(ctrl)
	inverted := invertedindex.NewMockInvertedFlusher(ctrl)
	newInvertedFlusherFunc = func(kvFlusher kv.Flusher) invertedindex.InvertedFlusher {
		return inverted
	}
	forwardFamily := kv.NewMockFamily(ctrl)
	forward := invertedindex.NewMockForwardFlusher(ctrl)
	newForwardFlusherFunc = func(kvFlusher kv.Flusher) invertedindex.ForwardFlusher {
		return forward
	}

	index := newInvertedIndex(nil, forwardFamily, invertedFamily)
	// case 1: flush not tiger
	err := index.Flush()
	assert.NoError(t, err)

	// mock data
	idx := index.(*invertedIndex)
	tagIndex := NewMockTagIndex(ctrl)
	idx.mutable.Put(5, tagIndex)

	// case 1: flush tag index flush err, immutable cannot set nil
	gomock.InOrder(
		forwardFamily.EXPECT().NewFlusher().Return(nil),
		invertedFamily.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")),
	)
	err = index.Flush()
	assert.Error(t, err)
	assert.NotNil(t, idx.immutable)
	// case 2: commit forward err
	gomock.InOrder(
		forwardFamily.EXPECT().NewFlusher().Return(nil),
		invertedFamily.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		forward.EXPECT().Commit().Return(fmt.Errorf("err")),
	)
	err = index.Flush()
	assert.Error(t, err)
	assert.NotNil(t, idx.immutable)
	// case 3: commit inverted err
	gomock.InOrder(
		forwardFamily.EXPECT().NewFlusher().Return(nil),
		invertedFamily.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		forward.EXPECT().Commit().Return(nil),
		inverted.EXPECT().Commit().Return(fmt.Errorf("err")),
	)
	err = index.Flush()
	assert.Error(t, err)
	assert.NotNil(t, idx.immutable)
	// case 4: commit success
	gomock.InOrder(
		forwardFamily.EXPECT().NewFlusher().Return(nil),
		invertedFamily.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		forward.EXPECT().Commit().Return(nil),
		inverted.EXPECT().Commit().Return(nil),
	)
	err = index.Flush()
	assert.NoError(t, err)
	assert.Nil(t, idx.immutable)
}

func prepareInvertedIndex(ctrl *gomock.Controller) InvertedIndex {
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	tagMetadata := metadb.NewMockTagMetadata(ctrl)
	metadata.EXPECT().DatabaseName().Return("test").AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMetadata).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(1), nil).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "zone").Return(uint32(2), nil).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "zone_err").Return(uint32(0), fmt.Errorf("err")).AnyTimes()
	tagMetadata.EXPECT().GenTagValueID(uint32(1), "1.1.1.1").Return(uint32(1), nil).Times(2)
	tagMetadata.EXPECT().GenTagValueID(uint32(1), "1.1.1.5").Return(uint32(0), fmt.Errorf("err"))
	tagMetadata.EXPECT().GenTagValueID(uint32(2), "sh").Return(uint32(1), nil)
	tagMetadata.EXPECT().GenTagValueID(uint32(2), "bj").Return(uint32(2), nil)
	index := newInvertedIndex(metadata, nil, nil)
	index.buildInvertIndex("ns", "name", map[string]string{
		"host": "1.1.1.1",
		"zone": "sh",
	}, 1)
	index.buildInvertIndex("ns", "name", map[string]string{
		"host": "1.1.1.1",
		"zone": "bj",
	}, 2)
	index.buildInvertIndex("ns", "name", map[string]string{
		"host":     "1.1.1.5",
		"zone_err": "bj",
	}, 3)
	return index
}

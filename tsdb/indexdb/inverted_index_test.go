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
		newReaderFunc = invertedindex.NewReader
		ctrl.Finish()
	}()
	reader := invertedindex.NewMockReader(ctrl)
	newReaderFunc = func(readers []table.Reader) invertedindex.Reader {
		return reader
	}

	index := prepareInvertedIndex(ctrl)
	family := kv.NewMockFamily(ctrl)
	idx := index.(*invertedIndex)
	idx.family = family
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
	seriesIDs, err = index.GetSeriesIDsForTag(4)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)
	// case 5: get all series ids under tag key
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	seriesIDs, err = index.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	// case 6: get kv readers err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err")).Times(2)
	seriesIDs, err = index.GetSeriesIDsForTag(1)
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1, 2, 3))
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 6: reader get data err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).AnyTimes()
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = index.GetSeriesIDsForTag(1)
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	reader.EXPECT().FindSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1, 2, 3))
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 6: reader get data success
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil)
	seriesIDs, err = index.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), seriesIDs)
	reader.EXPECT().FindSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), seriesIDs)

	// mock immutable data
	tagIndex := NewMockTagIndex(ctrl)
	idx.immutable = NewTagIndexStore()
	idx.immutable.Put(50, tagIndex)
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(roaring.New(), nil)
	tagIndex.EXPECT().getAllSeriesIDs().Return(roaring.BitmapOf(10, 200, 3000))
	seriesIDs, err = index.GetSeriesIDsForTag(50)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(10, 200, 3000), seriesIDs)
}

func TestInvertedIndex_GetSeriesIDsForTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newReaderFunc = invertedindex.NewReader
		ctrl.Finish()
	}()
	reader := invertedindex.NewMockReader(ctrl)
	newReaderFunc = func(readers []table.Reader) invertedindex.Reader {
		return reader
	}

	index := prepareInvertedIndex(ctrl)
	family := kv.NewMockFamily(ctrl)
	idx := index.(*invertedIndex)
	idx.family = family
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()

	// case 1: get reader err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err := index.GetSeriesIDsForTags([]uint32{1, 2, 3})
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 2: reader get data success
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).Times(3)
	reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	seriesIDs, err = index.GetSeriesIDsForTags([]uint32{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), seriesIDs)
}

func TestInvertedIndex_GetGroupingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := prepareInvertedIndex(ctrl)
	ctx, err := index.GetGroupingContext([]uint32{3, 4})
	assert.Error(t, err)
	assert.Nil(t, ctx)

	ctx, err = index.GetGroupingContext([]uint32{1, 2})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	ctx, err = index.GetGroupingContext([]uint32{1})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
}

func TestInvertedIndex_FlushInvertedIndexTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newFlusherFunc = invertedindex.NewFlusher
		ctrl.Finish()
	}()
	family := kv.NewMockFamily(ctrl)
	flusher := invertedindex.NewMockFlusher(ctrl)
	newFlusherFunc = func(kvFlusher kv.Flusher) invertedindex.Flusher {
		return flusher
	}
	index := newInvertedIndex(nil, family)
	// case 1: flush not tiger
	err := index.Flush()
	assert.NoError(t, err)

	// mock data
	idx := index.(*invertedIndex)
	tagIndex := NewMockTagIndex(ctrl)
	idx.mutable.Put(5, tagIndex)

	// case 2: flush tag key err, immutable cannot set nil
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any()).Return(nil),
		flusher.EXPECT().FlushTagKeyID(uint32(5)).Return(fmt.Errorf("err")),
	)
	err = index.Flush()
	assert.Error(t, err)
	assert.NotNil(t, idx.immutable)
	// case 3: flush tag index flush err, immutable cannot set nil
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any()).Return(fmt.Errorf("err")),
	)
	err = index.Flush()
	assert.Error(t, err)
	assert.NotNil(t, idx.immutable)
	// case 4: commit err
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any()).Return(nil),
		flusher.EXPECT().FlushTagKeyID(uint32(5)).Return(nil),
		flusher.EXPECT().Commit().Return(fmt.Errorf("err")),
	)
	err = index.Flush()
	assert.Error(t, err)
	assert.NotNil(t, idx.immutable)
	// case 4: commit success
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		tagIndex.EXPECT().flush(gomock.Any()).Return(nil),
		flusher.EXPECT().FlushTagKeyID(uint32(5)).Return(nil),
		flusher.EXPECT().Commit().Return(nil),
	)
	err = index.Flush()
	assert.NoError(t, err)
	assert.Nil(t, idx.immutable)
}

func prepareInvertedIndex(ctrl *gomock.Controller) InvertedIndex {
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	tagMetadata := metadb.NewMockTagMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMetadata).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(1), nil).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "zone").Return(uint32(2), nil).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "zone_err").Return(uint32(0), fmt.Errorf("err")).AnyTimes()
	tagMetadata.EXPECT().GenTagValueID(uint32(1), "1.1.1.1").Return(uint32(1), nil).Times(2)
	tagMetadata.EXPECT().GenTagValueID(uint32(1), "1.1.1.5").Return(uint32(0), fmt.Errorf("err"))
	tagMetadata.EXPECT().GenTagValueID(uint32(2), "sh").Return(uint32(1), nil)
	tagMetadata.EXPECT().GenTagValueID(uint32(2), "bj").Return(uint32(2), nil)
	index := newInvertedIndex(metadata, nil)
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

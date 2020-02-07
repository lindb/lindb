package metadb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore/metricsmeta"
)

func TestTagMetadata_GenTagValueID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagReaderFunc = metricsmeta.NewTagReader
		ctrl.Finish()
	}()

	meta, _, snapshot := mockTagMetadata(ctrl)

	tagReader := metricsmeta.NewMockTagReader(ctrl)
	newTagReaderFunc = func(readers []table.Reader) metricsmeta.TagReader {
		return tagReader
	}

	// case 1: gen tag value id
	snapshot.EXPECT().FindReaders(uint32(1)).Return(nil, nil)
	tagValueID, err := meta.GenTagValueID(1, "tag-value-1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), tagValueID)
	// case 2: get tag value id from mem
	tagValueID, err = meta.GenTagValueID(1, "tag-value-1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), tagValueID)
	// case 3: get kv readers err
	snapshot.EXPECT().FindReaders(uint32(1)).Return(nil, fmt.Errorf("err"))
	tagValueID, err = meta.GenTagValueID(1, "tag-value-err")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), tagValueID)
	// case 4: get tag value from kv store
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil).AnyTimes()
	tagReader.EXPECT().GetTagValueID(uint32(1), "tag-value-2").Return(uint32(2), nil)
	tagValueID, err = meta.GenTagValueID(1, "tag-value-2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), tagValueID)
	// case 5: get tag value from kv store err
	tagReader.EXPECT().GetTagValueID(uint32(1), "tag-value-2-err").Return(uint32(0), fmt.Errorf("err"))
	tagValueID, err = meta.GenTagValueID(1, "tag-value-2-err")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), tagValueID)
	// case 6: init tag entry from kv store err
	tagReader.EXPECT().GetTagValueID(uint32(5), "tag-value-2").Return(uint32(0), constants.ErrNotFound)
	tagReader.EXPECT().GetTagValueSeq(uint32(5)).Return(uint32(0), fmt.Errorf("err"))
	tagValueID, err = meta.GenTagValueID(5, "tag-value-2")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), tagValueID)
	// case 7: init tag entry from kv store
	tagReader.EXPECT().GetTagValueID(uint32(5), "tag-value-2").Return(uint32(0), constants.ErrNotFound)
	tagReader.EXPECT().GetTagValueSeq(uint32(5)).Return(uint32(20), nil)
	tagValueID, err = meta.GenTagValueID(5, "tag-value-2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(21), tagValueID)
	// case 8: get tag value id from immutable
	m := meta.(*tagMetadata)
	m.rwMutex.Lock()
	m.immutable = NewTagStore()
	tagEntry := newTagEntry(10)
	tagEntry.addTagValue("tag-value-5", 10)
	m.immutable.Put(5, tagEntry)
	m.rwMutex.Unlock()
	tagValueID, err = meta.GenTagValueID(5, "tag-value-5")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagValueID)
	// case 8: get tag value id from immutable not exist
	tagReader.EXPECT().GetTagValueID(uint32(5), "tag-value-6").Return(uint32(0), constants.ErrNotFound)
	tagValueID, err = meta.GenTagValueID(5, "tag-value-6")
	assert.NoError(t, err)
	assert.Equal(t, uint32(22), tagValueID)
}

func TestTagMetadata_SuggestTagValues(t *testing.T) {
	meta := NewTagMetadata(nil)
	assert.Panics(t, func() {
		_ = meta.SuggestTagValues(1, "11", 10)
	})
}

func TestTagMetadata_FindTagValueDsByExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagReaderFunc = metricsmeta.NewTagReader
		ctrl.Finish()
	}()

	meta, _, snapshot := mockTagMetadata(ctrl)

	tagReader := metricsmeta.NewMockTagReader(ctrl)
	newTagReaderFunc = func(readers []table.Reader) metricsmeta.TagReader {
		return tagReader
	}
	mockTagMetadataMemData(meta)

	// case 1: find from mutable
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	ids, err := meta.FindTagValueDsByExpr(uint32(5), &stmt.EqualsExpr{Value: "tag-value-5"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(10), ids)
	// case 2: find from mutable
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	ids, err = meta.FindTagValueDsByExpr(uint32(10), &stmt.EqualsExpr{Value: "tag-value-20"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(20), ids)
	// case 3: no data
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	ids, err = meta.FindTagValueDsByExpr(uint32(10), &stmt.EqualsExpr{Value: "tag-value-210"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), ids)
	// case 4: kv store find readers err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	ids, err = meta.FindTagValueDsByExpr(uint32(10), &stmt.EqualsExpr{Value: "tag-value-20"})
	assert.Error(t, err)
	assert.Nil(t, ids)
	// case 5: find ids from kv err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil)
	tagReader.EXPECT().FindValueIDsByExprForTagKeyID(uint32(10), gomock.Any()).Return(nil, fmt.Errorf("err"))
	ids, err = meta.FindTagValueDsByExpr(uint32(10), &stmt.EqualsExpr{Value: "tag-value-20"})
	assert.Error(t, err)
	assert.Nil(t, ids)
	// case 5: find ids from kv
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil)
	tagReader.EXPECT().FindValueIDsByExprForTagKeyID(uint32(10), gomock.Any()).Return(roaring.BitmapOf(30, 40), nil)
	ids, err = meta.FindTagValueDsByExpr(uint32(10), &stmt.EqualsExpr{Value: "tag-value-20"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(20, 30, 40), ids)
}

func TestTagMetadata_GetTagValueIDsForTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagReaderFunc = metricsmeta.NewTagReader
		ctrl.Finish()
	}()

	meta, _, snapshot := mockTagMetadata(ctrl)

	tagReader := metricsmeta.NewMockTagReader(ctrl)
	newTagReaderFunc = func(readers []table.Reader) metricsmeta.TagReader {
		return tagReader
	}
	mockTagMetadataMemData(meta)

	// case 1: get from mutable
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	ids, err := meta.GetTagValueIDsForTag(uint32(5))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(10), ids)
	// case 2: get from mutable
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	ids, err = meta.GetTagValueIDsForTag(uint32(10))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(20), ids)
	// case 3: no data
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	ids, err = meta.GetTagValueIDsForTag(uint32(100))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), ids)
	// case 4: kv store find readers err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	ids, err = meta.GetTagValueIDsForTag(uint32(10))
	assert.Error(t, err)
	assert.Nil(t, ids)
	// case 5: find ids from kv err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil)
	tagReader.EXPECT().GetTagValueIDsForTagKeyID(uint32(10)).Return(nil, fmt.Errorf("err"))
	ids, err = meta.GetTagValueIDsForTag(uint32(10))
	assert.Error(t, err)
	assert.Nil(t, ids)
	// case 5: find ids from kv
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{table.NewMockReader(ctrl)}, nil)
	tagReader.EXPECT().GetTagValueIDsForTagKeyID(uint32(10)).Return(roaring.BitmapOf(30, 40), nil)
	ids, err = meta.GetTagValueIDsForTag(uint32(10))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(20, 30, 40), ids)
}

func TestTagMetadata_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagFlusherFunc = metricsmeta.NewTagFlusher
		ctrl.Finish()
	}()

	meta, family, _ := mockTagMetadata(ctrl)
	flusher := metricsmeta.NewMockTagFlusher(ctrl)
	newTagFlusherFunc = func(kvFlusher kv.Flusher) metricsmeta.TagFlusher {
		return flusher
	}
	// case 1: flush not tiger
	err := meta.Flush()
	assert.NoError(t, err)

	// mock data
	m := meta.(*tagMetadata)
	m.rwMutex.Lock()
	tagEntry := newTagEntry(10)
	tagEntry.addTagValue("tag-value-5", 10)
	m.mutable.Put(5, tagEntry)
	m.rwMutex.Unlock()
	// case 2: flush tag key err, immutable cannot set nil
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		flusher.EXPECT().FlushTagValue("tag-value-5", uint32(10)),
		flusher.EXPECT().FlushTagKeyID(uint32(5), uint32(10)).Return(fmt.Errorf("err")),
	)
	err = meta.Flush()
	assert.Error(t, err)
	m.rwMutex.Lock()
	assert.NotNil(t, m.immutable)
	m.rwMutex.Unlock()
	// case 3: commit err, immutable cannot set nil
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		flusher.EXPECT().FlushTagValue("tag-value-5", uint32(10)),
		flusher.EXPECT().FlushTagKeyID(uint32(5), uint32(10)).Return(nil),
		flusher.EXPECT().Commit().Return(fmt.Errorf("err")),
	)
	err = meta.Flush()
	assert.Error(t, err)
	m.rwMutex.Lock()
	assert.NotNil(t, m.immutable)
	m.rwMutex.Unlock()
	// case 4: flush success, immutable is nil
	gomock.InOrder(
		family.EXPECT().NewFlusher().Return(nil),
		flusher.EXPECT().FlushTagValue("tag-value-5", uint32(10)),
		flusher.EXPECT().FlushTagKeyID(uint32(5), uint32(10)).Return(nil),
		flusher.EXPECT().Commit().Return(nil),
	)
	err = meta.Flush()
	assert.NoError(t, err)
	m.rwMutex.Lock()
	assert.Nil(t, m.immutable)
	m.rwMutex.Unlock()
}

func mockTagMetadata(ctrl *gomock.Controller) (TagMetadata, *kv.MockFamily, *version.MockSnapshot) {
	family := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	return NewTagMetadata(family), family, snapshot
}

func mockTagMetadataMemData(meta TagMetadata) {
	m := meta.(*tagMetadata)
	m.rwMutex.Lock()

	m.immutable = NewTagStore()
	tagEntry := newTagEntry(10)
	tagEntry.addTagValue("tag-value-5", 10)
	m.immutable.Put(5, tagEntry)

	m.mutable = NewTagStore()
	tagEntry = newTagEntry(20)
	tagEntry.addTagValue("tag-value-20", 20)
	m.mutable.Put(10, tagEntry)

	m.rwMutex.Unlock()
}

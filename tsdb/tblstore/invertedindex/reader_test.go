package invertedindex

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

var bitmapUnmarshal = encoding.BitmapUnmarshal

func buildInvertedIndexBlock() (zoneBlock []byte, ipBlock []byte, hostBlock []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher := NewFlusher(nopKVFlusher)
	zoneMapping := map[uint32]*roaring.Bitmap{
		1: roaring.BitmapOf(1),
		2: roaring.BitmapOf(2),
		3: roaring.BitmapOf(3),
	}
	hostMapping := map[uint32]*roaring.Bitmap{
		1: roaring.BitmapOf(1),
		2: roaring.BitmapOf(2),
		3: roaring.BitmapOf(3),
		4: roaring.BitmapOf(4),
		5: roaring.BitmapOf(5),
		6: roaring.BitmapOf(6),
		7: roaring.BitmapOf(7),
		8: roaring.BitmapOf(8),
		9: roaring.BitmapOf(9),
	}
	flush := func(tagValueIDs []uint32, mapping map[uint32]*roaring.Bitmap) {
		for _, tagValueID := range tagValueIDs {
			_ = seriesFlusher.FlushInvertedIndex(tagValueID, mapping[tagValueID])
		}
	}
	/////////////////////////
	// flush zone tag, tagID: 20
	/////////////////////////
	_ = seriesFlusher.FlushInvertedIndex(0, roaring.BitmapOf(1, 2, 3))
	flush([]uint32{1, 2, 3}, zoneMapping)
	// pick the zoneBlock buffer
	_ = seriesFlusher.FlushTagKeyID(20)
	zoneBlock = append(zoneBlock, nopKVFlusher.Bytes()...)

	/////////////////////////
	// flush ip tag, tagID: 21
	/////////////////////////
	//flush(ipMapping)
	_ = seriesFlusher.FlushInvertedIndex(0, roaring.BitmapOf(1, 2, 3, 4000000, 5000000, 6000000, 7000000, 8000000, 9000000))
	_ = seriesFlusher.FlushInvertedIndex(1, roaring.BitmapOf(1))
	_ = seriesFlusher.FlushInvertedIndex(2, roaring.BitmapOf(2))
	_ = seriesFlusher.FlushInvertedIndex(3, roaring.BitmapOf(3))
	_ = seriesFlusher.FlushInvertedIndex(4000000, roaring.BitmapOf(4000000))
	_ = seriesFlusher.FlushInvertedIndex(5000000, roaring.BitmapOf(5000000))
	_ = seriesFlusher.FlushInvertedIndex(6000000, roaring.BitmapOf(6000000))
	_ = seriesFlusher.FlushInvertedIndex(7000000, roaring.BitmapOf(7000000))
	_ = seriesFlusher.FlushInvertedIndex(8000000, roaring.BitmapOf(8000000))
	_ = seriesFlusher.FlushInvertedIndex(9000000, roaring.BitmapOf(9000000))

	// pick the ipBlock buffer
	_ = seriesFlusher.FlushTagKeyID(21)
	ipBlock = append(ipBlock, nopKVFlusher.Bytes()...)

	/////////////////////////
	// flush host tag, tagID: 22
	/////////////////////////
	_ = seriesFlusher.FlushInvertedIndex(0, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9))
	flush([]uint32{1, 2, 3, 4, 5, 6, 7, 8, 9}, hostMapping)
	// pick the hostBlock buffer
	_ = seriesFlusher.FlushTagKeyID(22)
	hostBlock = append(hostBlock, nopKVFlusher.Bytes()...)
	return zoneBlock, ipBlock, hostBlock
}

func buildInvertedIndexReader(ctrl *gomock.Controller) Reader {
	zoneBlock, ipBlock, hostBlock := buildInvertedIndexBlock()
	// mock readers
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(10)).Return(nil, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(19)).Return(nil, false).AnyTimes()
	mockReader.EXPECT().Get(uint32(20)).Return(zoneBlock, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(21)).Return(ipBlock, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(22)).Return(hostBlock, true).AnyTimes()
	// build series index reader
	return NewReader([]table.Reader{mockReader})
}

func TestReader_FindValueIDsForTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := buildInvertedIndexReader(ctrl)
	// read not tagID key
	idSet, err := reader.GetSeriesIDsForTagKeyID(19)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	idSet, err = reader.GetSeriesIDsForTagKeyID(10)
	assert.Error(t, err)
	assert.Nil(t, idSet)

	// read zone block
	idSet, err = reader.GetSeriesIDsForTagKeyID(21)
	a := roaring.BitmapOf(1, 2, 3, 4000000, 5000000, 6000000, 7000000, 8000000, 9000000)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ToArray(), idSet.ToArray())
}

func TestReader_FindSeriesIDsByTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := buildInvertedIndexReader(ctrl)
	// read not tag key id
	idSet, err := reader.FindSeriesIDsByTagValueIDs(19, roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	// tag value ids is empty
	idSet, err = reader.FindSeriesIDsByTagValueIDs(19, nil)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	// not found
	idSet, err = reader.FindSeriesIDsByTagValueIDs(10, roaring.BitmapOf(1))
	assert.Error(t, err)
	assert.Nil(t, idSet)

	// read zone block
	idSet, err = reader.FindSeriesIDsByTagValueIDs(21, roaring.BitmapOf(2, 49, 6000000, 6000033, 7000000))
	a := roaring.BitmapOf(2, 6000000, 7000000)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ToArray(), idSet.ToArray())

	idSet, err = reader.FindSeriesIDsByTagValueIDs(20, roaring.BitmapOf(1, 2))
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1, 2).ToArray(), idSet.ToArray())
}

func TestReader_InvertedIndex_reader_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		ctrl.Finish()
	}()

	zoneBlock, _, _ := buildInvertedIndexBlock()
	reader := newInvertedIndexReader(zoneBlock)

	idSet, err := reader.findSeriesIDsByTagValueIDs(roaring.BitmapOf(1, 2))
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1, 2).ToArray(), idSet.ToArray())
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	// case 1: unmarshal series id err
	idSet, err = reader.findSeriesIDsByTagValueIDs(roaring.BitmapOf(1, 2))
	assert.Error(t, err)
	assert.Nil(t, idSet)
	// case 2: init inverted reader err
	reader = newInvertedIndexReader(zoneBlock)
	idSet, err = reader.findSeriesIDsByTagValueIDs(roaring.BitmapOf(1, 2))
	assert.Error(t, err)
	assert.Nil(t, idSet)
	// case 3: validation offset err
	reader = newInvertedIndexReader([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	idSet, err = reader.findSeriesIDsByTagValueIDs(roaring.BitmapOf(1, 2))
	assert.Error(t, err)
	assert.Nil(t, idSet)
}

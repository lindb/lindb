package invertedindex

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestInvertedMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		ctrl.Finish()
	}()

	merge := NewInvertedMerger()
	// case 1: merge data success
	data, err := merge.Merge(1, mockInvertedMergeData())
	assert.NoError(t, err)
	reader, err := newTagInvertedReader(data)
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8000000, 9000000).ToArray(), reader.keys.ToArray())
	seriesIDs, _ := reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(1))
	assert.EqualValues(t, roaring.BitmapOf(1, 10).ToArray(), seriesIDs.ToArray())
	seriesIDs, _ = reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(2))
	assert.EqualValues(t, roaring.BitmapOf(2).ToArray(), seriesIDs.ToArray())
	seriesIDs, _ = reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(8000000))
	assert.EqualValues(t, roaring.BitmapOf(8000000).ToArray(), seriesIDs.ToArray())
	// case 2: new reader err
	data, err = merge.Merge(1, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 3: flush tag value data err
	flusher := NewMockInvertedFlusher(ctrl)
	m := merge.(*invertedMerger)
	m.invertedFlusher = flusher
	flusher.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	data, err = merge.Merge(1, mockInvertedMergeData())
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 4: flush tag data err
	flusher.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	flusher.EXPECT().FlushTagKeyID(gomock.Any()).Return(fmt.Errorf("err"))
	data, err = merge.Merge(1, mockInvertedMergeData())
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 5: scan data err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		d, _ := roaring.BitmapOf(1).ToBytes()
		if reflect.DeepEqual(d, data[:len(d)]) {
			// mock get series ids data err
			return fmt.Errorf("err")
		}
		// for other unmarshal
		return bitmap.UnmarshalBinary(data)
	}
	data, err = merge.Merge(1, mockInvertedMergeData())
	assert.Error(t, err)
	assert.Nil(t, data)
}

func mockInvertedMergeData() (data [][]byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher := NewInvertedFlusher(nopKVFlusher)
	zoneMapping := map[uint32]*roaring.Bitmap{
		1: roaring.BitmapOf(1),
		2: roaring.BitmapOf(2),
		3: roaring.BitmapOf(3),
	}
	hostMapping := map[uint32]*roaring.Bitmap{
		1:       roaring.BitmapOf(10),
		3:       roaring.BitmapOf(30),
		4:       roaring.BitmapOf(4),
		5:       roaring.BitmapOf(5),
		6:       roaring.BitmapOf(6),
		7:       roaring.BitmapOf(7),
		8000000: roaring.BitmapOf(8000000),
		9000000: roaring.BitmapOf(9000000),
	}
	flush := func(tagValueIDs []uint32, mapping map[uint32]*roaring.Bitmap) {
		for _, tagValueID := range tagValueIDs {
			_ = seriesFlusher.FlushInvertedIndex(tagValueID, mapping[tagValueID])
		}
	}
	flush([]uint32{1, 2, 3}, zoneMapping)
	_ = seriesFlusher.FlushTagKeyID(20)
	data = append(data, nopKVFlusher.Bytes())
	flush([]uint32{1, 3, 4, 5, 6, 7, 8000000, 9000000}, hostMapping)
	// pick the hostBlock buffer
	_ = seriesFlusher.FlushTagKeyID(22)
	data = append(data, nopKVFlusher.Bytes())
	return data
}

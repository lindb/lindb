package forwardindex

import (
	"fmt"
	"math"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Reader(t *testing.T) {
	reader := NewReader(nil)
	assert.NotNil(t, reader)
}

func buildForwardIndexBlock() []byte {
	var (
		ipMapping   = make(map[string]uint32)
		zoneMapping = make(map[string]*roaring.Bitmap)
		hostMapping = make(map[string]uint32)
	)
	for x := 0; x < math.MaxUint8; x++ {
		for y := 0; y < math.MaxUint8; y++ {
			// build ip
			seriesID := uint32(x*math.MaxUint8 + y)
			ip := fmt.Sprintf("192.168.%d.%d", x, y)
			ipMapping[ip] = seriesID
			// build zone
			var thisZone string
			switch x % 3 {
			case 0:
				thisZone = "nj"
			case 1:
				thisZone = "sh"
			default:
				thisZone = "bj"
			}
			bitmap, ok := zoneMapping[thisZone]
			if !ok {
				bitmap = roaring.NewBitmap()
			}
			bitmap.Add(seriesID)
			zoneMapping[thisZone] = bitmap
			// build host
			host := fmt.Sprintf("lindb-test-%s-%d", thisZone, seriesID)
			hostMapping[host] = seriesID
		}
	}

	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	for v := 0; v < 3; v++ {
		// flush tag ip
		for ip, seriesID := range ipMapping {
			if seriesID < 10000 {
				flusher.FlushTagValue(ip, roaring.BitmapOf(seriesID))
			}
		}
		flusher.FlushTagKey("ip")
		// flush tag zone
		for zone, bitmap := range zoneMapping {
			flusher.FlushTagValue(zone, bitmap)
		}
		flusher.FlushTagKey("zone")
		// flush tag host
		for host, seriesID := range hostMapping {
			bitmap := roaring.NewBitmap()
			bitmap.Add(seriesID)
			flusher.FlushTagValue(host, bitmap)
		}
		flusher.FlushTagKey("host")
		// flush version
		flusher.FlushVersion(series.Version(v), timeutil.TimeRange{Start: int64(v) * 100, End: (int64(v) + 1) * 100})
	}

	_ = flusher.FlushMetricID(1)
	data := nopKVFlusher.Bytes()
	return data
}

func buildForwardIndexReader(ctrl *gomock.Controller) *reader {
	data := buildForwardIndexBlock()
	// build mock reader
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(0)).Return(nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(1)).Return(data).AnyTimes()
	mockReaders := []table.Reader{mockReader}
	// build forward index reader
	indexReader := NewReader(mockReaders)
	return indexReader.(*reader)
}

func Test_ForwardIndexReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// build forward index reader
	indexReader := buildForwardIndexReader(ctrl)

	// test inexist version
	seriesID2TagValues, err := indexReader.GetTagValues(
		1,
		[]string{"host", "zone"},
		4,
		roaring.BitmapOf(1, 2, 3))
	assert.Len(t, seriesID2TagValues, 0)
	assert.NotNil(t, err)

	// test inexist metricID
	seriesID2TagValues, err = indexReader.GetTagValues(
		0,
		[]string{"host", "zone"},
		2,
		roaring.BitmapOf(1, 2, 3))
	assert.Len(t, seriesID2TagValues, 0)
	assert.NotNil(t, err)

	// test inexist tagKeys
	_, err = indexReader.GetTagValues(
		1,
		[]string{"notexisttag1", "notexisttag2"},
		2,
		roaring.BitmapOf(1, 2, 3))
	assert.NotNil(t, err)

	// test no keys
	_, err = indexReader.GetTagValues(1, nil, 2, roaring.BitmapOf(1, 2, 3))
	assert.NotNil(t, err)

	// test existed tagKeys
	seriesID2TagValues, err = indexReader.GetTagValues(
		1, []string{"host", "zone"}, 2, roaring.BitmapOf(1, 501, 1002, 999999999))
	assert.Nil(t, err)
	assert.NotNil(t, seriesID2TagValues)
	assert.Len(t, seriesID2TagValues, 3)
	assert.Equal(t, []string{"lindb-test-nj-1", "nj"}, seriesID2TagValues[1])
	assert.Equal(t, []string{"lindb-test-sh-501", "sh"}, seriesID2TagValues[501])
	assert.Equal(t, []string{"lindb-test-nj-1002", "nj"}, seriesID2TagValues[1002])
	// test empty tagKeys
	seriesID2TagValues, err = indexReader.GetTagValues(
		1, []string{"host", "ip", "zone"}, 2, roaring.BitmapOf(9999, 10000, 10001))
	assert.Nil(t, err)
	assert.NotNil(t, seriesID2TagValues)
	assert.Len(t, seriesID2TagValues, 3)
	assert.Equal(t, []string{"lindb-test-nj-9999", "192.168.39.54", "nj"}, seriesID2TagValues[9999])
	assert.Equal(t, []string{"lindb-test-nj-10000", "", "nj"}, seriesID2TagValues[10000])
	assert.Equal(t, []string{"lindb-test-nj-10001", "", "nj"}, seriesID2TagValues[10001])
}

func Test_forwardIndexVersionEntry_errorCases(t *testing.T) {

	// read footer error
	entry, err := newForwardIndexVersionEntry(nil)
	assert.NotNil(t, err)
	assert.Nil(t, entry)
	// red footer, position check failure
	block := []byte{
		1, 2, 3, 4, 5, 6, 7, 8, // time range
		0,                                  // tagKey
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // footer
	}
	entry, err = newForwardIndexVersionEntry(block)
	assert.NotNil(t, err)
	assert.Nil(t, entry)
	// read tagKeys error
	block = []byte{
		1, 2, 3, 4, 5, 6, 7, 8, // time range
		100,                                // tagKey
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // footer
	}
	entry, err = newForwardIndexVersionEntry(block)
	assert.NotNil(t, err)
	assert.Nil(t, entry)
	// unmarshal bitmap error
	block = []byte{
		1, 2, 3, 4, 5, 6, 7, 8, // time range
		0,                                  // tagKey
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // footer
	}
	entry, err = newForwardIndexVersionEntry(block)
	assert.NotNil(t, err)
	assert.Nil(t, entry)
}

package tblstore

import (
	"fmt"
	"math"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewForwardIndexReader(t *testing.T) {
	reader := NewForwardIndexReader(nil)
	assert.NotNil(t, reader)
}

func buildForwardIndexBlock(ctrl *gomock.Controller) []byte {
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
	mockKVFlusher := kv.NewMockFlusher(ctrl)
	mockKVFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	flusher := NewForwardIndexFlusher(mockKVFlusher)
	flusherImpl := flusher.(*forwardIndexFlusher)
	for version := 0; version < 3; version++ {
		// flush tag ip
		for ip, seriesID := range ipMapping {
			bitmap := roaring.NewBitmap()
			bitmap.Add(seriesID)
			flusher.FlushTagValue(ip, bitmap)
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
		flusher.FlushVersion(uint32(version), uint32(version*100), uint32(version+1)*100)
	}
	flusherImpl.resetDisabled = true
	_ = flusher.FlushMetricID(1)
	data, _ := flusherImpl.metricBlockWriter.Bytes()
	return data
}

func buildForwardIndexReader(ctrl *gomock.Controller) *forwardIndexReader {
	data := buildForwardIndexBlock(ctrl)
	// build mock reader
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(0)).Return(nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(1)).Return(data).AnyTimes()
	// build mock snapshot
	mockSnapShot := version.NewMockSnapshot(ctrl)
	mockSnapShot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{mockReader}, nil).AnyTimes()
	// build forward index reader
	indexReader := NewForwardIndexReader(mockSnapShot)
	return indexReader.(*forwardIndexReader)
}

func Test_ForwardIndexReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// build forward index reader
	indexReader := buildForwardIndexReader(ctrl)
	// test not exist metricID
	tagValues, err := indexReader.GetTagValues(0, []string{"host", "zone"}, 2)
	assert.Len(t, tagValues, 0)
	assert.NotNil(t, err)
	// test not exist tagKeys
	_, err = indexReader.GetTagValues(1, []string{"notexisttag1", "notexisttag2"}, 2)
	assert.NotNil(t, err)
	// test no keys
	_, err = indexReader.GetTagValues(1, nil, 2)
	assert.Nil(t, err)
	// test existed tagKeys
	tagValues, err = indexReader.GetTagValues(1, []string{"host", "not-exist-key", "zone"}, 2)
	assert.Len(t, tagValues, 3)
	assert.Len(t, tagValues[0], 65025)
	assert.Len(t, tagValues[1], 0)
	assert.Len(t, tagValues[2], 3)
	assert.Nil(t, err)
}

func Test_ForwardIndexReader_errorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// build forward index reader
	indexReader := buildForwardIndexReader(ctrl)
	// versionBlock is invalid
	_, _, err := indexReader.readKeysLUTBlock(nil, nil)
	assert.NotNil(t, err)
	err = indexReader.readDictBlockByIndexes(nil, nil)
	assert.NotNil(t, err)
	// tagKeysBlock corrupt
	_, err = indexReader.readTagKeysBlock([]byte{1, 2, 3, 4, 5, 6, 7, 8, 2}) // timeRange + count
	assert.NotNil(t, err)
	// index string block failure
	err = indexReader.readStringBlockByOffsets(nil, []int{1, 2}, []int{1, 3}, []int{1, 2, 3})
	assert.NotNil(t, err)
	// index cannot be found in dict block
	var strIndexes []int
	for i := 0; i < 1000; i++ {
		strIndexes = append(strIndexes, i)
	}
	err = indexReader.readStringBlockByOffsets(nil, nil, nil, strIndexes)
	assert.NotNil(t, err)
}

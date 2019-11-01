package invertedindex

import (
	"math"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func buildInvertedIndexBlock() (zoneBlock []byte, ipBlock []byte, hostBlock []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher := NewFlusher(nopKVFlusher)
	// disable auto reset to pick the entrySetBuffer
	/////////////////////////
	// seriesID mapping relation
	/////////////////////////
	ipMapping := map[uint32]string{
		1: "192.168.1.1",
		2: "192.168.1.2",
		3: "192.168.1.3",
		4: "192.168.2.4",
		5: "192.168.2.5",
		6: "192.168.2.6",
		7: "192.168.3.7",
		8: "192.168.3.8",
		9: "192.168.3.9"}
	zoneMapping := map[string][]uint32{
		"nj": {1, 2, 3},
		"sh": {4, 5, 6},
		"bj": {7, 8, 9}}
	hostMapping := map[uint32]string{
		1: "eleme-dev-nj-1",
		2: "eleme-dev-nj-2",
		3: "eleme-dev-nj-3",
		4: "eleme-dev-sh-4",
		5: "eleme-dev-sh-5",
		6: "eleme-dev-sh-6",
		7: "eleme-dev-bj-7",
		8: "eleme-dev-bj-8",
		9: "eleme-dev-bj-9"}
	/////////////////////////
	// flush zone tag, tagID: 20
	/////////////////////////
	for zone, ids := range zoneMapping {
		for v := series.Version(1500000000000); v < 1800000000000; v += 100000000000 {
			bitmap := roaring.New()
			bitmap.AddMany(ids)
			seriesFlusher.FlushVersion(v, timeutil.TimeRange{
				Start: v.Int64() + 10000*1000, End: v.Int64() + 20000*10000}, bitmap)
		}
		seriesFlusher.FlushTagValue(zone)
	}
	// pick the zoneBlock buffer
	_ = seriesFlusher.FlushTagKeyID(20)
	zoneBlock = append(zoneBlock, nopKVFlusher.Bytes()...)

	/////////////////////////
	// flush ip tag, tagID: 21
	/////////////////////////
	for seriesID, ip := range ipMapping {
		for v := series.Version(1500000000000); v < 1800000000000; v += 100000000000 {
			bitmap := roaring.New()
			bitmap.Add(seriesID)
			seriesFlusher.FlushVersion(v, timeutil.TimeRange{
				Start: v.Int64() + 10000*1000, End: v.Int64() + 20000*10000}, bitmap)
		}
		seriesFlusher.FlushTagValue(ip)
	}
	// pick the ipBlock buffer
	_ = seriesFlusher.FlushTagKeyID(21)
	ipBlock = append(ipBlock, nopKVFlusher.Bytes()...)

	/////////////////////////
	// flush host tag, tagID: 22
	/////////////////////////
	for seriesID, host := range hostMapping {
		for v := series.Version(1500000000000); v < 1800000000000; v += 100000000000 {
			bitmap := roaring.New()
			bitmap.Add(seriesID)
			seriesFlusher.FlushVersion(v, timeutil.TimeRange{
				Start: v.Int64() + 10000*1000, End: v.Int64() + 20000*10000}, bitmap)
		}
		seriesFlusher.FlushTagValue(host)
	}
	// pick the hostBlock buffer
	_ = seriesFlusher.FlushTagKeyID(22)
	hostBlock = append(hostBlock, nopKVFlusher.Bytes()...)
	return zoneBlock, ipBlock, hostBlock
}

func buildSeriesIndexReader(ctrl *gomock.Controller) Reader {
	zoneBlock, ipBlock, hostBlock := buildInvertedIndexBlock()
	// mock readers
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(19)).Return(nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(20)).Return(zoneBlock).AnyTimes()
	mockReader.EXPECT().Get(uint32(21)).Return(ipBlock).AnyTimes()
	mockReader.EXPECT().Get(uint32(22)).Return(hostBlock).AnyTimes()
	// build series index reader
	return NewReader([]table.Reader{mockReader})
}

func Test_InvertedIndexReader_GetSeriesIDsForTagID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := buildSeriesIndexReader(ctrl)
	// read not tagID key
	idSet, err := reader.GetSeriesIDsForTagKeyID(19, timeutil.TimeRange{})
	assert.NotNil(t, err)
	assert.Nil(t, idSet)

	// read zone block but not overlaps
	idSet, err = reader.GetSeriesIDsForTagKeyID(20,
		timeutil.TimeRange{
			Start: 1400000000 * 1000,
			End:   1500000000 * 1000})
	assert.NotNil(t, err)
	assert.Nil(t, idSet)

	// read zone block, overlaps
	idSet, err = reader.GetSeriesIDsForTagKeyID(20,
		timeutil.TimeRange{
			Start: 1500000000 * 1000,
			End:   1600000000 * 1000})
	assert.Nil(t, err)
	assert.NotNil(t, idSet)
	assert.Contains(t, idSet.Versions(), series.Version(1500000000000))
	assert.Equal(t, uint32(1), idSet.Versions()[series.Version(1500000000000)].Minimum())
	assert.Equal(t, uint32(9), idSet.Versions()[series.Version(1500000000000)].Maximum())
}

func Test_intSliceContains(t *testing.T) {
	assert.False(t, intSliceContains(nil, 1))
	assert.False(t, intSliceContains([]int{1, 3, 4, 5, 8}, 0))
	assert.True(t, intSliceContains([]int{1, 3, 4, 5, 8}, 1))
	assert.False(t, intSliceContains([]int{1, 3, 4, 5, 8}, 2))
	assert.True(t, intSliceContains([]int{1, 3, 4, 5, 8}, 3))
	assert.True(t, intSliceContains([]int{1, 3, 4, 5, 8}, 8))
	assert.False(t, intSliceContains([]int{1, 3, 4, 5, 8}, 9))
}

func Test_InvertedIndexReader_FindSeriesIDsByExprForTagID_badCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := buildSeriesIndexReader(ctrl)

	// tagID not exist
	idSet, err := reader.FindSeriesIDsByExprForTagKeyID(19, nil, timeutil.TimeRange{})
	assert.NotNil(t, err)
	assert.Nil(t, idSet)

	// find zone with bad expression
	idSet, err = reader.FindSeriesIDsByExprForTagKeyID(20, nil,
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.NotNil(t, err)
	assert.Nil(t, idSet)
	// find zone with bad time range
	idSet, err = reader.FindSeriesIDsByExprForTagKeyID(20, nil,
		timeutil.TimeRange{Start: 12 * 1000, End: 13 * 1000})
	assert.NotNil(t, err)
	assert.Nil(t, idSet)
}
func Test_InvertedIndexReader_FindSeriesIDsByExprForTagID_EqualExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := buildSeriesIndexReader(ctrl)

	idSet, err := reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.EqualsExpr{Key: "host", Value: "eleme-dev-sh-4"},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.Nil(t, err)
	assert.Contains(t, idSet.Versions(), series.Version(1500000000000))
	assert.Equal(t, uint64(1), idSet.Versions()[1500000000000].GetCardinality())
	assert.True(t, idSet.Versions()[series.Version(1500000000000)].Contains(4))
	// find not existed host
	_, err = reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.EqualsExpr{Key: "host", Value: "eleme-dev-sh-41"},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.NotNil(t, err)
}

func Test_InvertedIndexReader_FindSeriesIDsByExprForTagID_InExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := buildSeriesIndexReader(ctrl)

	// find existed host
	idSet, err := reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.InExpr{
		Key: "host", Values: []string{"eleme-dev-sh-4", "eleme-dev-sh-5", "eleme-dev-sh-55"}},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.Nil(t, err)
	assert.Contains(t, idSet.Versions(), series.Version(1500000000000))
	assert.Equal(t, uint64(2), idSet.Versions()[1500000000000].GetCardinality())
	assert.True(t, idSet.Versions()[series.Version(1500000000000)].Contains(4))
	assert.True(t, idSet.Versions()[series.Version(1500000000000)].Contains(5))
	// find not existed host
	_, err = reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.InExpr{
		Key: "host", Values: []string{"eleme-dev-sh-55"}},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.NotNil(t, err)
}

func Test_InvertedIndexReader_FindSeriesIDsByExprForTagID_LikeExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := buildSeriesIndexReader(ctrl)

	// find existed host
	idSet, err := reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.LikeExpr{
		Key: "host", Value: "eleme-dev-sh-"},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.Nil(t, err)
	assert.Contains(t, idSet.Versions(), series.Version(1500000000000))
	assert.Equal(t, uint64(3), idSet.Versions()[series.Version(1500000000000)].GetCardinality())
	assert.Equal(t, uint32(4), idSet.Versions()[series.Version(1500000000000)].Minimum())
	assert.Equal(t, uint32(6), idSet.Versions()[series.Version(1500000000000)].Maximum())
	// find not existed host
	_, err = reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.InExpr{
		Key: "host", Values: []string{"eleme-dev-sh---"}},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.NotNil(t, err)
}

func Test_InvertedIndexReader_FindSeriesIDsByExprForTagID_RegexExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := buildSeriesIndexReader(ctrl)

	_, err := reader.FindSeriesIDsByExprForTagKeyID(22, &stmt.RegexExpr{
		Key: "host", Regexp: "eleme-dev-sh-"},
		timeutil.TimeRange{Start: 1500000000 * 1000, End: 1600000000 * 1000})
	assert.Nil(t, err)
}

func Test_newTagKVEntrySet_error_cases(t *testing.T) {
	// block length too short, 8 bytes
	_, err := newTagKVEntrySet([]byte{16, 86, 104, 89, 32, 63, 84, 101})
	assert.NotNil(t, err)
	// validate offsets failure
	_, err = newTagKVEntrySet([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	assert.NotNil(t, err)
}

func Test_InvertedIndexReader_entrySetToIDSet_error_cases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	readerIntf := buildSeriesIndexReader(ctrl)
	readerImpl := readerIntf.(*reader)

	zoneBlock, _, _ := buildInvertedIndexBlock()
	zoneEntrySet, _ := newTagKVEntrySet(zoneBlock)
	// first offset not exist
	idSet, err := readerImpl.entrySetToIDSet(
		zoneEntrySet,
		timeutil.TimeRange{Start: 0, End: math.MaxInt64},
		[]int{1000, 1200})
	assert.Nil(t, idSet)
	assert.NotNil(t, err)
	// last offset not exist
	idSet, err = readerImpl.entrySetToIDSet(
		zoneEntrySet,
		timeutil.TimeRange{Start: 0, End: math.MaxInt64},
		[]int{0, 1200})
	assert.Nil(t, idSet)
	assert.NotNil(t, err)
}

func Test_tagKVEntrySet_TrieTree_error_cases(t *testing.T) {
	zoneBlock, _, _ := buildInvertedIndexBlock()
	entrySet, _ := newTagKVEntrySet(zoneBlock)

	// read stream eof
	entrySet.sr.Reset([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 1, 1, 1, 1})
	// read stream eof
	_, err := entrySet.TrieTree()
	assert.NotNil(t, err)

	// failed validation of trie tree
	entrySet.sr.Reset([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 1, 1, 1, 1, 1, 1})
	_, err = entrySet.TrieTree()
	assert.NotNil(t, err)

	// LOUDS block unmarshal failed
	entrySet.sr.Reset([]byte{1, 2, 3, 4, 5, 6, 7, 8, 6, 1, 1, 1, 1, 1, 1})
	_, err = entrySet.TrieTree()
	assert.NotNil(t, err)

	// isPrefixKey block unmarshal failed
	out, _ := NewRankSelect().MarshalBinary()
	badBLOCK := append([]byte{1, 2, 3, 4, 5, 6, 7, 8,
		18,   // trie tree length
		1, 1, // labels
		1, 1, // is prefix
		13}) // louds

	badBLOCK = append(badBLOCK, out...) // LOUDS block
	entrySet.sr.Reset(badBLOCK)
	_, err = entrySet.TrieTree()
	assert.NotNil(t, err)
}

func Test_tagKVEntrySet_Bitmap(t *testing.T) {
	data := versionedTagValueData{}
	_, err := data.Bitmap()
	assert.NotNil(t, err)
}

func Test_InvertedIndexReader_SuggestTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := buildSeriesIndexReader(ctrl)

	// tagID not exist
	assert.Nil(t, reader.SuggestTagValues(19, "", 10000000))
	// search ip
	assert.Len(t, reader.SuggestTagValues(21, "192", 1000), 9)
	assert.Len(t, reader.SuggestTagValues(21, "192", 3), 3)

	// mock corruption
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(18)).Return([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}).AnyTimes()
	corruptedReader := NewReader([]table.Reader{mockReader})
	assert.Nil(t, corruptedReader.SuggestTagValues(18, "", 10000000))
}

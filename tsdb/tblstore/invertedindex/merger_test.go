package invertedindex

//
//import (
//	"sort"
//	"testing"
//	"time"
//
//	"github.com/lindb/roaring"
//	"github.com/stretchr/testify/assert"
//
//	"github.com/lindb/lindb/kv"
//	"github.com/lindb/lindb/pkg/timeutil"
//	"github.com/lindb/lindb/series"
//)
//
//func buildInvertedIndexBlockToCompact() (data [][]byte) {
//	nopKVFlusher := kv.NewNopFlusher()
//	now := timeutil.Now()
//
//	// ttl: 30 day
//	tagValueVersions1 := map[string]map[series.Version]*roaring.Bitmap{
//		"192.168.1.1": {
//			series.Version(now - 25*timeutil.OneDay): roaring.BitmapOf(1, 2, 3),
//			series.Version(now - 35*timeutil.OneDay): roaring.BitmapOf(1, 2, 3)},
//		"192.168.1.2": {
//			series.Version(now - 12*timeutil.OneDay): roaring.BitmapOf(1, 2, 3),
//			series.Version(now - 32*timeutil.OneDay): roaring.BitmapOf(1, 2, 3)},
//		"192.168.1.3": {
//			series.Version(now - 35*timeutil.OneDay): roaring.BitmapOf(1, 2, 3),
//			series.Version(now - 36*timeutil.OneDay): roaring.BitmapOf(1, 2, 3)}}
//	tagValueVersions2 := map[string]map[series.Version]*roaring.Bitmap{
//		"192.168.1.1": {
//			series.Version(now - 15*timeutil.OneDay): roaring.BitmapOf(1, 2, 3),
//			series.Version(now - 35*timeutil.OneDay): roaring.BitmapOf(1, 2, 3)},
//		"192.168.1.3": {
//			series.Version(now - 20*timeutil.OneDay): roaring.BitmapOf(1, 2, 3)},
//		"192.168.1.4": {
//			series.Version(now - 36*timeutil.OneDay): roaring.BitmapOf(1, 2, 3)}}
//	getFlushedData := func(tagValueVersions map[string]map[series.Version]*roaring.Bitmap) []byte {
//		invertedFlusher := NewFlusher(nopKVFlusher).(*invertedFlusher)
//		for tagValue, versions := range tagValueVersions {
//			for version, bitmap := range versions {
//				invertedFlusher.FlushVersion(version, timeutil.TimeRange{Start: 1, End: 1}, bitmap)
//			}
//			invertedFlusher.FlushTagValue(tagValue)
//		}
//		_ = invertedFlusher.FlushTagKeyID(1)
//		return append([]byte{}, nopKVFlusher.Bytes()...)
//	}
//
//	data = append(data, getFlushedData(tagValueVersions1), getFlushedData(tagValueVersions2))
//	return data
//}
//
//func Test_Merge_TTL_30Day(t *testing.T) {
//	m := NewMerger(time.Hour * 24 * 30).(*invertedIndexMerger)
//	compacted, err := m.Merge(1, buildInvertedIndexBlockToCompact())
//	assert.Nil(t, err)
//	assert.NotNil(t, compacted)
//	// convert to entrySet
//	entrySet, err := newTagKVEntrySet(compacted)
//	assert.Nil(t, err)
//	tree, err := entrySet.TrieTree()
//	assert.Nil(t, err)
//	tagValues := tree.PrefixSearch("", 10)
//	sort.Slice(tagValues, func(i, j int) bool { return tagValues[i] < tagValues[j] })
//	assert.Equal(t, []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}, tagValues)
//}
//
//func TestInvertedIndexMerger_Merge_TTL_10Day(t *testing.T) {
//	m := NewMerger(time.Hour * 24 * 10).(*invertedIndexMerger)
//	compacted, err := m.Merge(1, buildInvertedIndexBlockToCompact())
//	assert.Nil(t, err)
//	assert.NotNil(t, compacted)
//	// convert to entrySet
//	entrySet, err := newTagKVEntrySet(compacted)
//	assert.Nil(t, err)
//	tree, err := entrySet.TrieTree()
//	assert.Nil(t, err)
//	tagValues := tree.PrefixSearch("", 10)
//	assert.Equal(t, []string{"192.168.1.2"}, tagValues)
//}
//
//func TestInvertedIndexMerger_Merge_BadBlock(t *testing.T) {
//	m := NewMerger(time.Hour * 24 * 10).(*invertedIndexMerger)
//	compacted, err := m.Merge(1, [][]byte{{1, 2, 3, 4}, {1, 2, 3, 4}})
//	assert.NotNil(t, err)
//	assert.Nil(t, compacted)
//}

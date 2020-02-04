package invertedindex

import (
	"time"

	"github.com/lindb/lindb/kv"
)

type invertedIndexMerger struct {
	flusher      *tagFlusher
	reader       *tagReader
	nopKVFlusher *kv.NopFlusher
	ttl          time.Duration
}

func NewMerger(ttl time.Duration) kv.Merger {
	nopKVFlusher := kv.NewNopFlusher()
	return &invertedIndexMerger{
		flusher:      NewTagFlusher(nopKVFlusher).(*tagFlusher),
		reader:       NewTagReader(nil).(*tagReader),
		nopKVFlusher: nopKVFlusher,
		ttl:          ttl}
}

func (m *invertedIndexMerger) reset() {
	m.flusher.reset()
}

func (m *invertedIndexMerger) Merge(
	key uint32,
	value [][]byte,
) (
	[]byte,
	error,
) {
	//defer m.reset()
	//var (
	//	tagValueData = make(map[string][]*versionedTagValueData)
	//)
	//// extract
	//for _, block := range value {
	//	var (
	//		offsetTagValues = make(map[int]string)
	//	)
	//	entrySet, err := newTagKVEntrySet(block)
	//	if err != nil {
	//		return nil, err
	//	}
	//	tree, err := entrySet.TrieTree()
	//	if err != nil {
	//		return nil, err
	//	}
	//	// read offsets
	//	itr := tree.Iterator("")
	//	for itr.HasNext() {
	//		value, offset := itr.Next()
	//		offsetTagValues[offset] = string(value)
	//		dataItr, err := entrySet.ReadTagValueDataBlock(offset)
	//		if err != nil {
	//			return nil, err
	//		}
	//		tagValue := offsetTagValues[offset]
	//
	//		dataUnionList, ok := tagValueData[tagValue]
	//		if !ok {
	//			dataUnionList = []*versionedTagValueData{}
	//		}
	//		for dataItr.HasNext() {
	//			dataUnionList = append(dataUnionList, &versionedTagValueData{
	//				version:   dataItr.DataVersion(),
	//				timeRange: dataItr.DataTimeRange(),
	//				data:      dataItr.Next(),
	//			})
	//		}
	//		tagValueData[tagValue] = dataUnionList
	//	}
	//}
	//// do ttl
	//m.evictOldVersion(tagValueData)
	//// do flush
	//m.flush(tagValueData, key)
	//return m.nopKVFlusher.Bytes(), nil
	return nil, nil
}

//
//func (m *invertedIndexMerger) evictOldVersion(
//	tagValueData map[string][]*versionedTagValueData,
//) {
//	var latestVersion series.Version
//	// sort and pick the latestVersion
//	for _, dataList := range tagValueData {
//		// desc order
//		dataList := dataList
//		sort.Slice(dataList, func(i, j int) bool {
//			return dataList[i].version.After(dataList[j].version)
//		})
//		thisLatestVersion := (dataList)[0].version
//		if thisLatestVersion.After(latestVersion) {
//			latestVersion = thisLatestVersion
//		}
//	}
//	for tagValue, dataList := range tagValueData {
//		var lastAliveIndex = len(dataList)
//		for index, data := range dataList {
//			// expire
//			if data.version.IsExpired(m.ttl) {
//				lastAliveIndex = index
//				break
//			}
//		}
//		// case1: all versions expired, but it's the latest version in use
//		if lastAliveIndex == 0 && dataList[0].version.Equal(latestVersion) {
//			// remove expired versions
//			lastAliveIndex = 1
//		}
//		// delete all versions
//		if lastAliveIndex == 0 {
//			delete(tagValueData, tagValue)
//			continue
//		}
//		// remove expired versions
//		dataList = dataList[:lastAliveIndex]
//		tagValueData[tagValue] = dataList
//	}
//}

//func (m *invertedIndexMerger) flush(
//	tagValueData map[string][]*versionedTagValueData,
//	tagKeyID uint32,
//) {
//	for tagValue, dataList := range tagValueData {
//		for _, data := range dataList {
//			m.flusher.flushVersion(
//				data.version,
//				data.timeRange,
//				data.data)
//		}
//		m.flusher.FlushTagValue(tagValue,)
//	}
//	_ = m.flusher.FlushTagKeyID(tagKeyID)
//}
//
//type versionedTagValueData struct {
//	data []byte
//}

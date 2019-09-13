package tblstore

import (
	"sort"
	"time"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series"
)

type invertedIndexMerger struct {
	flusher      *invertedIndexFlusher
	reader       *invertedIndexReader
	nopKVFlusher *kv.NopFlusher
	ttl          time.Duration
}

func NewInvertedIndexMerger(ttl time.Duration) kv.Merger {
	nopKVFlusher := kv.NewNopFlusher()
	return &invertedIndexMerger{
		flusher:      NewInvertedIndexFlusher(nopKVFlusher).(*invertedIndexFlusher),
		reader:       NewInvertedIndexReader(nil).(*invertedIndexReader),
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
	defer m.reset()
	var (
		tagValueData = make(map[string]*[]versionedTagValueData)
	)
	// extract
	for _, block := range value {
		var (
			offsets         []int
			offsetTagValues = make(map[int]string)
		)
		entrySet, err := newTagKVEntrySet(block)
		if err != nil {
			return nil, err
		}
		tree, err := entrySet.TrieTree()
		if err != nil {
			return nil, err
		}
		// read offsets
		itr := tree.Iterator("")
		for itr.HasNext() {
			value, offset := itr.Next()
			offsetTagValues[offset] = value
			offsets = append(offsets, offset)
		}
		// read all positions
		offsetPositions, err := entrySet.OffsetsToPosition(offsets)
		if err != nil {
			return nil, err
		}
		for offset, pos := range offsetPositions {
			dataList, err := entrySet.ReadTagValueDataBlock(pos)
			if err != nil {
				return nil, err
			}
			tagValue := offsetTagValues[offset]

			dataUnionList, ok := tagValueData[tagValue]
			if ok {
				*dataUnionList = append(*dataUnionList, dataList...)
			} else {
				dataUnionList = &[]versionedTagValueData{}
				*dataUnionList = append(*dataUnionList, dataList...)
			}
			tagValueData[tagValue] = dataUnionList
		}
	}
	// do ttl
	m.evictOldVersion(tagValueData)
	// do flush
	m.flush(tagValueData, key)
	return m.nopKVFlusher.Bytes(), nil
}

func (m *invertedIndexMerger) evictOldVersion(
	tagValueData map[string]*[]versionedTagValueData,
) {
	var latestVersion series.Version
	// sort and pick the latestVersion
	for _, dataList := range tagValueData {
		// desc order
		dataList := dataList
		sort.Slice(*dataList, func(i, j int) bool {
			return (*dataList)[i].version.After((*dataList)[j].version)
		})
		thisLatestVersion := (*dataList)[0].version
		if thisLatestVersion.After(latestVersion) {
			latestVersion = thisLatestVersion
		}
	}
	for tagValue, dataList := range tagValueData {
		var lastAliveIndex = len(*dataList)
		for index, data := range *dataList {
			// expire
			if data.version.IsExpired(m.ttl) {
				lastAliveIndex = index
				break
			}
		}
		// case1: all versions expired, but it's the latest version in use
		if lastAliveIndex == 0 && (*dataList)[0].version.Equal(latestVersion) {
			// remove expired versions
			lastAliveIndex = 1
		}
		// delete all versions
		if lastAliveIndex == 0 {
			delete(tagValueData, tagValue)
			continue
		}
		// remove expired versions
		*dataList = (*dataList)[:lastAliveIndex]
	}
}

func (m *invertedIndexMerger) flush(
	tagValueData map[string]*[]versionedTagValueData,
	tagKeyID uint32,
) {
	for tagValue, dataList := range tagValueData {
		for _, data := range *dataList {
			timeRange := data.TimeRange()
			m.flusher.flushVersion(
				data.version,
				uint32(timeRange.Start/1000),
				uint32(timeRange.End/1000),
				data.bitMapData)
		}
		m.flusher.FlushTagValue(tagValue)
	}
	_ = m.flusher.FlushTagKeyID(tagKeyID)
}

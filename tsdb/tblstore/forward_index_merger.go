package tblstore

import (
	"fmt"
	"sort"
	"time"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series"
)

type forwardIndexMerger struct {
	flusher          *forwardIndexFlusher
	reader           *forwardIndexReader
	nopKVFlusher     *kv.NopFlusher
	versionBlocksMap map[series.Version][][]byte // version->List<VersionBlock>
	ttl              time.Duration
}

func NewForwardIndexMerger(ttl time.Duration) kv.Merger {
	m := &forwardIndexMerger{
		reader:           NewForwardIndexReader(nil).(*forwardIndexReader),
		nopKVFlusher:     kv.NewNopFlusher(),
		versionBlocksMap: make(map[series.Version][][]byte),
		ttl:              ttl,
	}
	m.flusher = NewForwardIndexFlusher(m.nopKVFlusher).(*forwardIndexFlusher)
	return m
}

func (m *forwardIndexMerger) Reset() {
	m.flusher.resetMetricBlockContext()
	for version := range m.versionBlocksMap {
		delete(m.versionBlocksMap, version)
	}
}

func (m *forwardIndexMerger) latestVersionBlock(version series.Version) []byte {
	list := m.versionBlocksMap[version]
	var longestBlock []byte
	for _, block := range list {
		if len(block) > len(longestBlock) {
			longestBlock = block
		}
	}
	return longestBlock
}

// AliveVersions deletes the expired versions
func (m *forwardIndexMerger) AliveVersions() (alive []series.Version) {
	var versions []series.Version
	// collect a sorted versions list
	for version := range m.versionBlocksMap {
		versions = append(versions, version)
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i] < versions[j] })

	var lastVersion series.Version
	for _, version := range versions {
		lastVersion = version
		if version.Time().Add(m.ttl).After(time.Now()) {
			alive = append(alive, version)
		}
	}
	if len(alive) == 0 {
		alive = append(alive, lastVersion)
	}
	return alive
}

func (m *forwardIndexMerger) Merge(
	key uint32,
	value [][]byte,
) (
	[]byte,
	error,
) {
	defer m.Reset()

	for _, block := range value {
		versionBlockItr := newForwardIndexVersionBlockIterator(block)
		for versionBlockItr.HasNext() {
			version, versionBlock := versionBlockItr.Next()
			if versionBlock == nil {
				continue
			}
			list, ok := m.versionBlocksMap[version]
			if ok {
				list = append(list, versionBlock)
			} else {
				list = [][]byte{versionBlock}
			}
			m.versionBlocksMap[version] = list
		}
	}
	if len(m.versionBlocksMap) == 0 {
		return nil, fmt.Errorf("no available blocks for compacting")
	}
	for _, version := range m.AliveVersions() {
		latestVersionBlock := m.latestVersionBlock(version)
		startPos := m.flusher.metricBlockWriter.Len()
		m.flusher.metricBlockWriter.PutBytes(latestVersionBlock)
		m.flusher.RecordVersionOffset(version, startPos)
	}
	_ = m.flusher.FlushMetricID(key)
	return m.nopKVFlusher.Bytes(), nil
}

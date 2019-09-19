package tblstore

import (
	"fmt"
	"sort"
	"time"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
)

type forwardIndexMerger struct {
	flusher      *forwardIndexFlusher
	reader       *forwardIndexReader
	nopKVFlusher *kv.NopFlusher
	ttl          time.Duration
	sr           *stream.Reader
}

func NewForwardIndexMerger(ttl time.Duration) kv.Merger {
	nopKVFlusher := kv.NewNopFlusher()
	return &forwardIndexMerger{
		reader:       NewForwardIndexReader(nil).(*forwardIndexReader),
		nopKVFlusher: nopKVFlusher,
		flusher:      NewForwardIndexFlusher(nopKVFlusher).(*forwardIndexFlusher),
		ttl:          ttl,
		sr:           stream.NewReader(nil)}
}

func (m *forwardIndexMerger) Reset() {
	m.flusher.Reset()
}

func (m *forwardIndexMerger) latestVersionBlock(versionBlocks [][]byte) []byte {
	var (
		latestIndex int
		maxEndTime  uint32
	)
	for idx, block := range versionBlocks {
		m.sr.Reset(block)
		_, endTime := m.sr.ReadUint32(), m.sr.ReadUint32()
		if endTime > maxEndTime {
			maxEndTime = endTime
			latestIndex = idx
		}
	}
	m.sr.Reset(nil)
	return versionBlocks[latestIndex]
}

// AliveVersions deletes the expired versions
func (m *forwardIndexMerger) AliveVersions(
	versionBlocksMap map[series.Version][][]byte,
) (alive []series.Version,
) {
	var versions []series.Version
	// collect a sorted versions list
	for version := range versionBlocksMap {
		versions = append(versions, version)
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i].Before(versions[j]) })

	var lastVersion series.Version
	for _, version := range versions {
		lastVersion = version
		if !version.IsExpired(m.ttl) {
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

	// version->List<VersionBlock>
	var versionBlocksMap = make(map[series.Version][][]byte)
	for _, block := range value {
		versionBlockItr, err := newVersionBlockIterator(block)
		if err != nil {
			continue
		}
		for versionBlockItr.HasNext() {
			version, versionBlock := versionBlockItr.Next()
			list, ok := versionBlocksMap[version]
			if ok {
				list = append(list, versionBlock)
			} else {
				list = [][]byte{versionBlock}
			}
			versionBlocksMap[version] = list
		}
	}
	if len(versionBlocksMap) == 0 {
		return nil, fmt.Errorf("no available blocks for compacting")
	}
	for _, version := range m.AliveVersions(versionBlocksMap) {
		latestVersionBlock := m.latestVersionBlock(versionBlocksMap[version])
		startPos := m.flusher.metricBlockWriter.Len()
		m.flusher.metricBlockWriter.PutBytes(latestVersionBlock)
		m.flusher.RecordVersionOffset(version, startPos)
	}
	_ = m.flusher.FlushMetricID(key)
	return m.nopKVFlusher.Bytes(), nil
}

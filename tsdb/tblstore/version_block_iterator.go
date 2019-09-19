package tblstore

import (
	"fmt"

	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
)

const footerSizeAfterVersionEntries = 4 + // versionOffsetPos, uint32
	4 // CRC32 checksum, uint32

// versionBlockIterator a iterator for iterating version-block in common-use.
// see Level2 of MetricDataTable & Level2 of ForwardIndexTable in `tsdb/doc`
type versionBlockIterator struct {
	block            []byte
	offsetsReader    *stream.Reader // reading version offsets
	blockReader      *stream.Reader // reading blocks
	versionsCount    int            // total
	versionsRead     int            // have read
	lastVersion      series.Version // last read
	lastVersionBlock []byte         // last read
}

func newVersionBlockIterator(block []byte) (*versionBlockIterator, error) {
	if len(block) <= footerSizeAfterVersionEntries {
		return nil, fmt.Errorf("block length too short")
	}
	itr := &versionBlockIterator{
		block:         block,
		offsetsReader: stream.NewReader(block),
		blockReader:   stream.NewReader(block)}
	itr.readVersionsCount()
	return itr, itr.offsetsReader.Error()
}

func (itr *versionBlockIterator) readVersionsCount() {
	//////////////////////////////////////////////////
	// Read VersionOffSetsBlock
	//////////////////////////////////////////////////
	_ = itr.offsetsReader.ReadSlice(len(itr.block) - footerSizeAfterVersionEntries)
	versionOffsetPos := itr.offsetsReader.ReadUint32()
	// shift to Start Position of the VersionOffsetsBlock
	itr.offsetsReader.SeekStart()
	_ = itr.offsetsReader.ReadSlice(int(versionOffsetPos))
	// read version count
	itr.versionsCount = int(itr.offsetsReader.ReadUvarint64())
}

func (itr *versionBlockIterator) HasNext() bool {
	// read all versions
	if itr.versionsRead >= itr.versionsCount {
		return false
	}
	defer func() { itr.versionsRead++ }()
	// read version
	itr.lastVersion = series.Version(itr.offsetsReader.ReadInt64())
	// read version length
	versionLength := int(itr.offsetsReader.ReadUvarint64())
	itr.lastVersionBlock = itr.blockReader.ReadSlice(versionLength)
	return itr.blockReader.Error() == nil || itr.offsetsReader.Error() == nil
}

func (itr *versionBlockIterator) Next() (version series.Version, versionBlock []byte) {
	return itr.lastVersion, itr.lastVersionBlock
}

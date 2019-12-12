package tblstore

import (
	"fmt"

	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source ./version_block_iterator.go -destination=./version_block_iterator_mock.go -package tblstore

const footerSizeAfterVersionEntries = 4 + // versionOffsetPos, uint32
	4 // CRC32 checksum, uint32

// VersionBlockIterator represents a iterator for iterating version block as one metric data
type VersionBlockIterator interface {
	// HasNext returns if the iteration has more version block
	HasNext() bool
	// Next goes to next loop
	Next()
	// Peek peeks the current version and block
	Peek() (version series.Version, versionBlock []byte)
}

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

	needDoNext bool // mark if invoke next method
}

func NewVersionBlockIterator(block []byte) (VersionBlockIterator, error) {
	if len(block) <= footerSizeAfterVersionEntries {
		return nil, fmt.Errorf("block length too short")
	}
	itr := &versionBlockIterator{
		block:         block,
		offsetsReader: stream.NewReader(block),
		blockReader:   stream.NewReader(block),
		needDoNext:    true,
	}
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
	if !itr.needDoNext {
		return true
	}
	// read all versions
	if itr.versionsRead >= itr.versionsCount {
		return false
	}
	defer func() {
		itr.versionsRead++
		itr.needDoNext = false
	}()
	// read version
	itr.lastVersion = series.Version(itr.offsetsReader.ReadInt64())
	// read version length
	versionLength := int(itr.offsetsReader.ReadUvarint64())
	//FIXME need read the data if version not match
	itr.lastVersionBlock = itr.blockReader.ReadSlice(versionLength)
	return itr.blockReader.Error() == nil || itr.offsetsReader.Error() == nil
}

func (itr *versionBlockIterator) Next() {
	itr.needDoNext = true
}

func (itr *versionBlockIterator) Peek() (version series.Version, versionBlock []byte) {
	return itr.lastVersion, itr.lastVersionBlock
}

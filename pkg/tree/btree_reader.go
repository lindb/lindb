package tree

import (
	"bytes"

	"github.com/eleme/lindb/pkg/stream"
)

//Reader represents is used to parse B+Tree in the disk, support get and prefix query
type Reader struct {
	bufReader       *stream.ByteBufReader
	height          int         //height of the B-tree
	highPos         map[int]int //key:height of the tree  value:start position
	bodyPos         int         //The starting position of the data block
	hasChildrenNode bool
}

//ReaderIterator returns enumerating in the tree. It is returned from the Seek methods.
type ReaderIterator struct {
	reader *Reader //tree reader

	//Hit data block information
	leafNodes int    //Number of nodes
	idx       int    //index
	lcp       []byte //longest common prefix

	key   []byte // k
	value int    // v

	hit    bool
	filter Filter
	init   bool // The tag is read for the first time
}

//Next returns true if the iteration has elements
func (it *ReaderIterator) Next() bool {
	//
	if !it.init {
		it.init = true
		return true
	}
	//end of reading
	if it.reader.bufReader.IsEnd() {
		it.key = nil
		return false
	}
	//read next data block
	if !it.hit {
		it.leafNodes = int(it.reader.bufReader.ReadUvarint64())
		_, lcp := it.reader.bufReader.ReadLenBytes()
		it.lcp = lcp
		it.hit = true
	}
	if it.idx < it.leafNodes {
		_, key := it.reader.bufReader.ReadLenBytes()
		v := int(it.reader.bufReader.ReadUvarint64())

		it.idx++
		if len(it.lcp) > 0 {
			key = bytesCombine(it.lcp, key)
		}
		if it.filter.endMatch(key) {
			it.key = key
			it.value = v
			return true
		}
		it.key = nil
		return false
	}
	it.idx = 0
	it.hit = false
	return it.Next()
}

//GetKey returns the next key in the iteration.
func (it *ReaderIterator) GetKey() []byte {
	return it.key
}

//GetValue returns the next value in the iteration.
func (it *ReaderIterator) GetValue() int {
	return it.value
}

//NewReader returns create a B+Tree Reader
func NewReader(treeBytes []byte) *Reader {
	reader := &Reader{
		bufReader:       stream.NewBufReader(treeBytes),
		highPos:         make(map[int]int),
		hasChildrenNode: false,
	}

	//reader header
	if HasChildrenNode == reader.bufReader.ReadByte() {
		reader.hasChildrenNode = true
	}
	reader.height = int(reader.bufReader.ReadUvarint64())
	//Starting offset of each height
	for i := 0; i < reader.height; i++ {
		high := int(reader.bufReader.ReadUvarint64())
		start := int(reader.bufReader.ReadUvarint64())
		reader.highPos[high] = start
	}

	reader.bodyPos = reader.bufReader.GetPosition()
	return reader
}

//Duplicator return creates a new reader that shares this buffer's content.
func (r *Reader) Duplicator() *Reader {
	reader := &Reader{
		bufReader:       r.bufReader,
		highPos:         r.highPos,
		hasChildrenNode: r.hasChildrenNode,
		bodyPos:         r.bodyPos,
		height:          r.height,
	}
	return reader
}

// Get returns the value associated with k and true if it exists. Otherwise Get
// returns (zero-value, false).
func (r *Reader) Get(target []byte) (v int /*V*/, ok bool) {
	startPos := r.findTargetLeafNodePos(target)
	if startPos == NotFound {
		return NotFound, false
	}
	return r.linearSearchTarget(startPos, target)
}

//findTargetLeafNodePos returns find the target leaf node position
func (r *Reader) findTargetLeafNodePos(target []byte) int {
	var startPos = 0

	if r.hasChildrenNode {
		for high := 1; high < r.height; high++ {
			r.bufReader.NewPosition(r.bodyPos + r.highPos[high] + startPos)
			//read branch count
			count := int(r.bufReader.ReadUvarint64())
			//read common prefix
			lcpLen, lcp := r.bufReader.ReadLenBytes()

			if lcpLen >= len(target) {
				if bytes.Compare(lcp, target) < 0 {
					return NotFound
				}
			}
			startPos = r.linearSearchTargetPos(target, count, lcpLen, lcp)
		}
	} else {
		//If only leaf node, read from the leaf node 0 position
		return 0
	}
	return startPos
}

//extractHeaderAndTargetLcp returns to extract the longest common prefix
func extractHeaderAndTargetLcp(target, headerLcp []byte) int {
	var keyArray = make([][]byte, 2)
	keyArray[0] = target
	keyArray[1] = headerLcp

	return len(lcpByte(keyArray))
}

//linearSearchTargetPos returns the position of the leftmost element in the target,
//return -1 if not found.
func (r *Reader) linearSearchTargetPos(target []byte, count, lcpLen int, lcp []byte) int {
	if lcpLen >= len(target) {
		if bytes.Compare(lcp, target) < 0 {
			return NotFound
		}
	}

	targetLcpLen := extractHeaderAndTargetLcp(lcp, target)
	lcpDiff := lcpLen - targetLcpLen
	targetSuffix := target[targetLcpLen:]

	if lcpDiff > 0 {
		cmp := BytesCompare(lcp[targetLcpLen:], targetSuffix)
		if cmp < 0 {
			return NotFound
		}
		if cmp > 0 {
			//read next branch node
			r.bufReader.ReadByte()
			r.bufReader.ReadLenBytes()
			return int(r.bufReader.ReadUvarint64())
		}
	}

	for i := 0; i < count; i++ {
		//read a branch node
		hasParent := r.bufReader.ReadByte() == HasParent
		suffixLen := int(r.bufReader.ReadUvarint64())
		suffix := r.bufReader.ReadBytes(suffixLen)
		nextStartPos := int(r.bufReader.ReadUvarint64())

		cmp := BytesCompare(suffix, targetSuffix)
		if cmp > 0 {
			return nextStartPos
		} else if cmp == 0 {
			if !hasParent {
				return nextStartPos
			}
		}
	}
	return NotFound
}

//linearSearchTarget returns associated value; return false if no such key
func (r *Reader) linearSearchTarget(pos int, target []byte) (int /*V*/, bool) {
	r.bufReader.NewPosition(r.bodyPos + r.highPos[r.height] + pos)

	count := int(r.bufReader.ReadUvarint64())
	lcpLen, lcp := r.bufReader.ReadLenBytes()

	if lcpLen > 0 {
		if !bytes.HasPrefix(target, lcp) {
			return NotFound, false
		}
	}

	for i := 0; i < count; i++ {
		//read a leaf node
		suffix := r.bufReader.ReadBytes(int(r.bufReader.ReadUvarint64()))
		v := int(r.bufReader.ReadUvarint64())
		if bytes.Equal(target[lcpLen:], suffix) {
			return v, true
		}
	}
	return NotFound, false
}

//Range returns an Iterator for the given key range.
func (r *Reader) Range(startKey, endKey []byte) Iterator {
	startPos := r.findTargetLeafNodePos(startKey)
	if startPos == NotFound {
		return nil
	}
	r.bufReader.NewPosition(r.bodyPos + r.highPos[r.height] + startPos)

	rangeFilter := &RangeFilter{
		startKey: startKey,
		endKey:   endKey,
	}
	return r.seekLeafNodes(rangeFilter)
}

//SeekFirst returns an Iterator positioned on the first K-V pair in the tree
func (r *Reader) SeekToFirst() Iterator {
	r.bufReader.NewPosition(r.bodyPos + r.highPos[r.height])

	it := &ReaderIterator{
		reader: r,
		filter: &SkipFilter{},
		init:   true,
	}
	return it
}

//Seek returns an Iterator positioned on an item such that item's key is prefix key
func (r *Reader) Seek(prefix []byte) Iterator {
	startPos := r.findTargetLeafNodePos(prefix)
	if startPos == NotFound {
		return nil
	}
	r.bufReader.NewPosition(r.bodyPos + r.highPos[r.height] + startPos)

	seekFilter := &SeekFilter{
		prefix: prefix,
	}

	return r.seekLeafNodes(seekFilter)
}

//seekLeafNodes return a ReaderIterator.
//Leaf node for seek lookup.
func (r *Reader) seekLeafNodes(filter Filter) *ReaderIterator {
	leafNodes := int(r.bufReader.ReadUvarint64())
	leafLcpLen, leafLcp := r.bufReader.ReadLenBytes()

	it := &ReaderIterator{
		reader:    r,
		filter:    filter,
		leafNodes: leafNodes,
		lcp:       leafLcp,
	}

	for i := 0; i < leafNodes; i++ {
		it.idx++
		_, suffix := r.bufReader.ReadLenBytes()
		v := int(r.bufReader.ReadUvarint64())
		key := suffix
		if leafLcpLen > 0 {
			key = bytesCombine(leafLcp, suffix)
		}

		if filter.beginMatch(key) {
			it.key = key
			it.value = v
			it.hit = true
			break
		}
	}
	if it.hit {
		return it
	}
	return nil
}

package trie

import (
	"bytes"
)

type Iterator struct {
	valid        bool
	atTerminator bool
	tree         *trie
	level        uint32
	keyBuf       []byte
	fullKeyBuf   []byte
	posInTrie    []uint32
	nodeID       []uint32
	prefixLen    []uint32
}

func (it *Iterator) init(tree *trie) {
	it.tree = tree
	it.posInTrie = make([]uint32, tree.height)
	it.prefixLen = make([]uint32, tree.height)
	it.nodeID = make([]uint32, tree.height)
}

func (it *Iterator) Next() {
	if !it.valid {
		return
	}
	it.atTerminator = false
	pos := it.posInTrie[it.level] + 1

	for pos >= it.tree.loudsVec.numBits || it.tree.loudsVec.IsSet(pos) {
		if it.level == 0 {
			it.valid = false
			it.keyBuf = it.keyBuf[:0]
			return
		}
		it.level--
		pos = it.posInTrie[it.level] + 1
	}
	it.setAt(it.level, pos)
	it.moveToLeftMostKey()
}

func (it *Iterator) Valid() bool {
	return it.valid
}

func (it *Iterator) Prev() {
	if !it.valid {
		return
	}
	it.atTerminator = false
	pos := it.posInTrie[it.level]

	if pos == 0 {
		it.valid = false
		return
	}
	for it.tree.loudsVec.IsSet(pos) {
		if it.level == 0 {
			it.valid = false
			it.keyBuf = it.keyBuf[:0]
			return
		}
		it.level--
		pos = it.posInTrie[it.level]
	}
	it.setAt(it.level, pos-1)
	it.moveToRightMostKey()
}

func (it *Iterator) Seek(key []byte) bool {
	var fp bool
	it.Reset()

	if it.tree.height == 0 {
		return false
	}

	fp = it.seek(key)
	if !it.valid {
		it.moveToRightMostKey()
	}
	return fp
}

func (it *Iterator) seek(key []byte) bool {
	nodeID := uint32(0)
	pos := it.tree.firstLabelPos(nodeID)
	var ok bool
	depth := uint32(0)
	for it.level = 0; it.level < it.tree.sparseLevels(); it.level++ {
		prefix := it.tree.prefixVec.GetPrefix(it.tree.prefixID(nodeID))
		var prefixCmp int
		if len(prefix) != 0 {
			end := int(depth) + len(prefix)
			if end > len(key) {
				end = len(key)
			}
			prefixCmp = bytes.Compare(prefix, key[depth:end])
		}

		if prefixCmp < 0 {
			if it.level == 0 {
				it.valid = false
				return false
			}
			it.level--
			it.Next()
			return false
		}

		depth += uint32(len(prefix))
		if depth >= uint32(len(key)) || prefixCmp > 0 {
			it.append(it.tree.labelVec.GetLabel(pos), pos, nodeID)
			it.moveToLeftMostKey()
			return false
		}

		nodeSize := it.tree.nodeSize(pos)
		pos, ok = it.tree.labelVec.Search(key[depth], pos, nodeSize)
		if !ok {
			it.moveToLeftInNextSubTrie(pos, nodeID, nodeSize, key[depth])
			return false
		}

		it.append(key[depth], pos, nodeID)

		if !it.tree.hasChildVec.IsSet(pos) {
			return it.tree.suffixes.CheckSuffix(it.tree.suffixPos(pos), key, depth+1)
		}

		nodeID = it.tree.childNodeID(pos)
		pos = it.tree.firstLabelPos(nodeID)
		depth++
	}

	if it.tree.labelVec.GetLabel(pos) == labelTerminator && !it.tree.hasChildVec.IsSet(pos) && !it.tree.isEndOfNode(pos) {
		it.append(labelTerminator, pos, nodeID)
		it.atTerminator = true
		it.valid = true
		return false
	}

	if uint32(len(key)) <= depth {
		it.moveToLeftMostKey()
		return false
	}

	it.valid = true
	return true
}

func (it *Iterator) SeekToFirst() {
	it.Reset()

	if it.tree.height > 0 {
		it.setToFirstInRoot()
		it.moveToLeftMostKey()
	}
}

func (it *Iterator) SeekToLast() {
	it.Reset()

	if it.tree.height > 0 {
		it.setToLastInRoot()
		it.moveToRightMostKey()
	}
}

func (it *Iterator) uniqueKey() []byte {
	if it.atTerminator {
		return it.keyBuf[:len(it.keyBuf)-1]
	}
	return it.keyBuf
}

func (it *Iterator) Key() []byte {
	valPos := it.tree.suffixPos(it.posInTrie[it.level])
	suffix := it.tree.suffixes.GetSuffix(valPos)

	if len(suffix) == 0 {
		return it.uniqueKey()
	}

	it.fullKeyBuf = append(it.fullKeyBuf[:0], it.uniqueKey()...)
	it.fullKeyBuf = append(it.fullKeyBuf, suffix...)
	return it.fullKeyBuf
}

func (it *Iterator) Value() []byte {
	valPos := it.tree.suffixPos(it.posInTrie[it.level])
	return it.tree.values.Get(valPos)
}

func (it *Iterator) Reset() {
	it.valid = false
	it.level = 0
	it.atTerminator = false
	it.keyBuf = it.keyBuf[:0]
	it.fullKeyBuf = it.fullKeyBuf[:0]
}

func (it *Iterator) moveToLeftMostKey()  { it.moveToMostKey(true) }
func (it *Iterator) moveToRightMostKey() { it.moveToMostKey(false) }
func (it *Iterator) moveToMostKey(left bool) {
	var labelPosFunc func(nodeID uint32) uint32
	if left {
		labelPosFunc = it.tree.firstLabelPos
	} else {
		labelPosFunc = it.tree.lastLabelPos
	}

	if len(it.keyBuf) == 0 {
		pos := labelPosFunc(0)
		label := it.tree.labelVec.GetLabel(pos)
		it.append(label, pos, 0)
	}

	pos := it.posInTrie[it.level]
	label := it.tree.labelVec.GetLabel(pos)

	if !it.tree.hasChildVec.IsSet(pos) {
		if label == labelTerminator && !it.tree.isEndOfNode(pos) {
			it.atTerminator = true
		}
		it.valid = true
		return
	}

	for it.level < it.tree.sparseLevels() {
		it.level++
		nodeID := it.tree.childNodeID(pos)
		pos = labelPosFunc(nodeID)
		label = it.tree.labelVec.GetLabel(pos)

		if !it.tree.hasChildVec.IsSet(pos) {
			it.append(label, pos, nodeID)
			if label == labelTerminator && !it.tree.isEndOfNode(pos) {
				it.atTerminator = true
			}
			it.valid = true
			return
		}
		it.append(label, pos, nodeID)
	}
}

func (it *Iterator) setToFirstInRoot() {
	it.append(it.tree.labelVec.GetLabel(0), 0, 0)
}

func (it *Iterator) setToLastInRoot() {
	pos := it.tree.lastLabelPos(0)
	it.append(it.tree.labelVec.GetLabel(pos), pos, 0)
}

func (it *Iterator) append(label byte, pos, nodeID uint32) {
	prefix := it.tree.prefixVec.GetPrefix(it.tree.prefixID(nodeID))
	it.keyBuf = append(it.keyBuf, prefix...)
	it.keyBuf = append(it.keyBuf, label)
	it.posInTrie[it.level] = pos
	it.prefixLen[it.level] = uint32(len(prefix)) + 1
	if it.level != 0 {
		it.prefixLen[it.level] += it.prefixLen[it.level-1]
	}
	it.nodeID[it.level] = nodeID
}

func (it *Iterator) setAt(level, pos uint32) {
	it.keyBuf = append(it.keyBuf[:it.prefixLen[level]-1], it.tree.labelVec.GetLabel(pos))
	it.posInTrie[it.level] = pos
}

func (it *Iterator) moveToLeftInNextSubTrie(pos, nodeID, nodeSize uint32, label byte) {
	pos, ok := it.tree.labelVec.SearchGreaterThan(label, pos, nodeSize)
	it.append(it.tree.labelVec.GetLabel(pos), pos, nodeID)
	if ok {
		it.moveToLeftMostKey()
	} else {
		it.Next()
	}
}

package trie

import (
	"encoding/binary"
)

// builder is builder of Succinct Trie.
type builder struct {
	valueWidth uint32
	totalCount int

	// LOUDS-Sparse bitvecs, pooling
	lsLabels    [][]byte
	lsHasChild  [][]uint64
	lsLoudsBits [][]uint64

	// suffix keys
	scratch         []byte  // for variant encoding
	suffixesOffsets [][]int // pooling
	suffixesBlock   []byte

	// value
	values      [][]byte
	valueCounts []uint32

	// prefix
	hasPrefix [][]uint64
	prefixes  [][][]byte

	nodeCounts           []uint32
	isLastItemTerminator []bool

	// pooling data-structures
	poolByteSlice   []*[]byte
	poolUint64Slice []*[]uint64
	poolIntSlice    []*[]int
}

// NewBuilder returns a new Trie builder.
func NewBuilder() Builder {
	return &builder{
		suffixesBlock: make([]byte, 4096)[:0],
		scratch:       make([]byte, binary.MaxVarintLen32),
	}
}

func (b *builder) Build(keys, vals [][]byte, valueWidth uint32) SuccinctTrie {
	b.valueWidth = valueWidth

	b.totalCount = len(keys)
	b.buildNodes(keys, vals, 0, 0, 0)

	tree := new(trie)
	tree.Init(b)
	return tree
}

// buildNodes is recursive algorithm to bulk building Trie nodes.
//	* We divide keys into groups by the `key[depth]`, so keys in each group shares the same prefix
//	* If depth larger than the length if the first key in group, the key is prefix of others in group
//	  So we should append `labelTerminator` to labels and update `b.isLastItemTerminator`, then remove it from group.
//	* Scan over keys in current group when meets different label, use the new sub group call buildNodes with level+1 recursively
//	* If all keys in current group have the same label, this node can be compressed, use this group call buildNodes with level recursively.
//	* If current group contains only one key constract suffix of this key and return.
func (b *builder) buildNodes(keys, vals [][]byte, prefixDepth, depth, level int) {
	b.ensureLevel(level)
	nodeStartPos := b.numItems(level)

	groupStart := 0
	if depth >= len(keys[groupStart]) {
		b.lsLabels[level] = append(b.lsLabels[level], labelTerminator)
		b.isLastItemTerminator[level] = true
		b.insertSuffix(keys[groupStart], level, depth)
		b.insertValue(vals[groupStart], level)
		b.moveToNextItemSlot(level)
		groupStart++
	}

	for groupEnd := groupStart; groupEnd <= len(keys); groupEnd++ {
		if groupEnd < len(keys) && keys[groupStart][depth] == keys[groupEnd][depth] {
			continue
		}

		if groupEnd == len(keys) && groupStart == 0 && groupEnd-groupStart != 1 {
			// node at this level is one-way node, compress it to next node
			b.buildNodes(keys, vals, prefixDepth, depth+1, level)
			return
		}

		b.lsLabels[level] = append(b.lsLabels[level], keys[groupStart][depth])
		b.moveToNextItemSlot(level)
		if groupEnd-groupStart == 1 {
			b.insertSuffix(keys[groupStart], level, depth)
			b.insertValue(vals[groupStart], level)
		} else {
			setBit(b.lsHasChild[level], b.numItems(level)-1)
			b.buildNodes(keys[groupStart:groupEnd], vals[groupStart:groupEnd], depth+1, depth+1, level+1)
		}

		groupStart = groupEnd
	}

	// check if current node contains compressed path.
	if depth-prefixDepth > 0 {
		prefix := keys[0][prefixDepth:depth]
		setBit(b.hasPrefix[level], b.nodeCounts[level])
		b.insertPrefix(prefix, level)
	}
	setBit(b.lsLoudsBits[level], nodeStartPos)

	b.nodeCounts[level]++
	if b.nodeCounts[level]%wordSize == 0 {
		b.hasPrefix[level] = append(b.hasPrefix[level], 0)
	}
}

func (b *builder) ensureLevel(level int) {
	if level >= b.treeHeight() {
		b.addLevel()
	}
}

func (b *builder) treeHeight() int {
	return len(b.nodeCounts)
}

func (b *builder) numItems(level int) uint32 {
	return uint32(len(b.lsLabels[level]))
}

func (b *builder) addLevel() {
	// pooled
	b.lsLabels = append(b.lsLabels, *b.pickByteSlice())
	b.lsHasChild = append(b.lsHasChild, *b.pickUint64Slice())
	b.lsLoudsBits = append(b.lsLoudsBits, *b.pickUint64Slice())
	// not pooled
	b.hasPrefix = append(b.hasPrefix, []uint64{})
	// pooled
	b.suffixesOffsets = append(b.suffixesOffsets, *b.pickIntSlice())

	b.values = append(b.values, []byte{})
	b.valueCounts = append(b.valueCounts, 0)
	b.prefixes = append(b.prefixes, [][]byte{})

	b.nodeCounts = append(b.nodeCounts, 0)
	b.isLastItemTerminator = append(b.isLastItemTerminator, false)

	level := b.treeHeight() - 1
	b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	b.lsLoudsBits[level] = append(b.lsLoudsBits[level], 0)
	b.hasPrefix[level] = append(b.hasPrefix[level], 0)
}

func (b *builder) moveToNextItemSlot(level int) {
	if b.numItems(level)%wordSize == 0 {
		b.lsHasChild[level] = append(b.lsHasChild[level], 0)
		b.lsLoudsBits[level] = append(b.lsLoudsBits[level], 0)
	}
}

func (b *builder) insertSuffix(key []byte, level, depth int) {
	if level >= b.treeHeight() {
		b.addLevel()
	}

	var keySuffix []byte
	cutPos := depth + 1
	if cutPos > len(key) {
		keySuffix = nil
	} else {
		keySuffix = key[cutPos:]
	}
	offset := len(b.suffixesBlock)
	b.suffixesOffsets[level] = append(b.suffixesOffsets[level], offset)

	// put uvarint length of key-suffix
	width := binary.PutUvarint(b.scratch, uint64(len(keySuffix)))
	b.suffixesBlock = append(b.suffixesBlock, b.scratch[:width]...)
	// put key-suffix
	if len(keySuffix) != 0 {
		b.suffixesBlock = append(b.suffixesBlock, keySuffix...)
	}
}

func (b *builder) insertValue(value []byte, level int) {
	b.values[level] = append(b.values[level], value[:b.valueWidth]...)
	b.valueCounts[level]++
}

func (b *builder) insertPrefix(prefix []byte, level int) {
	b.prefixes[level] = append(b.prefixes[level], append([]byte{}, prefix...))
}

func (b *builder) Reset() {
	b.valueWidth = 0
	b.totalCount = 0

	// cache lsLabels
	for idx := range b.lsLabels {
		sl := b.lsLabels[idx][:0]
		b.poolByteSlice = append(b.poolByteSlice, &sl)
	}
	b.lsLabels = b.lsLabels[:0]

	// cache lsHasChild
	for idx := range b.lsHasChild {
		sl := b.lsHasChild[idx][:0]
		b.poolUint64Slice = append(b.poolUint64Slice, &sl)
	}
	b.lsHasChild = b.lsHasChild[:0]

	// cache lsLoudsBits
	for idx := range b.lsLoudsBits {
		sl := b.lsLoudsBits[idx][:0]
		b.poolUint64Slice = append(b.poolUint64Slice, &sl)
	}
	b.lsLoudsBits = b.lsLoudsBits[:0]

	// cache suffixOffsets
	for idx := range b.suffixesOffsets {
		sl := b.suffixesOffsets[idx][:0]
		b.poolIntSlice = append(b.poolIntSlice, &sl)
	}
	b.suffixesOffsets = b.suffixesOffsets[:0]

	// reset suffixesBlock
	b.suffixesBlock = b.suffixesBlock[:0]

	// reset values
	b.values = b.values[:0]
	b.valueCounts = b.valueCounts[:0]

	// reset prefixes
	b.hasPrefix = b.hasPrefix[:0]
	b.prefixes = b.prefixes[:0]

	// reset nodeCounts
	b.nodeCounts = b.nodeCounts[:0]
	b.isLastItemTerminator = b.isLastItemTerminator[:0]
}

func (b *builder) pickByteSlice() *[]byte {
	if len(b.poolByteSlice) == 0 {
		return &[]byte{}
	}
	tailIndex := len(b.poolByteSlice) - 1
	ptr := b.poolByteSlice[tailIndex]
	b.poolByteSlice = b.poolByteSlice[:tailIndex]
	return ptr
}

func (b *builder) pickUint64Slice() *[]uint64 {
	if len(b.poolUint64Slice) == 0 {
		return &[]uint64{}
	}
	tailIndex := len(b.poolUint64Slice) - 1
	ptr := b.poolUint64Slice[tailIndex]
	b.poolUint64Slice = b.poolUint64Slice[:tailIndex]
	return ptr
}
func (b *builder) pickIntSlice() *[]int {
	if len(b.poolIntSlice) == 0 {
		return &[]int{}
	}
	tailIndex := len(b.poolIntSlice) - 1
	ptr := b.poolIntSlice[tailIndex]
	b.poolIntSlice = b.poolIntSlice[:tailIndex]
	return ptr
}

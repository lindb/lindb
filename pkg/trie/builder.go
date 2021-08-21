// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package trie

// builder builds Succinct Trie.
type builder struct {
	valueWidth uint32
	totalCount int

	// LOUDS-Sparse bitvecs, pooling
	lsLabels    [][]byte
	lsHasChild  [][]uint64
	lsLoudsBits [][]uint64

	// value
	values      [][]byte
	valueCounts []uint32

	// prefix
	hasPrefix [][]uint64
	prefixes  [][][]byte

	// suffix
	hasSuffix [][]uint64
	suffixes  [][][]byte

	nodeCounts           []uint32
	isLastItemTerminator []bool

	// pooling data-structures
	cachedLabel   [][]byte
	cachedUint64s [][]uint64
}

// NewBuilder returns a new Trie builder.
func NewBuilder() Builder {
	return &builder{}
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
		b.ensureLevel(level)
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
			if depth+1 < len(keys[groupStart]) {
				b.ensureLevel(level)
				setBit(b.hasSuffix[level], b.numItems(level)-1)
				b.suffixes[level] = append(b.suffixes[level], keys[groupStart][depth+1:])
			}
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
		b.insertPrefix(prefix, level)
	}

	setBit(b.lsLoudsBits[level], nodeStartPos)

	b.nodeCounts[level]++
	if b.nodeCounts[level]%wordSize == 0 {
		b.hasPrefix[level] = append(b.hasPrefix[level], 0)
		b.hasSuffix[level] = append(b.hasSuffix[level], 0)
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
	// cached
	b.lsLabels = append(b.lsLabels, b.pickLabels())
	b.lsHasChild = append(b.lsHasChild, b.pickUint64Slice())
	b.lsLoudsBits = append(b.lsLoudsBits, b.pickUint64Slice())
	b.hasPrefix = append(b.hasPrefix, b.pickUint64Slice())
	b.hasSuffix = append(b.hasSuffix, b.pickUint64Slice())

	b.values = append(b.values, []byte{})
	b.valueCounts = append(b.valueCounts, 0)
	b.prefixes = append(b.prefixes, [][]byte{})
	b.suffixes = append(b.suffixes, [][]byte{})

	b.nodeCounts = append(b.nodeCounts, 0)
	b.isLastItemTerminator = append(b.isLastItemTerminator, false)

	level := b.treeHeight() - 1
	b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	b.lsLoudsBits[level] = append(b.lsLoudsBits[level], 0)
	b.hasPrefix[level] = append(b.hasPrefix[level], 0)
	b.hasSuffix[level] = append(b.hasSuffix[level], 0)
}

func (b *builder) moveToNextItemSlot(level int) {
	if b.numItems(level)%wordSize == 0 {
		b.hasSuffix[level] = append(b.hasSuffix[level], 0)
		b.lsHasChild[level] = append(b.lsHasChild[level], 0)
		b.lsLoudsBits[level] = append(b.lsLoudsBits[level], 0)
	}
}

func (b *builder) insertValue(value []byte, level int) {
	b.values[level] = append(b.values[level], value[:b.valueWidth]...)
	b.valueCounts[level]++
}

func (b *builder) insertPrefix(prefix []byte, level int) {
	setBit(b.hasPrefix[level], b.nodeCounts[level])
	b.prefixes[level] = append(b.prefixes[level], prefix)
}

func (b *builder) Reset() {
	b.valueWidth = 0
	b.totalCount = 0

	// cache lsLabels
	for idx := range b.lsLabels {
		b.cachedLabel = append(b.cachedLabel, b.lsLabels[idx][:0])
	}
	b.lsLabels = b.lsLabels[:0]

	// cache lsHasChild
	for idx := range b.lsHasChild {
		b.cachedUint64s = append(b.cachedUint64s, b.lsHasChild[idx][:0])
	}
	b.lsHasChild = b.lsHasChild[:0]

	// cache lsLoudsBits
	for idx := range b.lsLoudsBits {
		b.cachedUint64s = append(b.cachedUint64s, b.lsLoudsBits[idx][:0])
	}
	b.lsLoudsBits = b.lsLoudsBits[:0]

	// reset values
	b.values = b.values[:0]
	b.valueCounts = b.valueCounts[:0]

	// cache has prefix
	for idx := range b.hasPrefix {
		b.hasPrefix = append(b.hasPrefix, b.hasPrefix[idx][:0])
	}
	b.hasPrefix = b.hasPrefix[:0]

	// cache has suffix
	for idx := range b.hasSuffix {
		b.hasSuffix = append(b.hasSuffix, b.hasSuffix[idx][:0])
	}
	b.hasSuffix = b.hasSuffix[:0]

	// reset prefixes
	b.hasPrefix = b.hasPrefix[:0]
	b.prefixes = b.prefixes[:0]

	// reset suffixes
	b.hasSuffix = b.hasSuffix[:0]
	b.suffixes = b.suffixes[:0]

	// reset nodeCounts
	b.nodeCounts = b.nodeCounts[:0]
	b.isLastItemTerminator = b.isLastItemTerminator[:0]
}

func (b *builder) pickLabels() []byte {
	if len(b.cachedLabel) == 0 {
		return []byte{}
	}
	tailIndex := len(b.cachedLabel) - 1
	ptr := b.cachedLabel[tailIndex]
	b.cachedLabel = b.cachedLabel[:tailIndex]
	return ptr
}

func (b *builder) pickUint64Slice() []uint64 {
	if len(b.cachedUint64s) == 0 {
		return []uint64{}
	}
	tailIndex := len(b.cachedUint64s) - 1
	ptr := b.cachedUint64s[tailIndex]
	b.cachedUint64s = b.cachedUint64s[:tailIndex]
	return ptr
}

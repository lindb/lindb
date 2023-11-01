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

import (
	"io"
)

// builder builds Succinct Trie.
type builder struct {
	totalCount int
	height     int
	// LOUDS-Sparse bitvecs, pooling
	levels []*Level

	// pooling data-structures
	cachedLevels []*Level

	// write context(reuse)
	labelVec    labelVector
	hasChildVec rankVectorSparse
	loudsVec    selectVector
	prefixVec   prefixVector
	suffixVec   suffixVector
}

// NewBuilder returns a new Trie builder.
func NewBuilder() Builder {
	return &builder{}
}

func (b *builder) Build(keys [][]byte, vals []uint32) {
	b.totalCount = len(keys)

	b.buildNodes(keys, vals, 0, 0, 0)
}

func (b *builder) Trie() SuccinctTrie {
	tree := new(trie)
	tree.Init(b)
	return tree
}

// buildNodes is recursive algorithm to bulk building Trie nodes.
//  1. We divide keys into groups by the `key[depth]`, so keys in each group shares the same prefix
//  2. If depth larger than the length if the first key in group, the key is prefix of others in group
//     So we should append `labelTerminator` to labels and update `b.isLastItemTerminator`, then remove it from group.
//  3. Scan over keys in current group when meets different label, use the new sub group call buildNodes with level+1 recursively
//  4. If all keys in current group have the same label, this node can be compressed, use this group call buildNodes with level recursively.
//  5. If current group contains only one key constract suffix of this key and return.
func (b *builder) buildNodes(keys [][]byte, vals []uint32, prefixDepth, depth, level int) {
	b.ensureLevel(level)
	levelObj := b.levels[level]
	nodeStartPos := len(levelObj.lsLabels) // first label pos
	keysLen := len(keys)

	groupStart := 0
	if depth >= len(keys[groupStart]) {
		// first key is completed, append terminator label
		levelObj.lsLabels = append(levelObj.lsLabels, labelTerminator)
		b.moveToNextItemSlot(levelObj)
		levelObj.values = append(levelObj.values, vals[groupStart])
		groupStart++ // move to next key
	}

	currentKey := keys[groupStart][depth]
	for groupEnd := groupStart; groupEnd <= keysLen; groupEnd++ {
		// if groupEnd < keysLen && currentKey == keys[groupEnd][depth] {
		if groupEnd < keysLen {
			// try skip more same labels
			skipEnd := groupEnd + 4
			if skipEnd < keysLen && currentKey == keys[skipEnd][depth] {
				groupEnd = skipEnd
				continue
			}
			if currentKey == keys[groupEnd][depth] {
				// skip same lable
				continue
			}
		}
		width := groupEnd - groupStart
		nextDepth := depth + 1
		if groupEnd == keysLen && groupStart == 0 && width != 1 {
			// node at this level is one-way node, compress it to next node
			b.buildNodes(keys, vals, prefixDepth, nextDepth, level)
			return
		}

		levelObj.lsLabels = append(levelObj.lsLabels, currentKey)
		b.moveToNextItemSlot(levelObj)
		if width == 1 {
			// parent only have two sub tries, complete left child trie
			if nextDepth < len(keys[groupStart]) {
				// if has suffix, store suffix
				setBit(levelObj.hasSuffix, uint32(len(levelObj.lsLabels)-1))
				levelObj.suffixes = append(levelObj.suffixes, keys[groupStart][nextDepth:])
			}
			levelObj.values = append(levelObj.values, vals[groupStart])
		} else {
			// goto next level
			setBit(levelObj.lsHasChild, uint32(len(levelObj.lsLabels)-1))
			b.buildNodes(keys[groupStart:groupEnd], vals[groupStart:groupEnd], nextDepth, nextDepth, level+1)
		}

		groupStart = groupEnd // process right sub trie
		if groupStart < keysLen {
			// get new  current key
			currentKey = keys[groupStart][depth]
		}
	}

	// check if current node contains compressed path.
	if depth > prefixDepth {
		setBit(levelObj.hasPrefix, uint32(levelObj.nodeCount))
		levelObj.prefixes = append(levelObj.prefixes, keys[0][prefixDepth:depth])
	}

	// store start node pos in louds
	setBit(levelObj.lsLouds, uint32(nodeStartPos))

	levelObj.nodeCount++
}

func (b *builder) ensureLevel(level int) {
	if level >= b.height {
		b.addLevel()
	}
}

func (b *builder) addLevel() {
	b.height++
	levelObj := b.pickLevels()
	b.levels = append(b.levels, levelObj)

	levelObj.lsHasChild = append(levelObj.lsHasChild, 0)
	levelObj.lsLouds = append(levelObj.lsLouds, 0)
	levelObj.hasPrefix = append(levelObj.hasPrefix, 0)
	levelObj.hasSuffix = append(levelObj.hasSuffix, 0)
}

func (b *builder) moveToNextItemSlot(level *Level) {
	if len(level.lsLabels)%wordSize == 0 {
		level.lsHasChild = append(level.lsHasChild, 0)
		level.lsLouds = append(level.lsLouds, 0)
		level.hasPrefix = append(level.hasPrefix, 0)
		level.hasSuffix = append(level.hasSuffix, 0)
	}
}

func (b *builder) Reset() {
	b.totalCount = 0
	b.height = 0

	// cache level
	for idx := range b.levels {
		level := b.levels[idx]
		level.Reset()
		b.cachedLevels = append(b.cachedLevels, level)
	}
	b.levels = b.levels[:0]
}

func (b *builder) pickLevels() *Level {
	if len(b.cachedLevels) == 0 {
		return NewLevel()
	}
	tailIndex := len(b.cachedLevels) - 1
	ptr := b.cachedLevels[tailIndex]
	b.cachedLevels = b.cachedLevels[:tailIndex]
	return ptr
}

func (b *builder) Write(w io.Writer) error {
	var (
		bs [4]byte
	)
	// write total keys
	endian.PutUint32(bs[:], uint32(b.totalCount))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write height
	endian.PutUint32(bs[:], uint32(b.height))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write labels
	if err := b.labelVec.Write(w, b.levels); err != nil {
		return err
	}
	// write has child
	b.hasChildVec.init(rankSparseBlockSize, b.levels, HasChild)
	if err := b.hasChildVec.Write(w); err != nil {
		return err
	}
	// write louds
	b.loudsVec.Init(b.levels, Louds)
	if err := b.loudsVec.Write(w); err != nil {
		return err
	}
	// write prefix
	b.prefixVec.Init(b.levels, HasPrefix)
	if err := b.prefixVec.Write(w); err != nil {
		return err
	}
	// write suffix
	b.suffixVec.Init(b.levels, HasSuffix)
	if err := b.suffixVec.Write(w); err != nil {
		return err
	}

	// write values
	for level := range b.levels {
		values := b.levels[level].values
		if len(values) > 0 {
			if _, err := w.Write(u32SliceToBytes(values)); err != nil {
				return err
			}
		}
	}

	return nil
}

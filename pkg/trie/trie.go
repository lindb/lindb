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
	"bytes"
	"io"
)

type trie struct {
	height uint32

	labelVec    labelVector
	hasChildVec rankVectorSparse
	loudsVec    selectVector
	values      valueVector
	prefixVec   prefixVector
	suffixVec   suffixVector
}

// NewTrie returns a new empty SuccinctTrie
func NewTrie() SuccinctTrie {
	return &trie{}
}

func (tree *trie) Init(builder *builder) *trie {
	tree.height = uint32(len(builder.lsLabels))

	tree.labelVec.Init(builder.lsLabels, tree.sparseLevels())

	numItemsPerLevel := make([]uint32, tree.sparseLevels())
	for level := range numItemsPerLevel {
		numItemsPerLevel[level] = uint32(len(builder.lsLabels[level]))
	}
	tree.hasChildVec.Init(builder.lsHasChild, numItemsPerLevel)
	tree.loudsVec.Init(builder.lsLoudsBits, numItemsPerLevel)

	tree.values.Init(builder.values, builder.valueWidth)

	tree.prefixVec.Init(builder.hasPrefix, builder.nodeCounts, builder.prefixes)
	tree.suffixVec.Init(builder.hasSuffix, numItemsPerLevel, builder.suffixes)

	return tree
}

func (tree *trie) Get(key []byte) (value []byte, ok bool) {
	var (
		nodeID    uint32
		pos       = tree.firstLabelPos(nodeID)
		depth     uint32
		prefixLen uint32
		exhausted bool
	)
	for depth = 0; depth < uint32(len(key)); depth++ {
		prefixLen, ok = tree.prefixVec.CheckPrefix(key, depth, tree.prefixID(nodeID))
		if !ok {
			return nil, false
		}
		depth += prefixLen

		if depth >= uint32(len(key)) {
			exhausted = true
			break
		}

		if pos, ok = tree.labelVec.Search(key[depth], pos, tree.nodeSize(pos)); !ok {
			return nil, false
		}
		if !tree.hasChildVec.IsSet(pos) {
			valPos := tree.suffixPos(pos)
			if ok = tree.suffixVec.CheckSuffix(key, depth, pos); ok {
				value = tree.values.Get(valPos)
				ok = true
			}
			return value, ok
		}

		nodeID = tree.childNodeID(pos)
		pos = tree.firstLabelPos(nodeID)
	}
	// key is not exhausted, re-check the prefix
	if !exhausted {
		_, ok = tree.prefixVec.CheckPrefix(key, depth, tree.prefixID(nodeID))
		if !ok {
			return nil, false
		}
	}

	if tree.labelVec.GetLabel(pos) == labelTerminator && !tree.hasChildVec.IsSet(pos) {
		if ok = tree.suffixVec.CheckSuffix(key, depth, pos); ok {
			valPos := tree.suffixPos(pos)
			value = tree.values.Get(valPos)
		}
		return value, ok
	}

	return nil, false
}

func (tree *trie) MarshalSize() int64 {
	return align(tree.rawMarshalSize()) + tree.values.MarshalSize()
}

func (tree *trie) rawMarshalSize() int64 {
	return 4 + tree.labelVec.MarshalSize() + tree.hasChildVec.MarshalSize() + tree.loudsVec.MarshalSize() +
		tree.suffixVec.MarshalSize() + tree.prefixVec.MarshalSize()
}

func (tree *trie) MarshalBinary() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, tree.MarshalSize()))
	_ = tree.Write(w)
	return w.Bytes(), nil
}

func (tree *trie) Write(w io.Writer) error {
	var (
		bs [4]byte
	)
	endian.PutUint32(bs[:], tree.height)
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	if err := tree.labelVec.Write(w); err != nil {
		return err
	}
	if err := tree.hasChildVec.Write(w); err != nil {
		return err
	}
	if err := tree.loudsVec.Write(w); err != nil {
		return err
	}
	if err := tree.prefixVec.Write(w); err != nil {
		return err
	}
	if err := tree.suffixVec.Write(w); err != nil {
		return err
	}
	rawMarshalSize := tree.rawMarshalSize()
	// align
	padding := align(rawMarshalSize) - tree.rawMarshalSize()
	var zeros [8]byte
	if _, err := w.Write(zeros[:padding]); err != nil {
		return err
	}
	// write values
	return tree.values.Write(w)
}

func (tree *trie) UnmarshalBinary(buf []byte) (err error) {
	if len(buf) <= 4 {
		return io.EOF
	}
	buf1 := buf
	tree.height = endian.Uint32(buf1)
	buf1 = buf1[4:]

	if buf1, err = tree.labelVec.Unmarshal(buf1); err != nil {
		return err
	}
	if buf1, err = tree.hasChildVec.Unmarshal(buf1); err != nil {
		return err
	}
	if buf1, err = tree.loudsVec.Unmarshal(buf1); err != nil {
		return err
	}
	if buf1, err = tree.prefixVec.Unmarshal(buf1); err != nil {
		return err
	}
	if buf1, err = tree.suffixVec.Unmarshal(buf1); err != nil {
		return err
	}

	sz := align(int64(len(buf) - len(buf1)))
	if sz > int64(len(buf)) {
		return io.EOF
	}
	buf1 = buf[sz:]
	if _, err = tree.values.Unmarshal(buf1); err != nil {
		return err
	}
	return nil
}

func (tree *trie) NewIterator() *Iterator {
	itr := new(Iterator)
	itr.init(tree)
	return itr
}

func (tree *trie) NewPrefixIterator(prefix []byte) *PrefixIterator {
	rawItr := tree.NewIterator()
	rawItr.Seek(prefix)
	return &PrefixIterator{prefix: prefix, it: rawItr}
}

func (tree *trie) suffixPos(pos uint32) uint32 {
	return pos - tree.hasChildVec.Rank(pos)
}

func (tree *trie) firstLabelPos(nodeID uint32) uint32 {
	return tree.loudsVec.Select(nodeID + 1)
}

func (tree *trie) sparseLevels() uint32 {
	return tree.height
}
func (tree *trie) prefixID(nodeID uint32) uint32 {
	return nodeID
}

func (tree *trie) lastLabelPos(nodeID uint32) uint32 {
	nextRank := nodeID + 2
	if nextRank > tree.loudsVec.numOnes {
		return tree.loudsVec.numBits - 1
	}
	return tree.loudsVec.Select(nextRank) - 1
}

func (tree *trie) childNodeID(pos uint32) uint32 {
	return tree.hasChildVec.Rank(pos)
}

func (tree *trie) nodeSize(pos uint32) uint32 {
	return tree.loudsVec.DistanceToNextSetBit(pos)
}

func (tree *trie) isEndOfNode(pos uint32) bool {
	return pos == tree.loudsVec.numBits-1 || tree.loudsVec.IsSet(pos+1)
}

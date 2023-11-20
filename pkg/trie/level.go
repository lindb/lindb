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

type BitmapType int

const (
	HasChild BitmapType = iota + 1
	Louds
	HasPrefix
	HasSuffix
)

type Level struct {
	// LOUDS-Sparse context: labels/hasChild/louds
	//
	// store all the branching labels for each trie node
	// lsLabels [][]byte
	lsLabels []byte
	// one bit for each byte in labels to indicate whether
	// a child branch continues(i.e. points to a sub-trie)
	// or terminals(i.e. points to a value)
	lsHasChild []uint64
	// one bit for each byte in labels to indicate if a lable
	// is the first node in trie
	lsLouds []uint64

	// prefix
	hasPrefix []uint64
	prefixes  [][]byte
	// suffix
	hasSuffix []uint64
	suffixes  [][]byte
	// value
	values []uint32
	// level node count
	nodeCount int
}

func NewLevel() *Level {
	return &Level{
		lsLabels:   make([]byte, 0, 32),
		lsHasChild: make([]uint64, 0, 32),
		lsLouds:    make([]uint64, 0, 32),
		hasPrefix:  make([]uint64, 0, 32),
		hasSuffix:  make([]uint64, 0, 32),
		suffixes:   [][]byte{},
		values:     make([]uint32, 0, 32),
	}
}

func (l *Level) GetBitmap(t BitmapType) []uint64 {
	switch t {
	case HasChild:
		return l.lsHasChild
	case Louds:
		return l.lsLouds
	case HasPrefix:
		return l.hasPrefix
	case HasSuffix:
		return l.hasSuffix
	default:
		return []uint64{}
	}
}

func (l *Level) Reset() {
	l.lsLabels = l.lsLabels[:0]
	l.lsHasChild = l.lsHasChild[:0]
	l.lsLouds = l.lsLouds[:0]
	l.hasPrefix = l.hasPrefix[:0]
	l.prefixes = l.prefixes[:0]
	l.hasSuffix = l.hasSuffix[:0]
	l.suffixes = l.suffixes[:0]
	l.values = l.values[:0]
	l.nodeCount = 0
}

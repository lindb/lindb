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
	"encoding"
	"io"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package trie

// Builder represents the succinct trie builder.
type Builder interface {
	// Build returns the SuccinctTrie for added kv pairs.
	// Keys shall be sorted before building.
	Build(keys [][]byte, vals []uint32)
	// Write writes trie data.
	Write(w io.Writer) error
	// Trie returns a succinct trie
	Trie() SuccinctTrie
	// Reset resets the underlying data-structure for next use
	Reset()
}

// SuccinctTrie represents a succinct trie
type SuccinctTrie interface {
	// Get gets the value from trie
	Get(key []byte) (uint32, bool)
	// MarshalSize is the size after padding
	MarshalSize() int64
	encoding.BinaryUnmarshaler
	// NewIterator returns a iterator for arbitrarily iterating the trie
	NewIterator() *Iterator
	// NewPrefixIterator returns a iterator for prefix-iterating the trie
	NewPrefixIterator(prefix []byte) *PrefixIterator
}

package trie

import (
	"encoding"
	"io"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package trie

type Builder interface {
	// Build returns the SuccinctTrie for added kv pairs.
	// Keys shall be sorted before building.
	Build(keys, vals [][]byte, valueWidth uint32) SuccinctTrie

	// Reset resets the underlying data-structure for next use
	Reset()
}

// SuccinctTrie represents a succinct trie
type SuccinctTrie interface {
	// Get gets the value from trie
	Get(key []byte) ([]byte, bool)
	// MarshalSize is the size after padding
	MarshalSize() int64
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	WriteTo(w io.Writer) error
	// NewIterator returns a iterator for arbitrarily iterating the trie
	NewIterator() *Iterator
	// NewPrefixIterator returns a iterator for prefix-iterating the trie
	NewPrefixIterator(prefix []byte) *PrefixIterator
}

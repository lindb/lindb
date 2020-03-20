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
	Get(key []byte) ([]byte, bool)
	MarshalSize() int64
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	WriteTo(w io.Writer) error
	NewIterator() *Iterator
}

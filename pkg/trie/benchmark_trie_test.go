package trie_test

import (
	"testing"

	"github.com/lindb/lindb/pkg/trie"
)

// 79.2ms
func BenchmarkTrie_MarshalBinary(b *testing.B) {
	ips, ranks := newTestIPs()
	builder := trie.NewBuilder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := builder.Build(ips, ranks, 3)
		_, _ = tree.MarshalBinary()
		builder.Reset()
	}
}

// 13.0ms
func BenchmarkTrie_Iterator_NoRead(b *testing.B) {
	ips, ranks := newTestIPs()
	builder := trie.NewBuilder()
	tree := builder.Build(ips, ranks, 3)
	itr := tree.NewIterator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr.SeekToFirst()
		for itr.Valid() {
			itr.Next()
		}
	}
}

// 13.6ms
func BenchmarkTrie_Iterator_Read(b *testing.B) {
	ips, ranks := newTestIPs()
	builder := trie.NewBuilder()
	tree := builder.Build(ips, ranks, 3)
	itr := tree.NewIterator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr.SeekToFirst()
		for itr.Valid() {
			_, _ = itr.Key(), itr.Value()
			itr.Next()
		}
	}
}

// 364ns
func BenchmarkTrie_Get(b *testing.B) {
	ips, ranks := newTestIPs()
	builder := trie.NewBuilder()
	tree := builder.Build(ips, ranks, 3)
	key := []byte("1.1.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Get(key)
	}
}

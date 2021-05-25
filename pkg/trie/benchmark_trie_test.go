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

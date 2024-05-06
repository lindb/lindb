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
	"testing"
)

var ips, ranks = newTestIPs(1 << 8)

// after:  2982368 size 42.2ms (650k ip)
// before: 5488152 size 62.1ms (650k ip)
func BenchmarkTrie_Marshal(b *testing.B) {
	b.StopTimer()
	builder := NewBuilder()
	b.StartTimer()
	var buf = &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		builder.Build(ips, ranks)
		_ = builder.Write(buf)
		buf.Reset()
		builder.Reset()
	}
}

func BenchmarkTrie_Unmarshal(b *testing.B) {
	builder := NewBuilder()

	var buf = &bytes.Buffer{}
	builder.Build(ips, ranks)
	_ = builder.Write(buf)
	data := buf.Bytes()

	tree2 := NewTrie()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tree2.UnmarshalBinary(data)
	}
}

// 13.5ms
func BenchmarkTrie_Iterator_NoRead(b *testing.B) {
	builder := NewBuilder()
	builder.Build(ips, ranks)
	tree := builder.Trie()
	itr := tree.NewIterator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr.SeekToFirst()
		for itr.Valid() {
			itr.Next()
		}
	}
}

// 32.7ms
func BenchmarkTrie_Iterator_Read(b *testing.B) {
	builder := NewBuilder()
	builder.Build(ips, ranks)
	tree := builder.Trie()
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

// 320ns
func BenchmarkTrie_Get(b *testing.B) {
	builder := NewBuilder()
	builder.Build(ips, ranks)
	tree := builder.Trie()
	key := ips[len(ips)-1]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Get(key)
	}
}

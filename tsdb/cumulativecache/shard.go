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

package cumulativecache

import (
	"encoding/binary"
	"sync"

	"github.com/lindb/lindb/pkg/fasttime"
)

type cacheShard struct {
	hashmap map[uint64][]byte
	lock    sync.Mutex
	ttl     uint32 // seconds
	metrics *cacheMetrics
}

func newCacheShard(ttl uint32, metrics *cacheMetrics) *cacheShard {
	return &cacheShard{
		hashmap: map[uint64][]byte{},
		ttl:     ttl,
		metrics: metrics,
	}
}

func (s *cacheShard) getAndSet(hashKey uint64, newValue []byte) (entry []byte, exist bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	wrappedEntry, ok := s.hashmap[hashKey]
	if !ok {
		s.metrics.Misses.Incr()
		s.hashmap[hashKey] = newValue
		s.metrics.Nums.Incr()
		return nil, false
	}

	s.metrics.Hits.Incr()
	// read entry
	entry = readEntry(wrappedEntry)
	s.hashmap[hashKey] = newValue
	return entry, true
}

func (s *cacheShard) cleanUp() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for key, wrappedEntry := range s.hashmap {
		lastTimestamp := readTimestampFromEntry(wrappedEntry)
		if uint32(fasttime.UnixTimestamp())-lastTimestamp >= s.ttl {
			delete(s.hashmap, key)
			s.metrics.Evicts.Incr()
			s.metrics.Nums.Decr()
		}
	}
}

const (
	timestampSizeInBytes = 4 // Number of bytes used for timestamp
)

func readTimestampFromEntry(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func readEntry(data []byte) []byte {
	length := len(data) - timestampSizeInBytes

	// copy on read
	dst := make([]byte, length)
	copy(dst, data[timestampSizeInBytes:])
	return dst
}

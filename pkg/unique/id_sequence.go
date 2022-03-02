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

package unique

import "encoding/binary"

//go:generate mockgen -source ./id_sequence.go -destination=./id_sequence_mock.go -package=unique

// Sequence represents allocate sequence with batch cache, reduces write operations.
// NOTICE: not thread safe.
type Sequence interface {
	// HasNext checks if it has sequence in cache.
	HasNext() bool
	// Next returns next sequence from cache.
	Next() uint32
	// Current returns current sequence from cache.
	Current() uint32
	// Limit sets new limit value.
	Limit(limit uint32)
}

// sequence implements Sequence interface.
type sequence struct {
	val, limit uint32
}

// NewSequence creates a Sequence instance.
func NewSequence(val, limit uint32) Sequence {
	return &sequence{val: val, limit: limit}
}

// HasNext checks if it has sequence in cache.
func (s *sequence) HasNext() bool {
	return s.val < s.limit
}

// Next returns next sequence from cache.
func (s *sequence) Next() uint32 {
	s.val++
	return s.val
}

// Current returns current sequence from cache.
func (s *sequence) Current() uint32 {
	return s.val
}

// Limit sets new limit value.
func (s *sequence) Limit(limit uint32) {
	s.limit = limit
}

// SaveSequence persists current sequence value into store with key.
func SaveSequence(store IDStore, key []byte, seq uint32) error {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], seq)
	return store.Put(key, scratch[:])
}

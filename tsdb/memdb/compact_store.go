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

package memdb

import (
	"sync"
	"unsafe"

	"go.uber.org/atomic"
)

// CompressStore represents memory compress buffer store for field writing.
type CompressStore interface {
	// GetCompressBuffer returns memory compress buffer by memory time series id.
	GetCompressBuffer(memSeriesID uint32) []byte
	// StoreCompressBuffer stores memory compress buffer based on momery time series id.
	StoreCompressBuffer(memSeriesID uint32, buf []byte)
	// MemSize returns compress store memory approximate size.
	MemSize() int64
}

// compressStore implements CompressStore interface.
type compressStore struct {
	store sync.Map // memory series id => compress buffer

	memSize atomic.Int64
}

// NewCompressStore creates CompressStore instance.
func NewCompressStore() CompressStore {
	return &compressStore{}
}

// GetCompressBuffer returns memory compress buffer by memory time series id.
func (s *compressStore) GetCompressBuffer(memSeriesID uint32) []byte {
	buf, ok := s.store.Load(memSeriesID)
	if ok {
		return buf.([]byte)
	}
	return nil
}

// StoreCompressBuffer stores memory compress buffer based on momery time series id.
func (s *compressStore) StoreCompressBuffer(memSeriesID uint32, buf []byte) {
	oldBuf, ok := s.store.Load(memSeriesID)
	var diff int
	if ok {
		diff = len(buf) - len(oldBuf.([]byte))
	} else {
		diff = len(buf) + 4
	}
	s.store.Store(memSeriesID, buf)
	s.memSize.Add(int64(diff))
}

// MemSize returns compress store memory approximate size.
func (s *compressStore) MemSize() (memSize int64) {
	return int64(unsafe.Sizeof(s)) + s.memSize.Load()
}

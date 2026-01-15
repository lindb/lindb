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
	"fmt"
	"sync"

	"github.com/lindb/lindb/pkg/imap"
)

type exemplarBuffer struct {
	ids  *imap.IntMap[ExemplarPage] // store all time series ids(memory time series id => ExemplarPage)
	lock sync.RWMutex
}

func newExemplarBuffer() DataPointBuffer {
	return &exemplarBuffer{
		ids: imap.NewIntMap[ExemplarPage](),
	}
}

func (e *exemplarBuffer) BufferSize() int64 {
	return int64(e.ids.Size())
}

// Close implements [DataPointBuffer].
func (e *exemplarBuffer) Close() error {
	return nil
}

// GetExemplarPage implements [DataPointBuffer].
func (e *exemplarBuffer) GetExemplarPage(memSeriesID uint32) (ExemplarPage, bool) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.ids.Get(memSeriesID)
}

// GetOrCreateExemplarPage implements [DataPointBuffer].
func (e *exemplarBuffer) GetOrCreateExemplarPage(memSeriesID uint32) (ExemplarPage, error) {
	var (
		page ExemplarPage
		ok   bool
	)

	e.lock.RLock()
	page, ok = e.ids.Get(memSeriesID)
	e.lock.RUnlock()
	if ok {
		return page, nil
	}
	// generate a new page
	// NOTE: single goroutine write family data, so can read directly
	page = newExemplarPage()
	e.lock.Lock()
	e.ids.PutIfNotExist(memSeriesID, page)
	e.lock.Unlock()
	return page, nil
}

// GetOrCreatePage implements [DataPointBuffer].
func (e *exemplarBuffer) GetOrCreatePage(memSeriesID uint32) ([]byte, error) {
	panic("exemplar not support page")
}

// GetPage implements [DataPointBuffer].
func (e *exemplarBuffer) GetPage(memSeriesID uint32) ([]byte, bool) {
	panic("exemplar not support page")
}

func (e *exemplarBuffer) IsDirty() bool {
	return false
}

func (e *exemplarBuffer) Release() {
}

type ExemplarPage interface {
	write(slot uint16, traceID, spanID []byte, duration int64) error
}

type exemplar struct {
	traceID  []byte
	spanID   []byte
	duration int64
}

type exemplarPage struct {
	store map[uint16]exemplar
}

func newExemplarPage() ExemplarPage {
	return &exemplarPage{
		store: make(map[uint16]exemplar),
	}
}

func (e *exemplarPage) write(slot uint16, traceID, spanID []byte, duration int64) error {
	fmt.Println("write exemplar")
	e.store[slot] = exemplar{
		traceID:  traceID,
		spanID:   spanID,
		duration: duration,
	}
	return nil
}

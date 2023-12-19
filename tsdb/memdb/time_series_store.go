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

	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/series/field"
)

// tStoreINTF represents time series store for storing field data.
type tStoreINTF interface {
	// Write writes field data.
	Write(timeSeriesID uint32, fType field.Type, slot uint16, fValue float64,
		newFStore func() (fStoreINTF, error)) (writtenSize int, err error)
	// Get returns field store by time series id.
	Get(timeSeriesID uint32) (fStoreINTF, bool)
}

// tStore implements tStoreINTF interface.
type tStore struct {
	stores *imap.IntMap[fStoreINTF] // time series id(memory unique) => field store

	lock sync.RWMutex
}

// newTimeSeriesStore creates time series store.
func newTimeSeriesStore() tStoreINTF {
	return &tStore{
		stores: imap.NewIntMap[fStoreINTF](),
	}
}

// Write writes field data.
func (ts *tStore) Write(timeSeriesID uint32,
	fType field.Type, slot uint16, fValue float64,
	newFStore func() (fStoreINTF, error),
) (written int, err error) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	fStore, ok := ts.stores.Get(timeSeriesID)
	if !ok {
		fStore, err = newFStore()
		if err != nil {
			return 0, err
		}
		size := len(ts.stores.Values())
		ts.stores.Put(timeSeriesID, fStore)
		if len(ts.stores.Values())-size > 0 {
			written += int(IntMapStructValuesEntry)
		}
	}
	before := fStore.Capacity()
	fStore.Write(fType, slot, fValue)
	written += fStore.Capacity() - before
	return
}

// Get returns field store by time series id.
func (ts *tStore) Get(timeSeriesID uint32) (fStoreINTF, bool) {
	ts.lock.RLock()
	defer ts.lock.RUnlock()
	return ts.stores.Get(timeSeriesID)
}

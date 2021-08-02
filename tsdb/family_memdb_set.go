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

package tsdb

import (
	"sort"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/tsdb/memdb"
)

type memDBEntry struct {
	familyTime int64
	memDB      memdb.MemoryDatabase
}

type memDBEntries []memDBEntry

func (entries memDBEntries) Len() int           { return len(entries) }
func (entries memDBEntries) Less(i, j int) bool { return entries[i].familyTime < entries[j].familyTime }
func (entries memDBEntries) Swap(i, j int)      { entries[i], entries[j] = entries[j], entries[i] }

// familyMemDBSet is a immutable data structure in database to provide lock-free lookup
type familyMemDBSet struct {
	value atomic.Value // memDBEntries
}

// newFamilyMemDBSet returns a default empty familyMemDBSet
func newFamilyMemDBSet() *familyMemDBSet {
	set := &familyMemDBSet{}
	// initialize it with a new empty entry slice
	set.value.Store(memDBEntries{})
	return set
}

// InsertFamily inserts a new family into the set
func (ss *familyMemDBSet) InsertFamily(familyTime int64, memDB memdb.MemoryDatabase) {
	oldEntries := ss.value.Load().(memDBEntries)
	var (
		newEntries memDBEntries
		newEntry   = memDBEntry{familyTime: familyTime, memDB: memDB}
	)

	newEntries = make([]memDBEntry, oldEntries.Len()+1)
	copy(newEntries, oldEntries)
	newEntries[len(newEntries)-1] = newEntry
	sort.Sort(newEntries)
	ss.value.Store(newEntries)
}

// GetFamily searches the memDB by familyTime from the familyMemDBSet
func (ss *familyMemDBSet) GetFamily(familyTime int64) (memdb.MemoryDatabase, bool) {
	entries := ss.value.Load().(memDBEntries)
	// fast path when length < 20
	if entries.Len() <= 20 {
		for idx := range entries {
			if entries[idx].familyTime == familyTime {
				return entries[idx].memDB, true
			}
		}
		return nil, false
	}
	index := sort.Search(entries.Len(), func(i int) bool {
		return entries[i].familyTime >= familyTime
	})
	if index < 0 || index >= entries.Len() {
		return nil, false
	}
	return entries[index].memDB, entries[index].familyTime == familyTime
}

func (ss *familyMemDBSet) Entries() memDBEntries {
	return ss.value.Load().(memDBEntries)
}

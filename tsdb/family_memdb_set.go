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
	"go.uber.org/atomic"

	"github.com/lindb/lindb/tsdb/memdb"
)

type memDBEntry struct {
	familyTime int64
	memDB      memdb.MemoryDatabase
}

// memDBSet holds all memory databases
type memDBSet struct {
	mutable   []memDBEntry
	immutable []memDBEntry
}

// New returns a new data structure with samp capacity
func (s *memDBSet) New() *memDBSet {
	return &memDBSet{
		mutable:   make([]memDBEntry, 0, len(s.mutable)),
		immutable: make([]memDBEntry, 0, len(s.immutable)),
	}
}

func newMemDBSet() *memDBSet {
	return &memDBSet{}
}

// familyMemDBSet is a immutable data structure in database to provide lock-free lookup
// operation like inserting or switching memdb from mutable to immutable should be protected with lock
type familyMemDBSet struct {
	value atomic.Value // memDBSet
}

// newFamilyMemDBSet returns a default empty familyMemDBSet
func newFamilyMemDBSet() *familyMemDBSet {
	set := &familyMemDBSet{}
	set.value.Store(*newMemDBSet())
	return set
}

// InsertFamily inserts a new family into the set
func (ss *familyMemDBSet) InsertFamily(familyTime int64, memDB memdb.MemoryDatabase) {
	set := ss.value.Load().(memDBSet)
	newSet := set.New()

	var (
		newMutableEntry = memDBEntry{familyTime: familyTime, memDB: memDB}
	)
	newSet.mutable = append(newSet.mutable, set.mutable...)
	newSet.immutable = append(newSet.immutable, set.immutable...)

	newSet.mutable = append(newSet.mutable, newMutableEntry)

	ss.value.Store(*newSet)
}

// GetMutableFamily searches writable memDB by familyTime from the familyMemDBSet
func (ss *familyMemDBSet) GetMutableFamily(familyTime int64) (memdb.MemoryDatabase, bool) {
	set := ss.value.Load().(memDBSet)
	for idx := range set.mutable {
		if set.mutable[idx].familyTime == familyTime {
			return set.mutable[idx].memDB, true
		}
	}
	return nil, false
}

// SetLargestMutableMemDBImmutable choose a largest mutable memdb, then evict it
func (ss *familyMemDBSet) SetLargestMutableMemDBImmutable() bool {
	mutable := ss.MutableEntries()
	var (
		maxSize    int64
		familyTime int64
	)
	for _, entry := range mutable {
		memdbSize := entry.memDB.MemSize()
		if memdbSize > maxSize {
			familyTime = entry.familyTime
			maxSize = memdbSize
		}
	}
	if maxSize > 0 {
		ss.SetFamilyImmutable(familyTime)
		return true
	}
	return false
}

// SetFamilyImmutable moves a memdb from mutable to immutable
func (ss *familyMemDBSet) SetFamilyImmutable(familyTime int64) {
	oldSet := ss.value.Load().(memDBSet)
	newSet := oldSet.New()

	newSet.immutable = append(newSet.immutable, oldSet.immutable...)

	for _, entry := range oldSet.mutable {
		if entry.familyTime != familyTime {
			newSet.mutable = append(newSet.mutable, entry)
		} else {
			newSet.immutable = append(newSet.immutable, entry)
		}
	}

	// list keeps order as push order
	ss.value.Store(*newSet)
}

// RemoveHeadImmutable removes first immutable memdb after flushed
func (ss *familyMemDBSet) RemoveHeadImmutable() {
	oldSet := ss.value.Load().(memDBSet)
	// we must make sure there is already a immutable memdb in list
	newSet := oldSet.New()
	newSet.mutable = oldSet.mutable

	for idx := 1; idx < len(oldSet.immutable); idx++ {
		newSet.immutable = append(newSet.immutable, oldSet.immutable[idx])
	}
	ss.value.Store(*newSet)

}

// Entries returns all mutable and immutable memdb
func (ss *familyMemDBSet) Entries() []memDBEntry {
	set := ss.value.Load().(memDBSet)
	var dst = make([]memDBEntry, len(set.mutable)+len(set.immutable))
	copy(dst, set.immutable)
	copy(dst[len(set.immutable):], set.mutable)
	return dst
}

func (ss *familyMemDBSet) MutableEntries() []memDBEntry {
	set := ss.value.Load().(memDBSet)
	return set.mutable
}

func (ss *familyMemDBSet) ImmutableEntries() []memDBEntry {
	set := ss.value.Load().(memDBSet)
	return set.immutable
}

func (ss *familyMemDBSet) TotalSize() int64 {
	set := ss.value.Load().(memDBSet)
	var size int64
	for _, entry := range set.mutable {
		size += entry.memDB.MemSize()
	}
	for _, entry := range set.immutable {
		size += entry.memDB.MemSize()
	}
	return size
}

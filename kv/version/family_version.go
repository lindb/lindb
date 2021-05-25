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

package version

import (
	"sync"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./family_version.go -destination=./family_version_mock.go -package=version

type FamilyID int

// Int returns int value of family id
func (f FamilyID) Int() int {
	return int(f)
}

// Int returns int32 value of family id
func (f FamilyID) Int32() int32 {
	return int32(f)
}

type FamilyVersion interface {
	// GetID returns the family id
	GetID() FamilyID
	// GetVersionSet returns the store version set
	GetVersionSet() StoreVersionSet
	// GetAllActiveFiles returns all files based on all active versions
	GetAllActiveFiles() []*FileMeta
	// GetSnapshot returns the current version's snapshot
	GetSnapshot() Snapshot
	// GetLiveRollupFiles returns all need rollup files
	GetLiveRollupFiles() map[table.FileNumber]timeutil.Interval
	// GetLiveReferenceFiles returns all rollup reference files
	GetLiveReferenceFiles() map[FamilyID][]table.FileNumber
	// removeVersion removes version from active versions
	removeVersion(v Version)
	// appendVersion swaps family's current version, then releases previous version
	appendVersion(v Version)
}

// familyVersion maintains family level metadata
type familyVersion struct {
	ID         FamilyID
	familyName string
	versionSet StoreVersionSet

	current        Version           // current mutable version
	activeVersions map[int64]Version // all active versions include mutable/immutable versions

	mutex sync.RWMutex
}

// newFamilyVersion new FamilyVersion instance
func newFamilyVersion(familyID FamilyID, familyName string, versionSet StoreVersionSet) FamilyVersion {
	fv := &familyVersion{
		ID:             familyID,
		familyName:     familyName,
		versionSet:     versionSet,
		activeVersions: make(map[int64]Version),
	}
	// create new version for current mutable version
	current := newVersion(fv.versionSet.newVersionID(), fv)
	fv.activeVersions[current.ID()] = current
	fv.current = current
	return fv
}

// GetID returns the family id
func (fv *familyVersion) GetID() FamilyID {
	return fv.ID
}

// GetVersionSet returns the store version set
func (fv *familyVersion) GetVersionSet() StoreVersionSet {
	return fv.versionSet
}

// GetSnapshot returns the current version's snapshot
func (fv *familyVersion) GetSnapshot() Snapshot {
	fv.mutex.RLock()
	defer fv.mutex.RUnlock()
	return newSnapshot(fv.familyName, fv.current, fv.versionSet.getCache())
}

// GetAllActiveFiles returns all files based on all active versions
func (fv *familyVersion) GetAllActiveFiles() []*FileMeta {
	fv.mutex.RLock()
	defer fv.mutex.RUnlock()

	var files []*FileMeta
	var fileNumbers = make(map[table.FileNumber]table.FileNumber)
	for _, version := range fv.activeVersions {
		versionFiles := version.GetAllFiles()
		for _, file := range versionFiles {
			fileNumber := file.fileNumber
			// remove duplicate file in diff versions
			_, ok := fileNumbers[fileNumber]
			if !ok {
				files = append(files, file)
				fileNumbers[fileNumber] = fileNumber
			}
		}
	}
	return files
}

// GetLiveRollupFiles returns all need rollup files
func (fv *familyVersion) GetLiveRollupFiles() map[table.FileNumber]timeutil.Interval {
	fv.mutex.RLock()
	defer fv.mutex.RUnlock()
	return fv.current.GetRollupFiles()
}

// GetLiveReferenceFiles returns all rollup reference files
func (fv *familyVersion) GetLiveReferenceFiles() map[FamilyID][]table.FileNumber {
	fv.mutex.RLock()
	defer fv.mutex.RUnlock()
	return fv.current.GetReferenceFiles()
}

// removeVersion removes version from active versions,
// cannot remove current version from active versions.
func (fv *familyVersion) removeVersion(v Version) {
	fv.mutex.Lock()
	if v != fv.current {
		delete(fv.activeVersions, v.ID())
	}
	fv.mutex.Unlock()
}

// appendVersion swaps family's current version, then releases previous version
func (fv *familyVersion) appendVersion(v Version) {
	previous := fv.current

	fv.mutex.Lock()
	fv.activeVersions[v.ID()] = v
	fv.current = v
	fv.mutex.Unlock()

	if previous != nil && previous.NumOfRef() == 0 {
		// remove version from family active versions
		v.GetFamilyVersion().removeVersion(previous)
	}
}

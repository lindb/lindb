package version

import (
	"sync"
)

//go:generate mockgen -source=./family_version.go -destination=./family_version_mock.go -package=version

type FamilyVersion interface {
	// GetID returns the family id
	GetID() int
	// GetVersionSet returns the store version set
	GetVersionSet() StoreVersionSet
	// GetAllActiveFiles returns all files based on all active versions
	GetAllActiveFiles() []*FileMeta
	// GetSnapshot returns the current version's snapshot
	GetSnapshot() Snapshot

	// removeVersion removes version from active versions
	removeVersion(v *Version)
	// appendVersion swaps family's current version, then releases previous version
	appendVersion(v *Version)
}

// familyVersion maintains family level metadata
type familyVersion struct {
	ID         int
	familyName string
	versionSet StoreVersionSet

	current        *Version           // current mutable version
	activeVersions map[int64]*Version // all active versions include mutable/immutable versions

	mutex sync.RWMutex
}

// newFamilyVersion new FamilyVersion instance
func newFamilyVersion(familyID int, familyName string, versionSet StoreVersionSet) FamilyVersion {
	fv := &familyVersion{
		ID:             familyID,
		familyName:     familyName,
		versionSet:     versionSet,
		activeVersions: make(map[int64]*Version),
	}
	// create new version for current mutable version
	current := newVersion(fv.versionSet.newVersionID(), fv)
	fv.activeVersions[current.id] = current
	fv.current = current
	return fv
}

// GetID returns the family id
func (fv *familyVersion) GetID() int {
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
	var files []*FileMeta
	var fileNumbers = make(map[int64]int64)
	for _, version := range fv.activeVersions {
		versionFiles := version.getAllFiles()
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

// removeVersion removes version from active versions,
// cannot remove current version from active versions.
func (fv *familyVersion) removeVersion(v *Version) {
	fv.mutex.Lock()
	if v != fv.current {
		delete(fv.activeVersions, v.id)
	}
	fv.mutex.Unlock()
}

// appendVersion swaps family's current version, then releases previous version
func (fv *familyVersion) appendVersion(v *Version) {
	previous := fv.current

	fv.mutex.Lock()
	fv.activeVersions[v.id] = v
	fv.current = v
	fv.mutex.Unlock()

	if previous != nil && previous.numOfRef() == 0 {
		// remove version from family active versions
		v.fv.removeVersion(previous)
	}
}

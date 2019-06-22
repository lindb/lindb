package version

import "sync"

// FamilyVersion maintains family level metadata
type FamilyVersion struct {
	versionSet *StoreVersionSet

	current        *Version           // current mutable version
	activeVersions map[int64]*Version // all active versions include mutable/immutable versions

	mutex sync.RWMutex
}

// newFamilyVersion new FamilyVersion instance
func newFamilyVersion(versionSet *StoreVersionSet) *FamilyVersion {
	fv := &FamilyVersion{
		versionSet:     versionSet,
		activeVersions: make(map[int64]*Version),
	}
	// create new version for current mutable version
	current := newVersion(fv.versionSet.newVersionID(), fv)
	fv.activeVersions[current.id] = current
	fv.current = current
	return fv
}

// GetCurrent returns current mutable version
func (fv *FamilyVersion) GetCurrent() *Version {
	fv.mutex.RLock()
	defer fv.mutex.RUnlock()
	// inc ref count of version
	fv.current.retain()
	return fv.current
}

// FindFiles finds all files include key from current's level,
// must return files related version, and retain it, release version after read data.
func (fv *FamilyVersion) FindFiles(key uint32) (*Version, []*FileMeta) {
	fv.mutex.RLock()
	current := fv.current
	// must retain it, don't release util finish read, release it during snapshot's closing.
	current.retain()
	// find files related given key
	files := current.findFiles(key)
	fv.mutex.RUnlock()
	return current, files
}

// GetAllFiles returns all files based on all active versions
func (fv *FamilyVersion) GetAllFiles() []*FileMeta {
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

// removeVersion removes version from active versions
func (fv *FamilyVersion) removeVersion(v *Version) {
	fv.mutex.Lock()
	delete(fv.activeVersions, v.id)
	fv.mutex.Unlock()
}

// appendVersion swaps family's current version, then releases previous version
func (fv *FamilyVersion) appendVersion(v *Version) {
	previous := fv.current

	fv.mutex.Lock()
	fv.activeVersions[v.id] = v
	fv.current = v
	fv.mutex.Unlock()

	if previous != nil {
		previous.Release()
	}
}

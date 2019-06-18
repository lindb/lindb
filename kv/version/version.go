package version

import (
	"sync/atomic"
)

// Version is snapshot for current storage metadata includes levels/sst files
type Version struct {
	id          int64 // unique id in kv store level
	numOfLevels int   // num of levels
	fv          *FamilyVersion

	ref int32 // current version ref count for using

	levels []*level // each level sst files exclude level0
}

// newVersion new Version instance
func newVersion(id int64, fv *FamilyVersion) *Version {
	numOfLevel := fv.versionSet.numOfLevels
	if numOfLevel <= 0 {
		panic("num of levels cannot be <=0")
	}
	v := &Version{
		id:          id,
		fv:          fv,
		numOfLevels: numOfLevel,
	}
	v.levels = make([]*level, numOfLevel)
	for i := 0; i < numOfLevel; i++ {
		v.levels[i] = newLevel()
	}
	return v
}

// Release decrements version ref count,
// if ref==0, then remove current version from list of family level.
func (v *Version) Release() {
	val := atomic.AddInt32(&v.ref, -1)
	if val == 0 {
		// remove version from family active versions
		v.fv.removeVersion(v)
	}
}

// findFiles finds all files include key from each level
func (v *Version) findFiles(key uint32) []*FileMeta {
	var files []*FileMeta
	for _, level := range v.levels {
		for _, file := range level.getFiles() {
			if key >= file.minKey && key <= file.maxKey {
				files = append(files, file)
			}
		}
	}
	return files
}

// getAllFilesetAllFiles returns all ative files of each level
func (v *Version) getAllFiles() []*FileMeta {
	var files []*FileMeta
	for _, value := range v.levels {
		files = append(files, value.getFiles()...)
	}
	return files
}

// retain increments version ref count
func (v *Version) retain() {
	atomic.AddInt32(&v.ref, 1)
}

// cloneVersion builds new version based on current version
func (v *Version) cloneVersion() *Version {
	newVersion := newVersion(v.fv.versionSet.newVersionID(), v.fv)
	for level, value := range v.levels {
		for _, file := range value.files {
			newVersion.addFile(level, file)
		}
	}
	return newVersion
}

// addFiles adds file meta into spec level
func (v *Version) addFiles(level int, files []*FileMeta) {
	v.levels[level].addFiles(files...)
}

// addFile adds file meta into spec level
func (v *Version) addFile(level int, file *FileMeta) {
	v.levels[level].addFile(file)
}

// deleteFile delete file from spec level file list
func (v *Version) deleteFile(level int, fileNumber int64) {
	v.levels[level].deleteFile(fileNumber)
}

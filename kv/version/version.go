package version

import (
	"sync/atomic"
)

// Version is snapshot for current storage metadata includes levels/sst files
type Version struct {
	id          int64 // unique id in kv store level
	numOfLevels int   // num of levels
	fv          FamilyVersion

	ref int32 // current version ref count for using

	levels []*level // each level sst files exclude level0
}

// newVersion new Version instance
func newVersion(id int64, fv FamilyVersion) *Version {
	numOfLevel := fv.GetVersionSet().numberOfLevels()
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

// GetFamilyVersion return the family version
func (v *Version) GetFamilyVersion() FamilyVersion {
	return v.fv
}

// retain increments version ref count
func (v *Version) retain() {
	atomic.AddInt32(&v.ref, 1)
}

// release decrements version ref count,
// if ref==0, then remove current version from list of family level.
func (v *Version) release() {
	newVal := atomic.AddInt32(&v.ref, -1)
	if newVal == 0 {
		v.fv.removeVersion(v)
	}
}

// numOfRef returns the ref count
func (v *Version) numOfRef() int32 {
	return atomic.LoadInt32(&v.ref)
}

// NumberOfFilesInLevel returns the number of files by spec level,
// if level > numOfLevels return 0.
func (v *Version) NumberOfFilesInLevel(level int) int {
	if level < 0 || level > v.numOfLevels {
		return 0
	}
	return v.levels[level].numberOfFiles()
}

// PickL0Compaction picks level0 compaction context,
// if hasn't congruent compaction return nil.
func (v *Version) PickL0Compaction(compactThreshold int) *Compaction {
	// We prefer compactions triggered by too much data level 0 over the compactions triggered by seeks.
	if v.NumberOfFilesInLevel(0) < compactThreshold {
		return nil
	}
	var levelInputs []*FileMeta
	// pick the level one files to do compaction
	levelInputs = append(levelInputs, v.getFiles(0)...)
	/*
	 * Get over lapping input from level 1, based on level 0 key range.
	 * Only pick over lapping file, not use key range in all files for level 0, maybe happen overhead for reading.
	 * for example:
	 * Level 0:
	 * file 1: 1~10
	 * file 2: 1000~1001
	 *
	 * Level 1:
	 * file 3: 1~5
	 * file 4: 100~200
	 * file 5: 400~500
	 *
	 * if use key for all files in level 0, final key is 1~1001, pick level 1 files is 3,4,5.
	 * if use key for each files in level 0, final pick level 1 files is 3.
	 */
	levelUpInputMap := make(map[int64]*FileMeta)
	for _, lowInput := range levelInputs {
		upInputs := v.getOverlappingInputs(1, lowInput.GetMinKey(), lowInput.GetMaxKey())
		for _, upInput := range upInputs {
			levelUpInputMap[upInput.GetFileNumber()] = upInput
		}
	}
	var levelUpInputs []*FileMeta
	for _, upInput := range levelUpInputMap {
		levelUpInputs = append(levelUpInputs, upInput)
	}
	return NewCompaction(v.fv.GetID(), 0, levelInputs, levelUpInputs)
}

// findFiles finds all files include key from each level
func (v *Version) findFiles(key uint32) []*FileMeta {
	var files []*FileMeta
	for _, level := range v.levels {
		for _, file := range level.getFiles() {
			if key >= file.GetMinKey() && key <= file.GetMaxKey() {
				files = append(files, file)
			}
		}
	}
	return files
}

// getAllFilesetAllFiles returns all active files of each level
func (v *Version) getAllFiles() []*FileMeta {
	var files []*FileMeta
	for _, value := range v.levels {
		files = append(files, value.getFiles()...)
	}
	return files
}

// cloneVersion builds new version based on current version
func (v *Version) cloneVersion() *Version {
	newVersion := newVersion(v.fv.GetVersionSet().newVersionID(), v.fv)
	for level, value := range v.levels {
		for _, file := range value.files {
			newVersion.addFile(level, file)
		}
	}
	return newVersion
}

// getFiles returns files by level
func (v *Version) getFiles(level int) []*FileMeta {
	if level < 0 || level >= v.numOfLevels {
		return nil
	}
	return v.levels[level].getFiles()
}

// addFiles adds file meta into spec level
func (v *Version) addFiles(level int, files []*FileMeta) {
	if level >= 0 && level < v.numOfLevels {
		v.levels[level].addFiles(files...)
	}
}

// addFile adds file meta into spec level
func (v *Version) addFile(level int, file *FileMeta) {
	if level >= 0 && level < v.numOfLevels {
		v.levels[level].addFile(file)
	}
}

// deleteFile deletes file from spec level file list
func (v *Version) deleteFile(level int, fileNumber int64) {
	if level >= 0 && level < v.numOfLevels {
		v.levels[level].deleteFile(fileNumber)
	}
}

// getOverlappingInputs gets overlapping input based on level and key range,
// returns the over lapping th given level for key range.
func (v *Version) getOverlappingInputs(level int, minKey, maxKey uint32) []*FileMeta {
	var results []*FileMeta
	files := v.getFiles(level)
	for idx := range files {
		fileMeta := files[idx]
		if fileMeta.GetMaxKey() < minKey || fileMeta.GetMinKey() > maxKey {
			// Either completely before or after range; skip it
			continue
		}
		results = append(results, fileMeta)
	}
	return results
}

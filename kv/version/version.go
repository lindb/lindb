package version

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./version.go -destination=./version_mock.go -package=version

type Version interface {
	// ID returns version id
	ID() int64
	// AddFile adds file meta into spec level
	AddFile(level int, file *FileMeta)
	// AddFiles adds file meta into spec level
	AddFiles(level int, files []*FileMeta)
	// DeleteFile deletes file from spec level file list
	DeleteFile(level int, fileNumber table.FileNumber)
	// GetFiles returns files by level
	GetFiles(level int) []*FileMeta
	// GetFamilyVersion return the family version
	GetFamilyVersion() FamilyVersion
	// NumOfRef returns the number of reference which version be used by search/compact/rollup
	NumOfRef() int32
	// retain increments version ref count
	Retain()
	// Release decrements version ref count,
	// if ref==0, then remove current version from list of family level.
	Release()
	// FindFiles finds all files include key from each level
	FindFiles(key uint32) []*FileMeta
	// GetAllFiles returns all active files of each level
	GetAllFiles() []*FileMeta
	// Clone builds new version based on current version
	Clone() Version
	// Levels returns the files in each level
	Levels() []*level
	// NumberOfFilesInLevel returns the number of files by spec level,
	// if level > numOfLevels return 0.
	NumberOfFilesInLevel(level int) int
	// PickL0Compaction picks level0 compaction context,
	// if hasn't congruent compaction return nil.
	PickL0Compaction(compactThreshold int) *Compaction

	// AddRollupFile adds need rollup file and target interval
	AddRollupFile(fileNumber table.FileNumber, interval timeutil.Interval)
	// DeleteRollupFile removes rollup file after rollup job complete successfully
	DeleteRollupFile(fileNumber table.FileNumber)
	// AddReferenceFile adds rollup reference file under target family
	AddReferenceFile(familyID FamilyID, fileNumber table.FileNumber)
	// DeleteReferenceFile removes rollup reference file under target family
	DeleteReferenceFile(familyID FamilyID, fileNumber table.FileNumber)
	// GetRollupFiles returns all need rollup files
	GetRollupFiles() map[table.FileNumber]timeutil.Interval
	// GetReferenceFiles returns the reference files under target family
	GetReferenceFiles() map[FamilyID][]table.FileNumber
}

// version is snapshot for current storage metadata includes levels/sst files
type version struct {
	id          int64 // unique id in kv store level
	numOfLevels int   // num of levels
	fv          FamilyVersion
	ref         atomic.Int32 // current version ref count for using
	rollup      *rollup

	levels []*level // each level sst files exclude level0
}

// newVersion new version instance
func newVersion(id int64, fv FamilyVersion) Version {
	numOfLevel := fv.GetVersionSet().numberOfLevels()
	if numOfLevel <= 0 {
		panic("num of levels cannot be <=0")
	}
	v := &version{
		id:          id,
		fv:          fv,
		numOfLevels: numOfLevel,
		rollup:      newRollup(),
	}
	v.levels = make([]*level, numOfLevel)
	for i := 0; i < numOfLevel; i++ {
		v.levels[i] = newLevel()
	}
	return v
}

// ID returns version id
func (v *version) ID() int64 {
	return v.id
}

// Levels returns the files in each level
func (v *version) Levels() []*level {
	return v.levels
}

// GetFamilyVersion return the family version
func (v *version) GetFamilyVersion() FamilyVersion {
	return v.fv
}

// NumOfRef returns the number of reference which version be used by search/compact/rollup
func (v *version) NumOfRef() int32 {
	return v.ref.Load()
}

// retain increments version ref count
func (v *version) Retain() {
	v.ref.Inc()
}

// Release decrements version ref count,
// if ref==0, then remove current version from list of family level.
func (v *version) Release() {
	newVal := v.ref.Dec()
	if newVal == 0 {
		v.fv.removeVersion(v)
	}
}

// NumberOfFilesInLevel returns the number of files by spec level,
// if level > numOfLevels return 0.
func (v *version) NumberOfFilesInLevel(level int) int {
	if level < 0 || level > v.numOfLevels {
		return 0
	}
	return v.levels[level].numberOfFiles()
}

// PickL0Compaction picks level0 compaction context,
// if hasn't congruent compaction return nil.
func (v *version) PickL0Compaction(compactThreshold int) *Compaction {
	// We prefer compactions triggered by too much data level 0 over the compactions triggered by seeks.
	if v.NumberOfFilesInLevel(0) < compactThreshold {
		return nil
	}
	var levelInputs []*FileMeta
	// pick the level one files to do compaction
	levelInputs = append(levelInputs, v.GetFiles(0)...)
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
	levelUpInputMap := make(map[table.FileNumber]*FileMeta)
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

// FindFiles finds all files include key from each level
func (v *version) FindFiles(key uint32) []*FileMeta {
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

// GetAllFiles returns all active files of each level
func (v *version) GetAllFiles() []*FileMeta {
	var files []*FileMeta
	for _, value := range v.levels {
		files = append(files, value.getFiles()...)
	}
	return files
}

// Clone builds new version based on current version
func (v *version) Clone() Version {
	newVersion := newVersion(v.fv.GetVersionSet().newVersionID(), v.fv)
	for level, value := range v.levels {
		for _, file := range value.files {
			newVersion.AddFile(level, file)
		}
	}
	return newVersion
}

// GetFiles returns files by level
func (v *version) GetFiles(level int) []*FileMeta {
	if level < 0 || level >= v.numOfLevels {
		return nil
	}
	return v.levels[level].getFiles()
}

// AddFiles adds file meta into spec level
func (v *version) AddFiles(level int, files []*FileMeta) {
	if level >= 0 && level < v.numOfLevels {
		v.levels[level].addFiles(files...)
	}
}

// AddFile adds file meta into spec level
func (v *version) AddFile(level int, file *FileMeta) {
	if level >= 0 && level < v.numOfLevels {
		v.levels[level].addFile(file)
	}
}

// DeleteFile deletes file from spec level file list
func (v *version) DeleteFile(level int, fileNumber table.FileNumber) {
	if level >= 0 && level < v.numOfLevels {
		v.levels[level].deleteFile(fileNumber)
	}
}

// AddRollupFile adds need rollup file and target interval
func (v *version) AddRollupFile(fileNumber table.FileNumber, interval timeutil.Interval) {
	v.rollup.addRollupFile(fileNumber, interval)
}

// DeleteRollupFile removes rollup file after rollup job complete successfully
func (v *version) DeleteRollupFile(fileNumber table.FileNumber) {
	v.rollup.removeRollupFile(fileNumber)
}

// AddReferenceFile adds rollup reference file under target family
func (v *version) AddReferenceFile(familyID FamilyID, fileNumber table.FileNumber) {
	v.rollup.addReferenceFile(familyID, fileNumber)
}

// DeleteReferenceFile removes rollup reference file under target family
func (v *version) DeleteReferenceFile(familyID FamilyID, fileNumber table.FileNumber) {
	v.rollup.removeReferenceFile(familyID, fileNumber)
}

// GetRollupFiles returns all need rollup files
func (v *version) GetRollupFiles() map[table.FileNumber]timeutil.Interval {
	return v.rollup.getRollupFiles()
}

// GetReferenceFiles returns the reference files under target family
func (v *version) GetReferenceFiles() map[FamilyID][]table.FileNumber {
	return v.rollup.getReferenceFiles()
}

// getOverlappingInputs gets overlapping input based on level and key range,
// returns the over lapping th given level for key range.
func (v *version) getOverlappingInputs(level int, minKey, maxKey uint32) []*FileMeta {
	var results []*FileMeta
	files := v.GetFiles(level)
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

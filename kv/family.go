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

package kv

import (
	"fmt"
	"path/filepath"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./family.go -destination=./family_mock.go -package kv

// for testing
var (
	newCompactJobFunc = newCompactJob
	removeDirFunc     = fileutil.RemoveDir
)

// Family implements column family for data isolation each family.
type Family interface {
	// ID return family's id
	ID() version.FamilyID
	// Name return family's name
	Name() string
	// NewFlusher creates flusher for saving data to family.
	NewFlusher() Flusher
	// GetSnapshot returns current version's snapshot
	GetSnapshot() version.Snapshot
	// familyInfo return family info
	familyInfo() string

	// getFamilyVersion returns the family version
	getFamilyVersion() version.FamilyVersion
	// commitEditLog persists edit logs into manifest file.
	commitEditLog(editLog version.EditLog) bool
	// newTableBuilder creates table builder instance for storing kv data.
	newTableBuilder() (table.Builder, error)
	// needCompat returns level0 files if need do compact job
	needCompat() bool
	// compact does compaction job
	compact()
	// getNewMerger returns new merger function, merger need implement Merger interface
	getNewMerger() NewMerger
	// addPendingOutput add a file which current writing file number
	addPendingOutput(fileNumber table.FileNumber)
	// removePendingOutput removes pending output file after compact or flush
	removePendingOutput(fileNumber table.FileNumber)
	// doRollupWork does rollup job, merge source family data to target family
	doRollupWork(sourceFamily Family, rollup Rollup, sourceFiles []table.FileNumber) (err error)

	// deleteObsoleteFiles deletes obsolete files
	deleteObsoleteFiles()
}

// family implements Family interface
type family struct {
	store         Store
	name          string
	familyPath    string
	option        FamilyOption
	merger        NewMerger
	familyVersion version.FamilyVersion
	maxFileSize   int32

	pendingOutputs    sync.Map
	newCompactJobFunc func(family Family, state *compactionState, rollup Rollup) CompactJob

	rolluping  atomic.Bool
	compacting atomic.Bool
}

// newFamily creates new family or open existed family.
func newFamily(store Store, option FamilyOption) (Family, error) {
	name := option.Name

	familyPath := filepath.Join(store.Option().Path, name)

	if !fileutil.Exist(familyPath) {
		if err := mkDirFunc(familyPath); err != nil {
			return nil, fmt.Errorf("mkdir family path error:%w", err)
		}
	}
	merger, ok := mergers[MergerType(option.Merger)]
	if !ok {
		return nil, fmt.Errorf("merger of option not impelement Merger interface, merger is [%s]", option.Merger)
	}
	maxFileSize := defaultMaxFileSize
	if option.MaxFileSize > 0 {
		maxFileSize = option.MaxFileSize
	}

	f := &family{
		familyPath:        familyPath,
		store:             store,
		name:              name,
		option:            option,
		merger:            merger,
		maxFileSize:       maxFileSize,
		newCompactJobFunc: newCompactJobFunc,
		familyVersion:     store.createFamilyVersion(name, version.FamilyID(option.ID)),
	}

	kvLogger.Info("create new family successfully", logger.String("family", f.familyInfo()))
	return f, nil
}

// ID return family's id
func (f *family) ID() version.FamilyID {
	return version.FamilyID(f.option.ID)
}

// Name return family's name
func (f *family) Name() string {
	return f.name
}

// NewFlusher creates flusher for saving data to family.
func (f *family) NewFlusher() Flusher {
	return newStoreFlusher(f)
}

// GetSnapshot returns current version's snapshot
func (f *family) GetSnapshot() version.Snapshot {
	return f.familyVersion.GetSnapshot()
}

// familyInfo return family info
func (f *family) familyInfo() string {
	return f.familyPath
}

// newTableBuilder creates table builder instance for storing kv data.
func (f *family) newTableBuilder() (table.Builder, error) {
	fileNumber := f.store.nextFileNumber()
	fileName := filepath.Join(f.familyPath, version.Table(fileNumber))
	return table.NewStoreBuilder(fileNumber, fileName)
}

// commitEditLog persists edit logs into manifest file.
// returns true on committing successfully and false on failure
func (f *family) commitEditLog(editLog version.EditLog) bool {
	if editLog == nil || editLog.IsEmpty() {
		kvLogger.Warn("edit log is empty", logger.String("family", f.familyInfo()))
		return true
	}
	if err := f.store.commitFamilyEditLog(f.name, editLog); err != nil {
		kvLogger.Error("commit edit log error:", logger.String("family", f.familyInfo()), logger.Error(err))
		return false
	}
	return true
}

// needCompat returns level0 files if need do compact job
func (f *family) needCompat() bool {
	// has compaction job doing
	if f.compacting.Load() {
		return false
	}

	snapshot := f.GetSnapshot()
	defer snapshot.Close()
	threshold := f.option.CompactThreshold
	if threshold <= 0 {
		threshold = defaultCompactThreshold
	}

	numberOfFiles := snapshot.GetCurrent().NumberOfFilesInLevel(0)
	if numberOfFiles > 0 && numberOfFiles >= threshold {
		kvLogger.Info("need to compact level0 files", logger.String("family", f.familyInfo()),
			logger.Any("numOfFiles", numberOfFiles), logger.Any("threshold", f.option.CompactThreshold))
		return true
	}
	return false
}

// compact does compact job if hasn't compact job running
func (f *family) compact() {
	if f.compacting.CAS(false, true) {
		go func() {
			defer f.compacting.Store(false)

			if err := f.backgroundCompactionJob(); err != nil {
				kvLogger.Error("do compact job error",
					logger.String("family", f.familyInfo()), logger.Error(err), logger.Stack())
			}
		}()
	}
}

// backgroundCompactionJob runs compact job in background goroutine
func (f *family) backgroundCompactionJob() error {
	snapshot := f.GetSnapshot()
	defer func() {
		snapshot.Close()
		// clean up unused files, maybe some file not used
		f.deleteObsoleteFiles()
	}()

	compaction := snapshot.GetCurrent().PickL0Compaction(f.option.CompactThreshold)
	if compaction == nil {
		// no compaction job need to do
		return nil
	}
	compactionState := newCompactionState(f.maxFileSize, snapshot, compaction)
	compactJob := f.newCompactJobFunc(f, compactionState, nil)
	if err := compactJob.Run(); err != nil {
		return err
	}
	return nil
}

// addPendingOutput add a file which current writing file number
func (f *family) addPendingOutput(fileNumber table.FileNumber) {
	f.pendingOutputs.Store(fileNumber, dummy)
}

// removePendingOutput removes pending output file after compact or flush
func (f *family) removePendingOutput(fileNumber table.FileNumber) {
	f.pendingOutputs.Delete(fileNumber)
}

// deleteSST deletes the temp sst file, if flush or compact fail
func (f *family) deleteSST(fileNumber table.FileNumber) error {
	if err := removeDirFunc(filepath.Join(f.familyPath, version.Table(fileNumber))); err != nil {
		return err
	}
	return nil
}

// getFamilyVersion returns the family version
func (f *family) getFamilyVersion() version.FamilyVersion {
	return f.familyVersion
}

// getNewMerger returns new merger function, merger need implement Merger interface
func (f *family) getNewMerger() NewMerger {
	return f.merger
}

// deleteObsoleteFiles deletes obsolete files
func (f *family) deleteObsoleteFiles() {
	sstFiles, err := listDirFunc(f.familyPath)
	if err != nil {
		kvLogger.Error("list sst file fail when delete obsolete files", logger.String("family", f.familyInfo()))
		return
	}
	// make a map for all live files
	liveFiles := make(map[table.FileNumber]string)
	f.pendingOutputs.Range(func(key, value interface{}) bool {
		k, ok := key.(table.FileNumber)
		if ok {
			liveFiles[k] = dummy
		}
		return true
	})
	// add live files
	allLiveSSTFiles := f.familyVersion.GetAllActiveFiles()
	for idx := range allLiveSSTFiles {
		liveFiles[allLiveSSTFiles[idx].GetFileNumber()] = dummy
	}
	// add live rollup files, maybe some rollup files is not alive in current family version,
	// but those files cannot delete, because need read those files when do rollup job
	rollupFiles := f.familyVersion.GetLiveRollupFiles()
	for file := range rollupFiles {
		liveFiles[file] = dummy
	}
	for _, fileName := range sstFiles {
		fileDesc := version.ParseFileName(fileName)
		if fileDesc == nil {
			continue
		}
		keep := true
		fileNumber := fileDesc.FileNumber
		if fileDesc.FileType == version.TypeTable {
			_, keep = liveFiles[fileNumber]
		}
		if !keep {
			f.store.evictFamilyFile(f.name, fileNumber)
			if err := f.deleteSST(fileNumber); err != nil {
				kvLogger.Error("delete sst file fail",
					logger.String("family", f.familyInfo()), logger.Any("fileNumber", fileNumber))
			} else {
				kvLogger.Info("delete sst file successfully",
					logger.String("family", f.familyInfo()), logger.Any("fileNumber", fileNumber))
			}
		}
	}
}

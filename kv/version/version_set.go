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
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/bufioutil"
)

//go:generate mockgen -source=./version_set.go -destination=./version_set_mock.go -package=version

// for test
var (
	writeFileFunc       = os.WriteFile
	readFileFunc        = os.ReadFile
	renameFunc          = os.Rename
	newBufferReaderFunc = bufioutil.NewBufioEntryReader
	newBufferWriterFunc = bufioutil.NewBufioEntryWriter
	newEmptyEditLogFunc = newEmptyEditLog
)

// StoreVersionSet maintains all metadata for kv store
type StoreVersionSet interface {
	// Recover recover version set if exist, recover been invoked when kv store init.
	Recover() error
	// Destroy closes version set, release resource, such as journal writer etc.
	Destroy() error
	// NextFileNumber generates next file number
	NextFileNumber() table.FileNumber
	// ManifestFileNumber returns the current manifest file number
	ManifestFileNumber() table.FileNumber
	// CommitFamilyEditLog persists edit logs to manifest file, then apply new version to family version
	CommitFamilyEditLog(family string, editLog EditLog) error
	// CreateFamilyVersion creates family version using family name,
	// if family version exist, return exist one
	CreateFamilyVersion(family string, familyID FamilyID) FamilyVersion
	// GetFamilyVersion returns family version if it existed, else return nil
	GetFamilyVersion(family string) FamilyVersion

	// newVersionID generates new version id
	newVersionID() int64
	// setNextFileNumberWithoutLock set next file number, invoker must add lock
	setNextFileNumberWithoutLock(newNextFileNumber table.FileNumber)
	// numberOfLevels returns num. of levels
	numberOfLevels() int
	// getCache returns table cache for reading data
	getCache() table.Cache
}

// storeVersionSet maintains all metadata for kv store
type storeVersionSet struct {
	manifestFileNumber atomic.Int64
	nextFileNumber     atomic.Int64
	storePath          string
	familyVersions     map[string]FamilyVersion
	familyIDs          map[FamilyID]string
	versionID          atomic.Int64 // unique in for increasing version id
	storeCache         table.Cache

	numOfLevels int // num of levels

	manifest bufioutil.BufioWriter
	mutex    sync.RWMutex
}

// NewStoreVersionSet new VersionSet instance
func NewStoreVersionSet(storePath string, storeCache table.Cache, numOfLevels int) StoreVersionSet {
	return &storeVersionSet{
		manifestFileNumber: *atomic.NewInt64(1), // default value for initialize store
		nextFileNumber:     *atomic.NewInt64(2), // default value
		storePath:          storePath,
		storeCache:         storeCache,
		numOfLevels:        numOfLevels,
		familyVersions:     make(map[string]FamilyVersion),
		familyIDs:          make(map[FamilyID]string),
	}
}

// getCache returns table cache for reading data
func (vs *storeVersionSet) getCache() table.Cache {
	return vs.storeCache
}

// numberOfLevels returns num. of levels
func (vs *storeVersionSet) numberOfLevels() int {
	return vs.numOfLevels
}

// Destroy closes version set, release resource, such as journal writer etc.
func (vs *storeVersionSet) Destroy() error {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// close manifest journal writer if it's exist
	if vs.manifest != nil {
		if err := vs.manifest.Close(); err != nil {
			return err
		}
	}
	return nil
}

// NextFileNumber generates next file number
func (vs *storeVersionSet) NextFileNumber() table.FileNumber {
	// need add lock, because CommitFamilyEditLog will reset nextFileNumber
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	nextNumber := vs.nextFileNumber.Inc()
	return table.FileNumber(nextNumber - 1)
}

// ManifestFileNumber returns the current manifest file number
func (vs *storeVersionSet) ManifestFileNumber() table.FileNumber {
	return table.FileNumber(vs.manifestFileNumber.Load())
}

// CommitFamilyEditLog persists edit logs to manifest file, then apply new version to family version
func (vs *storeVersionSet) CommitFamilyEditLog(family string, editLog EditLog) error {
	// get family version based on family name
	familyVersion := vs.GetFamilyVersion(family)
	if familyVersion == nil {
		return fmt.Errorf("cannot find family version for name: %s", family)
	}

	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// add next file number init edit log for each delta edit log
	editLog.Add(NewNextFileNumber(table.FileNumber(vs.nextFileNumber.Load())))
	// persist edit log
	if err := vs.persistEditLogs(vs.manifest, []EditLog{editLog}); err != nil {
		return err
	}
	// get current snapshot
	snapshot := familyVersion.GetSnapshot()
	defer snapshot.Close()

	newVersion := snapshot.GetCurrent().Clone()

	// apply delta edit to new version
	editLog.apply(newVersion)

	// Install the new version for family level version edit log
	familyVersion.appendVersion(newVersion)
	versionLogger.Info("log and apply new version edit",
		logger.String("path", vs.storePath),
		logger.String("family", family),
		logger.Any("log", editLog))
	return nil
}

// CreateFamilyVersion creates family version using family name,
// if family version exist, return exist one
func (vs *storeVersionSet) CreateFamilyVersion(family string, familyID FamilyID) FamilyVersion {
	var familyVersion = vs.GetFamilyVersion(family)
	if familyVersion != nil {
		versionLogger.Warn("family version exist, use it.",
			logger.String("path", vs.storePath),
			logger.String("family", family))
		return familyVersion
	}
	familyVersion = newFamilyVersion(familyID, family, vs)
	vs.mutex.Lock()
	vs.familyVersions[family] = familyVersion
	vs.familyIDs[familyID] = family
	vs.mutex.Unlock()
	return familyVersion
}

// GetFamilyVersion returns family version if it exists, else return nil
func (vs *storeVersionSet) GetFamilyVersion(family string) FamilyVersion {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()

	if familyVersion, ok := vs.familyVersions[family]; ok {
		return familyVersion
	}
	return nil
}

// Recover recovers version set if it exist, recover been invoked when kv store init.
// Initialize if version file not exists, else recover old data then init journal writer.
func (vs *storeVersionSet) Recover() error {
	if !fileutil.Exist(filepath.Join(vs.storePath, current())) {
		versionLogger.Info("version set's current file not exist, initialize it",
			logger.String("path", vs.storePath))
		if err := vs.initJournal(); err != nil {
			return err
		}
		return nil
	}
	versionLogger.Info("recover version set data from journal file", logger.String("path", vs.storePath))
	if err := vs.recover(); err != nil {
		return err
	}
	if err := vs.initJournal(); err != nil {
		return err
	}
	return nil
}

// recover does recover logic, read journal wal record and recover it
func (vs *storeVersionSet) recover() error {
	manifestFileName, err := vs.readManifestFileName()
	if err != nil {
		return err
	}
	manifestPath := vs.getManifestFilePath(manifestFileName)
	reader, err := newBufferReaderFunc(manifestPath)
	defer func() {
		if reader != nil {
			if e := reader.Close(); e != nil {
				versionLogger.Error("close manifest reader error",
					logger.String("path", vs.storePath),
					logger.String("manifest", manifestPath))
			}
		}
	}()
	if err != nil {
		return fmt.Errorf("create journal reader error:%s", err)
	}
	// read edit log
	for reader.Next() {
		record, err := reader.Read()
		if err != nil {
			return fmt.Errorf("recover data from manifest file error:%s", err)
		}
		editLog := newEmptyEditLogFunc()
		unmarshalErr := editLog.unmarshal(record)
		if unmarshalErr != nil {
			return fmt.Errorf("unmarshal edit log data from manifest file error:%s", unmarshalErr)
		}

		familyID := editLog.FamilyID()
		if familyID == StoreFamilyID {
			editLog.applyVersionSet(vs)
		} else if err := vs.applyFamilyVersion(familyID, editLog); err != nil {
			return err
		}
	}
	return nil
}

// applyFamilyVersion applies edit log to family version
func (vs *storeVersionSet) applyFamilyVersion(familyID FamilyID, editLog EditLog) error {
	// find related family version
	familyVersion := vs.getFamilyVersion(familyID)
	if familyVersion == nil {
		return fmt.Errorf("cannot get family version by id:%d", familyID)
	}
	snapshot := familyVersion.GetSnapshot()
	defer snapshot.Close()
	// apply edit log to family current family
	editLog.apply(snapshot.GetCurrent())
	return nil
}

// setNextFileNumberWithoutLock set next file number, invoker must add lock
func (vs *storeVersionSet) setNextFileNumberWithoutLock(newNextFileNumber table.FileNumber) {
	next := int64(newNextFileNumber)
	vs.manifestFileNumber.Store(next)
	vs.nextFileNumber.Store(next + 1)
}

// readManifestFileName reads manifest file name from current file
func (vs *storeVersionSet) readManifestFileName() (string, error) {
	current := vs.getCurrentPath()
	v, err := readFileFunc(current)
	if err != nil {
		return "", fmt.Errorf("write manifest file name error:%s", err)
	}
	return string(v), nil
}

// initJournal creates journal writer,
// 1. must write version set's data into journal,
// 2. set current manifest file name into current file.
// 3. set version set's manifest writer
func (vs *storeVersionSet) initJournal() error {
	if vs.manifest == nil {
		manifestFileName := ManifestFileName(table.FileNumber(vs.manifestFileNumber.Load())) // manifest file name
		manifestPath := vs.getManifestFilePath(manifestFileName)
		writer, err := newBufferWriterFunc(manifestPath)
		if err != nil {
			return err
		}
		// need snapshot writes snapshot first
		editLogs := vs.createSnapshot()
		if err := vs.persistEditLogs(writer, editLogs); err != nil {
			return err
		}
		// make sure write snapshot success, important!!!!!!!
		// then set manifest file name into current file
		if err := vs.setCurrent(manifestFileName); err != nil {
			return err
		}
		// finally set version set's manifest writer
		vs.manifest = writer
	}
	return nil
}

// getFamilyVersion returns family version
func (vs *storeVersionSet) getFamilyVersion(familyID FamilyID) FamilyVersion {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()

	if familyName, ok := vs.familyIDs[familyID]; ok {
		return vs.familyVersions[familyName]
	}
	return nil
}

// newVersionID generates new version id
func (vs *storeVersionSet) newVersionID() int64 {
	newID := vs.versionID.Add(1)
	return newID - 1
}

// setCurrent writes manifest file name into CURRENT file
func (vs *storeVersionSet) setCurrent(manifestFile string) error {
	current := vs.getCurrentPath()
	tmp := fmt.Sprintf("%s.%s", current, TmpSuffix)
	// write manifest file name into current file
	if err := writeFileFunc(tmp, []byte(manifestFile), 0666); err != nil {
		return fmt.Errorf("write manifest file name into current tmp file error:%s", err)
	}
	if err := renameFunc(tmp, current); err != nil {
		return fmt.Errorf("rename current tmp file name to current error:%s", err)
	}
	return nil
}

// getCurrent returns current file path
func (vs *storeVersionSet) getCurrentPath() string {
	return filepath.Join(vs.storePath, current())
}

// getManifestFilePath returns manifest file path
func (vs *storeVersionSet) getManifestFilePath(manifestFileName string) string {
	return filepath.Join(vs.storePath, manifestFileName)
}

// createSnapshot builds current version edit log
func (vs *storeVersionSet) createSnapshot() (editLogs []EditLog) {
	// for family level edit log
	for id, name := range vs.familyIDs {
		editLog := vs.createFamilySnapshot(id, vs.familyVersions[name])
		editLogs = append(editLogs, editLog)
	}

	// for store level edit log
	editLogs = append(editLogs, vs.createStoreSnapshot())
	return
}

// createFamilySnapshot creates snapshot of edit log for family level.
// NOTE: IMPORTANT!!!!!, need write edit logs for all data of version.
func (vs *storeVersionSet) createFamilySnapshot(familyID FamilyID, familyVersion FamilyVersion) EditLog {
	editLog := NewEditLog(familyID)
	// save current version all active files
	snapshot := familyVersion.GetSnapshot()
	defer snapshot.Close()
	current := snapshot.GetCurrent()

	// write log for current file list under this family.
	levels := current.Levels()
	for numOfLevel, level := range levels {
		files := level.getFiles()
		for _, file := range files {
			// level -> file meta
			newFile := CreateNewFile(int32(numOfLevel), file)
			editLog.Add(newFile)
		}
	}
	// write log if family has replica sequences.
	sequences := current.GetSequences()
	for leader, seq := range sequences {
		// leader -> replica sequence
		editLog.Add(CreateSequence(leader, seq))
	}

	// write log if family has reference files
	refFiles := current.GetAllReferenceFiles()
	for store, families := range refFiles {
		for familyID, files := range families {
			for _, file := range files {
				editLog.Add(CreateNewReferenceFile(store, familyID, file))
			}
		}
	}

	// write log if family has rollup files
	rollupFiles := current.GetRollupFiles()
	for file, intervals := range rollupFiles {
		for _, interval := range intervals {
			editLog.Add(CreateNewRollupFile(file, interval))
		}
	}

	return editLog
}

// createStoreSnapshot creates snapshot of edit log for store level
func (vs *storeVersionSet) createStoreSnapshot() EditLog {
	editLog := NewEditLog(StoreFamilyID)
	// save next file number
	editLog.Add(NewNextFileNumber(table.FileNumber(vs.nextFileNumber.Load())))
	return editLog
}

// persistEditLogs persists edit logs into manifest file
func (vs *storeVersionSet) persistEditLogs(writer bufioutil.BufioWriter, editLogs []EditLog) error {
	for _, editLog := range editLogs {
		v, err := editLog.marshal()
		if err != nil {
			return fmt.Errorf("encode edit log error:%s", err)
		}
		if _, err := writer.Write(v); err != nil {
			return fmt.Errorf("write edit log error:%s", err)
		}
		if err := writer.Sync(); err != nil {
			return fmt.Errorf("sync edit log error:%s", err)
		}
	}
	return nil
}

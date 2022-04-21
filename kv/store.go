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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
)

//go:generate mockgen -source ./store.go -destination=./store_mock.go -package kv

// for testing
var (
	encodeTomlFunc    = ltoml.EncodeToml
	decodeTomlFunc    = ltoml.DecodeToml
	newFamilyFunc     = newFamily
	newVersionSetFunc = version.NewStoreVersionSet
	listDirFunc       = fileutil.ListDir
	mkDirFunc         = fileutil.MkDir
	removeFunc        = os.Remove
	newFileLockFunc   = lockers.NewFileLock
	newStoreFunc      = newStore
)

// Store is kv store, supporting column family, but is different from other LSM implementation.
// Current implementation doesn't contain memory table write logic.
type Store interface {
	// Name returns the store's name.
	Name() string
	// Path returns the store root path.
	Path() string
	// CreateFamily create/load column family.
	CreateFamily(familyName string, option FamilyOption) (Family, error)
	// GetFamily gets family based on name, return nil if not exist.
	GetFamily(familyName string) Family
	// ListFamilyNames returns the all family's name
	ListFamilyNames() []string
	// Option returns the store configuration options
	Option() StoreOption
	// ForceRollup does rollup job manual.
	ForceRollup()

	// compact the families under store.
	compact()
	// close store, then release some resource
	close() error
	// createFamilyVersion creates family version using family name and family id,
	// if family version exist, return exist one
	createFamilyVersion(name string, familyID version.FamilyID) version.FamilyVersion
	// nextFileNumber generates next file number
	nextFileNumber() table.FileNumber
	// commitFamilyEditLog persists edit logs to manifest file, then apply new version to family version
	commitFamilyEditLog(name string, editLog version.EditLog) error
	// evictFamilyFile evicts family file reader from cache
	evictFamilyFile(fileNumber table.FileNumber)
}

// store implements Store interface
type store struct {
	name   string
	path   string
	option StoreOption
	// file-lock restricts access to store by allowing only one instance
	lock     lockers.FileLock
	versions version.StoreVersionSet
	// each family instance need to be assigned a unique family id
	familySeq int
	families  map[string]Family
	// RWMutex for accessing family
	rwMutex sync.RWMutex

	storeInfo *storeInfo
	cache     table.Cache

	ctx    context.Context
	cancel context.CancelFunc
}

// newStore news store instance, need recover data if store existent
func newStore(name, path string, option StoreOption) (s Store, err error) {
	var info *storeInfo
	var isCreate bool
	if fileutil.Exist(path) {
		// exist store, open it, load store info and config from INFO
		info = &storeInfo{}
		optionsFile := filepath.Join(path, version.Options)
		if err = decodeTomlFunc(optionsFile, info); err != nil {
			return nil, fmt.Errorf("load store info file:%s, error:%s", optionsFile, err)
		}
	} else {
		// create store, initialize path and store info
		if err = mkDirFunc(path); err != nil {
			return nil, fmt.Errorf("create store path error:%s", err)
		}
		info = newStoreInfo(option)
		isCreate = true
	}
	// try lock
	lock := newFileLockFunc(filepath.Join(path, version.Lock))
	err = lock.Lock()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	store1 := &store{
		name:      name,
		path:      path,
		option:    option,
		lock:      lock,
		families:  make(map[string]Family),
		storeInfo: info,
		ctx:       ctx,
		cancel:    cancel,
	}

	defer func() {
		if err != nil {
			// if init err, need close store for release resource
			if err2 := store1.close(); err2 != nil {
				kvLogger.Error("close store err when create store fail",
					logger.String("store", path), logger.Error(err))
			}
		}

		// finally, try delete obsolete files
		store1.deleteObsoleteFiles()
		store1.deleteFamilyObsoleteFiles()
	}()

	// build store reader cache
	store1.cache = table.NewCache(path, option.TTL.Duration())
	// init version set
	store1.versions = newVersionSetFunc(path, store1.cache, store1.option.Levels)

	if isCreate {
		// if store is new created, need dump store info to INFO file
		err = store1.dumpStoreInfo()
		if err != nil {
			return nil, err
		}
	} else {
		// existed store need loading all families instance
		for familyName, familyOption := range info.Families {
			if store1.familySeq < familyOption.ID {
				store1.familySeq = familyOption.ID
			}
			// TODO lazy load??????
			var family Family
			// open existed family
			family, err = newFamily(store1, familyOption)
			if err != nil {
				return nil, fmt.Errorf("building family instance for existed store[%s] error:%s", path, err)
			}
			store1.families[familyName] = family
		}
	}
	// recover version set, after recovering family options
	err = store1.versions.Recover()
	if err != nil {
		return nil, fmt.Errorf("recover store version set error:%s", err)
	}

	return store1, nil
}

// Name returns the store's name.
func (s *store) Name() string {
	return s.name
}

// Path returns the store root path.
func (s *store) Path() string {
	return s.path
}

// CreateFamily create/load column family.
func (s *store) CreateFamily(familyName string, option FamilyOption) (family Family, err error) {
	s.rwMutex.RLock()
	family, ok := s.families[familyName]
	s.rwMutex.RUnlock()
	if ok {
		// return exist family
		return family, nil
	}

	familyPath := filepath.Join(s.path, familyName)

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if !fileutil.Exist(familyPath) {
		// create new family
		option.Name = familyName
		// assign unique family id
		s.familySeq++
		option.ID = s.familySeq
		s.storeInfo.Families[familyName] = option
		if err = s.dumpStoreInfo(); err != nil {
			// if dump store info error remove family option from store info
			delete(s.storeInfo.Families, familyName)
			return nil, err
		}
	}

	family, err = newFamilyFunc(s, s.storeInfo.Families[familyName])
	if err != nil {
		return nil, err
	}
	s.families[familyName] = family
	return family, nil
}

// GetFamily gets family based on name, return nil if not exist.
func (s *store) GetFamily(familyName string) Family {
	s.rwMutex.RLock()
	family := s.families[familyName]
	s.rwMutex.RUnlock()
	return family
}

// ListFamilyNames returns the all family's name
func (s *store) ListFamilyNames() []string {
	var result []string
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	for name := range s.families {
		result = append(result, name)
	}
	return result
}

// Option returns the store configuration options
func (s *store) Option() StoreOption {
	return s.option
}

// ForceRollup does rollup job manual.
func (s *store) ForceRollup() {
	families := s.getCurrentFamilies()
	for _, f := range families {
		f.rollup()
	}
}

// close the store, then release some resource
func (s *store) close() error {
	// close each family in kv store.
	families := s.getCurrentFamilies()
	for _, f := range families {
		f.close()
	}

	if err := s.cache.Close(); err != nil {
		kvLogger.Error("close store cache error", logger.String("store", s.path), logger.Error(err))
	}
	if err := s.versions.Destroy(); err != nil {
		kvLogger.Error("destroy store version set error",
			logger.String("store", s.path), logger.Error(err))
	}
	s.cancel()
	return s.lock.Unlock()
}

// evictFamilyFile evicts family file reader from cache
func (s *store) evictFamilyFile(fileNumber table.FileNumber) {
	s.cache.Evict(version.Table(fileNumber))
}

// createFamilyVersion creates family version using family name and family id,
// if family version exist, return exist one
func (s *store) createFamilyVersion(name string, familyID version.FamilyID) version.FamilyVersion {
	return s.versions.CreateFamilyVersion(name, familyID)
}

// nextFileNumber generates next file number
func (s *store) nextFileNumber() table.FileNumber {
	return s.versions.NextFileNumber()
}

// commitFamilyEditLog persists edit logs to manifest file, then apply new version to family version
func (s *store) commitFamilyEditLog(name string, editLog version.EditLog) error {
	return s.versions.CommitFamilyEditLog(name, editLog)
}

// dumpStoreInfo persists store info to OPTIONS file
func (s *store) dumpStoreInfo() error {
	infoPath := filepath.Join(s.path, version.Options)
	// write store info using toml format
	if err := encodeTomlFunc(infoPath, s.storeInfo); err != nil {
		return fmt.Errorf("write store info to file[%s] error:%s", infoPath, err)
	}
	return nil
}

// compact checks if family need to do compact,
// if it needs, does compaction job.
func (s *store) compact() {
	s.rwMutex.RLock()
	families := make([]Family, len(s.families))
	i := 0
	for _, family := range s.families {
		families[i] = family
		i++
	}
	s.rwMutex.RUnlock()

	for _, family := range families {
		if family.needCompact() {
			family.compact()
		}
		if family.needRollup() {
			family.rollup()
		}
	}

	// try to evict expired reader from cache.
	s.cache.Cleanup()
}

// deleteFamilyObsoleteFiles deletes the all families obsolete files when init kv store
func (s *store) deleteFamilyObsoleteFiles() {
	for _, family := range s.families {
		family.deleteObsoleteFiles()
	}
}

// deleteObsoleteFiles deletes the obsolete files
func (s *store) deleteObsoleteFiles() {
	files, err := listDirFunc(s.path)
	if err != nil {
		kvLogger.Error("list files fail when delete obsolete files",
			logger.String("store", s.path), logger.String("kv", s.name))
		return
	}

	currentManifest := version.ManifestFileName(s.versions.ManifestFileNumber())
	for _, fileName := range files {
		if !strings.HasPrefix(fileName, version.ManifestPrefix) {
			continue
		}
		if fileName == currentManifest {
			continue
		}
		if err := removeFunc(filepath.Join(s.path, fileName)); err != nil {
			kvLogger.Error("delete obsolete manifest file fail",
				logger.String("store", s.path), logger.String("file", fileName))
		}
	}
}

// getCurrentFamilies returns current families in kv store.
func (s *store) getCurrentFamilies() (families []Family) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	for _, family := range s.families {
		families = append(families, family)
	}
	return families
}

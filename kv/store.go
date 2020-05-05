package kv

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/timeutil"
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
)

// Store is kv store, supporting column family, but is different from other LSM implementation.
// Current implementation doesn't contain memory table write logic.
type Store interface {
	// CreateFamily create/load column family.
	CreateFamily(familyName string, option FamilyOption) (Family, error)
	// GetFamily gets family based on name, return nil if not exist.
	GetFamily(familyName string) Family
	// ListFamilyNames returns the all family's name
	ListFamilyNames() []string
	// Option returns the store configuration options
	Option() StoreOption
	// RegisterRollup registers the rollup source/target relation
	RegisterRollup(interval timeutil.Interval, rollup Rollup)
	// Close closes store, then release some resource
	Close() error

	// createFamilyVersion creates family version using family name and family id,
	// if family version exist, return exist one
	createFamilyVersion(name string, familyID version.FamilyID) version.FamilyVersion
	// nextFileNumber generates next file number
	nextFileNumber() table.FileNumber
	// commitFamilyEditLog persists edit logs to manifest file, then apply new version to family version
	commitFamilyEditLog(name string, editLog version.EditLog) error
	// evictFamilyFile evicts family file reader from cache
	evictFamilyFile(name string, fileNumber table.FileNumber)
	// getRollup returns the rollup relation by interval
	getRollup(interval timeutil.Interval) (Rollup, bool)
}

// store implements Store interface
type store struct {
	name   string
	option StoreOption
	// file-lock restricts access to store by allowing only one instance
	lock     lockers.FileLock
	versions version.StoreVersionSet
	// each family instance need to be assigned an unique family id
	familySeq int
	families  map[string]Family
	// RWMutex for accessing family
	rwMutex sync.RWMutex

	storeInfo *storeInfo
	cache     table.Cache

	rollupRelations map[timeutil.Interval]Rollup // save target kv store for rollup job

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStore new store instance, need recover data if store existent
func NewStore(name string, option StoreOption) (s Store, err error) {
	var info *storeInfo
	var isCreate bool
	if fileutil.Exist(option.Path) {
		// exist store, open it, load store info and config from INFO
		info = &storeInfo{}
		if err := decodeTomlFunc(filepath.Join(option.Path, version.Options), info); err != nil {
			return nil, fmt.Errorf("load store info error:%s", err)
		}
	} else {
		// create store, initialize path and store info
		if err := mkDirFunc(option.Path); err != nil {
			return nil, fmt.Errorf("create store path error:%s", err)
		}
		info = newStoreInfo(option)
		isCreate = true
	}
	// try lock
	lock := newFileLockFunc(filepath.Join(option.Path, version.Lock))
	err = lock.Lock()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	store1 := &store{
		name:      name,
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
			if err2 := store1.Close(); err2 != nil {
				kvLogger.Error("close store err when create store fail",
					logger.String("store", option.Path), logger.Error(err))
			}
		}

		// finally need try delete obsolete files
		store1.deleteObsoleteFiles()
		store1.deleteFamilyObsoleteFiles()
	}()

	// build store reader cache
	store1.cache = table.NewCache(store1.option.Path)
	// init version set
	store1.versions = newVersionSetFunc(store1.option.Path, store1.cache, store1.option.Levels)

	if isCreate {
		// if store is new created, need dump store info to INFO file
		if err := store1.dumpStoreInfo(); err != nil {
			return nil, err
		}
	} else {
		// existed store need loading all families instance
		for familyName, familyOption := range info.Families {
			if store1.familySeq < familyOption.ID {
				store1.familySeq = familyOption.ID
			}
			// open existed family
			family, err := newFamily(store1, familyOption)
			if err != nil {
				return nil, fmt.Errorf("building family instance for existed store[%s] error:%s", option.Path, err)
			}
			store1.families[familyName] = family
		}
	}
	// recover version set, after recovering family options
	if err = store1.versions.Recover(); err != nil {
		return nil, fmt.Errorf("recover store version set error:%s", err)
	}

	// schedule compact job
	store1.scheduleCompactJob()
	return store1, nil
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

	familyPath := filepath.Join(s.option.Path, familyName)

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if !fileutil.Exist(familyPath) {
		// create new family
		option.Name = familyName
		// assign unique family id
		s.familySeq++
		option.ID = s.familySeq
		s.storeInfo.Families[familyName] = option
		if err := s.dumpStoreInfo(); err != nil {
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

// RegisterRollup registers the rollup source/target relation
func (s *store) RegisterRollup(interval timeutil.Interval, rollup Rollup) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if s.rollupRelations == nil {
		s.rollupRelations = make(map[timeutil.Interval]Rollup)
	}
	_, ok := s.rollupRelations[interval]
	if ok {
		kvLogger.Error("rollup interval already register, ignore.",
			logger.String("store", s.option.Path),
			logger.Any("interval", interval))
		return
	}
	s.rollupRelations[interval] = rollup
}

// Close closes store, then release some resource
func (s *store) Close() error {
	//FIXME stone1100 need if has background job doing(family compact/flush etc.)
	if err := s.cache.Close(); err != nil {
		kvLogger.Error("close store cache error", logger.String("store", s.option.Path), logger.Error(err))
	}
	if err := s.versions.Destroy(); err != nil {
		kvLogger.Error("destroy store version set error",
			logger.String("store", s.option.Path), logger.Error(err))
	}
	s.cancel()
	return s.lock.Unlock()
}

// evictFamilyFile evicts family file reader from cache
func (s *store) evictFamilyFile(name string, fileNumber table.FileNumber) {
	s.cache.Evict(name, version.Table(fileNumber))
}

// getRollup returns the rollup relation by interval
func (s *store) getRollup(interval timeutil.Interval) (Rollup, bool) {
	rollup, ok := s.rollupRelations[interval]
	return rollup, ok
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
	infoPath := filepath.Join(s.option.Path, version.Options)
	// write store info using toml format
	if err := encodeTomlFunc(infoPath, s.storeInfo); err != nil {
		return fmt.Errorf("write store info to file[%s] error:%s", infoPath, err)
	}
	return nil
}

// scheduleCompactJob schedules a compaction background job
func (s *store) scheduleCompactJob() {
	interval := defaultCompactCheckInterval
	if s.option.CompactCheckInterval > 0 {
		interval = s.option.CompactCheckInterval
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.compact()
			case <-s.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// compact checks if family need do compact, if need, does compaction job
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
		if family.needCompat() {
			family.compact()
		}
	}
}

// deleteFamilyObsoleteFiles deletes the all families obsolete files when init kv store
func (s *store) deleteFamilyObsoleteFiles() {
	for _, family := range s.families {
		family.deleteObsoleteFiles()
	}
}

// deleteObsoleteFiles deletes the obsolete files
func (s *store) deleteObsoleteFiles() {
	files, err := listDirFunc(s.option.Path)
	if err != nil {
		kvLogger.Error("list files fail when delete obsolete files",
			logger.String("store", s.option.Path), logger.String("kv", s.name))
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
		if err := removeFunc(filepath.Join(s.option.Path, fileName)); err != nil {
			kvLogger.Error("delete obsolete manifest file fail",
				logger.String("store", s.option.Path), logger.String("file", fileName))
		}
	}
}

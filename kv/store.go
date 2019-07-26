package kv

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/kv/table"
	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/fileutil"
	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/pkg/logger"
)

// Store is kv store, supporting column family, but is different from other LSM implementation.
// Current implementation doesn't contain memory table write logic.
type Store interface {
	// CreateFamily create/load column family.
	CreateFamily(familyName string, option FamilyOption) (Family, error)
	// GetFamily gets family based on name, return nil if not exist.
	GetFamily(familyName string) Family
	// Close closes store, then release some resource
	Close() error
}

// store implements Store interface
type store struct {
	name   string
	option StoreOption
	// file-lock restricts access to store by allowing only one instance
	lock     *lockers.FileLock
	versions *version.StoreVersionSet
	// each family instance need to be assigned an unique family id
	familyID int
	families map[string]Family
	// RWMutex for accessing family
	rwMutex sync.RWMutex

	storeInfo *storeInfo
	cache     table.Cache

	logger *logger.Logger
}

// NewStore new store instance, need recover data if store existent
func NewStore(name string, option StoreOption) (Store, error) {
	var info *storeInfo
	var isCreate bool
	if fileutil.Exist(option.Path) {
		// exist store, open it, load store info and config from INFO
		info = &storeInfo{}
		if err := fileutil.DecodeToml(filepath.Join(option.Path, version.Options), info); err != nil {
			return nil, fmt.Errorf("load store info error:%s", err)
		}
	} else {
		// create store, initialize path and store info
		if err := fileutil.MkDir(option.Path); err != nil {
			return nil, fmt.Errorf("create store path error:%s", err)
		}
		info = newStoreInfo(option)
		isCreate = true
	}
	// try lock
	lock := lockers.NewFileLock(filepath.Join(option.Path, version.Lock))
	err := lock.Lock()
	if err != nil {
		return nil, err
	}

	log := logger.GetLogger(fmt.Sprintf("kv/store[%s]", option.Path))

	// unlock file lock if error
	defer func() {
		if err != nil {
			if e := lock.Unlock(); e != nil {
				log.Error("unlock file error:", logger.Error(e))
			}
		}
	}()

	store := &store{
		name:      name,
		option:    option,
		lock:      lock,
		families:  make(map[string]Family),
		logger:    log,
		storeInfo: info,
	}

	// init version set
	store.versions = version.NewStoreVersionSet(store.option.Path, store.option.Levels)

	if isCreate {
		// if store is new created, need dump store info to INFO file
		if err := store.dumpStoreInfo(); err != nil {
			return nil, err
		}
	} else {
		// existed store need loading all families instance
		for familyName, familyOption := range info.Families {
			if store.familyID < familyOption.ID {
				store.familyID = familyOption.ID
			}
			// open existed family
			family, err := newFamily(store, familyOption)
			if err != nil {
				return nil, fmt.Errorf("building family instance for existed store[%s] error:%s", option.Path, err)
			}
			store.families[familyName] = family
		}
	}
	// recover version set, after recovering family options
	if err := store.versions.Recover(); err != nil {
		return nil, fmt.Errorf("recover store version set error:%s", err)
	}

	// build store reader cache
	store.cache = table.NewCache(store.option.Path)
	return store, nil
}

// CreateFamily create/load column family.
func (s *store) CreateFamily(familyName string, option FamilyOption) (Family, error) {
	s.rwMutex.RLock()
	family, ok := s.families[familyName]
	s.rwMutex.RUnlock()
	if ok {
		return family, nil
	}

	// todo: check the lock granularity
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	familyPath := filepath.Join(s.option.Path, familyName)
	var err error
	if !fileutil.Exist(familyPath) {
		// create new family
		option.Name = familyName
		// assign unique family id
		s.familyID++
		option.ID = s.familyID
		s.storeInfo.Families[familyName] = option
		if err := s.dumpStoreInfo(); err != nil {
			// if dump store info error remove family option from store info
			delete(s.storeInfo.Families, familyName)
			return nil, err
		}
	}

	family, err = newFamily(s, s.storeInfo.Families[familyName])
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

// Close closes store, then release some resource
func (s *store) Close() error {
	if err := s.cache.Close(); err != nil {
		s.logger.Error("close store cache error", logger.Error(err))
	}
	if err := s.versions.Destroy(); err != nil {
		s.logger.Error("destroy store version set error", logger.Error(err))
	}
	return s.lock.Unlock()
}

// dumpStoreInfo persists store info to OPTIONS file
func (s *store) dumpStoreInfo() error {
	infoPath := filepath.Join(s.option.Path, version.Options)
	// write store info using toml format
	if err := fileutil.EncodeToml(infoPath, s.storeInfo); err != nil {
		return fmt.Errorf("write store info to file[%s] error:%s", infoPath, err)
	}
	return nil
}

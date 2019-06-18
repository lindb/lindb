package kv

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/kv/table"
	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"

	"go.uber.org/zap"
)

// Store is kv store, support column family, but different other LSM implements.
// current implement not include memory table write logic.
type Store struct {
	name     string
	option   StoreOption
	lock     *Lock // file lock make sure store only been open once instance
	versions *version.StoreVersionSet
	familyID int // each family instance need assign an unique family id
	families map[string]*Family

	storeInfo *storeInfo
	cache     table.Cache

	mutex sync.RWMutex

	logger *zap.Logger
}

// NewStore new store instance, need recover data if store existent
func NewStore(name string, option StoreOption) (*Store, error) {
	var info *storeInfo
	var isCreate bool
	if util.Exist(option.Path) {
		// exist store, open it, load store info and config from INFO
		info = &storeInfo{}
		if err := util.DecodeToml(filepath.Join(option.Path, version.Options), info); err != nil {
			return nil, fmt.Errorf("load store info error:%s", err)
		}
	} else {
		// create store, initialize path and store info
		if err := util.MkDir(option.Path); err != nil {
			return nil, fmt.Errorf("create store path error:%s", err)
		}
		info = newStoreInfo(option)
		isCreate = true
	}

	// first need do file lock, only allow open by a instance
	lock := NewLock(filepath.Join(option.Path, version.Lock))
	err := lock.Lock()
	if err != nil {
		return nil, err
	}

	log := logger.GetLogger()

	// unlock file lock if error
	defer func() {
		if err != nil {
			if e := lock.Unlock(); e != nil {
				log.Error("unlock file error:", zap.String("store", option.Path), zap.Error(e))
			}
		}
	}()

	store := &Store{
		name:      name,
		option:    option,
		lock:      lock,
		families:  make(map[string]*Family),
		logger:    log,
		storeInfo: info,
	}

	// init  version set
	store.versions = version.NewStoreVersionSet(store.option.Path, store.option.Levels)

	if isCreate {
		// if store is new created, need dump store info to INFO file
		if err := store.dumpStoreInfo(); err != nil {
			return nil, err
		}
	} else {
		// exist store need load all families instance
		for familyName, familyOption := range info.Familyies {
			if store.familyID < familyOption.ID {
				store.familyID = familyOption.ID
			}
			// open exist family
			family, err := newFamily(store, familyOption)
			if err != nil {
				return nil, fmt.Errorf("build family instance for exsit store[%s] error:%s", option.Path, err)
			}
			store.families[familyName] = family
		}
	}
	// recover version set, after recvoer family options
	if err := store.versions.Recover(); err != nil {
		return nil, fmt.Errorf("recover store version set error:%s", err)
	}

	// build store reader cache
	store.cache = table.NewCache(store.option.Path)
	return store, nil
}

// CreateFamily create/load column family.
func (s *Store) CreateFamily(familyName string, option FamilyOption) (*Family, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var family, ok = s.families[familyName]
	if !ok {
		familyPath := filepath.Join(s.option.Path, familyName)

		var err error
		if !util.Exist(familyPath) {
			// create new family
			option.Name = familyName
			// assign unqiue family id
			s.familyID++
			option.ID = s.familyID
			s.storeInfo.Familyies[familyName] = option
			if err := s.dumpStoreInfo(); err != nil {
				// if dump store info error remove family option from store info
				delete(s.storeInfo.Familyies, familyName)
				return nil, err
			}
		}

		family, err = newFamily(s, s.storeInfo.Familyies[familyName])

		if err != nil {
			return nil, err
		}
		s.families[familyName] = family
	}

	return family, nil
}

// GetFamily gets family based on name, if not exist return nil
func (s *Store) GetFamily(familyName string) (*Family, bool) {
	s.mutex.Lock()
	family, ok := s.families[familyName]
	s.mutex.Unlock()
	return family, ok
}

// Close closes store, then release some resoure
func (s *Store) Close() error {
	if err := s.cache.Close(); err != nil {
		s.logger.Error("close store cache error", zap.String("store", s.option.Path), zap.Error(err))
	}
	if err := s.versions.Destroy(); err != nil {
		s.logger.Error("destroy store version set error", zap.String("store", s.option.Path), zap.Error(err))
	}
	return s.lock.Unlock()
}

// dumpStoreInfo peresist store info to INFO file
func (s *Store) dumpStoreInfo() error {
	infoPath := filepath.Join(s.option.Path, version.Options)
	tmp := fmt.Sprintf("%s.%s", infoPath, version.TmpSuffix)
	// write store info using toml format
	if err := util.EncodeToml(tmp, s.storeInfo); err != nil {
		return fmt.Errorf("write store info error:%s", err)
	}
	if err := os.Rename(tmp, infoPath); err != nil {
		return fmt.Errorf("rename store info tmp file name error:%s", err)
	}
	return nil
}

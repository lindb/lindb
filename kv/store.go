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
)

//go:generate mockgen -source ./store.go -destination=./store_mock.go -package kv

var defaultCompactCheckInterval = 60

var mergers = make(map[string]Merger)

// RegisterMerger registers family merger
// NOTICE: must register before create family
func RegisterMerger(name string, merger Merger) {
	_, ok := mergers[name]
	if ok {
		panic("merger already register")
	}
	mergers[name] = merger
}

// Store is kv store, supporting column family, but is different from other LSM implementation.
// Current implementation doesn't contain memory table write logic.
type Store interface {
	// CreateFamily create/load column family.
	CreateFamily(familyName string, option FamilyOption) (Family, error)
	// GetFamily gets family based on name, return nil if not exist.
	GetFamily(familyName string) Family
	// ListFamilyNames returns the all family's name
	ListFamilyNames() []string
	// Close closes store, then release some resource
	Close() error
}

// store implements Store interface
type store struct {
	name   string
	option StoreOption
	// file-lock restricts access to store by allowing only one instance
	lock     *lockers.FileLock
	versions version.StoreVersionSet
	// each family instance need to be assigned an unique family id
	familyID int
	families map[string]Family
	// RWMutex for accessing family
	rwMutex sync.RWMutex

	storeInfo *storeInfo
	cache     table.Cache

	ctx    context.Context
	cancel context.CancelFunc

	logger *logger.Logger
}

// NewStore new store instance, need recover data if store existent
func NewStore(name string, option StoreOption) (Store, error) {
	var info *storeInfo
	var isCreate bool
	if fileutil.Exist(option.Path) {
		// exist store, open it, load store info and config from INFO
		info = &storeInfo{}
		if err := ltoml.DecodeToml(filepath.Join(option.Path, version.Options), info); err != nil {
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

	log := logger.GetLogger("kv", fmt.Sprintf("Store[%s]", option.Path))

	ctx, cancel := context.WithCancel(context.Background())
	// unlock file lock if error
	defer func() {
		if err != nil {
			if e := lock.Unlock(); e != nil {
				log.Error("unlock file error:", logger.Error(e))
			}
			cancel()
		}
	}()

	store := &store{
		name:      name,
		option:    option,
		lock:      lock,
		families:  make(map[string]Family),
		logger:    log,
		storeInfo: info,
		ctx:       ctx,
		cancel:    cancel,
	}
	// finally need try delete obsolete files
	defer store.deleteObsoleteFiles()

	// build store reader cache
	store.cache = table.NewCache(store.option.Path)
	// init version set
	store.versions = version.NewStoreVersionSet(store.option.Path, store.cache, store.option.Levels)

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

	// schedule compact job
	store.scheduleCompactJob()
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

// Close closes store, then release some resource
func (s *store) Close() error {
	//FIXME stone1100 need if has background job doing(family compact/flush etc.)
	if err := s.cache.Close(); err != nil {
		s.logger.Error("close store cache error", logger.Error(err))
	}
	if err := s.versions.Destroy(); err != nil {
		s.logger.Error("destroy store version set error", logger.Error(err))
	}
	s.cancel()
	return s.lock.Unlock()
}

// dumpStoreInfo persists store info to OPTIONS file
func (s *store) dumpStoreInfo() error {
	infoPath := filepath.Join(s.option.Path, version.Options)
	// write store info using toml format
	if err := ltoml.EncodeToml(infoPath, s.storeInfo); err != nil {
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

// deleteObsoleteFiles deletes the obsolete files
func (s *store) deleteObsoleteFiles() {
	files, err := fileutil.ListDir(s.option.Path)
	if err != nil {
		s.logger.Error("list files fail when delete obsolete files", logger.String("kv", s.name))
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
		if err := os.Remove(filepath.Join(s.option.Path, fileName)); err != nil {
			s.logger.Error("delete obsolete manifest file fail",
				logger.String("kv", s.name), logger.String("file", fileName))
		}
	}
}

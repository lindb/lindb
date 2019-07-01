package kv

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/eleme/lindb/kv/table"
	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"
)

// Family implements column family for data isolation each family.
type Family interface {
	// Name return family's name
	Name() string
	// NewFlusher creates flusher for saving data to family.
	NewFlusher() Flusher
	// GetSnapshot returns current version for given key, includes sst files
	GetSnapshot(key uint32) (Snapshot, error)
}

// family implements Family interface
type family struct {
	store         *store
	name          string
	familyPath    string
	option        FamilyOption
	familyVersion *version.FamilyVersion
	logger        *zap.Logger
}

// newFamily creates new family or open existed family.
func newFamily(store *store, option FamilyOption) (Family, error) {
	log := logger.GetLogger()
	name := option.Name

	familyPath := filepath.Join(store.option.Path, name)

	if !util.Exist(familyPath) {
		if err := util.MkDir(familyPath); err != nil {
			return nil, fmt.Errorf("mkdir family path error:%s", err)
		}
	}

	f := &family{
		familyPath:    familyPath,
		store:         store,
		name:          name,
		option:        option,
		familyVersion: store.versions.CreateFamilyVersion(name, option.ID),
		logger:        log,
	}

	log.Info("new family success", f.logStoreField(), f.logFamilyField())
	return f, nil
}

// Name return family's name
func (f *family) Name() string {
	return f.name
}

// NewFlusher creates flusher for saving data to family.
func (f *family) NewFlusher() Flusher {
	return newStoreFlusher(f)
}

// GetSnapshot returns current version for given key, includes sst files
func (f *family) GetSnapshot(key uint32) (Snapshot, error) {
	v, files := f.familyVersion.FindFiles(key)
	var readers []table.Reader
	for _, fileMeta := range files {
		// get store reader from cache
		reader, err := f.store.cache.GetReader(f.name, fileMeta.GetFileNumber())
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return newSnapshot(v, readers), nil
}

// newTableBuilder creates table builder instance for storing kv data.
func (f *family) newTableBuilder() (table.Builder, error) {
	fileNumber := f.store.versions.NextFileNumber()
	return table.NewStoreBuilder(f.familyPath, fileNumber)
}

// commitEditLog persists edit logs into manifest file.
// returns true on committing successfully and false on failure
func (f *family) commitEditLog(editLog *version.EditLog) bool {
	if editLog.IsEmpty() {
		f.logger.Warn("edit log is empty", f.logStoreField(), f.logFamilyField())
		return true
	}
	if err := f.store.versions.CommitFamilyEditLog(f.name, editLog); err != nil {
		f.logger.Error("commit edit log error:", f.logStoreField(), f.logFamilyField(), zap.Error(err))
		return false
	}
	return true
}

// logStoreField logs store info。
func (f *family) logStoreField() zap.Field {
	return zap.String("store", f.store.option.Path)
}

// logFamilyField logs family info.
func (f *family) logFamilyField() zap.Field {
	return zap.String("family", f.name)
}

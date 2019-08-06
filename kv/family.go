package kv

import (
	"fmt"
	"path/filepath"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
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
	logger        *logger.Logger
}

// newFamily creates new family or open existed family.
func newFamily(store *store, option FamilyOption) (Family, error) {
	name := option.Name

	familyPath := filepath.Join(store.option.Path, name)
	log := logger.GetLogger(fmt.Sprintf("kv/famliy[%s]", familyPath))

	if !fileutil.Exist(familyPath) {
		if err := fileutil.MkDir(familyPath); err != nil {
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

	log.Info("new family success")
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
		f.logger.Warn("edit log is empty")
		return true
	}
	if err := f.store.versions.CommitFamilyEditLog(f.name, editLog); err != nil {
		f.logger.Error("commit edit log error:", logger.Error(err))
		return false
	}
	return true
}

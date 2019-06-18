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

// Family implement column family for data isolation each family
type Family struct {
	store         *Store
	name          string
	familyPath    string
	option        FamilyOption
	familyVersion *version.FamilyVersion
	logger        *zap.Logger
}

// newFamily creates new family or open exsit family
func newFamily(store *Store, option FamilyOption) (*Family, error) {
	log := logger.GetLogger()
	name := option.Name

	familyPath := filepath.Join(store.option.Path, name)

	if !util.Exist(familyPath) {
		if err := util.MkDir(familyPath); err != nil {
			return nil, fmt.Errorf("mk family path error:%s", err)
		}
	}

	f := &Family{
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

// NewFlusher creates flusher for saving data to family
func (f *Family) NewFlusher() Flusher {
	return newStoreFlusher(f)
}

// GetSnapshot returns current version for given key, includes sst files
func (f *Family) GetSnapshot(key uint32) (*Snapshot, error) {
	version, files := f.familyVersion.FindFiles(key)
	var readers []table.Reader
	for _, fileMeta := range files {
		// get store reader from cache
		reader, err := f.store.cache.GetReader(f.name, fileMeta.GetFileNumber())
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return newSnapshot(version, readers), nil
}

// newTableBuilder create table builder instance for storing kv data
func (f *Family) newTableBuilder() (table.Builder, error) {
	fileNumber := f.store.versions.NextFileNumber()
	return table.NewStoreBuilder(f.familyPath, fileNumber)
}

// commitEditLog peresists eidt logs into mamanifest file
// returns ture commit successfully, else failure
func (f *Family) commitEditLog(editLog *version.EditLog) bool {
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

// logStoreField logging store info
func (f *Family) logStoreField() zap.Field {
	return zap.String("store", f.store.option.Path)
}

// logFamilyField logging family info
func (f *Family) logFamilyField() zap.Field {
	return zap.String("family", f.name)
}

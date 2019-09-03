package kv

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

const dummy = ""
const defaultMaxFileSize = int32(256 * 1024 * 1024)

//go:generate mockgen -source ./family.go -destination=./family_mock.go -package kv

// Family implements column family for data isolation each family.
type Family interface {
	// ID return family's id
	ID() int
	// Name return family's name
	Name() string
	// NewFlusher creates flusher for saving data to family.
	NewFlusher() Flusher
	// GetSnapshot returns current version's snapshot
	GetSnapshot() version.Snapshot

	// getFamilyVersion returns the family version
	getFamilyVersion() version.FamilyVersion
	// commitEditLog persists edit logs into manifest file.
	commitEditLog(editLog *version.EditLog) bool
	// newTableBuilder creates table builder instance for storing kv data.
	newTableBuilder() (table.Builder, error)
	// needCompat returns level0 files if need do compact job
	needCompat() bool
	// compact does compaction job
	compact()
	// getMerger returns user implement merger
	getMerger() Merger
	// addPendingOutput add a file which current writing file number
	addPendingOutput(fileNumber int64)
	// removePendingOutput removes pending output file after compact or flush
	removePendingOutput(fileNumber int64)
	// deleteSST deletes the temp sst file, if flush or compact fail
	deleteSST(fileNumber int64) error
	// getLogger returns the logger under family
	getLogger() *logger.Logger
}

// family implements Family interface
type family struct {
	store         *store
	name          string
	familyPath    string
	option        FamilyOption
	merger        Merger
	familyVersion version.FamilyVersion
	maxFileSize   int32

	pendingOutputs sync.Map

	compacting int32

	logger *logger.Logger
}

// newFamily creates new family or open existed family.
func newFamily(store *store, option FamilyOption) (Family, error) {
	name := option.Name

	familyPath := filepath.Join(store.option.Path, name)
	log := logger.GetLogger("kv", fmt.Sprintf("Family[%s]", familyPath))

	if !fileutil.Exist(familyPath) {
		if err := fileutil.MkDir(familyPath); err != nil {
			return nil, fmt.Errorf("mkdir family path error:%s", err)
		}
	}
	merger, ok := mergers[option.Merger]
	if !ok {
		return nil, fmt.Errorf("merger of option not impelement Merger interface, merger is [%s]", option.Merger)
	}
	maxFileSize := defaultMaxFileSize
	if option.MaxFileSize > 0 {
		maxFileSize = option.MaxFileSize
	}

	f := &family{
		familyPath:    familyPath,
		store:         store,
		name:          name,
		option:        option,
		compacting:    0,
		merger:        merger,
		maxFileSize:   maxFileSize,
		familyVersion: store.versions.CreateFamilyVersion(name, option.ID),
		logger:        log,
	}

	log.Info("new family success")
	return f, nil
}

func (f *family) ID() int {
	return f.option.ID
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

// newTableBuilder creates table builder instance for storing kv data.
func (f *family) newTableBuilder() (table.Builder, error) {
	fileNumber := f.store.versions.NextFileNumber()
	fileName := filepath.Join(f.familyPath, version.Table(fileNumber))
	return table.NewStoreBuilder(fileNumber, fileName)
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
	//FIXME deleteObsoleteFiles
	return true
}

// needCompat returns level0 files if need do compact job
func (f *family) needCompat() bool {
	// has compaction job doing
	if atomic.LoadInt32(&f.compacting) == 1 {
		return false
	}

	snapshot := f.GetSnapshot()
	defer snapshot.Close()

	numberOfFiles := snapshot.GetCurrent().NumberOfFilesInLevel(0)
	if numberOfFiles > 0 && numberOfFiles >= f.option.CompactThreshold {
		f.logger.Info("need to compact level0 files",
			logger.Any("numOfFiles", numberOfFiles), logger.Any("threshold", f.option.CompactThreshold))
		return true
	}
	return false
}

// compact does compact job if hasn't compact job running
func (f *family) compact() {
	if atomic.CompareAndSwapInt32(&f.compacting, 0, 1) {
		go func() {
			defer atomic.StoreInt32(&f.compacting, 0)

			if err := f.backgroundCompactionJob(); err != nil {
				f.logger.Error("do compact job error", logger.String("family", f.name))
			}
		}()
	}
}

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
	compactJob := newCompactJob(f, compactionState)
	if err := compactJob.run(); err != nil {
		return err
	}
	return nil
}

// addPendingOutput add a file which current writing file number
func (f *family) addPendingOutput(fileNumber int64) {
	f.pendingOutputs.Store(fileNumber, dummy)
}

// removePendingOutput removes pending output file after compact or flush
func (f *family) removePendingOutput(fileNumber int64) {
	f.pendingOutputs.Delete(fileNumber)
}

// deleteSST deletes the temp sst file, if flush or compact fail
func (f *family) deleteSST(fileNumber int64) error {
	if err := fileutil.RemoveDir(filepath.Join(f.familyPath, version.Table(fileNumber))); err != nil {
		return err
	}
	return nil
}

// getLogger returns the logger under family
func (f *family) getLogger() *logger.Logger {
	return f.logger
}

// getFamilyVersion returns the family version
func (f *family) getFamilyVersion() version.FamilyVersion {
	return f.familyVersion
}

// getMerger returns user implement merger
func (f *family) getMerger() Merger {
	return f.merger
}

// deleteObsoleteFiles deletes obsolete file
func (f *family) deleteObsoleteFiles() {
	sstFiles, err := fileutil.ListDir(f.familyPath)
	if err != nil {
		f.logger.Error("list sst file fail when delete obsolete files", logger.String("family", f.name))
		return
	}
	// make a map for all live files
	liveFiles := make(map[int64]string)
	f.pendingOutputs.Range(func(key, value interface{}) bool {
		k, ok := key.(int64)
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
	//TODO add rollup file ref??

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
			f.store.cache.Evict(f.name, version.Table(fileNumber))
			if err := f.deleteSST(fileNumber); err != nil {
				f.logger.Error("delete sst file fail",
					logger.String("family", f.name), logger.Int64("fileNumber", fileNumber))
			}
			f.logger.Info("delete sst file successfully",
				logger.String("family", f.name), logger.Int64("fileNumber", fileNumber))
		}
	}
}

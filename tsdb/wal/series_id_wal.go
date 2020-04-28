package wal

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source=./series_id_wal.go -destination=./series_id_wal_mock.go -package=wal

// for testing
var (
	mkDirFunc          = fileutil.MkDirIfNotExist
	newPageFactoryFunc = page.NewFactory
)

// RecoveryFunc represents the
type RecoveryFunc = func(metricID uint32, tagsHash uint64, seriesID uint32) error
type CommitFunc = func() error

var seriesWALLogger = logger.GetLogger("wal", "series")

const (
	seriesEntryLength = 4 + 8 + 4                      // metric id + tags hash + series id
	seriesPageSize    = seriesEntryLength * 512 * 1024 // series wal page size
	metricIDOffset    = 0                              // metric id offset
	tagsHashOffset    = metricIDOffset + 4             // tags hash offset
	seriesIDOffset    = tagsHashOffset + 8             // series id offset
)

// SeriesWAL represents write ahead log which stores series data for index database
type SeriesWAL interface {
	// Append appends metricID/tagsHash/seriesID into wal log
	Append(metricID uint32, tagsHash uint64, seriesID uint32) error
	// NeedRecovery checks if wal log need to recover
	NeedRecovery() bool
	// Recovery recoveries wal log, then writes data via recovery function
	Recovery(recovery RecoveryFunc, commit CommitFunc)
	// Sync flushes data into disk
	Sync() error
	// Close closes the wal log
	Close() error
}

// seriesWAL implements SeriesWAL interface
type seriesWAL struct {
	path        string
	pageSize    int
	walFactory  page.Factory
	currentPage page.MappedPage
	offset      int

	pageIndex       atomic.Int64
	commitPageIndex atomic.Int64
}

// NewSeriesWAL creates a new series write ahead log
func NewSeriesWAL(path string) (SeriesWAL, error) {
	var err error
	if err = mkDirFunc(path); err != nil {
		return nil, err
	}
	// init wal page factory
	fct, err := newPageFactoryFunc(path, seriesPageSize)
	if err != nil {
		return nil, err
	}
	pageIDs := fct.GetPageIDs()
	wal := &seriesWAL{path: path, walFactory: fct, pageSize: seriesPageSize}
	defer func() {
		if err != nil {
			if err1 := wal.walFactory.Close(); err1 != nil {
				seriesWALLogger.Error("close wal page factory error when init series wal",
					logger.String("wal", wal.path), logger.Error(err))
			}
		}
	}()
	if len(pageIDs) > 0 {
		wal.commitPageIndex.Store(pageIDs[0] - 1)
		wal.pageIndex.Store(pageIDs[len(pageIDs)-1])
	}
	// acquire new page for appending series data
	if wal.currentPage, err = wal.walFactory.AcquirePage(wal.pageIndex.Load() + 1); err != nil {
		return nil, err
	}
	wal.pageIndex.Inc()

	return wal, nil
}

// Append appends metricID/tagsHash/seriesID into wal log
func (wal *seriesWAL) Append(metricID uint32, tagsHash uint64, seriesID uint32) (err error) {
	// prepare the data pointer
	if wal.offset+seriesEntryLength > wal.pageSize {
		// sync previous data page
		if err := wal.currentPage.Sync(); err != nil {
			seriesWALLogger.Error("sync data page err when alloc",
				logger.String("wal", wal.path), logger.Error(err))
		}
		// not enough space in current data page, need create new page
		walPage, err := wal.walFactory.AcquirePage(wal.pageIndex.Load() + 1)
		if err != nil {
			return err
		}
		wal.currentPage = walPage
		wal.pageIndex.Inc()
		wal.offset = 0 // need reset message offset for new page append
	}

	wal.currentPage.PutUint32(metricID, wal.offset+metricIDOffset)
	wal.currentPage.PutUint64(tagsHash, wal.offset+tagsHashOffset)
	wal.currentPage.PutUint32(seriesID, wal.offset+seriesIDOffset)
	wal.offset += seriesEntryLength

	return nil
}

// NeedRecovery checks if wal log need to recover
func (wal *seriesWAL) NeedRecovery() bool {
	return wal.pageIndex.Load()-wal.commitPageIndex.Load() > 1
}

// Recovery recoveries wal log, then writes data via recovery function
func (wal *seriesWAL) Recovery(recovery RecoveryFunc, commit CommitFunc) {
	current := wal.pageIndex.Load()
	committed := wal.commitPageIndex.Load()
	for i := committed; i < current; i++ {
		walPage, ok := wal.walFactory.GetPage(i)
		if !ok {
			continue
		}
		offset := 0
		for offset < seriesPageSize {
			metricID := walPage.ReadUint32(offset + metricIDOffset)
			if metricID == 0 {
				break
			}

			if err := recovery(metricID,
				walPage.ReadUint64(offset+tagsHashOffset),
				walPage.ReadUint32(offset+seriesIDOffset)); err != nil {
				//TODO add metric?????
				seriesWALLogger.Error("invoke recovery func error",
					logger.String("wal", wal.path), logger.Error(err))
				return
			}
			offset += seriesEntryLength
		}

		if err := commit(); err != nil {
			//TODO add metric?????
			seriesWALLogger.Error("invoke commit func error",
				logger.String("wal", wal.path), logger.Error(err))
			return
		}

		if err := wal.walFactory.ReleasePage(i); err != nil {
			//TODO add metric?????
			seriesWALLogger.Error("release series wal page error",
				logger.String("wal", wal.path), logger.Error(err))
		}

		wal.commitPageIndex.Inc()
	}
}

// Sync flushes data into disk
func (wal *seriesWAL) Sync() error {
	return wal.currentPage.Sync()
}

// Close closes the wal log
func (wal *seriesWAL) Close() error {
	return wal.currentPage.Close()
}

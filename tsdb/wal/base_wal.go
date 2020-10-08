package wal

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
	"github.com/lindb/lindb/series/field"
)

var baseWALLogger = logger.GetLogger("wal", "base")

var (
	recoveryCommitFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "wal_recovery_commit_fail",
			Help: "Recovery commit fail when wal recovery.",
		},
	)
	releaseWALPageFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "wal_release_page_fail",
			Help: "Release wal page field fail when wal recovery.",
		},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(recoveryCommitFailCounter)
	monitoring.StorageRegistry.MustRegister(releaseWALPageFailCounter)
}

// SeriesRecoveryFunc represents the series recovery function
type SeriesRecoveryFunc = func(metricID uint32, tagsHash uint64, seriesID uint32) error

// MetricRecoveryFunc represents the metric recovery function
type MetricRecoveryFunc = func(namespace, metricName string, metricID uint32) error

// FieldRecoveryFunc represents the field recovery function
type FieldRecoveryFunc = func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error

// TagKeyRecoveryFunc represents the tag key recovery function
type TagKeyRecoveryFunc = func(metricID uint32, tagKeyID uint32, tagKey string) error

// CommitFunc represents the commit function after recovery
type CommitFunc = func() error

// baseWAL represents base write ahead log
type baseWAL struct {
	path     string
	pageSize int

	walFactory  page.Factory
	currentPage page.MappedPage

	offset int

	pageIndex       atomic.Int64
	commitPageIndex atomic.Int64
}

// newBaseWAL creates a new base write ahead log
func newBaseWAL(path string, pageSize int) (*baseWAL, error) {
	var err error
	if err = mkDirFunc(path); err != nil {
		return nil, err
	}

	// init wal page factory
	fct, err := newPageFactoryFunc(path, pageSize)
	if err != nil {
		return nil, err
	}

	pageIDs := fct.GetPageIDs()
	wal := &baseWAL{path: path, walFactory: fct, pageSize: pageSize}

	defer func() {
		if err != nil {
			if err1 := wal.walFactory.Close(); err1 != nil {
				baseWALLogger.Error("close wal page factory error when init base wal",
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

func (wal *baseWAL) checkPage(length int) error {
	// prepare the data pointer
	if wal.offset+length > wal.pageSize {
		// sync previous data page
		if err := wal.currentPage.Sync(); err != nil {
			baseWALLogger.Error("sync data page err when alloc",
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
	return nil
}

func (wal *baseWAL) putUint8(value uint8) {
	wal.currentPage.PutUint8(value, wal.offset)
	wal.offset++
}

func (wal *baseWAL) putUint32(value uint32) {
	wal.currentPage.PutUint32(value, wal.offset)
	wal.offset += 4
}

func (wal *baseWAL) putUint64(value uint64) {
	wal.currentPage.PutUint64(value, wal.offset)
	wal.offset += 8
}

func (wal *baseWAL) putString(value string) {
	length := len(value)
	wal.putUint8(uint8(length))
	wal.currentPage.WriteBytes([]byte(value), wal.offset)
	wal.offset += length
}

// sync flushes data into disk
func (wal *baseWAL) sync() error {
	return wal.currentPage.Sync()
}

// close closes the wal log
func (wal *baseWAL) close() error {
	return wal.currentPage.Close()
}

// needRecovery checks if wal log need to recover
func (wal *baseWAL) needRecovery() bool {
	return wal.pageIndex.Load()-wal.commitPageIndex.Load() > 1
}

func readString(dataPage page.MappedPage, offset int) (val string, n int) {
	length := int(dataPage.ReadUint8(offset))
	data := dataPage.ReadBytes(offset+1, length)
	return string(data), length + 1
}

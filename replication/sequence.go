package replication

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source=./sequence.go -destination=./sequence_mock.go -package=replication

// for testing
var (
	newPageFactoryFunc = page.NewFactory
)

var sequenceLogger = logger.GetLogger("replication", "sequence")

const (
	//sequenceMetaSize 8 bytes for int64
	sequenceMetaSize = 8
	metaPageID       = 0
)

// Sequence represents a persistence sequence recorder
// for on storage side when transferring data from broker to storage.
type Sequence interface {
	// GetHeadSeq returns the head sequence which is the latest sequence of replica received.
	GetHeadSeq() int64
	// SetHeadSeq sets the head sequence which is the latest sequence of replica received.
	SetHeadSeq(seq int64)
	// GetAckSeq returns the ack sequence which is the latest sequence of replica successfully flushed to disk.
	GetAckSeq() int64
	// GetAckSeq sets the ack sequence which is the latest sequence of replica successfully flushed to disk.
	SetAckSeq(seq int64)
	// Sync syncs the Sequence to storage.
	Sync() error
	// Synced checked if the Sequence has been synced.
	Synced() bool
	// ResetSynced resets Synced() to false.
	ResetSynced()

	//TODO need add close method??
}

// sequence implements Sequence.
type sequence struct {
	dirPath     string
	metaPageFct page.Factory
	// meta stores the ackSeq to page cache.
	metaPage page.MappedPage
	// headSeq represents the the max sequence num of replica received.
	headSeq atomic.Int64
	// ackSeq represents the the max sequence num of replica flushed to disk.
	ackSeq atomic.Int64
	// false -> not synced, true -> synced
	synced atomic.Bool
}

// NewSequence returns a sequence with page cache corresponding to dirPath.
func NewSequence(dirPath string) (Sequence, error) {
	var err error
	metaPageFct, err := newPageFactoryFunc(dirPath, sequenceMetaSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if err1 := metaPageFct.Close(); err1 != nil {
				sequenceLogger.Error("close meta page factory err",
					logger.String("path", dirPath), logger.Error(err1))
			}
		}
	}()

	metaPage, err := metaPageFct.AcquirePage(metaPageID)
	if err != nil {
		return nil, err
	}
	ackSeq := int64(metaPage.ReadUint64(0))

	return &sequence{
		dirPath:     dirPath,
		metaPageFct: metaPageFct,
		metaPage:    metaPage,
		headSeq:     *atomic.NewInt64(ackSeq),
		ackSeq:      *atomic.NewInt64(ackSeq),
	}, nil
}

// ResetSynced resets Synced() to false.
func (s *sequence) ResetSynced() {
	s.synced.Store(false)
}

// Synced checked if the Sequence has been synced.
func (s *sequence) Synced() bool {
	return s.synced.Load()
}

// GetHeadSeq returns the head sequence which is the latest sequence of replica received.
func (s *sequence) GetHeadSeq() int64 {
	return s.headSeq.Load()
}

// SetHeadSeq sets the head sequence which is the latest sequence of replica received.
func (s *sequence) SetHeadSeq(seq int64) {
	s.headSeq.Store(seq)
}

// GetAckSeq returns the ack sequence which is the latest sequence of replica successfully flushed to disk.
func (s *sequence) GetAckSeq() int64 {
	return s.ackSeq.Load()
}

// GetAckSeq sets the ack sequence which is the latest sequence of replica successfully flushed to disk.
func (s *sequence) SetAckSeq(seq int64) {
	s.ackSeq.Store(seq)
}

// Sync syncs the Sequence to storage.
func (s *sequence) Sync() error {
	s.metaPage.PutUint64(uint64(s.GetAckSeq()), 0)
	s.synced.Store(true)
	return s.metaPage.Sync()
}

package replication

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/queue"
)

//go:generate mockgen -source=./sequence.go -destination=./sequence_mock.go -package=replication

const (
	//sequenceMetaSize 8 bytes for int64
	sequenceMetaSize = 8
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
	dirPath string
	// meta stores the ackSeq to page cache.
	meta queue.Meta
	// headSeq represents the the max sequence num of replica received.
	headSeq atomic.Int64
	// ackSeq represents the the max sequence num of replica flushed to disk.
	ackSeq atomic.Int64
	// 0 -> not synced, 1 -> synced
	synced atomic.Int32
}

// ResetSynced resets Synced() to false.
func (s *sequence) ResetSynced() {
	s.synced.Store(0)
}

// Synced checked if the Sequence has been synced.
func (s *sequence) Synced() bool {
	return s.synced.Load() == 1
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
	s.meta.WriteInt64(0, s.GetAckSeq())
	s.synced.Store(1)
	return s.meta.Sync()
}

// NewSequence returns a sequence with page cache corresponding to dirPath.
func NewSequence(dirPath string) (Sequence, error) {
	meta, err := queue.NewMeta(dirPath, sequenceMetaSize)
	if err != nil {
		return nil, err
	}

	ackSeq := meta.ReadInt64(0)

	return &sequence{
		dirPath: dirPath,
		meta:    meta,
		headSeq: *atomic.NewInt64(ackSeq),
		ackSeq:  *atomic.NewInt64(ackSeq),
	}, nil
}

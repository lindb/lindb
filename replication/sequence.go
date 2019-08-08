package replication

import (
	"path"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
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
}

// sequence implements Sequence.
type sequence struct {
	dirPath string
	// meta stores the ackSeq to page cache.
	meta queue.Meta
	// headSeq represents the the max sequence num of replica received.
	headSeq int64
	// ackSeq represents the the max sequence num of replica flushed to disk.
	ackSeq int64
	// 0 -> not synced, 1 -> synced
	synced int32
}

// ResetSynced resets Synced() to false.
func (s *sequence) ResetSynced() {
	atomic.StoreInt32(&s.synced, 0)
}

// Synced checked if the Sequence has been synced.
func (s *sequence) Synced() bool {
	return atomic.LoadInt32(&s.synced) == 1
}

// GetHeadSeq returns the head sequence which is the latest sequence of replica received.
func (s *sequence) GetHeadSeq() int64 {
	return atomic.LoadInt64(&s.headSeq)
}

// SetHeadSeq sets the head sequence which is the latest sequence of replica received.
func (s *sequence) SetHeadSeq(seq int64) {
	atomic.StoreInt64(&s.headSeq, seq)
}

// GetAckSeq returns the ack sequence which is the latest sequence of replica successfully flushed to disk.
func (s *sequence) GetAckSeq() int64 {
	return atomic.LoadInt64(&s.ackSeq)
}

// GetAckSeq sets the ack sequence which is the latest sequence of replica successfully flushed to disk.
func (s *sequence) SetAckSeq(seq int64) {
	atomic.StoreInt64(&s.ackSeq, seq)
}

// Sync syncs the Sequence to storage.
func (s *sequence) Sync() error {
	s.meta.WriteInt64(0, s.GetAckSeq())
	atomic.StoreInt32(&s.synced, 1)
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
		headSeq: ackSeq,
		ackSeq:  ackSeq,
	}, nil
}

// SequenceManager manages the Sequences.
type SequenceManager interface {
	// GetSequence returns a sequence for given parameters.
	GetSequence(db string, shardID int32, node models.Node) (Sequence, bool)
	// CreateSequence creates a sequence for given parameters,
	// if the sequence already exists, directly return the existed one.
	// Concurrent safe.
	CreateSequence(db string, shardID int32, node models.Node) (Sequence, error)
}

// sequenceManager implements SequenceManager.
type sequenceManager struct {
	dirPath     string
	sequenceMap sync.Map
	lock4map    sync.Mutex
}

// GetSequence returns a sequence for given parameters.
func (sm *sequenceManager) GetSequence(db string, shardID int32, node models.Node) (Sequence, bool) {
	key := sm.buildKey(db, shardID, node)
	val, _ := sm.sequenceMap.Load(key)

	seq, ok := val.(Sequence)
	return seq, ok
}

// CreateSequence creates a sequence for given parameters,
// if the sequence already exists, directly return the existed one.
// Concurrent safe.
func (sm *sequenceManager) CreateSequence(db string, shardID int32, node models.Node) (Sequence, error) {
	key := sm.buildKey(db, shardID, node)
	val, ok := sm.sequenceMap.Load(key)
	if !ok {
		sm.lock4map.Lock()

		val, ok = sm.sequenceMap.Load(key)
		if !ok {
			dir := path.Join(sm.dirPath, db, strconv.Itoa(int(shardID)))
			if err := fileutil.MkDir(dir); err != nil {
				return nil, err
			}

			filePath := path.Join(dir, node.IP+"-"+strconv.Itoa(int(node.Port)))
			seq, err := NewSequence(filePath)
			if err != nil {
				sm.lock4map.Unlock()
				return nil, err
			}

			sm.sequenceMap.Store(key, seq)
			sm.lock4map.Unlock()
			return seq, nil
		}
	}

	seq := val.(Sequence)
	return seq, nil
}

func (sm *sequenceManager) buildKey(db string, shardID int32, node models.Node) string {
	return db + "/" + strconv.Itoa(int(shardID)) +
		"/" + node.IP + "-" + strconv.Itoa(int(node.Port))
}

func NewSequenceManager(dirPath string) (SequenceManager, error) {
	if err := fileutil.MkDir(dirPath); err != nil {
		return nil, err
	}
	return &sequenceManager{
		dirPath: dirPath,
	}, nil
}

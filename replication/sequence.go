package replication

import (
	"path"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/fileutil"
	"github.com/eleme/lindb/pkg/queue"
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
	meta    queue.Meta
	headSeq int64
	ackSeq  int64
	synced  int32
}

func (s *sequence) ResetSynced() {
	atomic.StoreInt32(&s.synced, 0)
}

func (s *sequence) Synced() bool {
	return atomic.LoadInt32(&s.synced) == 1
}

func (s *sequence) GetHeadSeq() int64 {
	return atomic.LoadInt64(&s.headSeq)
}

func (s *sequence) SetHeadSeq(seq int64) {
	atomic.StoreInt64(&s.headSeq, seq)
}

func (s *sequence) GetAckSeq() int64 {
	return atomic.LoadInt64(&s.ackSeq)
}

func (s *sequence) SetAckSeq(seq int64) {
	atomic.StoreInt64(&s.ackSeq, seq)
}

func (s *sequence) Sync() error {
	s.meta.WriteInt64(0, s.GetAckSeq())
	atomic.StoreInt32(&s.synced, 1)
	return s.meta.Sync()
}

func NewSequence(dirPath string) (Sequence, error) {
	meta, err := queue.NewMeta(dirPath, 8)
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
	// CreateSequence creates a sequence for given parameters, concurrent safe.
	CreateSequence(db string, shardID int32, node models.Node) (Sequence, error)
}

// sequenceManager implements SequenceManager.
type sequenceManager struct {
	dirPath     string
	sequenceMap sync.Map
	lock4map    sync.Mutex
}

func (sm *sequenceManager) GetSequence(db string, shardID int32, node models.Node) (Sequence, bool) {
	key := sm.buildKey(db, shardID, node)
	val, _ := sm.sequenceMap.Load(key)

	seq, ok := val.(Sequence)
	if ok {
		return seq, true
	}

	return nil, false
}

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

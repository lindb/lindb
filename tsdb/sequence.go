package tsdb

import (
	"path"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/replication"
)

//go:generate mockgen -source=./sequence.go -destination=./sequence_mock.go -package=tsdb

// for testing
var (
	newSequenceFunc = replication.NewSequence
)

// ReplicaSequence represents the shard level replica sequence
type ReplicaSequence interface {
	// getOrCreateSequence gets the replica sequence by remote replica peer if exist, else creates a new sequence
	getOrCreateSequence(remotePeer string) (replication.Sequence, error)
	// getAllHeads gets the current replica indexes for all replica remote peers
	getAllHeads() map[string]int64
	// ack acks the replica index that the data is persistent
	ack(heads map[string]int64) error
}

// replicaSequence implements ReplicaSequence
type replicaSequence struct {
	dirPath     string
	sequenceMap sync.Map
	lock4map    sync.Mutex
	syncing     atomic.Bool
}

// newReplicaSequence creates shard level replica sequence by dir path
func newReplicaSequence(dirPath string) (ReplicaSequence, error) {
	if fileutil.Exist(dirPath) {
		// if replica dir exist, load all exist replica sequences
		remotePeers, err := listDir(dirPath)
		if err != nil {
			return nil, err
		}
		ss := &replicaSequence{dirPath: dirPath}
		for _, remotePeer := range remotePeers {
			filePath := path.Join(dirPath, remotePeer)
			seq, err := newSequenceFunc(filePath)
			if err != nil {
				return nil, err
			}
			seq.SetHeadSeq(seq.GetAckSeq())
			ss.sequenceMap.Store(remotePeer, seq)
		}
		// persist new sequence
		if err := ss.syncSequence(); err != nil {
			return nil, err
		}
		return ss, nil
	}
	// create new sequence for creating shard
	if err := mkDirIfNotExist(dirPath); err != nil {
		return nil, err
	}
	return &replicaSequence{dirPath: dirPath}, nil
}

// getOrCreateSequence gets the replica sequence by remote replica peer if exist, else creates a new sequence
func (ss *replicaSequence) getOrCreateSequence(remotePeer string) (replication.Sequence, error) {
	val, ok := ss.sequenceMap.Load(remotePeer)
	if !ok {
		ss.lock4map.Lock()
		defer ss.lock4map.Unlock()
		// double check
		val, ok = ss.sequenceMap.Load(remotePeer)
		if !ok {
			filePath := path.Join(ss.dirPath, remotePeer)
			seq, err := newSequenceFunc(filePath)
			if err != nil {
				return nil, err
			}
			ss.sequenceMap.Store(remotePeer, seq)
			return seq, nil
		}
	}

	seq := val.(replication.Sequence)
	return seq, nil
}

// getAllHeads gets the current replica indexes for all replica remote peers
func (ss *replicaSequence) getAllHeads() map[string]int64 {
	result := make(map[string]int64)
	ss.sequenceMap.Range(func(key, value interface{}) bool {
		seq, ok := value.(replication.Sequence)
		if ok {
			replicaKey, ok := key.(string)
			if ok {
				result[replicaKey] = seq.GetHeadSeq()
			}
		}
		return true
	})
	return result
}

// ack acks the replica index that the data is persistent
func (ss *replicaSequence) ack(heads map[string]int64) error {
	for remotePeer, head := range heads {
		seq, ok := ss.sequenceMap.Load(remotePeer)
		if !ok {
			continue
		}
		s, ok := seq.(replication.Sequence)
		if !ok {
			continue
		}
		s.SetAckSeq(head)
	}
	return ss.syncSequence()
}

// sync syncs the all replica peer sequences
func (ss *replicaSequence) syncSequence() error {
	// make sure, just one worker does sync sequence
	var err error
	if ss.syncing.CAS(false, true) {
		ss.sequenceMap.Range(func(key, value interface{}) bool {
			seq, ok := value.(replication.Sequence)
			if ok {
				// sync one replica peer sequence
				err = seq.Sync()
			}
			return true
		})
		ss.syncing.Store(false)
	}
	return err
}

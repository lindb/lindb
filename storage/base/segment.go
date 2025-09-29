package base

import (
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/wal"
)

type Segment struct {
	Path string
	WALs map[models.NodeID]wal.WriteAheadLog // leader => write ahead log

	mutex sync.Mutex
}

func (s *Segment) GetOrCreateWAL(leader models.NodeID) (wal.WriteAheadLog, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if log, ok := s.WALs[leader]; ok {
		return log, nil
	}

	log, err := wal.NewWriteAheadLog(filepath.Join(s.Path, leader.String()))
	if err != nil {
		return nil, err
	}

	s.WALs[leader] = log

	return log, nil
}

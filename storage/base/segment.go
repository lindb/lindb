// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package base

import (
	"path/filepath"
	"sync"

	loggerpkg "github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

var logger = loggerpkg.GetLogger("storage", "base")

type Segment struct {
	TimeRange timeutil.TimeRange
	Path      string
	WALs      map[models.NodeID]store.WriteAheadLog // leader => write ahead log

	Sequence          map[models.NodeID]atomic.Int64 // leader => consume sequence(wal)
	ImmutableSequence map[models.NodeID]int64
	PersistSequence   map[models.NodeID]atomic.Int64 // leader=>ack sequence(wal)

	CreateWriteAheadLog func(p string) (store.WriteAheadLog, error)

	ref atomic.Int32

	mutex sync.Mutex
}

func (s *Segment) SegmentTimeRange() timeutil.TimeRange {
	return s.TimeRange
}

func (s *Segment) GetOrCreateWAL(leader models.NodeID) (store.WriteAheadLog, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if log, ok := s.WALs[leader]; ok {
		return log, nil
	}

	log, err := s.CreateWriteAheadLog(filepath.Join(s.Path, leader.String()))
	if err != nil {
		return nil, err
	}

	s.WALs[leader] = log

	return log, nil
}

// ValidateSequence validates replica sequence if valid.
func (s *Segment) ValidateSequence(leader models.NodeID, seq int64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if seqForLeader, ok := s.Sequence[leader]; ok {
		return seq > seqForLeader.Load()
	}
	return true
}

// CommitSequence commits written sequence after write data.
func (s *Segment) CommitSequence(leader models.NodeID, seq int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	seqForLeader := s.Sequence[leader]
	seqForLeader.Store(seq)
	s.Sequence[leader] = seqForLeader
}

// AckSequence acknowledges sequence after memory database flush successfully.
func (s *Segment) AckSequence(leader models.NodeID, fn func(seq int64)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO:??
	// s.callbacks[leader] = append(s.callbacks[leader], fn)

	seqForLeader, ok := s.PersistSequence[leader]
	logger.Info("register ack sequence callback",
		loggerpkg.String("path", s.Path), loggerpkg.Any("sequences", s.Sequence),
		loggerpkg.Any("leader", leader), loggerpkg.Any("exist", ok))
	if ok {
		// invoke ack sequence after register function, maybe some cases lost ack index.
		fn(seqForLeader.Load())
	}
}

// Retain increments write ref count
func (s *Segment) Retain() {
	s.ref.Inc()
}

// Release decrements write ref count,
// if ref==0, no data will write this family.
func (s *Segment) Release() {
	s.ref.Dec()
}

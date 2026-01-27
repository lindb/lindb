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
	"strconv"
	"sync"

	"github.com/lindb/common/pkg/fileutil"
	loggerpkg "github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

var logger = loggerpkg.GetLogger("storage", "base")

type Segment struct {
	TimeRange timeutil.TimeRange
	Interval  timeutil.Interval
	Path      string
	WALs      map[models.NodeID]store.WriteAheadLog // leader => write ahead log

	CreateWriteAheadLog func(p string) (store.WriteAheadLog, error)

	ref atomic.Int32

	mutex sync.Mutex
}

func (s *Segment) LoadWALs(ackSequences map[int32]int64) error {
	if !fileutil.Exist(s.Path) {
		return nil
	}
	leaders, err := fileutil.ListDir(s.Path)
	if err != nil {
		return err
	}
	for _, leader := range leaders {
		nodeID, err := strconv.ParseInt(leader, 10, 64)
		if err != nil {
			return err
		}
		wal, err := s.GetOrCreateWAL(models.NodeID(nodeID))
		if err != nil {
			return err
		}
		// set replica sequcence for local replicator
		ack, ok := ackSequences[int32(nodeID)]
		if !ok {
			wal.AckSequence(models.NodeID(nodeID), ack)
		}
	}
	return nil
}

func (s *Segment) SegmentTimeRange() timeutil.TimeRange {
	return s.TimeRange
}

func (s *Segment) NumOfPoints() int {
	return s.TimeRange.NumOfPoints(s.Interval)
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

	if leader == meta.CurrentNode() {
		// subscribe storage meta change events, if node is leader.
		// only leader node do replica or observe.
		// unsubscribe when write ahead log closed.
		meta.GetStorageMetaManager().Subscribe(log)
	}

	return log, nil
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

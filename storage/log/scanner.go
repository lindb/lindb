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

package log

import (
	"encoding/binary"

	logspkg "github.com/lindb/arrow/pkg/logs"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/store"
)

type Scanner interface {
	HasNext() bool
	// Next calls fn with the log reader, the row number within the WAL record, and
	// the per-segment log sequence number (logID). logID is required to construct
	// the composite pagination cursor (timestamp_ns, segment_start_ms, logID).
	Next(fn func(reader *logspkg.Reader, rowNum int, logID uint32)) error
	Close()
}

type scanner struct {
	segment *Segment

	logIDs *roaring.Bitmap
	it     roaring.IntPeekable

	wal map[models.NodeID]store.WriteAheadLog // leader => write ahead log

	current struct {
		leader   models.NodeID
		sequence uint32

		data []byte
	}

	reader *logspkg.Reader
}

func NewScanner(segment *Segment, logIDs *roaring.Bitmap) Scanner {
	return &scanner{
		segment: segment,
		logIDs:  logIDs,
		it:      logIDs.Iterator(),
		wal:     segment.GetWALs(),
	}
}

func (s *scanner) HasNext() bool {
	return s.it.HasNext()
}

func (s *scanner) Next(fn func(reader *logspkg.Reader, rowNum int, logID uint32)) error {
	logID := s.it.Next()
	index, err := s.segment.GetIndex(logID)
	if err != nil {
		return err
	}
	leader := models.NodeID(index[0])
	sequence := binary.LittleEndian.Uint32(index[1:])
	rowNum := binary.LittleEndian.Uint32(index[5:])
	if leader != s.current.leader || sequence != s.current.sequence || len(s.current.data) == 0 {
		// update current log info, if leader or sequence changed, read log data from wal.
		// if leader and sequence are same, it means log from same wal record.
		s.current.leader = leader
		s.current.sequence = sequence

		// read log data from wal
		data, err := s.wal[leader].Get(int64(sequence))
		if err != nil {
			return err
		}
		s.current.data = data

		if err := s.initializeReader(); err != nil {
			return err
		}
	}

	fn(s.reader, int(rowNum), logID)

	return nil
}

func (s *scanner) initializeReader() (err error) {
	if s.reader == nil {
		reader, err := logspkg.NewReader(s.current.data)
		if err != nil {
			return err
		}
		s.reader = reader
	} else {
		// release previous log and attribute record batch before read new data
		err = s.reader.Reset(s.current.data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *scanner) Close() {
	if s.reader != nil {
		s.reader.Release()
	}
}

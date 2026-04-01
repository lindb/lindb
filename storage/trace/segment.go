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

package trace

import (
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cockroachdb/pebble/v2"
	"github.com/lindb/arrow/pkg/traces"
	"github.com/lindb/common/pkg/fileutil"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/flush"
	"github.com/lindb/lindb/storage/store"
)

// estimatedBytesPerRow is a rough estimate of pebble memory per trace row
// (traceID key + 5-byte WAL index value).
const estimatedBytesPerRow = 256

// writeOpts disables the write-ahead log for better write performance.
var writeOpts = pebble.NoSync

type Segment struct {
	base.Segment

	partition *partition

	db *pebble.DB

	buf []byte

	reader *traces.TraceIDReader

	isFlushing  atomic.Bool  // guard against concurrent flush (Flush() closes the DB)
	createdTime int64        // unix nanoseconds, set at construction time
	rowsWritten atomic.Int64 // total rows written; used for MemSize estimation
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	segmentTime := store.MinuteIntervalCalc.CalcSegmentTime(timestamp)
	familySlot := store.MinuteIntervalCalc.CalcFamily(timestamp, segmentTime)
	family := fmt.Sprintf("%d", familySlot)
	segmentPath := filepath.Join(partition.Path(), family)

	indexPath := path.Join(segmentPath, "index")
	if err := fileutil.MkDirIfNotExist(indexPath); err != nil {
		return nil, err
	}

	opts := &pebble.Options{
		// Disable the WAL — data is recovered from the segment's own WAL files.
		DisableWAL: true,
	}
	db, err := pebble.Open(indexPath, opts)
	if err != nil {
		return nil, err
	}
	start := store.MinuteIntervalCalc.CalcFamilyStartTime(segmentTime, familySlot)

	seg := &Segment{
		Segment: base.Segment{
			TimeRange: timeutil.TimeRange{
				Start: start,
				End:   store.MinuteIntervalCalc.CalcFamilyEndTime(start),
			},
			Interval: store.MinuteInterval,
			Path:     segmentPath,
			WALs:     make(map[models.NodeID]store.WriteAheadLog),
		},
		partition:   partition,
		db:          db,
		createdTime: time.Now().UnixNano(),

		buf: make([]byte, 5),
	}

	wals, err := fileutil.ListDir(segmentPath)
	if err != nil {
		return nil, err
	}

	seg.CreateWriteAheadLog = func(p string) (store.WriteAheadLog, error) {
		return store.CreateWriteAheadLog(p, seg)
	}

	for _, leader := range wals {
		if leader == "index" {
			continue
		}
		l, err := strconv.Atoi(leader)
		if err != nil {
			// TODO: add metric
			continue
		}
		// check leader
		_, err = seg.GetOrCreateWAL(models.NodeID(l))
		if err != nil {
			return nil, err
		}
	}

	// register with flush tracker so the checker can schedule flush for this segment
	flush.GetMemDBTracker().Register(seg)

	return seg, nil
}

func (seg *Segment) Partition() store.Partition {
	return seg.partition
}

// MutableMemDBInfo implements flush.FlushableSegment.
// Trace segments write directly to pebble; we use rowsWritten as a proxy for memory pressure.
// Returns nil after Flush() has been called (db is closed, segment is terminal).
func (seg *Segment) MutableMemDBInfo() *flush.MemDBInfo {
	if seg.isFlushing.Load() {
		// flush in progress or already flushed (terminal state)
		return nil
	}
	rows := seg.rowsWritten.Load()
	if rows == 0 {
		return nil
	}
	created := time.Unix(0, seg.createdTime)
	return &flush.MemDBInfo{
		MemSize:     rows * estimatedBytesPerRow,
		CreatedTime: seg.createdTime,
		SegmentTime: seg.TimeRange.Start,
		NumOfRows:   int(rows),
		Uptime:      time.Since(created),
	}
}

// SegmentKey implements flush.FlushableSegment.
// Returns a globally unique key: "{dbName}/{shardID}/{segmentTime}".
func (seg *Segment) SegmentKey() string {
	db := seg.partition.shard.Database().(*Database)
	return flush.MakeSegmentKey(db.Name(), seg.partition.shard.ShardID().String(), seg.TimeRange.Start)
}

// SegmentMeta implements flush.FlushableSegment.
// For trace segments, size/TTL thresholds are intentionally 0 (disabled);
// only SegmentRange triggering is expected.
func (seg *Segment) SegmentMeta() flush.SegmentMeta {
	db := seg.partition.shard.Database().(*Database)
	dbOpt := db.GetOption().Option
	ahead, behind := dbOpt.GetAcceptWritableRange()

	return flush.SegmentMeta{
		DatabaseName:         db.Name(),
		ShardID:              seg.partition.shard.ShardID().String(),
		SizeThresholdBytes:   0, // disabled for trace: flush is a terminal operation
		TimeThresholdNano:    0, // disabled for trace
		Ahead:                ahead,
		Behind:               behind,
		SegmentOutRangeDelay: dbOpt.Data.SegmentOutRangeDelay,
	}
}

func (seg *Segment) Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error) {
	// OPT: thread safe reader, avoid new reader for each write
	if seg.reader == nil {
		seg.reader, err = traces.NewTraceIDReader(msg)
		if err != nil {
			return 0, err
		}
	} else {
		if err := seg.reader.Reset(msg); err != nil {
			return 0, err
		}
	}

	// Write the 5-byte WAL pointer into seg.buf: [leaderID(1)] + [seq(4)]
	seg.buf[0] = byte(leader)
	stream.PutUint32(seg.buf, 1, uint32(seq))

	traceIDs := make(map[string]struct{})
	numOfRows := seg.reader.NumOfRows()
	for i := 0; i < numOfRows; i++ {
		traceID := seg.reader.TraceID(i)
		traceIDStr := strutil.ByteSlice2String(traceID)
		if _, ok := traceIDs[traceIDStr]; !ok {
			// Use pebble's native Merge (backed by DefaultMerger / AppendValueMerger)
			// to atomically append the 5-byte WAL pointer
			if err := seg.db.Merge(traceID, seg.buf, writeOpts); err != nil {
				return 0, err
			}
			traceIDs[traceIDStr] = struct{}{}
		}
	}
	seg.rowsWritten.Add(int64(numOfRows))
	return numOfRows, nil
}

func (seg *Segment) GetTrace(traceID string) (rs [][]byte, err error) {
	data, closer, err := seg.db.Get([]byte(traceID))
	if err == pebble.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer closer.Close() //nolint:errcheck

	for i := range len(data) / 5 {
		index := data[i*5 : (i+1)*5]
		leader := models.NodeID(index[0])
		sequence := binary.LittleEndian.Uint32(index[1:])
		trace, err := seg.WALs[leader].Get(int64(sequence))
		if err != nil {
			return nil, err
		}
		rs = append(rs, trace)
	}
	return rs, nil
}

func (seg *Segment) Close() error {
	// unregister from flush tracker before closing
	flush.GetMemDBTracker().Unregister(seg)

	seg.Flush()

	for _, d := range seg.WALs {
		d.Close()
	}

	return nil
}

// Flush implements store.Segment.
// NOTE: for trace segments this is a terminal operation — it flushes and closes the pebble instance.
// The CAS guard ensures it only executes once even if called concurrently.
func (seg *Segment) Flush() error {
	if !seg.isFlushing.CompareAndSwap(false, true) {
		// already flushed or flush in progress
		return nil
	}
	// do not restore isFlushing — once flushed the segment is terminal
	seg.db.Flush() //nolint:errcheck
	seg.db.Close() //nolint:errcheck
	// reset rowsWritten so MutableMemDBInfo returns nil after flush
	seg.rowsWritten.Store(0)
	return nil
}

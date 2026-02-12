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
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	logspkg "github.com/lindb/arrow/pkg/logs"
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/log/memdb"
	"github.com/lindb/lindb/storage/log/tblstore"
	"github.com/lindb/lindb/storage/store"
)

type Segment struct {
	base.Segment

	shard     store.Shard
	partition *partition
	index     queue.Queue

	family kv.Family

	timestampIndexes sync.Map

	immutable memdb.Database
	mutable   memdb.Database

	buf []byte

	numOfPoints int

	reader    *logspkg.BinaryReader
	logReader *logspkg.Reader

	mutex sync.RWMutex
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	segmentTime := store.MinuteIntervalCalc.CalcSegmentTime(timestamp)
	familySlot := store.MinuteIntervalCalc.CalcFamily(timestamp, segmentTime)
	family := fmt.Sprintf("%d", familySlot)
	segmentPath := filepath.Join(partition.Path(), family)
	index, err := queue.NewQueue(path.Join(segmentPath, "index"), 128*1024*1024)
	if err != nil {
		return nil, err
	}
	kvFamily := partition.kvStore.GetFamily(family)
	if kvFamily == nil {
		// create kv family
		var err error
		familyOption := kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(tblstore.LogIndexMerger),
		}
		kvFamily, err = partition.kvStore.CreateFamily(family, familyOption)
		if err != nil {
			return nil, err
		}
	}
	segmentStartTime := store.MinuteIntervalCalc.CalcFamilyStartTime(segmentTime, familySlot)
	db := partition.shard.Database().(*Database)
	seg := &Segment{
		Segment: base.Segment{
			TimeRange: timeutil.TimeRange{
				Start: segmentStartTime,
				End:   store.MinuteIntervalCalc.CalcFamilyEndTime(segmentStartTime),
			},
			Interval: store.MinuteInterval,
			Path:     segmentPath,
			WALs:     make(map[models.NodeID]store.WriteAheadLog),
		},
		partition: partition,
		family:    kvFamily,
		index:     index,

		mutable: memdb.NewDatabase(db.indexDB),

		buf: make([]byte, 9),
	}

	seg.numOfPoints = seg.TimeRange.NumOfPoints(store.MinuteInterval)

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
	return seg, nil
}

func (s *Segment) Partition() store.Partition {
	return s.partition
}

func (s *Segment) NumOfPoints() int {
	return s.numOfPoints
}

func (s *Segment) FindLogIDsByTimeRange(timeRange timeutil.TimeRange, callback func(timestamp int64, logIDs *roaring.Bitmap)) {
	snapshot := s.partition.timestampIndex.GetSnapshot()
	defer snapshot.Close()

	interval := store.MinuteInterval.Int64()

	target := (&timeRange).Intersect(s.SegmentTimeRange())
	start := target.Start - target.Start%interval
	end := target.End - target.End%interval
	partitionTime := s.partition.PartitionTime()
	temp := roaring.New()
	result := roaring.New()
	if err := walkTimeRange(start, end, interval, func(timestamp int64) error {
		idsObj, _ := s.timestampIndexes.Load(timestamp)
		if ids, ok := idsObj.(*roaring.Bitmap); ok {
			result.Or(ids)
		} else {
			fmt.Println("not found.....")
		}
		if err := snapshot.Load(uint32(timestamp-partitionTime), func(value []byte) error {
			_, err := encoding.BitmapUnmarshal(temp, value)
			if err != nil {
				return err
			}
			result.Or(temp)
			return nil
		}); err != nil {
			panic(err)
		}

		if !result.IsEmpty() {
			callback(timestamp, result)
		}

		result.Clear()
		return nil
	}); err != nil {
		panic(err)
	}
}

func (s *Segment) FindLogIDsByFields(fieldIDs []uint32) *roaring.Bitmap {
	snapshot := s.family.GetSnapshot()
	defer snapshot.Close()

	logIDs := roaring.New()

	temp := roaring.New()
	for _, fieldID := range fieldIDs {
		if err := snapshot.Load(fieldID, func(value []byte) error {
			_, err := encoding.BitmapUnmarshal(temp, value)
			if err != nil {
				return err
			}
			logIDs.Or(temp)
			return nil
		}); err != nil {
			panic(err)
		}
		if memLogIDs := s.mutable.FindLogIDsByField(fieldID); memLogIDs != nil {
			logIDs.Or(memLogIDs)
		}
	}
	return logIDs
}

func (s *Segment) GetIndex(logID uint32) ([]byte, error) {
	return s.index.Get(int64(logID))
}

func (s *Segment) GetWALs() map[models.NodeID]store.WriteAheadLog {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make(map[models.NodeID]store.WriteAheadLog, len(s.WALs))
	for k, v := range s.WALs {
		result[k] = v
	}
	return result
}

func (s *Segment) Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error) {
	var logs, attributes arrow.RecordBatch

	// FIXME: need optimize, multiple leader write to same segment(thread not safe)
	if s.logReader == nil {
		reader, err := logspkg.NewBinaryReader()
		if err != nil {
			return 0, err
		}
		s.reader = reader
		logs, attributes, err = reader.ReadFrom(msg)
		if err != nil {
			return 0, err
		}
		s.logReader = logspkg.NewReader(logs, attributes)
	} else {
		logs, attributes, err = s.reader.ReadFrom(msg)
		if err != nil {
			return 0, err
		}
		s.logReader.Reset(logs, attributes)
	}

	defer func() {
		logs.Release()
		attributes.Release()
	}()

	it := s.logReader.Iterator()

	for it.HasNext() {
		row := it.Next()
		// generate log id then index it
		logID := uint32(s.index.AppendedSeq() + 1)
		s.buf[0] = byte(leader)
		stream.PutUint32(s.buf, 1, uint32(seq))
		stream.PutUint32(s.buf, 5, uint32(row))
		s.index.Put(s.buf)

		s.index.AppendedSeq()

		// build secondary index for log timestmap
		s.indexTimestamp(s.logReader.Timestamp(row)/1000_000, logID)

		// build secondary index for log fields
		s.mutable.Write([]byte("ns"), logID, s.logReader, row)
	}

	return s.logReader.NumRows(), nil
}

func (s *Segment) Flush() error {
	// flush timestamp index
	if err := s.FlushTimestampIndex(); err != nil {
		return err
	}
	// flush log fields index
	flusher := s.family.NewFlusher()
	if err := s.mutable.Flush(flusher); err != nil {
		return err
	}
	return flusher.Commit()
}

func (s *Segment) Close() error {
	s.Flush()

	for _, d := range s.WALs {
		d.Close()
	}

	s.index.Close()

	return nil
}

func (s *Segment) indexTimestamp(timestamp int64, logID uint32) {
	// truncate timestamp based on interval
	targetTimestamp := timestamp - timestamp%store.MinuteInterval.Int64()
	fmt.Println("index timestamp", targetTimestamp)

	index, ok := s.timestampIndexes.Load(targetTimestamp)
	if ok {
		(index.(*roaring.Bitmap)).Add(logID)
	} else {
		s.timestampIndexes.Store(targetTimestamp, roaring.BitmapOf(logID))
	}
}

func (s *Segment) FlushTimestampIndex() error {
	flusher := s.partition.timestampIndex.NewFlusher()
	start := s.TimeRange.Start
	end := s.TimeRange.End
	interval := store.MinuteInterval.Int64()
	partitionTime := s.partition.PartitionTime()
	if err := walkTimeRange(start, end, interval, func(timestamp int64) error {
		idsObj, _ := s.timestampIndexes.Load(timestamp)
		if ids, ok := idsObj.(*roaring.Bitmap); ok {
			data, err := ids.ToBytes()
			if err != nil {
				return err
			}
			// store timestamp index(key=offset based partition timestamp,value=log ids)
			flusher.Add(uint32(timestamp-partitionTime), data)
		}
		return nil
	}); err != nil {
		return err
	}
	return flusher.Commit()
}

func walkTimeRange(start, end, interval int64, fn func(timestamp int64) error) error {
	step := start
	for i := int64(0); step <= end; i++ {
		step = start + i*interval
		if err := fn(step); err != nil {
			return err
		}
	}
	return nil
}

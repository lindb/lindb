package log

import (
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/proto/gen/v1/flatLogV1"
	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	logproto "github.com/lindb/lindb/proto/log"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/log/memdb"
	"github.com/lindb/lindb/storage/log/tblstore"
	"github.com/lindb/lindb/storage/store"
	"github.com/lindb/lindb/storage/wal"
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

	running atomic.Bool

	buf []byte

	mutex sync.Mutex
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)
	familySlot := intervalCalc.CalcFamily(timestamp, segmentTime)
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
	segmentStartTime := intervalCalc.CalcFamilyStartTime(segmentTime, familySlot)
	seg := &Segment{
		Segment: base.Segment{
			TimeRange: timeutil.TimeRange{
				Start: segmentStartTime,
				End:   intervalCalc.CalcFamilyEndTime(segmentStartTime),
			},
			Path: segmentPath,
			WALs: make(map[models.NodeID]wal.WriteAheadLog),
		},
		partition: partition,
		family:    kvFamily,
		index:     index,
		running:   *atomic.NewBool(true),

		mutable: memdb.NewDatabase(partition.shard.database.indexDB),

		buf: make([]byte, 8),
	}

	wals, err := fileutil.ListDir(segmentPath)
	if err != nil {
		return nil, err
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

	// start build index goroutine
	go seg.buildIndex()

	return seg, nil
}

func (s *Segment) FindLogIDsByTimeRange(timeRange timeutil.TimeRange) *roaring.Bitmap {
	snapshot := s.partition.timestampIndex.GetSnapshot()
	defer snapshot.Close()

	target := (&timeRange).Intersect(s.SegmentTimeRange())
	start := target.Start - target.Start%minuteInterval.Int64()
	end := target.End - target.End%minuteInterval.Int64()
	interval := minuteInterval.Int64()
	partitionTime := s.partition.PartitionTime()
	logIDs := roaring.New()
	temp := roaring.New()
	for i := int64(1); start <= end; i++ {
		idsObj, _ := s.timestampIndexes.Load(start)
		if ids, ok := idsObj.(*roaring.Bitmap); ok {
			logIDs.Or(ids)
		}
		fmt.Printf("find mem start: %d, end: %d, i: %v\n", start, end, idsObj)
		if err := snapshot.Load(uint32(start-partitionTime), func(value []byte) error {
			_, err := encoding.BitmapUnmarshal(temp, value)
			if err != nil {
				return err
			}
			fmt.Printf("find disk start: %d, i: %v\n", start, temp)
			logIDs.Or(temp)
			return nil
		}); err != nil {
			panic(err)
		}
		start += i * interval
	}
	s.timestampIndexes.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", key, value)
		return true
	})
	fmt.Printf("logIDs: %v\n", logIDs)
	return logIDs
}

func (s *Segment) FindLogIDsByFields(fieldIDs []uint32) *roaring.Bitmap {
	snapshot := s.family.GetSnapshot()
	defer snapshot.Close()

	logIDs := roaring.New()

	temp := roaring.New()
	for _, fieldID := range fieldIDs {
		if err := snapshot.Load(fieldID, func(value []byte) error {
			fmt.Printf("==>>>>>fieldID: %d, value: %s\n", fieldID, value)
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

func (s *Segment) GetLog(logID uint32) ([]byte, error) {
	id, err := s.index.Get(int64(logID))
	if err != nil {
		return nil, err
	}
	index := binary.LittleEndian.Uint32(id[1:])
	return s.WALs[models.NodeID(id[0])].Get(int64(index))
}

func (s *Segment) Close() error {
	s.Flush()

	for _, d := range s.WALs {
		d.Close()
	}

	s.index.Close()
	s.running.Store(false)

	return nil
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

func (s *Segment) buildIndex() {
	for s.running.Load() {
		for leader, log := range s.WALs {
			seq, data, err := log.Consume()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if data != nil {
				s.indexLog(leader, seq, data)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *Segment) indexLog(leader models.NodeID, index int64, msg []byte) {
	log := &flatLogV1.Log{}
	log.Init(msg, flatbuffers.GetUOffsetT(msg))

	// generate log id then index it
	logID := uint32(s.index.AppendedSeq() + 1)
	s.buf[0] = byte(leader)
	stream.PutUint32(s.buf, 1, uint32(logID))
	s.index.Put(s.buf)

	// build secondary index for log timestmap
	s.indexTimestamp(log.Timestamp(), logID)

	// build secondary index for log fields
	s.mutable.Write([]byte("ns"), logID, logproto.NewFieldIterator(log))
}

func (s *Segment) indexTimestamp(timestamp int64, logID uint32) {
	// truncate timestamp based on interval
	targetTimestamp := timestamp - timestamp%minuteInterval.Int64()

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
	interval := minuteInterval.Int64()
	partitionTime := s.partition.PartitionTime()
	for i := int64(1); start <= end; i++ {
		idsObj, _ := s.timestampIndexes.Load(start)
		if ids, ok := idsObj.(*roaring.Bitmap); ok {
			data, err := ids.ToBytes()
			if err != nil {
				return err
			}
			// store timestamp index(key=offset based partition timestamp,value=log ids)
			flusher.Add(uint32(start-partitionTime), data)
		}

		fmt.Printf("start: %d, end: %d, i: %v\n", start, end, idsObj)

		start += i * interval
	}
	s.timestampIndexes.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", key, value)
		return true
	})
	return flusher.Commit()
}

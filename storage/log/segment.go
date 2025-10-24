package log

import (
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"sync"

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
			Path:              segmentPath,
			WALs:              make(map[models.NodeID]store.WriteAheadLog),
			Sequence:          make(map[models.NodeID]atomic.Int64),
			ImmutableSequence: make(map[models.NodeID]int64),
			PersistSequence:   make(map[models.NodeID]atomic.Int64),
		},
		partition: partition,
		family:    kvFamily,
		index:     index,

		mutable: memdb.NewDatabase(partition.shard.database.indexDB),

		buf: make([]byte, 8),
	}

	seg.numOfPoints = seg.TimeRange.NumOfPoints(minuteInterval)

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

	target := (&timeRange).Intersect(s.SegmentTimeRange())
	start := target.Start - target.Start%minuteInterval.Int64()
	end := target.End - target.End%minuteInterval.Int64()
	interval := minuteInterval.Int64()
	partitionTime := s.partition.PartitionTime()
	temp := roaring.New()
	result := roaring.New()
	if err := walkTimeRange(start, end, interval, func(timestamp int64) error {
		idsObj, _ := s.timestampIndexes.Load(timestamp)
		if ids, ok := idsObj.(*roaring.Bitmap); ok {
			result.Or(ids)
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
			fmt.Printf("timestamp:%d,log ids:%v\n", timestamp, result.GetCardinality())
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

func (s *Segment) GetLog(logID uint32) ([]byte, error) {
	id, err := s.index.Get(int64(logID))
	if err != nil {
		return nil, err
	}
	index := binary.LittleEndian.Uint32(id[1:])
	return s.WALs[models.NodeID(id[0])].Get(int64(index))
}

func (s *Segment) Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error) {
	log := &flatLogV1.Log{}
	log.Init(msg, flatbuffers.GetUOffsetT(msg))

	// generate log id then index it
	logID := uint32(s.index.AppendedSeq() + 1)
	s.buf[0] = byte(leader)
	stream.PutUint32(s.buf, 1, uint32(logID))
	s.index.Put(s.buf)

	s.index.AppendedSeq()

	// build secondary index for log timestmap
	s.indexTimestamp(log.Timestamp(), logID)

	// build secondary index for log fields
	s.mutable.Write([]byte("ns"), logID, logproto.NewFieldIterator(log))
	return 1, nil
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

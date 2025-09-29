package log

import (
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
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
	logproto "github.com/lindb/lindb/proto/log"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/log/memdb"
	"github.com/lindb/lindb/storage/log/tblstore"
	"github.com/lindb/lindb/storage/store"
	"github.com/lindb/lindb/storage/wal"
)

type Segment struct {
	base.Segment

	shard store.Shard
	index queue.Queue

	family kv.Family

	immutable memdb.Database
	mutable   memdb.Database

	running atomic.Bool

	buf []byte
}

func NewSegment(timestmap int64, partition *partition) (store.Segment, error) {
	family := fmt.Sprintf("%d", intervalCalc.CalcFamily(timestmap, intervalCalc.CalcSegmentTime(timestmap)))
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
	seg := &Segment{
		Segment: base.Segment{
			Path: segmentPath,
			WALs: make(map[models.NodeID]wal.WriteAheadLog),
		},
		family:  kvFamily,
		index:   index,
		running: *atomic.NewBool(true),

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

	go seg.buildIndex()
	return seg, nil
}

func (s *Segment) GetLogIDs(fieldID uint32) *roaring.Bitmap {
	snapshot := s.family.GetSnapshot()
	defer snapshot.Close()
	logIDs := roaring.New()
	temp := roaring.New()
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
	if memLogIDs := s.mutable.GetLogIDs(fieldID); memLogIDs != nil {
		logIDs.Or(memLogIDs)
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

	// build secondary index for log(timestamp/fields)
	s.mutable.Write([]byte("ns"), logID, log.Timestamp(), logproto.NewFieldIterator(log))
}

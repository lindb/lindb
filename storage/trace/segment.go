package trace

import (
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"strconv"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/linxGnu/grocksdb"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/store"
)

var (
	wo *grocksdb.WriteOptions
	ro *grocksdb.ReadOptions
)

func init() {
	wo = grocksdb.NewDefaultWriteOptions()
	wo.DisableWAL(true)

	ro = grocksdb.NewDefaultReadOptions()
}

type Segment struct {
	base.Segment

	partition *partition

	db *grocksdb.DB

	buf []byte
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)
	familySlot := intervalCalc.CalcFamily(timestamp, segmentTime)
	family := fmt.Sprintf("%d", familySlot)
	segmentPath := filepath.Join(partition.Path(), family)

	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetMergeOperator(&TraceMergeOperator{})
	indexPath := path.Join(segmentPath, "index")
	if err := fileutil.MkDirIfNotExist(indexPath); err != nil {
		return nil, err
	}
	db, err := grocksdb.OpenDb(opts, indexPath)
	fmt.Println(indexPath)
	if err != nil {
		return nil, err
	}
	start := intervalCalc.CalcFamilyStartTime(segmentTime, familySlot)

	seg := &Segment{
		Segment: base.Segment{
			TimeRange: timeutil.TimeRange{
				Start: start,
				End:   intervalCalc.CalcFamilyEndTime(start),
			},
			Path:              segmentPath,
			WALs:              make(map[models.NodeID]store.WriteAheadLog),
			Sequence:          make(map[models.NodeID]atomic.Int64),
			ImmutableSequence: make(map[models.NodeID]int64),
			PersistSequence:   make(map[models.NodeID]atomic.Int64),
		},
		partition: partition,
		db:        db,

		buf: make([]byte, 8),
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

	return seg, nil
}

func (seg *Segment) Partition() store.Partition {
	return seg.partition
}

func (seg *Segment) Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error) {
	req := ptraceotlp.NewExportRequest()
	if err = req.UnmarshalProto(msg); err != nil {
		fmt.Println(err)
		return
	}
	traceIDs := make(map[string]struct{})
	traces := req.Traces()
	spans := traces.ResourceSpans()
	for i := range spans.Len() {
		s := spans.At(i)
		scopeSpans := s.ScopeSpans()
		for j := range scopeSpans.Len() {
			span := scopeSpans.At(j)
			sSpans := span.Spans()
			for k := range sSpans.Len() {
				ss := sSpans.At(k)
				if _, ok := traceIDs[ss.TraceID().String()]; !ok {
					seg.db.Merge(wo, []byte(ss.TraceID().String()), encoding.U32ToBytes(uint32(seq)))
					traceIDs[ss.TraceID().String()] = struct{}{}
					fmt.Printf("traceID:%s, index:%d\n", ss.TraceID().String(), seq)
				}
			}
		}
	}
	return 1, nil
}

func (seg *Segment) GetTrace(traceID string) (rs [][]byte, err error) {
	indexes, err := seg.db.Get(ro, []byte(traceID))
	if err != nil {
		return nil, err
	}
	if !indexes.Exists() {
		fmt.Println("trace not found")
		return nil, nil
	}
	data := indexes.Data()
	fmt.Printf("get data len=%d\n", len(data))
	for i := range len(data) / 4 {
		index := binary.BigEndian.Uint32(data[i*4:])
		trace, err := seg.WALs[models.NodeID(1)].Get(int64(index))
		if err != nil {
			return nil, err
		}
		fmt.Printf("logid===%d,len=%d\n", index, len(trace))
		rs = append(rs, trace)
	}
	fmt.Printf("get data len=%d\n", len(rs))
	return rs, nil
}

func (seg *Segment) Close() error {
	seg.Flush()

	for _, d := range seg.WALs {
		d.Close()
	}

	return nil
}

func (seg *Segment) Flush() error {
	seg.db.Flush(grocksdb.NewDefaultFlushOptions())
	seg.db.Close()
	return nil
}

func (seg *Segment) indexTrace(leader models.NodeID, index int64, msg []byte) {
	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(msg); err != nil {
		fmt.Println(err)
		return
	}
	traceIDs := make(map[string]struct{})
	traces := req.Traces()
	spans := traces.ResourceSpans()
	for i := range spans.Len() {
		s := spans.At(i)
		scopeSpans := s.ScopeSpans()
		for j := range scopeSpans.Len() {
			span := scopeSpans.At(j)
			sSpans := span.Spans()
			for k := range sSpans.Len() {
				ss := sSpans.At(k)
				if _, ok := traceIDs[ss.TraceID().String()]; !ok {
					seg.db.Merge(wo, []byte(ss.TraceID().String()), encoding.U32ToBytes(uint32(index)))
					traceIDs[ss.TraceID().String()] = struct{}{}
					fmt.Printf("traceID:%s, index:%d\n", ss.TraceID().String(), index)
				}
			}
		}
	}
}

package trace

import (
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/linxGnu/grocksdb"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/store"
	"github.com/lindb/lindb/storage/wal"
)

type Segment struct {
	base.Segment

	shard store.Shard

	db *grocksdb.DB

	running atomic.Bool

	buf []byte
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	family := fmt.Sprintf("%d", intervalCalc.CalcFamily(timestamp,
		intervalCalc.CalcSegmentTime(timestamp)))
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

	seg := &Segment{
		Segment: base.Segment{
			Path: segmentPath,
			WALs: make(map[models.NodeID]wal.WriteAheadLog),
		},
		db:      db,
		running: *atomic.NewBool(true),

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

func (segment *Segment) GetTrace(traceID string) (rs [][]byte, err error) {
	options := grocksdb.NewDefaultReadOptions()
	indexes, err := segment.db.Get(options, []byte(traceID))
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
		trace, err := segment.WALs[models.NodeID(1)].Get(int64(index))
		if err != nil {
			return nil, err
		}
		fmt.Printf("logid===%d,len=%d\n", index, len(trace))
		rs = append(rs, trace)
	}
	fmt.Printf("get data len=%d\n", len(rs))
	return rs, nil
}

func (segment *Segment) Close() error {
	segment.Flush()

	for _, d := range segment.WALs {
		d.Close()
	}

	segment.running.Store(false)

	return nil
}

func (segment *Segment) Flush() error {
	segment.db.Flush(grocksdb.NewDefaultFlushOptions())
	segment.db.Close()
	return nil
}

func (segment *Segment) buildIndex() {
	for segment.running.Load() {
		for leader, log := range segment.WALs {
			seq, data, err := log.Consume()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if data != nil {
				segment.indexTrace(leader, seq, data)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (segment *Segment) indexTrace(leader models.NodeID, index int64, msg []byte) {
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
					wo := grocksdb.NewDefaultWriteOptions()
					wo.DisableWAL(true)

					segment.db.Merge(wo, []byte(ss.TraceID().String()), encoding.U32ToBytes(uint32(index)))
					traceIDs[ss.TraceID().String()] = struct{}{}
					fmt.Printf("traceID:%s, index:%d\n", ss.TraceID().String(), index)
				}
			}
		}
	}
}

type TraceMergeOperator struct{}

func (op *TraceMergeOperator) Name() string {
	return "TraceMergeOperator"
}

func (op *TraceMergeOperator) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	fmt.Println("full merge")
	total := len(existingValue)
	for _, v := range operands {
		total += len(v)
	}
	dest := make([]byte, total)
	offset := copy(dest, existingValue)
	for _, operand := range operands {
		offset += copy(dest[offset:], operand)
	}
	return dest, true
}

func (op *TraceMergeOperator) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	fmt.Println("merge...")
	dest := make([]byte, (len(rightOperand) + len(leftOperand)))
	copy(dest, leftOperand)
	copy(dest[len(leftOperand):], rightOperand)
	return dest, true
}

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

	"github.com/lindb/common/pkg/fileutil"
	"github.com/linxGnu/grocksdb"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

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
	segmentTime := store.MinuteIntervalCalc.CalcSegmentTime(timestamp)
	familySlot := store.MinuteIntervalCalc.CalcFamily(timestamp, segmentTime)
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
				// TODO: using trace id directly
				if _, ok := traceIDs[ss.TraceID().String()]; !ok {
					seg.db.Merge(wo, []byte(ss.TraceID().String()), encoding.U32ToBytes(uint32(seq)))
					traceIDs[ss.TraceID().String()] = struct{}{}
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
		return nil, nil
	}
	data := indexes.Data()
	for i := range len(data) / 4 {
		index := binary.BigEndian.Uint32(data[i*4:])
		trace, err := seg.WALs[models.NodeID(1)].Get(int64(index))
		if err != nil {
			return nil, err
		}
		rs = append(rs, trace)
	}
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
				}
			}
		}
	}
}

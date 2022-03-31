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

package metricsdata

import (
	"encoding/binary"
	"io"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metricsdata

// Flusher is a wrapper of kv.Builder, provides ability to flush a metric-table file to disk.
// The layout is available in `tsdb/doc.go`
// Level1: metric-block
// Level2: series entry
// Level3: compressed field data
//
// flush step:
// 1. flush field metas of metric level
// 2. flush field store of one series
// 3. flush series id
// 4. flush metric data include field metadata and all series ids data
type Flusher interface {
	// PrepareMetric prepares to write a new metric block
	PrepareMetric(
		metricID uint32,
		fieldMetas field.Metas,
	)
	// FlushField writes a compressed field data to writer.
	// It will be called in order with field metas even if field data is empty
	FlushField(data []byte) error
	// FlushSeries writes a full series, this will be called after writing all fields of this entry.
	FlushSeries(seriesID uint32) error
	// CommitMetric ends writing a full metric block
	// this will be called after writing all entries of this metric.
	CommitMetric(slotRange timeutil.SlotRange) error
	// GetFieldMetas returns current field metas of metric.
	GetFieldMetas() field.Metas

	// Closer closes the writer, syncs all data to the file.
	io.Closer
}

// flusher implements Flusher.
type flusher struct {
	// Level1 flusher
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━Layout of Metric Data Table━━━━━━━━━━━━━━━━━━━━━━
	//                     Level1
	//                    +---------+---------+---------+---------+---------+
	//                    │ Metric  │ Metric  │ Metric  │ Metric  │ Footer  │
	//                    │ Block   │ Block   │ Offsets │ Bitmap  │         │
	//                    +---------+---------+---------+---------+---------+
	//                   /           \
	//                  /             \
	//                 /               \
	//                /                 \
	//   +-----------+                   +-----------------+
	//  /                 Level2                            \
	// v--------+--------+--------+--------+--------+--------v
	// │ Series │ Series │  Field | Series │ HighKey│ Footer │
	// │ Bucket │ Bucket │  Metas | Bitmap │ Offsets│        │
	// +--------+--------+--------+--------+--------+--------+
	//
	//
	// Level2(Fields Meta)
	// ┌─────────────────────────────────────────────────────────────────┐
	// │                      Fields Meta                                │
	// ├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
	// │   Count  │ FieldID  │  Field   │ FieldID  │  Field   │          │
	// │          │ (uint16) │  Type    │ (uint16) │  Type    │  ......  │
	// ├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
	// │  1 Byte  │  1 Bytes │ 1 Byte   │  1 Bytes │ 1 Byte   │          │
	// └──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘
	//
	// Level2 (KV table: Series Bucket Footer)
	// ┌──────────────────────────────────────────────────────┐
	// │                    Footer                            │
	// ├──────────┬──────────┬──────────┬──────────┬──────────┤
	// │   time   │ position │ position │ position │  CRC32   │
	// │   range  │ OfMetas  │ OfBitmap │ OfOffsets│ CheckSum │
	// ├──────────┼──────────┼──────────┼──────────┼──────────┤
	// │  4 Byte  │ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │  4 Bytes │
	// └──────────┴──────────┴──────────┴──────────┴──────────┘
	//
	// Level2 is a context of the second level in kv table, used for a writing a full metric
	// each entry is a series bucket ordered by roaring high key
	// Resets it after completed writing a metric
	Level2 struct {
		fieldMetas     field.Metas
		seriesIDs      *roaring.Bitmap
		highKeyOffsets *encoding.FixedOffsetEncoder
		footer         [dataFooterSize]byte
	}
	// +--------+--------+--------+--------+--------+--------v
	// │ Series │ Series │  Field | Series │ HighKey│ Footer │
	// │ Bucket │ Bucket │  Metas | Bitmap │ Offsets│        │
	// +--------+--------+--------+--------+--------+--------+
	// │        │         Level3
	// v--------v--------+--------+--------+--------+
	// │ Series │ Series │ Series │ LowKey │ PosOf  |
	// │ Entry  │ Entry  │ Entry  │ Offsets│Offsets |
	// +--------+--------+--------+--------+--------+
	//
	// Level3 is the third level of kv table, context for writing series,
	// each entry is a different series ordered by low key
	// Resets it after completed writing bulk of series in a low container of roaring
	// each entry is a series
	Level3 struct {
		// startAt is the absolute position in the Level2's Series Bucket
		startAt int
		// highKey is a higher 16 bits of seriesIDs.
		// query will be performed concurrently by high key.
		highKey uint16
		// highKeySetEver symbols if highKey has been set before
		// sets it to true after flushing series
		isHighKeySetEver bool
		// low container of series ids
		// offset = seriesEntryPosition - start position of level3
		lowKeyOffsets *encoding.FixedOffsetEncoder
	}
	// v--------v--------+--------+--------+--------+
	// │ Series │ Series │ Series │ LowKey │ PosOf  |
	// │ Entry  │ Entry  │ Entry  │ Offsets│Offsets |
	// +--------+--------+--------+--------+--------+
	// │         \        \        \
	// │          \        \        \
	// │           \        \        |
	// │            \        \       +--------------------------+
	// │             \        +--------------------------+       \
	// │              +------------------+                \       \
	// │                  Level4          \                \        \
	// v--------+--------+--------+--------+                +--------+
	// │ Field  │ Field  │ Field  | LenOf  │                │ Field  |
	// │ Data   │ Data   │ Offsets│ Offsets│                │ Data   |
	// +--------+--------+--------+--------+                +--------+
	// uvariant64 encoding LenOfOffsets, reversed
	//
	// Level4 is the fourth level of kv table, context for writing field data.
	// Each entry is a different field data ordered by field ids
	// The following are two different scenarios：
	// Case1: Single Field Metric
	//        FieldData without other information
	// Case2: Multi Field Metric
	//        FieldData + FieldData + FieldData + FieldOffsets
	Level4 struct {
		// startAt is the absolute position in Level2's SeriesEntry
		startAt int
		// scratch for variant encoding field offsets marshal size's length
		scratch [binary.MaxVarintLen64]byte
		// fieldsOffsets holds distances between startAt and position of fieldData
		fieldDataOffsets *encoding.FixedOffsetEncoder
		fieldBuffer      [][]byte
		fieldAppendIdx   int
	}
}

// NewFlusher returns a new Flusher,
// interval is used to calculate the time-range of field data slots.`
func NewFlusher(kvFlusher kv.Flusher) (Flusher, error) {
	sw, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	flusher := &flusher{
		kvFlusher: kvFlusher,
		kvWriter:  sw,
	}
	// level2 context
	flusher.Level2.seriesIDs = roaring.New()
	flusher.Level2.highKeyOffsets = encoding.NewFixedOffsetEncoder(true)
	// level3 context
	flusher.Level3.lowKeyOffsets = encoding.NewFixedOffsetEncoder(true)
	// level4 context
	flusher.Level4.fieldDataOffsets = encoding.NewFixedOffsetEncoder(true)
	return flusher, nil
}

func (w *flusher) PrepareMetric(
	metricID uint32,
	fieldMetas field.Metas,
) {
	w.kvWriter.Prepare(metricID)
	w.Level2.fieldMetas = fieldMetas
	w.Level2.highKeyOffsets.Add(0)

	w.Level4.fieldBuffer = make([][]byte, len(fieldMetas))
	w.Level4.fieldAppendIdx = 0
}

func (w *flusher) FlushField(data []byte) error {
	// just buffer field data
	w.Level4.fieldBuffer[w.Level4.fieldAppendIdx] = data
	w.Level4.fieldAppendIdx++
	return nil
}

func (w *flusher) flushField() error {
	defer func() {
		w.Level4.fieldAppendIdx = 0
	}()
	isMultiField := w.Level2.fieldMetas.Len() > 1
	for fieldIdx := range w.Level4.fieldBuffer {
		data := w.Level4.fieldBuffer[fieldIdx]
		// if metric only has one field, just writes field data
		fieldDataAt := int(w.kvWriter.Size()) - w.Level4.startAt
		if _, err := w.kvWriter.Write(data); err != nil {
			return err
		}
		// if metric only has one field, just writes field data
		// multi fields, write the field offset
		if isMultiField {
			w.Level4.fieldDataOffsets.Add(fieldDataAt)
		}
	}
	// flush field offsets in necessary(multi field).
	if isMultiField {
		if err := w.writeLevel4OffsetsFooter(); err != nil {
			return err
		}
	}
	return nil
}

func (w *flusher) writeLevel4OffsetsFooter() error {
	// pick level4's start position of Offsets
	beforeLen := w.kvWriter.Size()
	// write level4's FieldOffsets
	if err := w.Level4.fieldDataOffsets.Write(w.kvWriter); err != nil {
		return err
	}

	// write level4's length of Offsets
	writtenLen := stream.PutUvariantLittleEndian(w.Level4.scratch[:], uint64(w.kvWriter.Size()-beforeLen))
	// reverse uvaiant little endian encoding
	_, err := w.kvWriter.Write(w.Level4.scratch[:writtenLen])
	return err
}

// FlushSeries writes a full series,
// this will be called after writing all fields of this entry.
// 1. only one field: series data = field data
// 2. multi-fields: series data = field offsets + fields data
func (w *flusher) FlushSeries(seriesID uint32) error {
	// reset level4 context
	defer func() {
		w.Level4.startAt = int(w.kvWriter.Size())
		w.Level4.fieldDataOffsets.Reset()
	}()

	seriesHasData := w.Level4.fieldAppendIdx > 0
	if !seriesHasData {
		// if not field data, drop this series
		return nil
	}
	// isHighKeySetEver means this is the first high key
	// If different high key arrives, level3 is done.
	// otherwise, we are still at level3, just keeps the low-key's offset
	highKey := encoding.HighBits(seriesID)
	if !w.Level3.isHighKeySetEver {
		w.Level3.isHighKeySetEver = true
		w.Level3.highKey = highKey
	}
	// first need check high key, because if first write offset, first field of low series(next high key) will be lost.
	if highKey != w.Level3.highKey {
		// flush data by diff high key
		if err := w.flushLevel2SeriesBucket(); err != nil {
			return err
		}
		// set high key, for next container storage
		w.Level3.highKey = highKey
		// reset low keys for next container of a different high key
		w.Level3.lowKeyOffsets.Reset()
		w.Level3.startAt = int(w.kvWriter.Size())
		// set high key offset to current series bucket
		w.Level2.highKeyOffsets.Add(int(w.kvWriter.Size()))
		w.Level4.startAt = int(w.kvWriter.Size())
	}

	// write field's offset for current series id
	w.Level3.lowKeyOffsets.Add(int(w.kvWriter.Size()) - w.Level3.startAt)
	// write field data
	if err := w.flushField(); err != nil {
		return err
	}
	// add series id into index block of metric
	w.Level2.seriesIDs.Add(seriesID)
	return nil
}

func (w *flusher) flushLevel2SeriesBucket() error {
	posOfLowKeyOffsets := int(w.kvWriter.Size()) - w.Level3.startAt
	if !(posOfLowKeyOffsets > 0) {
		return nil
	}
	// data in this series bucket has been flushed.
	// flush LowKey-Offsets in Level3
	if err := w.Level3.lowKeyOffsets.Write(w.kvWriter); err != nil {
		return err
	}
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(posOfLowKeyOffsets))
	_, err := w.kvWriter.Write(scratch[:])
	return err
}

func (w *flusher) reset() {
	w.Level2.fieldMetas = w.Level2.fieldMetas[:0]
	w.Level2.seriesIDs.Clear()
	w.Level2.highKeyOffsets.Reset()

	w.Level3.startAt = 0
	w.Level3.isHighKeySetEver = false
	w.Level3.lowKeyOffsets.Reset()

	w.Level4.startAt = 0
	w.Level4.fieldDataOffsets.Reset()
}

// CommitMetric writes a full metric-block,
// this will be called after writing all entries of this metric.
func (w *flusher) CommitMetric(slotRange timeutil.SlotRange) error {
	defer w.reset()

	// no metric data written ever
	if w.Level2.seriesIDs.IsEmpty() {
		return nil
	}
	if err := w.flushLevel2SeriesBucket(); err != nil {
		return err
	}

	// write fields-meta
	fieldMetasAt := w.kvWriter.Size()
	// write field-count
	if _, err := w.kvWriter.Write([]byte{byte(len(w.Level2.fieldMetas))}); err != nil {
		return err
	}
	// write field-id, field-type list
	for _, fm := range w.Level2.fieldMetas {
		// write field-id, field-type
		if _, err := w.kvWriter.Write([]byte{
			byte(fm.ID),
			byte(fm.Type),
		}); err != nil {
			return err
		}
	}
	// write series ids bitmap
	seriesIDAt := w.kvWriter.Size()
	if _, err := w.Level2.seriesIDs.WriteTo(w.kvWriter); err != nil {
		return err
	}
	// write high offsets
	highKeyOffsetsAt := w.kvWriter.Size()

	if err := w.Level2.highKeyOffsets.Write(w.kvWriter); err != nil {
		return err
	}

	//////////////////////////////////////////////////
	// build footer (field meta's offset+series ids' offset+high level offsets+crc32 checksum)
	// (2 bytes + 2 bytes +4 bytes + 4 bytes + 4 bytes + 4 bytes)
	//////////////////////////////////////////////////
	// write time range of metric level
	binary.LittleEndian.PutUint16(w.Level2.footer[:2], slotRange.Start)
	binary.LittleEndian.PutUint16(w.Level2.footer[2:4], slotRange.End)
	// write field metas' start position
	binary.LittleEndian.PutUint32(w.Level2.footer[4:8], fieldMetasAt)
	// write series ids' start position
	binary.LittleEndian.PutUint32(w.Level2.footer[8:12], seriesIDAt)
	// write offset block start position
	binary.LittleEndian.PutUint32(w.Level2.footer[12:16], highKeyOffsetsAt)
	// write CRC32 checksum
	binary.LittleEndian.PutUint32(w.Level2.footer[16:20], w.kvWriter.CRC32CheckSum())

	if _, err := w.kvWriter.Write(w.Level2.footer[:]); err != nil {
		return err
	}
	return w.kvWriter.Commit()
}

// Close adds the footer and then closes the kv builder,
// this will be called after writing all metric-blocks.
func (w *flusher) Close() error {
	return w.kvFlusher.Commit()
}

// GetFieldMetas returns the file metas of current metric.
func (w *flusher) GetFieldMetas() field.Metas {
	return w.Level2.fieldMetas
}

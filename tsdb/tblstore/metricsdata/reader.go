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
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package metricsdata

// for testing
var (
	getOffsetFunc = getOffset
)

const (
	dataFooterSize = 2 + // start time slot
		2 + // end time slot
		4 + // field metas position
		4 + // series ids position
		4 + // high offsets position
		4 // crc32 checksum
)

// MetricReader represents the metric block metricReader
type MetricReader interface {
	// Path returns file path
	Path() string
	// GetSeriesIDs returns the series ids in this sst file
	GetSeriesIDs() *roaring.Bitmap
	// GetFields returns the field metas in this sst file
	GetFields() field.Metas
	// GetTimeRange returns the time range in this sst file
	GetTimeRange() timeutil.SlotRange
	// Load loads the data from sst file, then returns the file metric scanner.
	Load(highKey uint16, seriesID roaring.Container, fields field.Metas) flow.DataLoader
	// readSeriesData reads series data from file by given position.
	readSeriesData(position int) [][]byte
}

// metricReader implements MetricReader interface that reads metric block
type metricReader struct {
	path          string
	buf           []byte
	highOffsets   *encoding.FixedOffsetDecoder
	seriesIDs     *roaring.Bitmap
	fields        field.Metas
	crc32CheckSum uint32
	timeRange     timeutil.SlotRange

	readFieldIndexes []int // read field indexes be used when query metric data
}

// NewReader creates a metric block metricReader
func NewReader(path string, buf []byte) (MetricReader, error) {
	r := &metricReader{
		path: path,
		buf:  buf,
	}
	if err := r.initReader(); err != nil {
		return nil, err
	}
	return r, nil
}

// Path returns the file path
func (r *metricReader) Path() string {
	return r.path
}

// GetSeriesIDs returns the series ids in this sst file
func (r *metricReader) GetSeriesIDs() *roaring.Bitmap {
	return r.seriesIDs
}

// GetFields returns the field metas in this sst file
func (r *metricReader) GetFields() field.Metas {
	return r.fields
}

// GetTimeRange returns the time range in this sst file
func (r *metricReader) GetTimeRange() timeutil.SlotRange {
	return r.timeRange
}

// prepare prepares the field aggregator based on query condition
func (r *metricReader) prepare(fields field.Metas) (found bool) {
	fieldMap := make(map[field.ID]int)
	for idx, fieldMeta := range r.fields {
		fieldMap[fieldMeta.ID] = idx
	}
	r.readFieldIndexes = make([]int, len(fields))
	for idx, f := range fields { // sort by field ids
		fieldIdx, ok := fieldMap[f.ID]
		if !ok {
			r.readFieldIndexes[idx] = -1
		} else {
			r.readFieldIndexes[idx] = fieldIdx
			found = true
		}
	}
	return
}

// Load loads the data from sst file, then returns the file metric scanner.
func (r *metricReader) Load(highKey uint16, seriesID roaring.Container, fields field.Metas) flow.DataLoader {
	// 1. get high container index by the high key of series ID
	highContainerIdx := r.seriesIDs.GetContainerIndex(highKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series IDs not exist) return it
		return nil
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := r.seriesIDs.GetContainerAtIndex(highContainerIdx)
	foundSeriesIDs := lowContainer.And(seriesID)
	if foundSeriesIDs.GetCardinality() == 0 {
		return nil
	}
	offset, _ := r.highOffsets.Get(highContainerIdx)
	seriesOffsets := encoding.NewFixedOffsetDecoder(r.buf[offset:])

	if !r.prepare(fields) {
		// field not found
		return nil
	}
	// must use lowContainer from store, because get series index based on container
	return newMetricLoader(r, lowContainer, seriesOffsets)
}

// readSeriesData reads series data from file by given position.
func (r *metricReader) readSeriesData(position int) [][]byte {
	fieldCount := r.fields.Len()
	if fieldCount == 1 {
		// metric has one field, just read the data
		return [][]byte{r.buf[position:]}
	}
	// read data for multi-fields
	seriesData := r.buf[position:]
	fieldOffsets := encoding.NewFixedOffsetDecoder(seriesData)
	fieldsData := seriesData[fieldOffsets.Header()+fieldCount*fieldOffsets.ValueWidth():]
	rs := make([][]byte, len(r.readFieldIndexes))
	for i, idx := range r.readFieldIndexes {
		if idx == -1 {
			continue
		}
		offset, ok := fieldOffsets.Get(idx)
		if ok {
			// read field data
			rs[i] = fieldsData[offset:]
		}
	}
	return rs
}

// initReader initializes the metricReader context includes tag value ids/high offsets
func (r *metricReader) initReader() error {
	if len(r.buf) <= dataFooterSize {
		return fmt.Errorf("block length not ok")
	}
	// read footer(2+2+4+4+4+4)
	footerPos := len(r.buf) - dataFooterSize
	r.timeRange.Start = stream.ReadUint16(r.buf, footerPos)
	r.timeRange.End = stream.ReadUint16(r.buf, footerPos+2)

	fieldMetaStartPos := int(stream.ReadUint32(r.buf, footerPos+4))
	seriesIDsStartPos := int(stream.ReadUint32(r.buf, footerPos+8))
	highOffsetsPos := int(stream.ReadUint32(r.buf, footerPos+12))
	r.crc32CheckSum = stream.ReadUint32(r.buf, footerPos+16)
	// validate offsets
	if fieldMetaStartPos > footerPos || seriesIDsStartPos > highOffsetsPos {
		return fmt.Errorf("bad offsets")
	}

	// read field metas
	offset := fieldMetaStartPos
	fieldCount := r.buf[offset]
	offset++
	r.fields = make(field.Metas, fieldCount)
	for i := byte(0); i < fieldCount; i++ {
		r.fields[i] = field.Meta{
			ID:   field.ID(r.buf[offset]),
			Type: field.Type(r.buf[offset+1]),
		}
		offset += 2
	}

	// read series ids
	seriesIDs := roaring.New()
	if err := encoding.BitmapUnmarshal(seriesIDs, r.buf[seriesIDsStartPos:]); err != nil {
		return err
	}
	r.seriesIDs = seriesIDs
	// read high offsets
	r.highOffsets = encoding.NewFixedOffsetDecoder(r.buf[highOffsetsPos:])
	return nil
}

// fieldIndexes returns field indexes of metric level
func (r *metricReader) fieldIndexes() map[field.ID]int {
	result := make(map[field.ID]int)
	for idx, f := range r.fields {
		result[f.ID] = idx
	}
	return result
}

// dataScanner represents the metric data scanner which scans the series data when merge operation
type dataScanner struct {
	reader        *metricReader
	container     roaring.Container
	seriesOffsets *encoding.FixedOffsetDecoder

	highKeys  []uint16
	highKey   uint16
	seriesPos int
}

// newDataScanner creates a data scanner for data merge
func newDataScanner(r MetricReader) *dataScanner {
	reader := r.(*metricReader)
	s := &dataScanner{
		reader:   reader,
		highKeys: reader.seriesIDs.GetHighKeys(),
	}
	s.nextContainer()
	return s
}

// fieldIndexes returns field indexes of metric level
func (s *dataScanner) fieldIndexes() map[field.ID]int {
	return s.reader.fieldIndexes()
}

// nextContainer goes next container context for scanner
func (s *dataScanner) nextContainer() {
	s.highKey = s.highKeys[s.seriesPos]
	s.container = s.reader.seriesIDs.GetContainerAtIndex(s.seriesPos)
	offset, _ := s.reader.highOffsets.Get(s.seriesPos)
	s.seriesOffsets = encoding.NewFixedOffsetDecoder(s.reader.buf[offset:])
	s.seriesPos++
}

// slotRange returns the slot range of metric level in current sst file
func (s *dataScanner) slotRange() timeutil.SlotRange {
	return s.reader.GetTimeRange()
}

// scan scans the data and returns series position if series id exist, else returns -1
func (s *dataScanner) scan(highKey, lowSeriesID uint16) int {
	if s.highKey < highKey {
		if s.seriesPos >= len(s.highKeys) {
			// current tag inverted no data can read
			return -1
		}
		s.nextContainer()
	}
	if highKey != s.highKey {
		// high key not match, return it
		return -1
	}
	// find data by low series id
	if s.container.Contains(lowSeriesID) {
		// get the index of low series id in container
		idx := s.container.Rank(lowSeriesID)
		// get series data data position
		offset, ok := getOffsetFunc(s.seriesOffsets, idx-1)
		if !ok {
			return -1
		}
		return offset
	}
	return -1
}

// getOffset returns the offset by idx
func getOffset(seriesOffsets *encoding.FixedOffsetDecoder, idx int) (int, bool) {
	return seriesOffsets.Get(idx)
}

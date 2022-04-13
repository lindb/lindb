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

package tagindex

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/tag"

	"github.com/lindb/roaring"
)

//go:generate mockgen -source ./forward_reader.go -destination=./forward_reader_mock.go -package tagindex

// ForwardReader represents read forward index data(series id=>tag value id)
type ForwardReader interface {
	flow.Grouping
	// GetSeriesIDsForTagKeyID returns series ids for spec tag key id of metric.
	GetSeriesIDsForTagKeyID(tagKeyID tag.KeyID) (*roaring.Bitmap, error)
}

// forwardReader implements ForwardReader
type forwardReader struct {
	readers []table.Reader
}

// NewForwardReader creates a Reader for reading forward index
func NewForwardReader(readers []table.Reader) ForwardReader {
	return &forwardReader{
		readers: readers,
	}
}

// GetSeriesIDsForTagKeyID get series ids for spec tag key id of metric.
func (r *forwardReader) GetSeriesIDsForTagKeyID(tagKeyID tag.KeyID) (*roaring.Bitmap, error) {
	seriesIDs := roaring.New()
	if err := r.findReader(tagKeyID, func(reader TagForwardReader) {
		seriesIDs.Or(reader.getSeriesIDs())
	}); err != nil {
		return nil, err
	}
	return seriesIDs, nil
}

// GetGroupingScanner returns the grouping scanners based on tag key ids and series ids
func (r *forwardReader) GetGroupingScanner(tagKeyID tag.KeyID, seriesIDs *roaring.Bitmap) ([]flow.GroupingScanner, error) {
	var scanners []flow.GroupingScanner
	if err := r.findReader(tagKeyID, func(reader TagForwardReader) {
		// check reader if it has series ids(after filtering)
		seriesIDs.And(reader.getSeriesIDs())
		if seriesIDs.IsEmpty() {
			// not found
			return
		}
		// found series ids in the sst file
		scanners = append(scanners, reader)
	}); err != nil {
		return nil, err
	}
	return scanners, nil
}

// findReader finds the tag forward reader by tag key id, if reader exist, will invoke callback function
func (r *forwardReader) findReader(tagKeyID tag.KeyID, callback func(reader TagForwardReader)) error {
	for _, reader := range r.readers {
		value, err := reader.Get(uint32(tagKeyID))
		if err != nil {
			continue
		}
		indexReader, err := NewTagForwardReader(value)
		if err != nil {
			return err
		}

		callback(indexReader)
	}
	return nil
}

// TagForwardReader represents the forward index inverterReader for one tag(series id=>tag value id)
type TagForwardReader interface {
	flow.GroupingScanner
	// getSeriesIDs gets all series ids under this tag key
	getSeriesIDs() *roaring.Bitmap
}

// tagForwardReader implements TagForwardReader interface
type tagForwardReader struct {
	baseReader
}

// NewTagForwardReader creates a forward index inverterReader
func NewTagForwardReader(buf []byte) (TagForwardReader, error) {
	r := &tagForwardReader{baseReader: baseReader{
		buf: buf,
	}}
	if err := r.initReader(); err != nil {
		return nil, err
	}
	return r, nil
}

// GetSeriesAndTagValue returns group by container and tag value ids
func (r *tagForwardReader) GetSeriesAndTagValue(highKey uint16) (lowSeriesIDs roaring.Container, tagValueIDs []uint32) {
	index := r.keys.GetContainerIndex(highKey)
	if index < 0 {
		// data not found
		return nil, nil
	}
	// tag value ids cannot reuse, because
	offset, _ := r.offsets.Get(index)
	tagValueIDsFromFile := encoding.NewDeltaBitPackingDecoder(r.buf[offset:])

	lowSeriesIDs = r.keys.GetContainerAtIndex(index)
	tagValueIDsCount := lowSeriesIDs.GetCardinality()
	tagValueIDs = make([]uint32, tagValueIDsCount)
	i := 0
	for tagValueIDsFromFile.HasNext() {
		tagValueIDs[i] = uint32(tagValueIDsFromFile.Next())
		i++
	}
	return lowSeriesIDs, tagValueIDs
}

// getSeriesIDs gets all series ids under this tag key
func (r *tagForwardReader) getSeriesIDs() *roaring.Bitmap {
	return r.keys
}

// tagForwardScanner represents the tag forward index scanner which scans the index data when merge operation
type tagForwardScanner struct {
	reader          *tagForwardReader
	container       roaring.Container
	tagValueOffsets *encoding.FixedOffsetDecoder
	highKeys        []uint16
	highKey         uint16
	keyPos          int
	tagValueIDs     *encoding.DeltaBitPackingDecoder
}

// newTagForwardScanner creates a tag forward index scanner
func newTagForwardScanner(reader TagForwardReader) *tagForwardScanner {
	forwardReader := reader.(*tagForwardReader)
	s := &tagForwardScanner{
		reader:   forwardReader,
		highKeys: forwardReader.keys.GetHighKeys(),
	}
	s.nextContainer()
	return s
}

// nextContainer goes next container context for scanner
func (s *tagForwardScanner) nextContainer() {
	s.highKey = s.highKeys[s.keyPos]
	s.container = s.reader.keys.GetContainerAtIndex(s.keyPos)
	offset, _ := s.reader.offsets.Get(s.keyPos)
	s.tagValueOffsets = encoding.NewFixedOffsetDecoder()
	_, _ = s.tagValueOffsets.Unmarshal(s.reader.buf[offset:])
	if s.tagValueIDs == nil {
		offset, _ := s.reader.offsets.Get(s.keyPos)
		s.tagValueIDs = encoding.NewDeltaBitPackingDecoder(s.reader.buf[offset:])
	} else {
		offset, _ := s.reader.offsets.Get(s.keyPos)
		s.tagValueIDs.Reset(s.reader.buf[offset:])
	}
	s.keyPos++
}

// scan scans the data then merges the tag value ids into target tag value ids
func (s *tagForwardScanner) scan(highKey, lowSeriesID uint16, tagValueIDs []uint32) []uint32 {
	if s.highKey < highKey {
		if s.keyPos >= len(s.highKeys) {
			// current tag inverted no data can read
			return tagValueIDs
		}
		s.nextContainer()
	}
	if highKey != s.highKey {
		// high key not match, return it
		return tagValueIDs
	}
	// find data by low tag value id
	if s.container.Contains(lowSeriesID) {
		tagValueIDs = append(tagValueIDs, uint32(s.tagValueIDs.Next()))
	}
	return tagValueIDs
}

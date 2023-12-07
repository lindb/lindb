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

package v1

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./forward_reader.go -destination=./forward_reader_mock.go -package v1

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
		seriesIDs.Or(reader.GetSeriesIDs())
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
		finalSeriesIDs := roaring.FastAnd(seriesIDs, reader.GetSeriesIDs())
		if finalSeriesIDs.IsEmpty() {
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
		indexReader, err := newTagForwardReader(value)
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
}

// tagForwardReader implements TagForwardReader interface
type tagForwardReader struct {
	buf       []byte
	lut       []int
	seriesIDs *roaring.Bitmap
}

// NewTagForwardReader creates a forward index inverterReader
func NewTagForwardReader(buf []byte) (TagForwardReader, error) {
	seriesIDs := roaring.New()
	size, err := bitmapUnmarshal(seriesIDs, buf)
	if err != nil {
		return nil, err
	}
	// calc lookup table
	highKeys := seriesIDs.GetHighKeys()
	lut := make([]int, len(highKeys)+1)
	lut[0] = 0
	for idx := range highKeys {
		lowContainer := seriesIDs.GetContainerAtIndex(idx)
		lut[idx+1] = lowContainer.GetCardinality()
	}
	return &tagForwardReader{
		buf:       buf[size:],
		lut:       lut,
		seriesIDs: seriesIDs,
	}, nil
}

// GetSeriesAndTagValue returns group by container and tag value ids
func (r *tagForwardReader) GetSeriesAndTagValue(highKey uint16) (lowSeriesIDs roaring.Container, tagValueIDs []uint32) {
	index := r.seriesIDs.GetContainerIndex(highKey)
	if index < 0 {
		// data not found
		return nil, nil
	}
	// tag value ids cannot reuse, because
	c := r.lut[index]
	lowSeriesIDs = r.seriesIDs.GetContainerAtIndex(index)
	tagValueIDs = encoding.BytesToU32Slice(r.buf[c*4 : c*4+lowSeriesIDs.GetCardinality()*4])
	return lowSeriesIDs, tagValueIDs
}

// GetSeriesIDs gets all series ids under this tag key
func (r *tagForwardReader) GetSeriesIDs() *roaring.Bitmap {
	return r.seriesIDs
}

// tagForwardScanner represents the tag forward index scanner which scans the index data when merge operation
type tagForwardScanner struct {
	reader      TagForwardReader
	highKey     uint16
	container   roaring.Container
	tagValueIDs []uint32

	tagValueIdx int
}

// newTagForwardScanner creates a tag forward index scanner
func newTagForwardScanner(reader TagForwardReader) *tagForwardScanner {
	min := reader.GetSeriesIDs().Minimum()
	s := &tagForwardScanner{
		reader: reader,
	}
	s.highKey = encoding.HighBits(min)
	s.nextContainer(s.highKey)
	return s
}

// nextContainer goes next container context for scanner
func (s *tagForwardScanner) nextContainer(highKey uint16) {
	lowContainer, tagValueIDs := s.reader.GetSeriesAndTagValue(highKey)
	s.container = lowContainer
	s.tagValueIDs = tagValueIDs
	s.tagValueIdx = 0
	s.highKey = highKey
}

// scan scans the data then merges the tag value ids into target tag value ids
func (s *tagForwardScanner) scan(highKey, lowSeriesID uint16, tagValueIDs []uint32) []uint32 {
	if s.highKey < highKey {
		s.nextContainer(highKey)
	}
	if highKey != s.highKey {
		// high key not match, return it
		return tagValueIDs
	}
	if s.container == nil {
		// not tag value ids found, return it
		return tagValueIDs
	}
	// find data by low tag value id
	if s.container.Contains(lowSeriesID) {
		tagValueIDs = append(tagValueIDs, s.tagValueIDs[s.tagValueIdx])
		s.tagValueIdx++
	}
	return tagValueIDs
}

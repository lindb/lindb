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

package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

//go:generate mockgen -source ./inverted_reader.go -destination=./inverted_reader_mock.go -package invertedindex

// InvertedReader reads seriesID bitmap from series-index-table
type InvertedReader interface {
	// GetSeriesIDsByTagValueIDs finds series ids by tag key id and tag value ids
	GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
}

// inverterReader implements InvertedReader
type inverterReader struct {
	readers []table.Reader
}

// NewInvertedReader creates a InvertedReader for reading inverted index
func NewInvertedReader(readers []table.Reader) InvertedReader {
	return &inverterReader{
		readers: readers,
	}
}

// GetSeriesIDsByTagValueIDs finds series ids by tag key id and tag value ids
func (r *inverterReader) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		return roaring.New(), nil
	}
	fn := func(indexReader *tagInvertedReader) (*roaring.Bitmap, error) {
		return indexReader.getSeriesIDsByTagValueIDs(tagValueIDs)
	}
	return r.loadSeriesIDs(tagKeyID, fn)
}

// loadSeriesIDs loads the series ids by tag key id, function need implement condition
func (r *inverterReader) loadSeriesIDs(tagKeyID uint32, fn func(indexReader *tagInvertedReader) (*roaring.Bitmap, error)) (*roaring.Bitmap, error) {
	seriesIDs := roaring.New()
	for _, reader := range r.readers {
		value, err := reader.Get(tagKeyID)
		if err != nil {
			continue
		}
		indexReader, err := newTagInvertedReader(value)
		if err != nil {
			return nil, err
		}
		ids, err := fn(indexReader)
		if err != nil {
			return nil, err
		}
		seriesIDs.Or(ids)
	}
	return seriesIDs, nil
}

// tagInvertedReader represents the inverted index inverterReader for one tag(tag value ids=>series ids)
type tagInvertedReader struct {
	baseReader
}

// newTagInvertedReader creates an inverted index tagInvertedReader
func newTagInvertedReader(buf []byte) (*tagInvertedReader, error) {
	r := &tagInvertedReader{
		baseReader: baseReader{buf: buf},
	}
	if err := r.initReader(); err != nil {
		return nil, err
	}
	return r, nil
}

// getSeriesIDsByTagValueIDs finds series ids by tag value ids under this tag key
func (r *tagInvertedReader) getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	result := roaring.New()
	// get final tag value ids need to load
	finalTagValueIDs := roaring.And(tagValueIDs, r.keys)
	highKeys := finalTagValueIDs.GetHighKeys()
	for idx, highKey := range highKeys {
		loadLowContainer := finalTagValueIDs.GetContainerAtIndex(idx)
		lowContainerIdx := r.keys.GetContainerIndex(highKey)
		lowContainer := r.keys.GetContainerAtIndex(lowContainerIdx)
		offset, _ := r.offsets.Get(lowContainerIdx)
		seriesOffsets := encoding.NewFixedOffsetDecoder()
		_, err := seriesOffsets.Unmarshal(r.buf[offset:])
		if err != nil {
			return nil, err
		}
		it := loadLowContainer.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// get the index of low tag value id in container
			lowIdx := lowContainer.Rank(lowTagValueID)
			seriesPos, _ := seriesOffsets.Get(lowIdx - 1)

			// unmarshal series ids
			seriesIDs := roaring.New()
			if err := encoding.BitmapUnmarshal(seriesIDs, r.buf[seriesPos:]); err != nil {
				return nil, err
			}
			result.Or(seriesIDs)
		}
	}
	return result, nil
}

// tagInvertedScanner represents the tag inverted index scanner which scans the index data when merge operation
type tagInvertedScanner struct {
	reader        *tagInvertedReader
	container     roaring.Container
	seriesOffsets *encoding.FixedOffsetDecoder
	highKeys      []uint16
	highKey       uint16
	keyPos        int
}

// newTagInvertedScanner creates a tag inverted index scanner
func newTagInvertedScanner(reader *tagInvertedReader) *tagInvertedScanner {
	s := &tagInvertedScanner{
		reader:   reader,
		highKeys: reader.keys.GetHighKeys(),
	}
	s.nextContainer()
	return s
}

// nextContainer goes next container context for scanner
func (s *tagInvertedScanner) nextContainer() {
	s.highKey = s.highKeys[s.keyPos]
	s.container = s.reader.keys.GetContainerAtIndex(s.keyPos)
	offset, _ := s.reader.offsets.Get(s.keyPos)
	s.seriesOffsets = encoding.NewFixedOffsetDecoder()
	_, _ = s.seriesOffsets.Unmarshal(s.reader.buf[offset:])
	s.keyPos++
}

// scan scans the data then merges the series ids into target series ids
func (s *tagInvertedScanner) scan(highKey, lowTagValueID uint16, targetSeriesIDs *roaring.Bitmap) error {
	if s.highKey < highKey {
		if s.keyPos >= len(s.highKeys) {
			// current tag inverted no data can read
			return nil
		}
		s.nextContainer()
	}
	if highKey != s.highKey {
		// high key not match, return it
		return nil
	}
	// find data by low tag value id
	if s.container.Contains(lowTagValueID) {
		lowIdx := s.container.Rank(lowTagValueID)
		seriesPos, _ := s.seriesOffsets.Get(lowIdx - 1)

		// unmarshal series ids
		seriesIDs := roaring.New()
		if err := encoding.BitmapUnmarshal(seriesIDs, s.reader.buf[seriesPos:]); err != nil {
			return err
		}
		// merge the data into target series ids
		targetSeriesIDs.Or(seriesIDs)
	}
	return nil
}

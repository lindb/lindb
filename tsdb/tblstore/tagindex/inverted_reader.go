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
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./inverted_reader.go -destination=./inverted_reader_mock.go -package tagindex

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
func (r *inverterReader) loadSeriesIDs(tagKeyID uint32,
	fn func(indexReader *tagInvertedReader) (*roaring.Bitmap, error)) (*roaring.Bitmap, error) {
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

	highOffsets := encoding.NewFixedOffsetDecoder()
	lowOffsets := encoding.NewFixedOffsetDecoder()

	if _, err := highOffsets.Unmarshal(r.buf[r.baseReader.offsetsAt:]); err != nil {
		return nil, err
	}
	entries := r.buf[:r.baseReader.tagValueBitmapAt]

	for idx, highKey := range highKeys {
		loadLowContainer := finalTagValueIDs.GetContainerAtIndex(idx)
		lowContainerIdx := r.keys.GetContainerIndex(highKey)
		lowContainer := r.keys.GetContainerAtIndex(lowContainerIdx)

		tagValueBucket, err := highOffsets.GetBlock(lowContainerIdx, entries)
		if err != nil {
			return nil, err
		}
		lowKeyOffsetsBlockLen, uVariantEncodingLen := stream.UvarintLittleEndian(tagValueBucket)
		lowKeyOffsetsAt := len(tagValueBucket) - int(lowKeyOffsetsBlockLen) - uVariantEncodingLen
		if uVariantEncodingLen <= 0 || lowKeyOffsetsAt <= 0 || lowKeyOffsetsAt >= len(tagValueBucket) {
			return nil, fmt.Errorf("read lowkey offsets error")
		}
		if _, err = lowOffsets.Unmarshal(tagValueBucket[lowKeyOffsetsAt:]); err != nil {
			return nil, err
		}
		level3Block := tagValueBucket[:lowKeyOffsetsAt]

		it := loadLowContainer.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// get the index of low tag value id in container
			lowIdx := lowContainer.Rank(lowTagValueID)

			block, err := lowOffsets.GetBlock(lowIdx-1, level3Block)
			if err != nil {
				continue
			}
			// unmarshal series ids
			seriesIDs := roaring.New()
			if err := encoding.BitmapUnmarshal(seriesIDs, block); err != nil {
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
	lowKeyOffsets *encoding.FixedOffsetDecoder
	entries       []byte
	level3Block   []byte

	highKeys         []uint16
	highKey          uint16
	highContainerIdx int
}

// newTagInvertedScanner creates a tag inverted index scanner
func newTagInvertedScanner(reader *tagInvertedReader) (*tagInvertedScanner, error) {
	s := &tagInvertedScanner{
		reader:        reader,
		highKeys:      reader.keys.GetHighKeys(),
		lowKeyOffsets: encoding.NewFixedOffsetDecoder(),
		entries:       reader.buf[:reader.tagValueBitmapAt],
	}
	if len(s.highKeys) == 0 {
		return nil, fmt.Errorf("tagValue bitmap is empty")
	}
	if err := s.nextContainer(); err != nil {
		return nil, err
	}
	return s, nil
}

// nextContainer goes next container context for scanner
func (s *tagInvertedScanner) nextContainer() error {
	s.highKey = s.highKeys[s.highContainerIdx]
	s.container = s.reader.keys.GetContainerAtIndex(s.highContainerIdx)

	tagValueBucket, err := s.reader.offsets.GetBlock(s.highContainerIdx, s.entries)
	if err != nil {
		return err
	}
	lowKeyOffsetsBlockLen, uVariantEncodingLen := stream.UvarintLittleEndian(tagValueBucket)
	lowKeyOffsetsAt := len(tagValueBucket) - int(lowKeyOffsetsBlockLen) - uVariantEncodingLen
	if uVariantEncodingLen <= 0 || lowKeyOffsetsAt <= 0 || lowKeyOffsetsAt >= len(tagValueBucket) {
		return fmt.Errorf("read lowkey offsets error")
	}
	if _, err = s.lowKeyOffsets.Unmarshal(tagValueBucket[lowKeyOffsetsAt:]); err != nil {
		return err
	}
	s.level3Block = tagValueBucket[:lowKeyOffsetsAt]
	s.highContainerIdx++
	return nil
}

// scan scans the data then merges the series ids into target series ids
func (s *tagInvertedScanner) scan(highKey, lowTagValueID uint16, targetSeriesIDs *roaring.Bitmap) error {
	if s.highKey < highKey {
		if s.highContainerIdx >= len(s.highKeys) {
			// current tag inverted no data can read
			return nil
		}
		if err := s.nextContainer(); err != nil {
			return err
		}
	}
	if highKey != s.highKey {
		// high key not match, return it
		return nil
	}
	// find data by low tag value id
	if s.container.Contains(lowTagValueID) {
		lowIdx := s.container.Rank(lowTagValueID)
		bitmapBlock, err := s.lowKeyOffsets.GetBlock(lowIdx-1, s.level3Block)
		if err != nil {
			return err
		}
		// unmarshal series ids
		seriesIDs := roaring.New()
		if err := encoding.BitmapUnmarshal(seriesIDs, bitmapBlock); err != nil {
			return err
		}
		// merge the data into target series ids
		targetSeriesIDs.Or(seriesIDs)
	}
	return nil
}

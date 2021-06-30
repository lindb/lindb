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

package tagkeymeta

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/lindb/roaring"
)

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package tagkeymeta

// Reader reads tag value data from tag-index-table
type Reader interface {
	// GetTagValueSeq returns the auto sequence of tag value under the tag key,
	// if not exist return constants.ErrTagValueSeqNotFound
	GetTagValueSeq(tagKeyID uint32) (tagValueSeq uint32, err error)

	// GetTagValueID returns the tag value id for spec metric's tag key id,
	// if not exist return constants.ErrTagValueIDNotFound
	GetTagValueID(tagKeyID uint32, tagValue string) (tagValueID uint32, err error)

	// GetTagValueIDsForTagKeyID get tag value ids for spec metric's tag key id
	GetTagValueIDsForTagKeyID(tagKeyID uint32) (tagValueIDs *roaring.Bitmap, err error)

	// FindValueIDsByExprForTagKeyID finds tag values ids by tag filter expr and tag key id
	FindValueIDsByExprForTagKeyID(tagKeyID uint32, expr stmt.TagFilter) (tagValueIDs *roaring.Bitmap, err error)

	// SuggestTagValues finds tag values by prefix search
	SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string

	// WalkTagValues walks each tag value and tag value id via fn.
	// If fn returns false, the iteration is stopped.
	// The values are the raw byte slices and not the converted types.
	WalkTagValues(
		tagKeyID uint32,
		tagValuePrefix string,
		fn func(tagValue []byte, tagValueID uint32) bool,
	) error

	// CollectTagValues collects the tag values by tag value ids,
	CollectTagValues(tagKeyID uint32, tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) error
}

// tagReader implements TagReader
type tagReader struct {
	readers []table.Reader
}

// NewReader returns a new TagReader
func NewReader(readers []table.Reader) Reader {
	return &tagReader{readers: readers}
}

// GetTagValueSeq returns the auto sequence of tag value under the tag key,
// if not exist return constants.ErrTagValueSeqNotFound
// kv store returns the table.readers in order,
// so the max sequence will be stored in the first table.reader that is tag key store.
func (r *tagReader) GetTagValueSeq(tagKeyID uint32) (tagValueSeq uint32, err error) {
	for _, reader := range r.readers {
		tagKeyMetaBlock, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		//FIXME stone1100 opt need cache entry set
		meta, err := newTagKeyMeta(tagKeyMetaBlock)
		if err != nil {
			return 0, fmt.Errorf("%w, %s: ", constants.ErrTagValueSeqNotFound, err)
		}
		return meta.TagValueIDSeq(), nil
	}
	return 0, fmt.Errorf("%w, tagKeyID:%d", constants.ErrTagValueSeqNotFound, tagKeyID)
}

// GetTagValueID returns the tag value id for spec metric's tag key id,
// if not exist return constants.ErrTagValueIDNotFound
func (r *tagReader) GetTagValueID(tagID uint32, tagValue string) (tagValueID uint32, err error) {
	for _, reader := range r.readers {
		tagKeyMetaBlock, ok := reader.Get(tagID)
		if !ok {
			continue
		}
		meta, err := newTagKeyMeta(tagKeyMetaBlock)
		if err != nil {
			return 0, fmt.Errorf("%w, tagValue: %s with error: %s",
				constants.ErrTagValueIDNotFound, tagValue, err)
		}
		tagValueIDs := meta.FindTagValueID(tagValue)
		if len(tagValueIDs) == 0 {
			continue
		}
		return tagValueIDs[0], nil
	}
	return 0, fmt.Errorf("%w, tagValue: %s", constants.ErrTagValueIDNotFound, tagValue)
}

// FindValueIDsByExprForTagKeyID finds tag values ids by tag filter expr and tag key id
func (r *tagReader) FindValueIDsByExprForTagKeyID(tagID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	tagKeyMetas := r.filterTagKeyMetas(tagID)
	if len(tagKeyMetas) == 0 {
		return nil, fmt.Errorf("%w, tagID: %d", constants.ErrTagKeyMetaNotFound, tagID)
	}
	tagValueIDs := roaring.New()
	for _, tagKeyMeta := range tagKeyMetas {
		switch expression := expr.(type) {
		case *stmt.EqualsExpr:
			tagValueIDs.AddMany(tagKeyMeta.FindTagValueID(expression.Value))
		case *stmt.InExpr:
			tagValueIDs.AddMany(tagKeyMeta.FindTagValueIDs(expression.Values))
		case *stmt.LikeExpr:
			tagValueIDs.AddMany(tagKeyMeta.FindTagValueIDsByLike(expression.Value))
		case *stmt.RegexExpr:
			tagValueIDs.AddMany(tagKeyMeta.FindTagValueIDsByRegex(expression.Regexp))
		default:
			return nil, fmt.Errorf("%w, unsupported expr, tagID: %d",
				constants.ErrTagKeyMetaNotFound, tagID)
		}
	}
	if tagValueIDs.IsEmpty() {
		return nil, fmt.Errorf("%w, tagID: %d", constants.ErrTagValueIDNotFound, tagID)
	}
	return tagValueIDs, nil
}

// GetTagValueIDsForTagKeyID get tag value ids for spec metric's tag key id
func (r *tagReader) GetTagValueIDsForTagKeyID(tagID uint32) (*roaring.Bitmap, error) {
	tagKeyMetas := r.filterTagKeyMetas(tagID)
	if len(tagKeyMetas) == 0 {
		return nil, fmt.Errorf("%w, tagID: %d", constants.ErrTagKeyMetaNotFound, tagID)
	}
	return tagKeyMetas.GetTagValueIDs()
}

// filterTagKeyMetas filters the tag-key-metas by tag key id
func (r *tagReader) filterTagKeyMetas(tagID uint32) (metas TagKeyMetas) {
	for _, reader := range r.readers {
		tagKeyMetaBlock, ok := reader.Get(tagID)
		if !ok {
			continue
		}
		tagKeyMeta, err := newTagKeyMeta(tagKeyMetaBlock)
		if err != nil {
			continue
		}
		metas = append(metas, tagKeyMeta)
	}
	return
}

// SuggestTagValues finds tagValues by prefix search
func (r *tagReader) SuggestTagValues(
	tagKeyID uint32,
	tagValuePrefix string,
	limit int,
) (
	tagValues []string,
) {
	if limit > constants.MaxSuggestions {
		limit = constants.MaxSuggestions
	}
	for _, reader := range r.readers {
		tagKeyMetaBlock, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		tagKeyMeta, err := newTagKeyMeta(tagKeyMetaBlock)
		if err != nil {
			continue
		}
		itr, err := tagKeyMeta.PrefixIterator(strutil.String2ByteSlice(tagValuePrefix))
		if err != nil {
			continue
		}
		for itr.Valid() {
			tagValues = append(tagValues, string(itr.Key()))
			if len(tagValues) >= limit {
				return tagValues
			}
			itr.Next()
		}
	}
	return tagValues
}

// WalkTagValues walks each tag value and tag value id via fn.
// If fn returns false, the iteration is stopped.
// The values are the raw byte slices and not the converted types.
func (r *tagReader) WalkTagValues(
	tagKeyID uint32,
	tagValuePrefix string,
	fn func(tagValue []byte, tagValueID uint32) bool,
) error {
	for _, reader := range r.readers {
		tagKeyMetaBlock, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		tagKeyMeta, err := newTagKeyMeta(tagKeyMetaBlock)
		if err != nil {
			continue
		}
		itr, err := tagKeyMeta.PrefixIterator(strutil.String2ByteSlice(tagValuePrefix))
		if err != nil {
			continue
		}
		for itr.Valid() {
			tagValue, tagValueID := itr.Key(), encoding.ByteSlice2Uint32(itr.Value())
			if fn != nil && !fn(tagValue, tagValueID) {
				return nil
			}
			itr.Next()
		}
	}
	return nil
}

// CollectTagValues collects the tag values by tag value ids
func (r *tagReader) CollectTagValues(
	tagKeyID uint32,
	tagValueIDs *roaring.Bitmap,
	tagValues map[uint32]string,
) error {
	for _, reader := range r.readers {
		if tagValueIDs.IsEmpty() {
			return nil
		}
		tagKeyMetaBlock, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		tagKeyMeta, err := newTagKeyMeta(tagKeyMetaBlock)
		if err != nil {
			continue
		}
		if err := tagKeyMeta.CollectTagValues(tagValueIDs, tagValues); err != nil {
			return err
		}
	}
	return nil
}

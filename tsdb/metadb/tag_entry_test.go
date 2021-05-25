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

package metadb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestTagEntry_genTagValueID(t *testing.T) {
	tagEntry := newTagEntry(0)
	assert.Equal(t, uint32(1), tagEntry.genTagValueID())
	assert.Equal(t, uint32(1), tagEntry.getTagValueIDSeq())

	tagEntry = newTagEntry(100)
	assert.Equal(t, uint32(101), tagEntry.genTagValueID())
	assert.Equal(t, uint32(101), tagEntry.getTagValueIDSeq())
}

func TestTagEntry_getTagValueID(t *testing.T) {
	tagIndex := prepareTagEntry()
	id, ok := tagIndex.getTagValueID("abc")
	assert.True(t, ok)
	assert.Equal(t, uint32(2), id)
	id, ok = tagIndex.getTagValueID("abcddd")
	assert.False(t, ok)
	assert.Equal(t, uint32(0), id)
	values := tagIndex.getTagValues()
	assert.Len(t, values, 8)
}

func TestTagEntry_findSeriesIDsByEquals(t *testing.T) {
	tagIndex := prepareTagEntry()
	// tag-value not exist
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "alpha"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(4), tagIndex.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "c"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(5), tagIndex.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "bc"}))
}

func TestTagEntry_findSeriesIDsByLike(t *testing.T) {
	tagIndex := prepareTagEntry()

	// tag-value is empty
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(2, 5, 8), tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "*bc*"}))
	// tag-value not exist
	assert.Equal(t, roaring.New(), tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "zz*"}))
	// tag-value is *
	assert.Equal(t, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8), tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "*"}))
	// tag-value is "abc" ==> equals
	assert.Equal(t, roaring.BitmapOf(2), tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "abc"}))
	// tag-value is "*cd"
	assert.Equal(t, roaring.BitmapOf(8), tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "*cd"}))
	// tag-value is "b*"
	assert.Equal(t, roaring.BitmapOf(3, 5, 6, 7, 8), tagIndex.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "b*"}))
}

func TestTagEntry_findSeriesIDsByIn(t *testing.T) {
	tagIndex := prepareTagEntry()
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(3, 5, 8), tagIndex.findSeriesIDsByExpr(&stmt.InExpr{Key: "host", Values: []string{"b", "bc", "bcd", "ahi"}}))
}

func TestTagEntry_findSeriesIDsByExpr_not_tagFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagIndex := prepareTagEntry()
	tagFilter := stmt.NewMockTagFilter(ctrl)
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(tagFilter))
}

func TestTagEntry_findSeriesIDsByRegex(t *testing.T) {
	tagIndex := prepareTagEntry()
	// pattern not match
	assert.Equal(t, roaring.New(), tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "bbbbbbbbbbb"}))
	// pattern error
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "b.32*++++\n"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(6, 7), tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `b2[0-9]+`}))
	// literal prefix:22 not exist
	assert.Equal(t, roaring.New(), tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `22+`}))
}

func TestTagEntry_collectTagValues(t *testing.T) {
	tagIndex := prepareTagEntry()
	tagValueIDs := roaring.BitmapOf(1, 2, 3, 100)
	result := make(map[uint32]string)
	// case 1: collect tag value
	tagIndex.collectTagValues(tagValueIDs, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "a", result[1])
	assert.Equal(t, "abc", result[2])
	assert.Equal(t, "b", result[3])
	assert.EqualValues(t, tagValueIDs.ToArray(), roaring.BitmapOf(100).ToArray())
	tagIndex.collectTagValues(tagValueIDs, result)
	assert.Len(t, result, 3)
	assert.EqualValues(t, tagValueIDs.ToArray(), roaring.BitmapOf(100).ToArray())
	// case 2: collect tag value ids empty
	result = make(map[uint32]string)
	tagValueIDs = roaring.BitmapOf(2, 3)
	tagIndex.collectTagValues(tagValueIDs, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "abc", result[2])
	assert.Equal(t, "b", result[3])
}

func prepareTagEntry() TagEntry {
	tagIndex := newTagEntry(0)
	tagIndex.addTagValue("a", 1)
	tagIndex.addTagValue("abc", 2)
	tagIndex.addTagValue("b", 3)
	tagIndex.addTagValue("c", 4)
	tagIndex.addTagValue("bc", 5)
	tagIndex.addTagValue("b21", 6)
	tagIndex.addTagValue("b22", 7)
	tagIndex.addTagValue("bcd", 8)
	return tagIndex
}

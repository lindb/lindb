package indexdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestTagIndex_buildInvertedIndex(t *testing.T) {
	index := newTagIndex()
	values := index.getValues()
	assert.Len(t, values, 0)
	index.buildInvertedIndex("tagValue1", 1)
	index.buildInvertedIndex("tagValue1", 2)
	index.buildInvertedIndex("tagValue2", 1)
	index.buildInvertedIndex("tagValue2", 2)
	index.buildInvertedIndex("tagValue2", 3)
	values = index.getValues()
	assert.Equal(t, roaring.BitmapOf(1, 2), values["tagValue1"])
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), values["tagValue2"])
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), index.getAllSeriesIDs())
}

func TestTagIndex_findSeriesIDsByEquals(t *testing.T) {
	tagIndex := prepareTagIdx()
	// tag-value not exist
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "alpha"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(4), tagIndex.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "c"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(5), tagIndex.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "bc"}))
}

func TestTagIndex_findSeriesIDsByLike(t *testing.T) {
	tagIndex := prepareTagIdx()

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

func TestTagIndex_findSeriesIDsByIn(t *testing.T) {
	tagIndex := prepareTagIdx()
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(3, 5, 8), tagIndex.findSeriesIDsByExpr(&stmt.InExpr{Key: "host", Values: []string{"b", "bc", "bcd", "ahi"}}))
}

func TestTagIndex_findSeriesIDsByExpr_not_tagFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagIndex := prepareTagIdx()
	tagFilter := stmt.NewMockTagFilter(ctrl)
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(tagFilter))
}

func TestTagIndex_findSeriesIDsByRegex(t *testing.T) {
	tagIndex := prepareTagIdx()
	// pattern not match
	assert.Equal(t, roaring.New(), tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "bbbbbbbbbbb"}))
	// pattern error
	assert.Nil(t, tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "b.32*++++\n"}))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(6, 7), tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `b2[0-9]+`}))
	// literal prefix:22 not exist
	assert.Equal(t, roaring.New(), tagIndex.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `22+`}))
}

func TestTagIndex_suggestTagValues(t *testing.T) {
	tagIndex := prepareTagIdx()
	assert.Nil(t, tagIndex.suggestTagValues("de", 0))
	assert.Zero(t, tagIndex.suggestTagValues("de", 10))
	assert.Len(t, tagIndex.suggestTagValues("b", 10), 5)
	assert.Len(t, tagIndex.suggestTagValues("b", 1), 1)
}

func prepareTagIdx() TagIndex {
	tagIndex := newTagIndex()
	tagIndex.buildInvertedIndex("a", 1)
	tagIndex.buildInvertedIndex("abc", 2)
	tagIndex.buildInvertedIndex("b", 3)
	tagIndex.buildInvertedIndex("c", 4)
	tagIndex.buildInvertedIndex("bc", 5)
	tagIndex.buildInvertedIndex("b21", 6)
	tagIndex.buildInvertedIndex("b22", 7)
	tagIndex.buildInvertedIndex("bcd", 8)
	return tagIndex
}

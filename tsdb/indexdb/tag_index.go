package indexdb

import (
	"regexp"
	"strings"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/sql/stmt"
)

// TagIndex represents the tag inverted index
type TagIndex interface {
	// buildInvertedIndex builds inverted index for tag value
	buildInvertedIndex(tagValue string, seriesID uint32)
	// findSeriesIDsByExpr finds series ids by tag filter expr
	findSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap
	// getValues returns the all tag values and series ids
	getValues() map[string]*roaring.Bitmap
	// getAllSeriesIDs returns all series ids
	getAllSeriesIDs() *roaring.Bitmap
	// suggestTagValues returns tagValues by prefix-search
	suggestTagValues(tagValuePrefix string, limit int) (tagValuesList []string)
}

// tagIndex is a inverted mapping relation of tag-value and seriesID group.
type tagIndex struct {
	values map[string]*roaring.Bitmap
}

// newTagKVEntrySet returns a new tagKVEntrySet
func newTagIndex() TagIndex {
	return &tagIndex{
		values: make(map[string]*roaring.Bitmap),
	}
}

func (index *tagIndex) buildInvertedIndex(tagValue string, seriesID uint32) {
	seriesIDs, ok := index.values[tagValue]
	if !ok {
		// create new series ids for new tag value
		seriesIDs = roaring.NewBitmap()
		index.values[tagValue] = seriesIDs
	}
	seriesIDs.Add(seriesID)
}

// findSeriesIDsByExpr finds series ids by tag filter expr
func (index *tagIndex) findSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap {
	switch expression := expr.(type) {
	case *stmt.EqualsExpr:
		return index.findSeriesIDsByEqual(expression.Value)
	case *stmt.InExpr:
		return index.findSeriesIDsByIn(expression)
	case *stmt.LikeExpr:
		return index.findSeriesIDsByLike(expression)
	case *stmt.RegexExpr:
		return index.findSeriesIDsByRegex(expression)
	}
	indexLogger.Warn("expr type is not tag filter when find series ids by expr")
	return nil
}

// findSeriesIDsByEqual finds series ids by tag value - equal
func (index *tagIndex) findSeriesIDsByEqual(value string) *roaring.Bitmap {
	bitmap, ok := index.values[value]
	if !ok {
		return nil
	}
	return bitmap.Clone()
}

// findSeriesIDsByIn finds series ids by tag value - in
func (index *tagIndex) findSeriesIDsByIn(expr *stmt.InExpr) *roaring.Bitmap {
	union := roaring.New()
	for _, value := range expr.Values {
		bitmap, ok := index.values[value]
		if !ok {
			continue
		}
		union.Or(bitmap)
	}
	return union
}

// findSeriesIDsByLike finds series ids by tag value - like
// case 1: value is empty, return nil
// case 2: value is "*", return all series ids
// case 3: value is "*xxx*", do contains
// case 4: value is "*xxx", do suffix
// case 5: value is "xxx*", do prefix
// case 6: value is "xxx", do equal
func (index *tagIndex) findSeriesIDsByLike(expr *stmt.LikeExpr) *roaring.Bitmap {
	likeTo := expr.Value
	length := len(likeTo)
	if length == 0 {
		return nil
	}
	if likeTo == "*" {
		return index.getAllSeriesIDs()
	}
	union := roaring.New()
	prefix := strings.HasPrefix(likeTo, "*")
	suffix := strings.HasSuffix(likeTo, "*")
	switch {
	case prefix && suffix:
		like := likeTo[1 : length-1]
		for value, bitmap := range index.values {
			if strings.Contains(value, like) {
				union.Or(bitmap)
			}
		}
	case prefix:
		like := likeTo[1:]
		for value, bitmap := range index.values {
			if strings.HasSuffix(value, like) {
				union.Or(bitmap)
			}
		}
	case suffix:
		like := likeTo[:length-1]
		for value, bitmap := range index.values {
			if strings.HasPrefix(value, like) {
				union.Or(bitmap)
			}
		}
	default:
		// like == equal
		return index.findSeriesIDsByEqual(likeTo)
	}
	return union
}

// findSeriesIDsByRegex finds series ids by tag value - regex
func (index *tagIndex) findSeriesIDsByRegex(expr *stmt.RegexExpr) *roaring.Bitmap {
	pattern, err := regexp.Compile(expr.Regexp)
	if err != nil {
		return nil
	}
	// the regex pattern is regarded as a prefix string + pattern
	literalPrefix, _ := pattern.LiteralPrefix()
	union := roaring.New()
	for value, bitmap := range index.values {
		if !strings.HasPrefix(value, literalPrefix) {
			continue
		}
		if pattern.MatchString(value) {
			union.Or(bitmap)
		}
	}
	return union
}

// getAllSeriesIDs returns all series ids
func (index *tagIndex) getAllSeriesIDs() *roaring.Bitmap {
	union := roaring.New()
	for _, bitMap := range index.values {
		union.Or(bitMap)
	}
	return union
}

// getValues returns the all tag values and series ids
func (index *tagIndex) getValues() map[string]*roaring.Bitmap {
	return index.values
}

// suggestTagValues returns tagValues by prefix-search
func (index *tagIndex) suggestTagValues(tagValuePrefix string, limit int) (tagValuesList []string) {
	for tagValue := range index.values {
		if strings.HasPrefix(tagValue, tagValuePrefix) {
			if len(tagValuesList) >= limit {
				return
			}
			tagValuesList = append(tagValuesList, tagValue)
		}
	}
	return
}

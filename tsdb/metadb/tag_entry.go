package metadb

import (
	"regexp"
	"strings"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/sql/stmt"
)

// TagEntry represents the tag value=>id under tag key
type TagEntry interface {
	// genTagValueID generates a new tag value id for new tag value, start with 1
	genTagValueID() uint32
	// getTagValueIDSeq returns the current tag value id sequence
	getTagValueIDSeq() uint32
	// addTagValue adds tag value=>id mapping
	addTagValue(tagValue string, tagValueID uint32)
	// findSeriesIDsByExpr finds tag value ids by tag filter expr
	findSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap
	// getTagValueID gets the tag value id by tag value under the tag key
	getTagValueID(tagValue string) (uint32, bool)
	// getTagValueIDs returns all tag value ids under the tag key
	getTagValueIDs() *roaring.Bitmap
	// getTagValues returns the all tag values
	getTagValues() map[string]uint32
	// collectTagValues collects the tag values by tag value ids,
	collectTagValues(tagValueIDs *roaring.Bitmap, tagValues map[uint32]string)
}

// tagEntry implements TagEntry interface
type tagEntry struct {
	tagValueSeq atomic.Uint32
	tagValues   map[string]uint32
}

// newTagEntry creates tag entry with tag value id auto sequence
func newTagEntry(tagValueSeq uint32) TagEntry {
	t := &tagEntry{
		tagValues: make(map[string]uint32),
	}
	t.tagValueSeq.Store(tagValueSeq)
	return t
}

// genTagValueID generates a new tag value id for new tag value, start with 1
func (t *tagEntry) genTagValueID() uint32 {
	return t.tagValueSeq.Inc()
}

// getTagValueIDSeq returns the current tag value id sequence
func (t *tagEntry) getTagValueIDSeq() uint32 {
	return t.tagValueSeq.Load()
}

// addTagValue adds tag value=>id mapping
func (t *tagEntry) addTagValue(tagValue string, tagValueID uint32) {
	t.tagValues[tagValue] = tagValueID
}

// getTagValueIDs returns all tag value ids under the tag key
func (t *tagEntry) getTagValueID(tagValue string) (uint32, bool) {
	tagValueID, ok := t.tagValues[tagValue]
	return tagValueID, ok
}

// getTagValues returns the all tag values
func (t *tagEntry) getTagValues() map[string]uint32 {
	return t.tagValues
}

// getTagValueIDs returns all tag value ids under the tag key
func (t *tagEntry) getTagValueIDs() *roaring.Bitmap {
	tagValueIDs := roaring.New()
	for _, tagValueID := range t.tagValues {
		tagValueIDs.Add(tagValueID)
	}
	return tagValueIDs
}

// findSeriesIDsByExpr finds tag value ids by tag filter expr
func (t *tagEntry) findSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap {
	switch expression := expr.(type) {
	case *stmt.EqualsExpr:
		return t.findSeriesIDsByEqual(expression.Value)
	case *stmt.InExpr:
		return t.findSeriesIDsByIn(expression)
	case *stmt.LikeExpr:
		return t.findSeriesIDsByLike(expression)
	case *stmt.RegexExpr:
		return t.findSeriesIDsByRegex(expression)
	}
	metaLogger.Warn("expr type is not tag filter when find tag value ids by expr")
	return nil
}

// findSeriesIDsByEqual finds tag value ids by tag value - equal
func (t *tagEntry) findSeriesIDsByEqual(value string) *roaring.Bitmap {
	tagValueID, ok := t.tagValues[value]
	if !ok {
		return nil
	}
	return roaring.BitmapOf(tagValueID)
}

// findSeriesIDsByIn finds series ids by tag value - in
func (t *tagEntry) findSeriesIDsByIn(expr *stmt.InExpr) *roaring.Bitmap {
	union := roaring.New()
	for _, value := range expr.Values {
		tagValueID, ok := t.tagValues[value]
		if !ok {
			continue
		}
		union.Add(tagValueID)
	}
	return union
}

// findSeriesIDsByLike finds tag values ids by tag value - like
// case 1: value is empty, return nil
// case 2: value is "*", return all tag value ids
// case 3: value is "*xxx*", do contains
// case 4: value is "*xxx", do suffix
// case 5: value is "xxx*", do prefix
// case 6: value is "xxx", do equal
func (t *tagEntry) findSeriesIDsByLike(expr *stmt.LikeExpr) *roaring.Bitmap {
	likeTo := expr.Value
	length := len(likeTo)
	if length == 0 {
		return nil
	}
	if likeTo == "*" {
		return t.getTagValueIDs()
	}
	result := roaring.New()
	prefix := strings.HasPrefix(likeTo, "*")
	suffix := strings.HasSuffix(likeTo, "*")
	switch {
	case prefix && suffix:
		like := likeTo[1 : length-1]
		for value, tagValueID := range t.tagValues {
			if strings.Contains(value, like) {
				result.Add(tagValueID)
			}
		}
	case prefix:
		like := likeTo[1:]
		for value, tagValueID := range t.tagValues {
			if strings.HasSuffix(value, like) {
				result.Add(tagValueID)
			}
		}
	case suffix:
		like := likeTo[:length-1]
		for value, tagValueID := range t.tagValues {
			if strings.HasPrefix(value, like) {
				result.Add(tagValueID)
			}
		}
	default:
		// like == equal
		return t.findSeriesIDsByEqual(likeTo)
	}
	return result
}

// findSeriesIDsByRegex finds tag value ids by tag value - regex
func (t *tagEntry) findSeriesIDsByRegex(expr *stmt.RegexExpr) *roaring.Bitmap {
	pattern, err := regexp.Compile(expr.Regexp)
	if err != nil {
		return nil
	}
	// the regex pattern is regarded as a prefix string + pattern
	literalPrefix, _ := pattern.LiteralPrefix()
	result := roaring.New()
	for value, tagValueID := range t.tagValues {
		if !strings.HasPrefix(value, literalPrefix) {
			continue
		}
		if pattern.MatchString(value) {
			result.Add(tagValueID)
		}
	}
	return result
}

// collectTagValues collects the tag values by tag value ids,
func (t *tagEntry) collectTagValues(tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) {
	for value, tagValueID := range t.tagValues {
		if tagValueIDs.IsEmpty() {
			break
		}
		if tagValueIDs.Contains(tagValueID) {
			tagValueIDs.Remove(tagValueID)
			tagValues[tagValueID] = value
		}
	}
}

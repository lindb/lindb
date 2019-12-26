package query

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series/tag"
)

// SeriesID2Tags represents the tag values of series id
type SeriesID2Tags struct {
	tagValues       []string
	size, tagsCount int
}

// NewSeriesID2Tags creates a SeriesID2Tags
func NewSeriesID2Tags(tagsCount int) *SeriesID2Tags {
	return &SeriesID2Tags{
		tagValues: make([]string, tagsCount),
		tagsCount: tagsCount,
	}
}

// Size returns the size of tag values
func (entry *SeriesID2Tags) Size() int {
	return entry.size
}

// TagValues returns the tag values
func (entry *SeriesID2Tags) TagValues() []string {
	return entry.tagValues
}

// AddTagValue adds the tag value for series id
func (entry *SeriesID2Tags) AddTagValue(tagValue string) {
	if entry.size < entry.tagsCount {
		entry.tagValues[entry.size] = tagValue
		entry.size++
	}
}

// TagValuesEntrySet represents the tag values and series ids mapping for a tag key
type TagValuesEntrySet struct {
	values map[string]*roaring.Bitmap
}

// NewTagValuesEntrySet creates a TagValuesEntrySet
func NewTagValuesEntrySet() *TagValuesEntrySet {
	return &TagValuesEntrySet{values: make(map[string]*roaring.Bitmap)}
}

// Values returns the tag values data
func (tes *TagValuesEntrySet) Values() map[string]*roaring.Bitmap {
	return tes.values
}

// SetTagValues sets the tag values data
func (tes *TagValuesEntrySet) SetTagValues(values map[string]*roaring.Bitmap) {
	tes.values = values
}

// AddTagValue adds tag value and series ids
func (tes *TagValuesEntrySet) AddTagValue(tagValue string, seriesIDs *roaring.Bitmap) {
	oldSeriesIDs, ok := tes.values[tagValue]
	if ok {
		oldSeriesIDs.Or(seriesIDs)
	} else {
		tes.values[tagValue] = seriesIDs
	}
}

// GroupingContext represents the context of group by query for tag keys
type GroupingContext struct {
	tagValuesEntrySets []*TagValuesEntrySet
}

// NewGroupContext creates a GroupingContext
func NewGroupContext(tagKeyCount int) *GroupingContext {
	return &GroupingContext{
		tagValuesEntrySets: make([]*TagValuesEntrySet, tagKeyCount),
	}
}

// SetTagValuesEntrySet sets the tag values entry set for group by tag keys
func (g *GroupingContext) SetTagValuesEntrySet(idx int, tagValuesEntrySet *TagValuesEntrySet) {
	g.tagValuesEntrySets[idx] = tagValuesEntrySet
}

// Len returns the group by tag key's length
func (g *GroupingContext) Len() int {
	return len(g.tagValuesEntrySets)
}

// BuildGroup builds the grouped series ids by the high key of series id
// and the container includes low keys of series id
func (g *GroupingContext) BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16 {
	groupTagKeysCount := len(g.tagValuesEntrySets)
	if groupTagKeysCount == 1 {
		return g.buildForSingleTagKey(highKey, container)
	}

	// new seriesIDs2Tags array based on range of max ~ min
	seriesIDs2Tags := g.buildSeriesIDs2Tags(highKey, container)

	// finds group tags => series IDs, and builds result
	it := container.PeekableIterator()
	min := container.Minimum()
	result := make(map[string][]uint16)
	for it.HasNext() {
		lowKey := it.Next()
		idx := lowKey - min
		seriesID2Tags := seriesIDs2Tags[idx]
		if seriesID2Tags != nil && seriesID2Tags.Size() == groupTagKeysCount {
			tagValuesStr := tag.ConcatTagValues(seriesID2Tags.TagValues())
			values, ok := result[tagValuesStr]
			if !ok {
				result[tagValuesStr] = []uint16{lowKey}
			} else {
				result[tagValuesStr] = append(values, lowKey)
			}
		}
	}
	return result
}

// buildSeriesIDs2Tags builds for multi group by keys
func (g *GroupingContext) buildSeriesIDs2Tags(highKey uint16, container roaring.Container) []*SeriesID2Tags {
	groupTagKeysCount := len(g.tagValuesEntrySets)
	// new seriesIDs2Tags array based on range of max ~ min
	min := container.Minimum()
	max := container.Maximum()
	seriesIDs2Tags := make([]*SeriesID2Tags, int(max-min)+1)

	// builds seriesIDs => tags mapping, using counting sort
	// https://en.wikipedia.org/wiki/Counting_sort
	for _, tagKV := range g.tagValuesEntrySets {
		for tagValue, lowSeriesIDs := range tagKV.Values() {
			lowContainer := lowSeriesIDs.GetContainer(highKey)
			if lowContainer != nil {
				matchContainer := lowContainer.And(container)
				it := matchContainer.PeekableIterator()
				for it.HasNext() {
					idx := it.Next() - min // index = lowKey - min
					if seriesIDs2Tags[idx] == nil {
						seriesIDs2Tags[idx] = NewSeriesID2Tags(groupTagKeysCount)
					}
					seriesIDs2Tags[idx].AddTagValue(tagValue)
				}
			}
		}
	}
	return seriesIDs2Tags
}

// buildForSingleTagKey builds for single group by tag key
func (g *GroupingContext) buildForSingleTagKey(highKey uint16, container roaring.Container) map[string][]uint16 {
	result := make(map[string][]uint16)
	for _, tagKV := range g.tagValuesEntrySets {
		for tagValue, lowSeriesIDs := range tagKV.Values() {
			lowContainer := lowSeriesIDs.GetContainer(highKey)
			if lowContainer != nil {
				matchContainer := lowContainer.And(container)
				result[tagValue] = matchContainer.ToArray()
			}
		}
	}
	return result
}

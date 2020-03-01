package query

import (
	"encoding/binary"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series"
)

// GroupingContext represents the context of group by query for tag keys
// builds tags => series ids mapping, using such as counting sort
// https://en.wikipedia.org/wiki/Counting_sort
type GroupingContext struct {
	tagKeys     []uint32
	scanners    map[uint32][]series.GroupingScanner
	tagValueIDs []*roaring.Bitmap // collect tag value ids for each group by tag key
}

// NewGroupContext creates a GroupingContext
func NewGroupContext(tagKeys []uint32, scanners map[uint32][]series.GroupingScanner) series.GroupingContext {
	return &GroupingContext{
		tagKeys:     tagKeys,
		scanners:    scanners,
		tagValueIDs: make([]*roaring.Bitmap, len(tagKeys)),
	}
}

// GetGroupByTagValueIDs returns the group by tag value ids for each tag key
func (g *GroupingContext) GetGroupByTagValueIDs() []*roaring.Bitmap {
	return g.tagValueIDs
}

// BuildGroup builds the grouped series ids by the high key of series id
// and the container includes low keys of series id.
func (g *GroupingContext) BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16 {
	// new tag value ids array for each group by tag key
	tagValueIDsForTags := g.buildTagValueIDs2SeriesIDs(highKey, container)

	min := container.Minimum()
	result := make(map[string][]uint16)
	tagValueIDs := make([]byte, len(g.tagKeys)*4)
	// iterator all series ids after filtering
	it := container.PeekableIterator()
	for it.HasNext() {
		seriesID := it.Next()
		found := true
		for idx := range g.tagKeys {
			scanners := tagValueIDsForTags[idx]
			tagValueID := scanners[seriesID-min]
			if tagValueID == 0 {
				found = false
				break
			}
			// collect group by tag value id
			groupByTagValueIDs := g.tagValueIDs[idx]
			if groupByTagValueIDs == nil {
				groupByTagValueIDs = roaring.New()
				g.tagValueIDs[idx] = groupByTagValueIDs
			}
			groupByTagValueIDs.Add(tagValueID)
			// build group key with group by tag value ids
			offset := idx * 4
			binary.LittleEndian.PutUint32(tagValueIDs[offset:], tagValueID)
		}
		if found {
			tagValuesStr := string(tagValueIDs)
			values, ok := result[tagValuesStr]
			if !ok {
				result[tagValuesStr] = []uint16{seriesID}
			} else {
				result[tagValuesStr] = append(values, seriesID)
			}
		}
	}
	return result
}

// buildTagValueIDs2SeriesIDs builds tag value id => series id mapping
func (g *GroupingContext) buildTagValueIDs2SeriesIDs(highKey uint16, container roaring.Container) [][]uint32 {
	// new seriesIDs2Tags array based on range of min ~ max
	min := container.Minimum()
	max := container.Maximum()
	seriesIDsLength := int(max-min) + 1
	tagValueIDsForTags := make([][]uint32, len(g.tagKeys))
	for i, tagKey := range g.tagKeys {
		scanners := g.scanners[tagKey]
		v := make([]uint32, seriesIDsLength)
		for _, scanner := range scanners {
			lowContainer, tagValueIDs := scanner.GetSeriesAndTagValue(highKey)
			it := lowContainer.PeekableIterator()
			idx := 0
			for it.HasNext() {
				seriesID := it.Next()
				if seriesID >= min && seriesID <= max {
					v[seriesID-min] = tagValueIDs[idx] // put tag value by index(series ids)
				}
				idx++
			}
		}
		tagValueIDsForTags[i] = v
	}
	return tagValueIDsForTags
}

package memdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series/tag"
)

type seriesID2Tags struct {
	tagValues       []string
	size, tagsCount int
}

func newSeriesID2Tags(tagsCount int) *seriesID2Tags {
	return &seriesID2Tags{
		tagValues: make([]string, tagsCount),
		tagsCount: tagsCount,
	}
}
func (entry *seriesID2Tags) addTagValue(tagValue string) {
	if entry.size < entry.tagsCount {
		entry.tagValues[entry.size] = tagValue
		entry.size++
	}
}

type groupingContext struct {
	ms             *metricStore
	tagKVEntrySets []*tagKVEntrySet
}

func (g *groupingContext) BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16 {
	groupTagKeysCount := len(g.tagKVEntrySets)
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
		if seriesID2Tags != nil && seriesID2Tags.size == groupTagKeysCount {
			tagValuesStr := tag.ConcatTagValues(seriesID2Tags.tagValues)
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

func (g *groupingContext) buildSeriesIDs2Tags(highKey uint16, container roaring.Container) []*seriesID2Tags {
	// need add read lock
	g.ms.mux.RLock()
	defer g.ms.mux.RUnlock()

	groupTagKeysCount := len(g.tagKVEntrySets)
	// new seriesIDs2Tags array based on range of max ~ min
	min := container.Minimum()
	max := container.Maximum()
	seriesIDs2Tags := make([]*seriesID2Tags, int(max-min)+1)

	// builds seriesIDs => tags mapping, using counting sort
	// https://en.wikipedia.org/wiki/Counting_sort
	for _, tagKV := range g.tagKVEntrySets {
		for tagValue, lowSeriesIDs := range tagKV.values {
			lowContainer := lowSeriesIDs.GetContainer(highKey)
			if lowContainer != nil {
				matchContainer := lowContainer.And(container)
				it := matchContainer.PeekableIterator()
				for it.HasNext() {
					idx := it.Next() - min // index = lowKey - min
					if seriesIDs2Tags[idx] == nil {
						seriesIDs2Tags[idx] = newSeriesID2Tags(groupTagKeysCount)
					}
					seriesIDs2Tags[idx].addTagValue(tagValue)
				}
			}
		}
	}
	return seriesIDs2Tags
}

func (g *groupingContext) buildForSingleTagKey(highKey uint16, container roaring.Container) map[string][]uint16 {
	// need add read lock
	g.ms.mux.RLock()
	defer g.ms.mux.RUnlock()

	result := make(map[string][]uint16)
	for _, tagKV := range g.tagKVEntrySets {
		for tagValue, lowSeriesIDs := range tagKV.values {
			lowContainer := lowSeriesIDs.GetContainer(highKey)
			if lowContainer != nil {
				matchContainer := lowContainer.And(container)
				result[tagValue] = matchContainer.ToArray()
			}
		}
	}
	return result
}

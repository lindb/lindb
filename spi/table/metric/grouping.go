package metric

import (
	"encoding/binary"
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
)

const (
	tagValueNotFound = "tag_value_not_found"
)

type Grouping struct {
	db   tsdb.Database
	tags tag.Metas

	// tag value ids for each grouping tag key
	groupingTagValueIDs []*roaring.Bitmap
	tagValuesMap        []map[uint32]string // tag value id=> tag value for each group by tag key
}

func NewGrouping(db tsdb.Database, tags tag.Metas) *Grouping {
	lenOfTags := tags.Len()
	return &Grouping{
		db:                  db,
		tags:                tags,
		groupingTagValueIDs: make([]*roaring.Bitmap, lenOfTags),
		tagValuesMap:        make([]map[uint32]string, lenOfTags),
	}
}

func (g *Grouping) CollectTagValueIDs(tagValueIDs []*roaring.Bitmap) {
	// TODO: add lock?
	for idx, ids := range tagValueIDs {
		if g.groupingTagValueIDs[idx] == nil {
			g.groupingTagValueIDs[idx] = ids
		} else {
			g.groupingTagValueIDs[idx].Or(ids)
		}
	}
}

func (g *Grouping) CollectTagValues() {
	metaDB := g.db.MetaDB()

	for idx := range g.groupingTagValueIDs {
		tagKey := g.tags[idx]
		tagValueIDs := g.groupingTagValueIDs[idx]

		if tagValueIDs == nil || tagValueIDs.IsEmpty() {
			continue
		}

		tagValues := make(map[uint32]string) // tag value id => tag value
		err := metaDB.CollectTagValues(tagKey.ID, tagValueIDs, tagValues)
		if err != nil {
			panic(err)
		}
		fmt.Printf("collect tag values...%v\n", tagValues)
		g.tagValuesMap[idx] = tagValues
	}
}

func (g *Grouping) GetTagValues(tagValueIDs string) []string {
	// if tagValues, ok := ctx.tagsMap[tagValueIDs]; ok {
	// 	return tagValues
	// }
	tagValues := make([]string, g.tags.Len())
	tagsData := []byte(tagValueIDs)
	for idx := range g.tagValuesMap {
		tagValuesForKey := g.tagValuesMap[idx]
		offset := idx * 4
		tagValueID := binary.LittleEndian.Uint32(tagsData[offset:])
		if tagValue, ok := tagValuesForKey[tagValueID]; ok {
			tagValues[idx] = tagValue
		} else {
			tagValues[idx] = tagValueNotFound
		}
	}
	return tagValues
}

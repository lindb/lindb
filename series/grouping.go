package series

import (
	"github.com/lindb/roaring"
)

//go:generate mockgen -source=./grouping.go -destination=./grouping_mock.go -package=series

// GroupingContext represents the context of group by query for tag keys
type GroupingContext interface {
	// BuildGroup builds the grouped series ids by the high key of series id
	// and the container includes low keys of series id
	BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16
	// GetGroupByTagValueIDs returns the group by tag value ids for each tag key
	GetGroupByTagValueIDs() []*roaring.Bitmap
}

// GroupingScanner represents the scanner which scans the group by data by high key of series id
type GroupingScanner interface {
	// GetSeriesAndTagValue returns group by container and tag value ids
	GetSeriesAndTagValue(highKey uint16) (roaring.Container, []uint32)
}

// Grouping represents the getter grouping scanners for tag key group by query
type Grouping interface {
	// GetGroupingScanner returns the grouping scanners based on tag key ids and series ids
	GetGroupingScanner(tagKeyID uint32, seriesIDs *roaring.Bitmap) ([]GroupingScanner, error)
}

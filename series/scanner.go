package series

import (
	"github.com/lindb/roaring"
)

//go:generate mockgen -source=./scanner.go -destination=./scanner_mock.go -package=series

// GroupingContext represents the context of group by query for tag keys
type GroupingContext interface {
	// BuildGroup builds the grouped series ids by the high key of series id
	// and the container includes low keys of series id
	BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16
}

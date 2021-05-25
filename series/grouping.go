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
	// ScanTagValueIDs scans grouping context by high key/container of series ids,
	// then returns grouped tag value ids for each tag key
	ScanTagValueIDs(highKey uint16, container roaring.Container) []*roaring.Bitmap
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

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

package indexdb

import "github.com/lindb/roaring"

// memGroupingScanner implements series.GroupingScanner for memory tag index
type memGroupingScanner struct {
	forward  *ForwardStore
	withLock func() (release func())
}

// GetSeriesAndTagValue returns group by container and tag value ids
func (g *memGroupingScanner) GetSeriesAndTagValue(highKey uint16) (lowSeriesIDs roaring.Container, tagValueIDs []uint32) {
	release := g.withLock()
	defer release()

	index := g.forward.keys.GetContainerIndex(highKey)
	if index < 0 {
		// data not found
		return nil, nil
	}
	return g.forward.keys.GetContainerAtIndex(index), g.forward.values[index]
}

// GetSeriesIDs returns the series ids in current memory scanner.
func (g *memGroupingScanner) GetSeriesIDs() *roaring.Bitmap {
	release := g.withLock()
	defer release()

	return g.forward.keys
}

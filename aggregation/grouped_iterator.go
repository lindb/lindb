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

package aggregation

import "github.com/lindb/lindb/series"

//////////////////////////////////////////////////////
// groupedIterator implements GroupedIterator
//////////////////////////////////////////////////////
type groupedIterator struct {
	tags       string // tag values
	aggregates FieldAggregates
	len        int
	idx        int
}

// newGroupedIterator creates a grouped iterator for field aggregates
func newGroupedIterator(tags string, aggregates FieldAggregates) series.GroupedIterator {
	return &groupedIterator{tags: tags, aggregates: aggregates, len: len(aggregates)}
}

// Tags returns the tags of series
func (g *groupedIterator) Tags() string {
	return g.tags
}

// HasNext returns if the iteration has more field's iterator
func (g *groupedIterator) HasNext() bool {
	if g.idx >= g.len {
		return false
	}
	g.idx++
	return true
}

// Next returns the result set of aggregator
func (g *groupedIterator) Next() series.Iterator {
	return g.aggregates[g.idx-1].ResultSet()
}

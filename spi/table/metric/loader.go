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

package metric

import (
	"github.com/lindb/common/models"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

type familyLoader struct {
	familyIndex     int
	filterResultSet flow.FilterResultSet
	loader          flow.DataLoader
	columns         []Column
}

func (fl *familyLoader) load(field field.Meta, getter encoding.TSDValueGetter) {
	slotRange := fl.filterResultSet.SlotRange()
	fc := fl.columns[field.Index]
	switch c := (fc).(type) {
	case *column[float64]:
		c.load(fl.familyIndex, slotRange, getter.GetValue)
	case *column[*models.Exemplar]:
		c.load(fl.familyIndex, slotRange, getter.GetExemplar)
	}
}

type loader struct {
	interval timeutil.Interval // storage interval
	families []*familyLoader

	familyTime int64
	timeRange  timeutil.SlotRange // time range
}

func (l *loader) load(lowSeriesID uint16) {
	for _, family := range l.families {
		family.loader.Load(lowSeriesID, family.load)
	}
}

type Streams[V float64 | *models.Exemplar] struct {
	streams []Stream[V] // stream of fields
}

func newStreams[V float64 | *models.Exemplar](numOfFamilies int) *Streams[V] {
	return &Streams[V]{
		streams: make([]Stream[V], numOfFamilies),
	}
}

func (s *Streams[V]) GetStreamByIndex(familyIndex int, fn func() Stream[V]) Stream[V] {
	stream := s.streams[familyIndex]
	if stream == nil {
		stream = fn()
		s.streams[familyIndex] = stream
	}
	return stream
}

func (s *Streams[V]) reset() {
	for i := range s.streams {
		if s.streams[i] != nil {
			s.streams[i].Reset()
		}
	}
}

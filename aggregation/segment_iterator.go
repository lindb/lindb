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

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// seriesIterator implements series.Iterator.
type seriesIterator struct {
	fieldName  field.Name
	fieldType  field.Type
	aggregates []FieldAggregator
	idx        int
	len        int
}

// newSeriesIterator creates the time series iterator
func newSeriesIterator(agg SeriesAggregator) series.Iterator {
	it := &seriesIterator{fieldName: agg.FieldName(), fieldType: agg.GetFieldType(), aggregates: agg.GetAggregates()}
	it.len = len(it.aggregates)
	return it
}

// FieldName returns field name
func (s *seriesIterator) FieldName() field.Name {
	return s.fieldName
}

// FieldType returns field type
func (s *seriesIterator) FieldType() field.Type {
	return s.fieldType
}

// HasNext returns if the iteration has more field's iterator
func (s *seriesIterator) HasNext() bool {
	for s.idx < s.len {
		if s.aggregates[s.idx] != nil {
			s.idx++
			return true
		}
		s.idx++
	}
	return false
}

// Next returns the field's iterator and segment start time
func (s *seriesIterator) Next() (startTime int64, fieldIt series.FieldIterator) {
	return s.aggregates[s.idx-1].ResultSet()
}

func (s *seriesIterator) MarshalBinary() ([]byte, error) {
	return series.MarshalIterator(s)
}

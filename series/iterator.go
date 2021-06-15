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
	enc "encoding"

	"github.com/lindb/lindb/models"

	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./iterator.go -destination=./iterator_mock.go -package=series

// TimeSeriesEvent represents time series event for query.
type TimeSeriesEvent struct {
	SeriesList []GroupedIterator

	Stats *models.QueryStats
	Err   error
}

// GroupedIterator represents a iterator for the grouped time series data.
type GroupedIterator interface {
	// HasNext returns if the iteration has more field's iterator.
	HasNext() bool
	// Next returns the field's iterator.
	Next() Iterator
	// Tags returns group tags, tags is tag values concat string.
	Tags() string
}

// Iterator represents an iterator for the time series data.
type Iterator interface {
	// FieldName returns the field name.
	FieldName() field.Name
	// FieldType returns the field type.
	FieldType() field.Type
	// HasNext returns if the iteration has more field's iterator.
	HasNext() bool
	// Next returns the field's iterator.
	Next() (startTime int64, fieldIt FieldIterator)
	// MarshalBinary marshals the data.
	enc.BinaryMarshaler
}

// FieldIterator represents a field's data iterator, support multi field for one series.
type FieldIterator interface {
	// HasNext returns if the iteration has more fields.
	HasNext() bool
	// Next returns the data point in the iteration.
	Next() PrimitiveIterator
	// MarshalBinary marshals the data.
	enc.BinaryMarshaler
}

// PrimitiveIterator represents an iterator over a primitive field, iterator points data of primitive field.
type PrimitiveIterator interface {
	// AggType returns the primitive field's agg type.
	AggType() field.AggType
	// HasNext returns if the iteration has more data points.
	HasNext() bool
	// Next returns the data point in the iteration.
	Next() (timeSlot int, value float64)
}

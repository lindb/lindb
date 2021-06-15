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

package fields

import (
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./field.go -destination=./field_mock.go -package=fields

// Field represents the field series for the time series.
type Field interface {
	// SetValue sets field value using series iterator.
	SetValue(fieldSeries series.Iterator)
	// GetValues returns the values which function call need by given function type.
	GetValues(funcType function.FuncType) (result []collections.FloatArray)
	// GetDefaultValues returns the field default values which aggregation need if user not input function type.
	GetDefaultValues() (result []collections.FloatArray)
	// Reset resets field's value for reusing.
	Reset()
}

// dynamicField represents the dynamic field for storing multi-agg types.
type dynamicField struct {
	fieldType field.Type
	startTime int64
	interval  int64
	capacity  int

	fields map[field.AggType]collections.FloatArray
}

// NewDynamicField creates a dynamic field series.
func NewDynamicField(fieldType field.Type, startTime int64, interval int64, capacity int) Field {
	return &dynamicField{
		fieldType: fieldType,
		startTime: startTime,
		interval:  interval,
		capacity:  capacity,
		fields:    make(map[field.AggType]collections.FloatArray),
	}
}

// SetValue sets the field's value by time slot
func (f *dynamicField) SetValue(fieldSeries series.Iterator) {
	if fieldSeries == nil {
		return
	}
	var fieldValues collections.FloatArray
	ok := false
	for fieldSeries.HasNext() {
		startTime, it := fieldSeries.Next()
		if it == nil {
			continue
		}
		for it.HasNext() {
			pIt := it.Next()
			aggType := pIt.AggType()
			fieldValues, ok = f.fields[aggType]
			if !ok {
				fieldValues = collections.NewFloatArray(f.capacity)
				f.fields[aggType] = fieldValues
			}
			for pIt.HasNext() {
				slot, val := pIt.Next()
				idx := ((int64(slot)*f.interval + startTime) - f.startTime) / f.interval
				fieldValues.SetValue(int(idx), val)
			}
		}
	}
}

// GetValues returns the values which function call need by given function type and field type
func (f *dynamicField) GetValues(funcType function.FuncType) (result []collections.FloatArray) {
	pFields := f.fieldType.GetFuncFieldParams(funcType)
	return f.getFieldValues(pFields)
}

// GetDefaultValues returns the field default values which aggregation need by field type
func (f *dynamicField) GetDefaultValues() []collections.FloatArray {
	return f.getFieldValues(f.fieldType.GetDefaultFuncFieldParams())
}

func (f *dynamicField) Reset() {
	for _, pField := range f.fields {
		pField.Reset()
	}
}

// getFieldValues returns the values by field name and agg type.
func (f *dynamicField) getFieldValues(aggTypes []field.AggType) (result []collections.FloatArray) {
	if len(aggTypes) == 0 {
		return
	}
	for _, aggType := range aggTypes {
		pField, ok := f.fields[aggType]
		if ok {
			result = append(result, pField)
		}
	}
	return
}

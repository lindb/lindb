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
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

// Aggregator represents aggregator spec for down sampling/aggregator.
type Aggregator struct {
	DownSampling AggregatorSpec
	Aggregator   AggregatorSpec
}

// AggregatorSpecs represents aggregator spec slice.
type AggregatorSpecs []AggregatorSpec

// AggregatorSpec represents aggregator spec.
type AggregatorSpec interface {
	// FieldName returns field name.
	FieldName() field.Name
	// GetFieldType sets field type
	GetFieldType() field.Type
	// AddFunctionType adds function type for down sampling.
	AddFunctionType(funcType function.FuncType)
	// Functions returns function types for down sampling.
	Functions() map[function.FuncType]function.FuncType
}

// aggregatorSpec implements AggregatorSpec interface.
type aggregatorSpec struct {
	fieldName field.Name
	fieldType field.Type
	functions map[function.FuncType]function.FuncType
}

// NewAggregatorSpec creates a AggregatorSpec.
func NewAggregatorSpec(fieldName field.Name, fieldType field.Type) AggregatorSpec {
	return &aggregatorSpec{
		fieldName: fieldName,
		fieldType: fieldType,
		functions: make(map[function.FuncType]function.FuncType),
	}
}

// GetFieldType sets field type
func (a *aggregatorSpec) GetFieldType() field.Type {
	return a.fieldType
}

// FieldName returns field name.
func (a *aggregatorSpec) FieldName() field.Name {
	return a.fieldName
}

// AddFunctionType adds function type for down sampling.
func (a *aggregatorSpec) AddFunctionType(funcType function.FuncType) {
	_, exist := a.functions[funcType]
	if !exist {
		a.functions[funcType] = funcType
	}
}

// Functions returns function types for down sampling.
func (a *aggregatorSpec) Functions() map[function.FuncType]function.FuncType {
	return a.functions
}

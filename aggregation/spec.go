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

type Aggregator struct {
	DownSampling AggregatorSpec
	Aggregator   AggregatorSpec
}

type AggregatorSpecs []AggregatorSpec

type AggregatorSpec interface {
	FieldName() field.Name
	GetFieldType() field.Type
	AddFunctionType(funcType function.FuncType)
	Functions() map[function.FuncType]function.FuncType
}

type aggregatorSpec struct {
	fieldName field.Name
	fieldType field.Type
	functions map[function.FuncType]function.FuncType
}

func NewAggregatorSpec(fieldName field.Name, fieldType field.Type) AggregatorSpec {
	return &aggregatorSpec{
		fieldName: fieldName,
		fieldType: fieldType,
		functions: make(map[function.FuncType]function.FuncType),
	}
}

func (a *aggregatorSpec) GetFieldType() field.Type {
	return a.fieldType
}

func (a *aggregatorSpec) SetFieldType(fieldType field.Type) {
	a.fieldType = fieldType
}

func (a *aggregatorSpec) FieldName() field.Name {
	return a.fieldName
}

func (a *aggregatorSpec) AddFunctionType(funcType function.FuncType) {
	_, exist := a.functions[funcType]
	if !exist {
		a.functions[funcType] = funcType
	}
}

func (a *aggregatorSpec) Functions() map[function.FuncType]function.FuncType {
	return a.functions
}

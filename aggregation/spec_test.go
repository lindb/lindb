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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

func TestAggregatorSpec_FieldName(t *testing.T) {
	agg := NewAggregatorSpec("f1", field.SumField)
	assert.Equal(t, field.Name("f1"), agg.FieldName())
	assert.Equal(t, field.SumField, agg.GetFieldType())
}

func TestAggregatorSpec_AddFunctionType(t *testing.T) {
	agg := NewAggregatorSpec("f1", field.SumField)
	agg.AddFunctionType(function.Sum)
	agg.AddFunctionType(function.Sum)
	assert.Equal(t, 1, len(agg.Functions()))
}

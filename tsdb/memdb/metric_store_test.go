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

package memdb

import (
	"testing"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
)

func TestMetricStore_genField(t *testing.T) {
	ms := &metricStore{}
	f, isNew := ms.genField("test", field.SumField)
	assert.True(t, isNew)
	assert.Equal(t, field.Meta{Name: "test", Type: field.SumField, Index: 0}, f)

	f, isNew = ms.genField("test", field.SumField)
	assert.False(t, isNew)
	assert.Equal(t, field.Meta{Name: "test", Type: field.SumField, Index: 0}, f)
}

func TestMetricStore_IsAction(t *testing.T) {
	ms := newMetricStore()
	_, isNew := ms.GenField("test", field.SumField)
	assert.True(t, isNew)
	assert.True(t, ms.IsActive(fasttime.UnixMilliseconds()))
	assert.False(t, ms.IsActive(fasttime.UnixMilliseconds()+timeutil.OneDay))
}

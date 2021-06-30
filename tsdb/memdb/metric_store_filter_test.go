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
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series/field"
)

func TestMetricStore_Filter(t *testing.T) {
	metricStore := mockMetricStore()

	// case 1: field not found
	rs, err := metricStore.Filter(1, nil, field.Metas{{ID: 1}, {ID: 2}})
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, rs)
	// case 3: series ids not found
	rs, err = metricStore.Filter(1, roaring.BitmapOf(1, 2), field.Metas{{ID: 1}, {ID: 20, Type: field.SumField}})
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, rs)
	// case 3: found data
	rs, err = metricStore.Filter(1, roaring.BitmapOf(1, 100, 200), field.Metas{{ID: 1}, {ID: 20, Type: field.SumField}})
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	mrs := rs[0].(*memFilterResultSet)
	assert.EqualValues(t, roaring.BitmapOf(100, 200).ToArray(), mrs.SeriesIDs().ToArray())
	assert.Equal(t,
		field.Metas{
			{ID: 1}, {
				ID:   20,
				Type: field.SumField,
			}}, mrs.fields)
	assert.Equal(t, "memory", rs[0].Identifier())
}

func TestMemFilterResultSet_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mStore := mockMetricStore()

	rs, err := mStore.Filter(1, roaring.BitmapOf(1, 100, 200), field.Metas{{ID: 1}, {ID: 20}})
	assert.NoError(t, err)
	// case 1: load data success
	loader := rs[0].Load(0, roaring.BitmapOf(100, 200).GetContainer(0))
	assert.NotNil(t, loader)
	loader.Load(100)
	loader.Load(200)
	// case 2: series ids not found
	rs, _ = mStore.Filter(1, roaring.BitmapOf(1, 100, 200), field.Metas{{ID: 1}, {ID: 20}})
	loader = rs[0].Load(0, roaring.BitmapOf(1, 2).GetContainer(0))
	assert.Nil(t, loader)
	// case 3: high key not exist
	rs, _ = mStore.Filter(1, roaring.BitmapOf(1, 100, 200), field.Metas{{ID: 1}, {ID: 20}})
	loader = rs[0].Load(10, roaring.BitmapOf(1, 2).GetContainer(0))
	assert.Nil(t, loader)
	// case 4: field not exist
	rs, err = mStore.Filter(1, roaring.BitmapOf(1, 100, 200), field.Metas{{ID: 100}, {ID: 200}})
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, rs)
}

func mockMetricStore() *metricStore {
	mStore := newMetricStore()
	mStore.AddField(field.ID(10), field.SumField)
	mStore.AddField(field.ID(20), field.SumField)
	mStore.SetSlot(10)
	mStore.SetSlot(20)
	mStore.GetOrCreateTStore(100)
	mStore.GetOrCreateTStore(120)
	mStore.GetOrCreateTStore(200)
	return mStore.(*metricStore)
}

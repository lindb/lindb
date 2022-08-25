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

package context

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

func TestLeafReduceContext_Reduce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageCtx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{},
	}
	ctx := NewLeafReduceContext(storageCtx, &LeafGroupingContext{})

	it := series.NewMockGroupedIterator(ctrl)
	it.EXPECT().Tags().Return("")
	it.EXPECT().HasNext().Return(false)
	ctx.Reduce(it)
}

func TestLeafReduceContext_BuildResultSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := aggregation.NewAggregatorSpec("f", field.SumField)
	spec.AddFunctionType(function.Sum)
	storageCtx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{
			GroupBy: []string{"key"},
		},
		AggregatorSpecs: aggregation.AggregatorSpecs{spec},
	}
	ctx := NewLeafReduceContext(storageCtx, &LeafGroupingContext{
		tagsMap: map[string]string{},
	})
	cases := []struct {
		name    string
		in      *models.Leaf
		prepare func()
		assert  func(rs [][]byte)
	}{
		{
			name: "send to root",
			in: &models.Leaf{
				Receivers: []models.StatelessNode{{}},
			},
			assert: func(rs [][]byte) {
				assert.Len(t, rs, 1)
			},
		},
		{
			name: "marshal series data failure",
			in: &models.Leaf{
				Receivers: []models.StatelessNode{{}},
			},
			prepare: func() {
				agg := aggregation.NewMockGroupingAggregator(ctrl)
				ctx.reduceAgg = agg
				gIt := series.NewMockGroupedIterator(ctrl)
				agg.EXPECT().ResultSet().Return(series.GroupedIterators{gIt})
				gIt.EXPECT().HasNext().Return(true)
				it := series.NewMockIterator(ctrl)
				gIt.EXPECT().Next().Return(it)
				it.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
				gIt.EXPECT().HasNext().Return(false)
			},
			assert: func(rs [][]byte) {
				assert.Len(t, rs, 1)
			},
		},
		{
			name: "need hash rs",
			in: &models.Leaf{
				Receivers: []models.StatelessNode{{}, {}},
			},
			prepare: func() {
				agg := aggregation.NewMockGroupingAggregator(ctrl)
				ctx.reduceAgg = agg
				gIt := series.NewMockGroupedIterator(ctrl)
				gIt.EXPECT().Tags().Return("key")
				agg.EXPECT().ResultSet().Return(series.GroupedIterators{gIt})
				gIt.EXPECT().HasNext().Return(true)
				it := series.NewMockIterator(ctrl)
				it.EXPECT().FieldName().Return(field.Name("f"))
				gIt.EXPECT().Next().Return(it)
				it.EXPECT().MarshalBinary().Return([]byte{1, 2, 3}, nil)
				gIt.EXPECT().HasNext().Return(false)
			},
			assert: func(rs [][]byte) {
				assert.Len(t, rs, 2)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				ctx.reduceAgg = nil
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			tt.assert(ctx.BuildResultSet(tt.in))
		})
	}
}

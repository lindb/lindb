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

package operator

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

func TestSeriesFiltering_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	storageCtx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{},
		TagFilterResult: map[string]*flow.TagFilterResult{
			"key1=value1": {
				TagKeyID:    tag.KeyID(1),
				TagValueIDs: roaring.BitmapOf(1, 2, 3),
			},
		},
	}
	shardCtx := flow.NewShardExecuteContext(storageCtx)
	cases := []struct {
		name    string
		in      stmtpkg.Expr
		prepare func()
		wantErr bool
	}{
		{
			name: "condition is empty",
		},
		{
			name: "tag values not found from context",
			in: &stmtpkg.EqualsExpr{
				Key:   "key2",
				Value: "value2",
			},
			wantErr: true,
		},
		{
			name: "find series failure",
			in: &stmtpkg.EqualsExpr{
				Key:   "key1",
				Value: "value1",
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "find series successfully",
			in: &stmtpkg.EqualsExpr{
				Key:   "key1",
				Value: "value1",
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2), nil)
			},
		},
		{
			name: "paren expr successfully",
			in: &stmtpkg.ParenExpr{
				Expr: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2), nil)
			},
		},
		{
			name: "not expr successfully",
			in: &stmtpkg.NotExpr{
				Expr: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsForTag(gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil)
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2), nil)
			},
		},
		{
			name: "not expr failure",
			in: &stmtpkg.NotExpr{
				Expr: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsForTag(gomock.Any()).Return(nil, fmt.Errorf("err"))
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2), nil)
			},
			wantErr: true,
		},
		{
			name: "binary expr failure",
			in: &stmtpkg.BinaryExpr{
				Right: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
				Left: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "binary expr(OR) successfully",
			in: &stmtpkg.BinaryExpr{
				Right: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
				Operator: stmtpkg.OR,
				Left: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).
					Return(roaring.BitmapOf(1, 2), nil).MaxTimes(2)
			},
		},
		{
			name: "binary expr(AND) successfully",
			in: &stmtpkg.BinaryExpr{
				Right: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
				Operator: stmtpkg.AND,
				Left: &stmtpkg.EqualsExpr{
					Key:   "key1",
					Value: "value1",
				},
			},
			prepare: func() {
				indexDB.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).
					Return(roaring.BitmapOf(1, 2), nil).MaxTimes(2)
			},
		},
		{
			name: "unknown condition expr",
			in: &stmtpkg.FieldExpr{
				Name: "f",
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := NewSeriesFiltering(shardCtx, shard)
			storageCtx.Query.Condition = tt.in
			if tt.prepare != nil {
				tt.prepare()
			}
			err := op.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

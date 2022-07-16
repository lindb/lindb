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
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestTagValuesLookup_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	meta.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	meta.EXPECT().MetadataDatabase().Return(metaDB).AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	ctx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{},
		TagKeys: map[string]tag.KeyID{
			"key-10": tag.KeyID(10),
		},
	}
	cases := []struct {
		name    string
		in      stmtpkg.Expr
		prepare func()
		wantErr bool
	}{
		{
			name:    "empty condition",
			wantErr: false,
		},
		{
			name: "get tag key failure",
			in: &stmtpkg.EqualsExpr{
				Key:   "key",
				Value: "value",
			},
			prepare: func() {
				metaDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.EmptyTagKeyID, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get tag values failure",
			in: &stmtpkg.EqualsExpr{
				Key:   "key-10",
				Value: "value",
			},
			prepare: func() {
				tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get tag values successfully",
			in: &stmtpkg.EqualsExpr{
				Key:   "key-1",
				Value: "value",
			},
			prepare: func() {
				metaDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.EmptyTagKeyID, nil)
				tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(1, 2, 3), nil)
			},
			wantErr: false,
		},
		{
			name: "wrong op type",
			in: &stmtpkg.BinaryExpr{
				Operator: stmtpkg.DIV,
			},
			wantErr: true,
		},
		{
			name: "binary expr successfully",
			in: &stmtpkg.BinaryExpr{
				Left: &stmtpkg.EqualsExpr{
					Key:   "key-10",
					Value: "value",
				},
				Operator: stmtpkg.AND,
				Right: &stmtpkg.EqualsExpr{
					Key:   "key-10",
					Value: "value",
				},
			},
			prepare: func() {
				tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).
					Return(roaring.BitmapOf(1, 2, 3), nil).MaxTimes(2)
			},
			wantErr: false,
		},
		{
			name: "binary expr failure",
			in: &stmtpkg.BinaryExpr{
				Left: &stmtpkg.EqualsExpr{
					Key:   "key-10",
					Value: "value",
				},
				Operator: stmtpkg.AND,
				Right: &stmtpkg.EqualsExpr{
					Key:   "key-10",
					Value: "value",
				},
			},
			prepare: func() {
				tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "not expr successfully",
			in: &stmtpkg.NotExpr{
				Expr: &stmtpkg.EqualsExpr{
					Key:   "key-10",
					Value: "value",
				},
			},
			prepare: func() {
				tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).
					Return(roaring.BitmapOf(1, 2, 3), nil)
			},
			wantErr: false,
		},
		{
			name: "paren expr successfully",
			in: &stmtpkg.ParenExpr{
				Expr: &stmtpkg.EqualsExpr{
					Key:   "key-10",
					Value: "value",
				},
			},
			prepare: func() {
				tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).
					Return(roaring.BitmapOf(1, 2, 3), nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := NewTagValuesLookup(ctx, db)
			ctx.Query.Condition = tt.in
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

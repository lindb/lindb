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
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestMetadataLookup_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	meta.EXPECT().MetadataDatabase().Return(metaDB).AnyTimes()

	ctx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{},
	}
	t.Run("find metric id failure", func(t *testing.T) {
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("get group tag failure", func(t *testing.T) {
		defer func() {
			ctx.Query.GroupBy = nil
		}()
		ctx.Query.GroupBy = []string{"a"}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.EmptyTagKeyID, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("select item empty", func(t *testing.T) {
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		assert.Error(t, op.Execute())
	})
	t.Run("get field failure", func(t *testing.T) {
		ctx.Query.SelectItems = []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f"}}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).Return(field.Meta{}, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("execute successfully", func(t *testing.T) {
		ctx.Query.SelectItems = []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f"}}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).Return(field.Meta{
			ID:   10,
			Type: field.SumField,
			Name: "f",
		}, nil)
		assert.NoError(t, op.Execute())
	})
}

func TestMetadataLookup_groupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := metadb.NewMockMetadataDatabase(ctrl)

	ctx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{
			GroupBy: []string{"k"},
		},
		TagKeys: make(map[string]tag.KeyID),
	}

	op := &metadataLookup{
		executeCtx: ctx,
		metadata:   metaDB,
	}

	metaDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "k").Return(tag.KeyID(10), nil)
	assert.NoError(t, op.groupBy())
}

func TestMetadataLookup_field(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	ctx := &flow.StorageExecuteContext{
		Query:   &stmtpkg.Query{},
		TagKeys: make(map[string]tag.KeyID),
	}
	metaDB.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).Return(field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f",
	}, nil).AnyTimes()

	t.Run("has err", func(t *testing.T) {
		op := &metadataLookup{err: fmt.Errorf("err")}
		op.field(nil, nil)
	})

	t.Run("not support field", func(t *testing.T) {
		metaDB2 := metadb.NewMockMetadataDatabase(ctrl)
		metaDB2.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).Return(field.Meta{
			Type: field.Unknown,
			Name: "f",
		}, nil)
		op := &metadataLookup{
			executeCtx: ctx,
			metadata:   metaDB2,
			fields:     make(map[field.ID]*aggregation.Aggregator),
		}
		op.field(nil, &stmtpkg.SelectItem{
			Expr:  &stmtpkg.FieldExpr{Name: "f"},
			Alias: "a",
		})
		assert.Error(t, op.err)
	})
	t.Run("function not support", func(t *testing.T) {
		metaDB2 := metadb.NewMockMetadataDatabase(ctrl)
		metaDB2.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).Return(field.Meta{
			ID:   field.ID(10),
			Type: field.MinField,
			Name: "f",
		}, nil)
		op := &metadataLookup{
			executeCtx: ctx,
			metadata:   metaDB2,
			fields:     make(map[field.ID]*aggregation.Aggregator),
		}
		op.field(nil, &stmtpkg.CallExpr{
			FuncType: function.Sum,
			Params:   []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f1"}},
		})
		assert.Error(t, op.err)
	})

	cases := []struct {
		name    string
		in      stmtpkg.Expr
		wantErr bool
	}{
		{
			name: "handle select expr",
			in: &stmtpkg.SelectItem{
				Expr:  &stmtpkg.FieldExpr{Name: "f"},
				Alias: "a",
			},
		},
		{
			name: "handle binary expr",
			in: &stmtpkg.BinaryExpr{
				Left:  &stmtpkg.FieldExpr{Name: "f1"},
				Right: &stmtpkg.FieldExpr{Name: "f2"},
			},
		},
		{
			name: "handle sum function",
			in: &stmtpkg.CallExpr{
				FuncType: function.Sum,
				Params:   []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f1"}},
			},
		},
		{
			name: "handle quantile function",
			in: &stmtpkg.CallExpr{
				FuncType: function.Quantile,
				Params:   []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f1"}},
			},
			wantErr: true,
		},
		{
			name: "handle paren",
			in: &stmtpkg.ParenExpr{
				Expr: &stmtpkg.BinaryExpr{
					Left:     &stmtpkg.CallExpr{FuncType: function.Sum, Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f"}}},
					Operator: stmtpkg.ADD,
					Right:    &stmtpkg.FieldExpr{Name: "a"},
				}},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := &metadataLookup{
				executeCtx: ctx,
				metadata:   metaDB,
				fields:     make(map[field.ID]*aggregation.Aggregator),
			}
			op.field(nil, tt.in)
			if (op.err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

func TestMetadataLookup_planHistogramFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	cases := []struct {
		name    string
		in      *stmtpkg.CallExpr
		prepare func()
		wantErr bool
	}{
		{
			name: "invalid params",
			in: &stmtpkg.CallExpr{
				Params: nil,
			},
			wantErr: true,
		},
		{
			name: "parse params failure",
			in: &stmtpkg.CallExpr{
				Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f1"}},
			},
			wantErr: true,
		},
		{
			name: "parse params out of range",
			in: &stmtpkg.CallExpr{
				Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "10.3"}},
			},
			wantErr: true,
		},
		{
			name: "find filed failure",
			in: &stmtpkg.CallExpr{
				Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "0.95"}},
			},
			prepare: func() {
				metaDB.EXPECT().GetAllHistogramFields(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "find filed successfully",
			in: &stmtpkg.CallExpr{
				Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "0.95"}},
			},
			prepare: func() {
				metaDB.EXPECT().GetAllHistogramFields(gomock.Any(), gomock.Any()).
					Return(field.Metas{{
						Type: field.SumField,
						Name: "11",
					}}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := &metadataLookup{
				executeCtx: &flow.StorageExecuteContext{
					Query: &stmtpkg.Query{},
				},
				metadata: metaDB,
				fields:   make(map[field.ID]*aggregation.Aggregator),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			op.planHistogramFields(tt.in)
			if (op.err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

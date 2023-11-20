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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestMetadataLookup_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	db.EXPECT().MetaDB().Return(metaDB).AnyTimes()

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
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("schema field empty", func(t *testing.T) {
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{}, nil)
		assert.Error(t, op.Execute())
	})
	t.Run("select item empty", func(t *testing.T) {
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{Fields: field.Metas{{}}}, nil)
		assert.Error(t, op.Execute())
	})
	t.Run("select item field not exist", func(t *testing.T) {
		ctx.Query.SelectItems = []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f"}}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{Fields: field.Metas{{}}}, nil)
		assert.Error(t, op.Execute())
	})
	t.Run("get field failure", func(t *testing.T) {
		ctx.Query.SelectItems = []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f"}}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("execute successfully", func(t *testing.T) {
		ctx.Query.SelectItems = []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f"}}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{
			Fields: field.Metas{{
				ID:   10,
				Type: field.SumField,
				Name: "f",
			}},
		}, nil)
		assert.NoError(t, op.Execute())
	})
	t.Run("get all fields failure", func(t *testing.T) {
		ctx.Query.AllFields = true
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, op.Execute())
	})
	t.Run("get all fields successfully", func(t *testing.T) {
		ctx.Query.AllFields = true
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{
			Fields: field.Metas{{
				ID:   1,
				Type: field.SumField,
				Name: "f",
			}},
		}, nil)
		assert.NoError(t, op.Execute())
	})
	t.Run("metric schema not exist", func(t *testing.T) {
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(nil, nil)
		assert.Error(t, op.Execute())
	})
	t.Run("group by key not exist", func(t *testing.T) {
		defer func() {
			ctx.Query.GroupBy = nil
		}()
		ctx.Query.GroupBy = []string{"a"}
		op := NewMetadataLookup(ctx, db)
		metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
		metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{Fields: field.Metas{{}}}, nil)
		assert.Error(t, op.Execute())
	})
}

func TestMetadataLookup_groupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := index.NewMockMetricMetaDatabase(ctrl)

	ctx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{
			GroupBy: []string{"k"},
		},
		Schema: &metric.Schema{
			TagKeys: tag.Metas{{ID: tag.KeyID(10), Key: "k"}},
		},
	}

	op := &metadataLookup{
		executeCtx: ctx,
		metaDB:     metaDB,
	}

	assert.NoError(t, op.groupBy())
}

func TestMetadataLookup_field(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	ctx := &flow.StorageExecuteContext{
		Query: &stmtpkg.Query{},
		Schema: &metric.Schema{
			Fields: field.Metas{
				{
					ID:   10,
					Type: field.SumField,
					Name: "f1",
				},
				{
					ID:   12,
					Type: field.SumField,
					Name: "f2",
				},
				{
					Name: "f_not_func",
					ID:   field.ID(10),
					Type: field.MinField,
				},
				{
					Type: field.Unknown,
					Name: "f_not",
				},
			},
		},
	}

	t.Run("has err", func(_ *testing.T) {
		op := &metadataLookup{err: fmt.Errorf("err")}
		op.field(nil, nil)
	})

	t.Run("not support field", func(t *testing.T) {
		metaDB2 := index.NewMockMetricMetaDatabase(ctrl)
		op := &metadataLookup{
			executeCtx: ctx,
			metaDB:     metaDB2,
			fields:     make(map[field.ID]*aggregation.Aggregator),
		}
		op.field(nil, &stmtpkg.SelectItem{
			Expr:  &stmtpkg.FieldExpr{Name: "f_not"},
			Alias: "a",
		})
		assert.Error(t, op.err)
	})
	t.Run("function not support", func(t *testing.T) {
		op := &metadataLookup{
			executeCtx: ctx,
			fields:     make(map[field.ID]*aggregation.Aggregator),
		}
		op.field(nil, &stmtpkg.CallExpr{
			FuncType: function.Sum,
			Params:   []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f_not_func"}},
		})
		assert.Error(t, op.err)
	})

	cases := []struct {
		name    string
		in      stmtpkg.Expr
		wantErr bool
	}{
		{
			name: "field not exist",
			in: &stmtpkg.SelectItem{
				Expr:  &stmtpkg.FieldExpr{Name: "f1_xx"},
				Alias: "a",
			},
			wantErr: true,
		},
		{
			name: "handle select expr",
			in: &stmtpkg.SelectItem{
				Expr:  &stmtpkg.FieldExpr{Name: "f1"},
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
					Left:     &stmtpkg.CallExpr{FuncType: function.Sum, Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "f1"}}},
					Operator: stmtpkg.ADD,
					Right:    &stmtpkg.FieldExpr{Name: "f2"},
				}},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := &metadataLookup{
				executeCtx: ctx,
				metaDB:     metaDB,
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
			name: "successfully",
			in: &stmtpkg.CallExpr{
				Params: []stmtpkg.Expr{&stmtpkg.FieldExpr{Name: "0.99"}},
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := &metadataLookup{
				executeCtx: &flow.StorageExecuteContext{
					Query: &stmtpkg.Query{},
					Schema: &metric.Schema{
						Fields: field.Metas{{
							Type: field.HistogramField,
							Name: "11",
						}},
					},
				},
				fields: make(map[field.ID]*aggregation.Aggregator),
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

func TestMetadataLookup_Identifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	db.EXPECT().MetaDB().Return(metaDB).AnyTimes()
	assert.Equal(t, "Metadata Lookup", NewMetadataLookup(nil, db).Identifier())
}

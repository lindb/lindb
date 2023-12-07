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

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/series/metric"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestFieldSuggest_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	db.EXPECT().MetaDB().Return(metaDB).AnyTimes()

	ctx := &context.LeafMetadataContext{
		Database: db,
		Request:  &stmtpkg.MetricMetadata{},
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get metric id failure",
			prepare: func() {
				metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "find schema failure",
			prepare: func() {
				metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), nil)
				metaDB.EXPECT().GetSchema(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "find fields successfully",
			prepare: func() {
				metaDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), nil)
				metaDB.EXPECT().GetSchema(gomock.Any()).Return(&metric.Schema{}, nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := NewFieldSuggest(ctx)
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

func TestFieldSuggest_Identifier(t *testing.T) {
	assert.Equal(t, "Field Suggest", NewFieldSuggest(nil).Identifier())
}

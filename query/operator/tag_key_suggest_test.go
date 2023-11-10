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

	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/series/tag"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestTagKeySuggest_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	meta.EXPECT().MetadataDatabase().Return(metaDB).AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()

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
			name: "tag key suggest failure",
			prepare: func() {
				metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "tag key suggest successfully",
			prepare: func() {
				metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).
					Return(tag.Metas{{}}, nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := NewTagKeySuggest(ctx)
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

func TestTagKeySuggest_Identifier(t *testing.T) {
	assert.Equal(t, "Tag Key Suggest", NewTagKeySuggest(nil).Identifier())
}

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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/query/context"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestTagValueSuggest_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	meta.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()

	ctx := &context.LeafMetadataContext{
		Database: db,
		Request:  &stmtpkg.MetricMetadata{},
	}
	op := NewTagValueSuggest(ctx)
	tagMeta.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"name"})
	assert.NoError(t, op.Execute())
}

func TestTagValueSuggest_Identifier(t *testing.T) {
	assert.Equal(t, "Tag Value Suggest", NewTagValueSuggest(nil).Identifier())
}

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

package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	trackerpkg "github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestSearchFail(t *testing.T) {
	rs, err := MetricMetadataSearch(context.TODO(), &models.ExecuteParam{}, &stmt.MetricMetadata{}, &SearchMgr{})
	assert.Error(t, err)
	assert.Nil(t, rs)
	rs, err = MetricMetadataSearchWithResult(context.TODO(), &models.ExecuteParam{}, &stmt.MetricMetadata{}, &SearchMgr{})
	assert.Error(t, err)
	assert.Nil(t, rs)
	rs, err = MetricDataSearch(context.TODO(), &models.ExecuteParam{}, &stmt.Query{}, &SearchMgr{})
	assert.Error(t, err)
	assert.Nil(t, rs)
}

func TestMetricMetadataSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newExecutePipelineFn = NewExecutePipeline
		ctrl.Finish()
	}()

	pipeline := NewMockPipeline(ctrl)
	newExecutePipelineFn = func(_ *trackerpkg.StageTracker,
		completeCallback func(err error)) Pipeline {
		completeCallback(nil) // just mock invoke
		return pipeline
	}
	pipeline.EXPECT().Execute(gomock.Any())
	taskMgr := NewMockTaskManager(ctrl)
	taskMgr.EXPECT().AddTask(gomock.Any(), gomock.Any())
	taskMgr.EXPECT().RemoveTask(gomock.Any())
	rs, err := MetricMetadataSearchWithResult(context.TODO(), &models.ExecuteParam{Database: "test"}, &stmt.MetricMetadata{}, &SearchMgr{
		RequestID: "xxxx-1bc",
		TaskMgr:   taskMgr,
	})
	assert.NoError(t, err)
	assert.NotNil(t, rs)
}

func TestBuildMetadataResultSet(t *testing.T) {
	rs, err := buildMetadataResultSet(&stmt.MetricMetadata{Type: stmt.Field}, []string{"avc"})
	assert.Error(t, err)
	assert.Nil(t, rs)

	rs, err = buildMetadataResultSet(
		&stmt.MetricMetadata{Type: stmt.Field},
		[]string{
			string(encoding.JSONMarshal(
				&field.Metas{{Name: "f", Type: field.FirstField}, {Name: "1", Type: field.HistogramField}},
			)),
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rs)
}

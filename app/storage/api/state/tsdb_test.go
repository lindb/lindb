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

package state

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/tsdb"
)

func TestTSDBAPI_GetMemoryDatabaseState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	f := tsdb.NewMockDataFamily(ctrl)
	f.EXPECT().Indicator().Return("f")
	f.EXPECT().GetState().Return(models.DataFamilyState{})
	s := tsdb.NewMockShard(ctrl)
	f.EXPECT().Shard().Return(s).AnyTimes()
	db := tsdb.NewMockDatabase(ctrl)
	s.EXPECT().Database().Return(db)
	db.EXPECT().Name().Return("test")
	tsdb.GetFamilyManager().AddFamily(f)

	api := NewTSDBAPI()
	r := gin.New()
	api.Register(r)

	// case 1: params invalid
	resp := mock.DoRequest(t, r, http.MethodGet, MemoryDatabase, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 2: get replica state ok
	resp = mock.DoRequest(t, r, http.MethodGet, MemoryDatabase+"?db=test", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}

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

package monitoring

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/series/tag"
)

func TestExploreAPI_ExploreCurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewExploreAPI(tag.Tags{
		{Key: []byte("role"), Value: []byte(constants.BrokerRole)},
	})
	r := gin.New()
	api.Register(r)
	resp := mock.DoRequest(t, r, http.MethodGet, ExploreCurrentPath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	metric := linmetric.
		NewScope("lindb.ut").
		NewGauge("path")
	metric.Add(1)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreCurrentPath+"?names=lindb.ut&tags[a]=b", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}

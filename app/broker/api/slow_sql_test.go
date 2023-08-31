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

package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	commonlogger "github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

func TestSlowLogMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(SlowSQLLog(
		&deps.HTTPDeps{
			BrokerCfg: &config.Broker{
				BrokerBase: config.BrokerBase{
					SlowSQL: ltoml.Duration(time.Millisecond),
				},
			},
		},
		commonlogger.GetLogger(logger.SlowSQLModule, "SQL"),
	))
	r.GET("/home", func(c *gin.Context) {
		c.Set(constants.CurrentSQL, &models.ExecuteParam{
			SQL: "show databases",
		})
		time.Sleep(time.Millisecond * 10)
		c.JSON(http.StatusOK, "ok")
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.Set(constants.CurrentSQL, &models.ExecuteParam{
			Database: "test",
			SQL:      "show metrics",
		})
		time.Sleep(time.Millisecond * 10)
		c.JSON(http.StatusOK, "ok")
	})
	_ = mock.DoRequest(t, r, http.MethodGet, "/home", `{"database": "db", "sql": "show databases"}`)
	_ = mock.DoRequest(t, r, http.MethodGet, "/metrics", `{"database": "db", "sql": "show databases"}`)
}

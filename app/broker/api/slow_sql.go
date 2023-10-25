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
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

// SlowSQLLog returns show sql log middleware.
func SlowSQLLog(deps *depspkg.HTTPDeps, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			sql, exist := c.Get(constants.CurrentSQL)
			if exist {
				throttle := deps.BrokerCfg.BrokerBase.SlowSQL
				sqlParam := sql.(*models.ExecuteParam)
				end := time.Now()
				duration := end.Sub(start)
				if int64(duration) >= int64(throttle) {
					sqlInfo := fmt.Sprintf("# Time: %s \n%s# Execute time: %s\n%s",
						timeutil.FormatTimestamp(start.UnixMilli(), timeutil.DataTimeFormat2),
						getDatabaseName(sqlParam),
						duration.String(),
						sqlParam.SQL,
					)
					log.Error(sqlInfo)
				}
			}
		}()
		c.Next()
	}
}

func getDatabaseName(param *models.ExecuteParam) string {
	if param.Database == "" {
		return ""
	}
	return fmt.Sprintf("# Database: %s\n", param.Database)
}

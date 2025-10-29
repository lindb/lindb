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

package ingest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/log"
	"github.com/lindb/common/models"
	"github.com/lindb/common/pkg/http"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

const Path = "/logs"

type Log struct {
	deps *depspkg.HTTPDeps

	statistics struct {
		flat   *linmetric.BoundHistogram
		proto  *linmetric.BoundHistogram
		influx *linmetric.BoundHistogram
	}
}

func NewLog(deps *depspkg.HTTPDeps) *Log {
	ingestStatistics := metrics.NewCommonIngestionStatistics()

	return &Log{
		deps: deps,
		statistics: struct {
			flat   *linmetric.BoundHistogram
			proto  *linmetric.BoundHistogram
			influx *linmetric.BoundHistogram
		}{
			flat:   ingestStatistics.Duration.WithTagValues("flat"),
			proto:  ingestStatistics.Duration.WithTagValues("proto"),
			influx: ingestStatistics.Duration.WithTagValues("influx"),
		},
	}
}

// Register adds the log ingest url route.
func (w *Log) Register(route gin.IRoutes) {
	route.POST(Path, w.Write)
	route.PUT(Path, w.Write)
}

func (w *Log) Write(c *gin.Context) {
	if err := w.write(c); err != nil {
		http.Error(c, err)
	}
}

func (w *Log) write(c *gin.Context) error {
	var logs []models.Log
	fmt.Println("write log...")
	if err := c.ShouldBind(&logs); err != nil {
		return err
	}
	fmt.Println(logs)
	rb := log.CreateRowBuilder()
	for _, l := range logs {
		rb.AddMessage([]byte(l.Message)).
			AddTimestamp(l.Timestamp)
		for k, v := range l.Fields {
			rb.AddField([]byte(k), []byte(v))
		}
		data, _ := rb.Build()
		fmt.Println(string(data))
		if err := w.deps.CM.WriteMsg(c.Request.Context(), "log_test", data); err != nil {
			return err
		}
		rb.Reset()
	}
	return nil
}

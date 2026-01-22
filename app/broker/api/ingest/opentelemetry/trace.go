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

package opentelemetry

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/http"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

const TracePath = "/opentelemetry/traces"

type Trace struct {
	deps *depspkg.HTTPDeps

	statistics struct {
		flat   *linmetric.BoundHistogram
		proto  *linmetric.BoundHistogram
		influx *linmetric.BoundHistogram
	}
}

func NewTrace(deps *depspkg.HTTPDeps) *Trace {
	ingestStatistics := metrics.NewCommonIngestionStatistics()

	return &Trace{
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
func (w *Trace) Register(route gin.IRoutes) {
	route.POST(TracePath, w.Write)
	route.PUT(TracePath, w.Write)
}

func (w *Trace) Write(c *gin.Context) {
	if err := w.write(c); err != nil {
		fmt.Println(err)
		http.Error(c, err)
	}
}

func (w *Trace) write(c *gin.Context) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(body); err != nil {
		return err
	}
	// TODO: set database name
	if err := w.deps.CM.WriteMsg(c.Request.Context(), "trace_test", body); err != nil {
		return err
	}

	return nil
}

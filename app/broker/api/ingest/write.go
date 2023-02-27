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
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"

	commonconstants "github.com/lindb/common/constants"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/ingestion/flat"
	"github.com/lindb/lindb/ingestion/influx"
	"github.com/lindb/lindb/ingestion/proto"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/series/metric"
)

var (
	// WritePath represents write http api router path.
	WritePath = "/write"
)

// Write represents write api that processes flat/proto/influx protocol data.
type Write struct {
	deps *depspkg.HTTPDeps

	statistics struct {
		flat   *linmetric.BoundHistogram
		proto  *linmetric.BoundHistogram
		influx *linmetric.BoundHistogram
	}
}

// NewWrite creates a writer instance.
func NewWrite(deps *depspkg.HTTPDeps) *Write {
	ingestStatistics := metrics.NewCommonIngestionStatistics()
	return &Write{
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

// Register adds the writer url route.
func (w *Write) Register(route gin.IRoutes) {
	route.POST(WritePath, w.Write)
	route.PUT(WritePath, w.Write)
}

// Write processes flat/proto/influx protocol data with ingest limit.
//
// @BasePath /api/v1
// @Summary write metric data
// @Schemes
// @Description receive metric data, then parse the data based on content type(flat buffer/proto buffer/influx).
// @Description write data via database channel, support content-type as below:
// @Description 1. application/flatbuffer
// @Description 2. application/protobuf
// @Description 3. application/influx
// @Tags Write
// @Accept application/flatbuffer
// @Accept application/protobuf
// @Accept application/influx
// @Param db query string true "database name"
// @Param ns query string false "namespace, default value: default-ns"
// @Param string body string ture "metric data"
// @Produce plain
// @Success 204 {string} string ""
// @Failure 500 {string} string "internal error"
// @Router /write [put]
// @Router /write [post]
func (w *Write) Write(c *gin.Context) {
	if err := w.deps.IngestLimiter.Do(func() error {
		return w.write(c)
	}); err != nil {
		http.Error(c, err)
	} else {
		http.NoContent(c)
	}
}

// parse flat/proto/influx protocol data, then write parsed data to database's write channel.
func (w *Write) write(c *gin.Context) (err error) {
	var param struct {
		Database  string `form:"db" binding:"required"`
		Namespace string `form:"ns"`
	}
	err = c.ShouldBindQuery(&param)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		w.deps.BrokerCfg.BrokerBase.Ingestion.IngestTimeout.Duration())
	defer cancel()

	if param.Namespace == "" {
		param.Namespace = commonconstants.DefaultNamespace
	}
	enrichedTags, err := ingestCommon.ExtractEnrichTags(c.Request)
	if err != nil {
		return err
	}

	limits := w.deps.StateMgr.GetDatabaseLimits(param.Database)
	for _, tag := range enrichedTags {
		if len(tag.Key) > limits.MaxTagNameLength {
			return constants.ErrTagKeyTooLong
		}
		if len(tag.Value) > limits.MaxTagValueLength {
			return constants.ErrTagValueTooLong
		}
	}
	if len(param.Namespace) > limits.MaxNamespaceLength {
		return constants.ErrNamespaceTooLong
	}
	contentType := strings.ToLower(strings.Trim(c.Request.Header.Get(headers.ContentType), " "))
	var rows *metric.BrokerBatchRows
	switch {
	case strings.HasPrefix(contentType, constants.ContentTypeFlat):
		rows, err = flat.Parse(c.Request, enrichedTags, param.Namespace, limits)
	case strings.HasPrefix(contentType, constants.ContentTypeInflux):
		rows, err = influx.Parse(c.Request, enrichedTags, param.Namespace, limits)
	case strings.HasPrefix(contentType, constants.ContentTypeProto):
		rows, err = proto.Parse(c.Request, enrichedTags, param.Namespace, limits)
	default:
		err = fmt.Errorf("not support content type: %s, only support %s/%s/%s", contentType,
			constants.ContentTypeFlat, constants.ContentTypeProto, constants.ContentTypeInflux)
	}
	if err != nil {
		return err
	}
	if err := w.deps.CM.Write(ctx, param.Database, rows); err != nil {
		return err
	}
	return nil
}

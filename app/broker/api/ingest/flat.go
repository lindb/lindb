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
	"time"

	"github.com/gin-gonic/gin"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/ingestion/flat"
	"github.com/lindb/lindb/internal/linmetric"
)

var (
	FlatWritePath = "/flat/write"
)

var (
	HTTPHandlerTimerVec = linmetric.BrokerRegistry.
		NewScope(
			"lindb.http_server.handle_duration",
		).
		NewHistogramVec("path").
		WithExponentBuckets(time.Millisecond, time.Second*5, 20)
)

func WithHistogram(histogram *linmetric.BoundHistogram) gin.HandlerFunc {
	// TODO need move??
	return func(c *gin.Context) {
		start := time.Now()
		defer histogram.UpdateSince(start)
		c.Next()
	}
}

// FlatWriter processes native proto metrics.
type FlatWriter struct {
	commonWriter
}

// NewFlatWriter creates native proto metrics writer
func NewFlatWriter(deps *depspkg.HTTPDeps) *FlatWriter {
	return &FlatWriter{
		commonWriter: commonWriter{
			deps:   deps,
			parser: flat.Parse,
		},
	}
}

// Register adds native write url route.
func (nw *FlatWriter) Register(route gin.IRoutes) {
	route.POST(
		FlatWritePath,
		WithHistogram(HTTPHandlerTimerVec.WithTagValues(FlatWritePath)),
		nw.Write,
	)
	route.PUT(
		FlatWritePath,
		WithHistogram(HTTPHandlerTimerVec.WithTagValues(FlatWritePath)),
		nw.Write,
	)
}

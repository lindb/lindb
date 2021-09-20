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
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/app/broker/middleware"
	"github.com/lindb/lindb/ingestion/flat"
)

var (
	FlatWritePath = "/flat/write"
)

// FlatWriter processes native proto metrics.
type FlatWriter struct {
	commonWriter
}

// NewFlatWriter creates native proto metrics writer
func NewFlatWriter(deps *deps.HTTPDeps) *FlatWriter {
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
		middleware.WithHistogram(middleware.HttHandlerTimerVec.WithTagValues(FlatWritePath)),
		nw.Write,
	)
	route.PUT(
		FlatWritePath,
		middleware.WithHistogram(middleware.HttHandlerTimerVec.WithTagValues(FlatWritePath)),
		nw.Write,
	)
}

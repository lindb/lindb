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

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/ingestion/proto"
)

var (
	ProtoWritePath = "/proto/write"
)

// ProtoWriter processes native proto metrics.
type ProtoWriter struct {
	commonWriter
}

// NewProtoWriter creates native proto metrics writer
func NewProtoWriter(deps *depspkg.HTTPDeps) *ProtoWriter {
	return &ProtoWriter{
		commonWriter: commonWriter{
			deps:   deps,
			parser: proto.Parse,
		},
	}
}

// Register adds native write url route.
func (nw *ProtoWriter) Register(route gin.IRoutes) {
	route.POST(
		ProtoWritePath,
		WithHistogram(ingestStatistics.Duration.WithTagValues(ProtoWritePath)),
		nw.Write,
	)
	route.PUT(
		ProtoWritePath,
		WithHistogram(ingestStatistics.Duration.WithTagValues(ProtoWritePath)),
		nw.Write,
	)
}

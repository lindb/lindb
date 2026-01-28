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
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/log"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"

	depspkg "github.com/lindb/lindb/app/broker/deps"
)

const LogPath = "/opentelemetry/logs"

// Log represents the OpenTelemetry log ingest api handler.
type Log struct {
	writer
}

// NewLog creates a new OpenTelemetry log ingest api handler.
func NewLog(deps *depspkg.HTTPDeps) *Log {
	l := &Log{
		writer: writer{
			deps: deps,
		},
	}
	l.writer.processProto = l.processProto
	return l
}

// Register adds the log ingest url route.
func (w *Log) Register(route gin.IRoutes) {
	route.POST(LogPath, w.Write)
	route.PUT(LogPath, w.Write)
}

// processProto processes the OpenTelemetry log proto data.
func (w *Log) processProto(ctx context.Context, database string, data []byte) error {
	req := plogotlp.NewExportRequest()
	if err := req.UnmarshalProto(data); err != nil {
		return err
	}
	// TODO:
	rb := log.CreateRowBuilder()

	logs := req.Logs()

	rLogs := logs.ResourceLogs()
	for i := range rLogs.Len() {
		log := rLogs.At(i)
		attr := log.Resource().Attributes()
		scopeLogs := log.ScopeLogs()
		for j := range scopeLogs.Len() {
			sl := scopeLogs.At(j)
			lrs := sl.LogRecords()
			for k := range lrs.Len() {
				lr := lrs.At(k)
				rb.AddMessage([]byte(lr.Body().AsString())).
					AddTimestamp(lr.Timestamp().AsTime().UnixMilli())
				rb.AddField([]byte("level"), []byte(lr.SeverityText()))
				attr.Range(func(k string, v pcommon.Value) bool {
					if v.AsString() != "" {
						rb.AddField([]byte(k), []byte(v.AsString()))
					}
					return true
				})
				lr.Attributes().Range(func(k string, v pcommon.Value) bool {
					if v.AsString() != "" {
						rb.AddField([]byte(k), []byte(v.AsString()))
					}
					return true
				})

				data, _ := rb.Build()
				dData := make([]byte, len(data))
				copy(dData, data)
				if err := w.deps.CM.WriteMsg(ctx, database, dData); err != nil {
					return err
				}
				rb.Reset()
			}
		}
	}

	return nil
}

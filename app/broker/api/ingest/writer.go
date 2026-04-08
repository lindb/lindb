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
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/http"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
)

// ProcessFunc processes ingest data for a specific content-type.
type ProcessFunc func(ctx context.Context, database string, data []byte) error

// BaseWriter is the shared base handler for all Arrow/OTLP ingest endpoints.
// It centralizes: IngestLimiter throttling, database header extraction,
// request body reading, timeout context construction, and content-type dispatch.
type BaseWriter struct {
	deps       *depspkg.HTTPDeps
	processors map[string]ProcessFunc
}

// NewBaseWriter creates a BaseWriter with the given content-type processor map.
func NewBaseWriter(deps *depspkg.HTTPDeps, processors map[string]ProcessFunc) BaseWriter {
	return BaseWriter{deps: deps, processors: processors}
}

// Write is the gin handler; wraps write() with IngestLimiter rate limiting.
func (w *BaseWriter) Write(c *gin.Context) {
	if err := w.deps.IngestLimiter.Do(func() error {
		return w.write(c)
	}); err != nil {
		http.Error(c, err)
	} else {
		http.NoContent(c)
	}
}

// write reads the request, dispatches to the matching ProcessFunc by Content-Type,
// and returns any error for the caller (Write) to handle.
func (w *BaseWriter) write(c *gin.Context) error {
	database := c.Request.Header.Get(constants.DatabaseHeader)
	if database == "" {
		return errors.New("database header not found[X-LinDB-Database]")
	}

	contentType := c.Request.Header.Get("Content-Type")
	process, ok := w.processors[contentType]
	if !ok {
		return fmt.Errorf("unsupported media type: %s", contentType)
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(),
		w.deps.BrokerCfg.BrokerBase.Ingestion.IngestTimeout.Duration())
	defer cancel()

	return process(ctx, database, body)
}

// StandardProcess returns a ProcessFunc that forwards data to the WriteManager.
// Trace, log, and metric handlers share identical dispatch logic:
// look up the db writer and call Write with the given encoding.
func StandardProcess(deps *depspkg.HTTPDeps, encoding constants.EncodingType) ProcessFunc {
	return func(ctx context.Context, database string, data []byte) error {
		w, ok := deps.WriteManager.GetWriter(database)
		if !ok {
			return constants.ErrDatabaseNotFound
		}
		return w.Write(ctx, data, encoding)
	}
}

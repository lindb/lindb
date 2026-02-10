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
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/http"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
)

// writer represents OpenTelemetry data write handler.
type writer struct {
	deps *depspkg.HTTPDeps

	// need set for processing protobuf data(trace/log etc.)
	processProto func(ctx context.Context, database string, data []byte) error
}

// Write handles OpenTelemetry data write request.
func (w *writer) Write(c *gin.Context) {
	if err := w.deps.IngestLimiter.Do(func() error {
		return w.write(c)
	}); err != nil {
		http.Error(c, err)
	} else {
		http.NoContent(c)
	}
}

// write processes the OpenTelemetry data based on Content-Type.
func (w *writer) write(c *gin.Context) error {
	database := c.Request.Header.Get(constants.DatabaseHeader)
	if database == "" {
		return errors.New("database header not found[X-LinDB-Database]")
	}
	// 1. Extract the Content-Type header
	contentType := c.Request.Header.Get("Content-Type")
	// 2. Determine the encoding logic
	switch contentType {
	case "application/x-protobuf":
		// TODO: add compression logic
		// Handle binary protobuf
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(),
			w.deps.BrokerCfg.BrokerBase.Ingestion.IngestTimeout.Duration())

		defer func() {
			c.Request.Body.Close()
			cancel()
		}()

		return w.processProto(ctx, database, body)
	default:
		return errors.New("unsupported media type: expected application/x-protobuf")
	}
}

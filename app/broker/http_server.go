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

package broker

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/felixge/fgprof"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb"
	"github.com/lindb/lindb/app/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/conntrack"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
)

const _apiRootPath = "/api"

// HTTPServer represents http server with gin framework.
type HTTPServer struct {
	addr   string
	server http.Server
	gin    *gin.Engine

	logger *logger.Logger
}

// NewHTTPServer creates http server.
func NewHTTPServer(cfg config.HTTP) *HTTPServer {
	s := &HTTPServer{
		addr: fmt.Sprintf(":%d", cfg.Port),
		gin:  gin.New(),
		server: http.Server{
			// use extra timeout for ingestion and query timeout
			WriteTimeout: cfg.WriteTimeout.Duration(),
			ReadTimeout:  cfg.ReadTimeout.Duration(),
			IdleTimeout:  cfg.IdleTimeout.Duration(),
		},
		logger: logger.GetLogger("broker", "HTTPServer"),
	}
	s.init()
	return s
}

// init initializes http server default router/handle/middleware.
func (s *HTTPServer) init() {
	// Using middlewares on group.
	// use AccessLogMiddleware to log panic error with zap
	s.gin.Use(middleware.AccessLogMiddleware())
	s.gin.Use(cors.Default())

	// self monitoring handler
	s.gin.GET(constants.HealthPath, monitoring.HealthHandler)

	if logger.IsDebug() {
		s.logger.Info("/debug/pprof is enabled")
		pprof.Register(s.gin)
		s.logger.Info("/debug/fgprof is enabled")
		s.gin.GET("/debug/fgprof", gin.WrapH(fgprof.Handler()))
	}
	// server static file
	staticFS, err := fs.Sub(lindb.StaticContent, "web/static")
	staticHome := "/console"
	if err != nil {
		s.logger.Error("cannot find static resource", logger.Error(err))
	} else {
		s.gin.StaticFS(staticHome, http.FS(staticFS))
		// redirects to admin console
		s.gin.GET("/", func(c *gin.Context) {
			c.Request.URL.Path = staticHome
			s.gin.HandleContext(c)
		})
	}
}

// GetAPIRouter returns api router.
func (s *HTTPServer) GetAPIRouter() *gin.RouterGroup {
	return s.gin.Group(_apiRootPath)
}

// Run runs the HTTP server.
func (s *HTTPServer) Run() error {
	s.logger.Info("starting http server", logger.String("addr", s.server.Addr))
	s.server.Handler = s.gin
	// Open listener.
	trackedListener, err := conntrack.NewTrackedListener("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.server.Serve(trackedListener)
}

// Close closes the server.
func (s *HTTPServer) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

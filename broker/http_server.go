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
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	promreporter "github.com/uber-go/tally/prometheus"

	"github.com/lindb/lindb"
	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
)

const _apiRootPath = "/api/v1"

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
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60, //TODO add config?
		},
		logger: logger.GetLogger("broker", "HTTPServer"),
	}
	s.init()
	return s
}

// init initializes http server default router/handle/middleware.
func (s *HTTPServer) init() {
	// Using middlewares on group.
	s.gin.Use(middleware.AccessLogMiddleware())
	s.gin.Use(cors.Default())
	s.gin.Use(gin.Recovery())

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

	// add prometheus metric report
	reporter := promreporter.NewReporter(promreporter.Options{})
	h := reporter.HTTPHandler()
	s.gin.GET("/metrics", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})
}

// GetAPIRouter returns api router.
func (s *HTTPServer) GetAPIRouter() *gin.RouterGroup {
	return s.gin.Group(_apiRootPath)
}

// Run runs the HTTP server.
func (s *HTTPServer) Run() error {
	s.logger.Info("starting http server", logger.Any("addr", s.server.Addr))
	s.server.Handler = s.gin
	// Open listener.
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.server.Serve(ln)
}

// Close closes the server.
func (s *HTTPServer) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

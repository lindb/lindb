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

package http

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/felixge/fgprof"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/lindb/common/pkg/http/middleware"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	_ "github.com/lindb/lindb/docs"
	"github.com/lindb/lindb/internal/conntrack"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/hostutil"
)

//go:generate mockgen -source ./http_server.go -destination=./http_server_mock.go -package=http

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// Server represents http server with gin framework.
type Server interface {
	// GetAPIRouter returns api router.
	GetAPIRouter() *gin.RouterGroup
	// Run runs the HTTP server.
	Run() error
	// Close closes the server.
	Close(ctx context.Context) error
}

// server implements Server interface.
type server struct {
	addr           string
	server         http.Server
	gin            *gin.Engine
	staticResource bool
	cfg            config.HTTP

	r      *linmetric.Registry
	logger logger.Logger
}

// NewServer creates http server.
func NewServer(cfg config.HTTP, staticResource bool, r *linmetric.Registry) Server {
	s := &server{
		cfg:            cfg,
		addr:           fmt.Sprintf(":%d", cfg.Port),
		gin:            gin.New(),
		staticResource: staticResource,
		server: http.Server{
			// use extra timeout for ingestion and query timeout
			// if write timeout will return ERR_EMPTY_RESPONSE, chrome will does auto retry.
			// https://www.bennadel.com/blog/3257-google-chrome-will-automatically-retry-requests-on-certain-error-responses.htm
			// https://mariocarrion.com/2021/09/17/golang-software-architecture-resilience-http-servers.html
			// WriteTimeout: cfg.WriteTimeout.Duration(),
			ReadTimeout: cfg.ReadTimeout.Duration(),
			IdleTimeout: cfg.IdleTimeout.Duration(),
		},
		r:      r,
		logger: logger.GetLogger("HTTP", "Server"),
	}
	s.init()
	return s
}

// init initializes http server default router/handle/middleware.
func (s *server) init() {
	// Using middlewares on group.
	s.gin.Use(middleware.Recovery())
	// use AccessLog to log panic error with zap
	s.gin.Use(middleware.AccessLog(logger.GetLogger(logger.AccessLogModule, "HTTP")))
	s.gin.Use(cors.Default())

	if config.Profile {
		s.logger.Info("/debug/pprof is enabled")
		pprof.Register(s.gin)
		s.logger.Info("/debug/fgprof is enabled")
		s.gin.GET("/debug/fgprof", gin.WrapH(fgprof.Handler()))
	}

	if config.Doc {
		// swagger-ui: http://localhost:port/swagger/index.html
		ip, _ := hostutil.GetHostIP()
		s.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", ip, s.cfg.Port)),
			ginSwagger.DefaultModelsExpandDepth(-1)))
	}
	if s.staticResource {
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
}

// GetAPIRouter returns api router.
func (s *server) GetAPIRouter() *gin.RouterGroup {
	return s.gin.Group(constants.APIRoot)
}

// Run runs the HTTP server.
func (s *server) Run() error {
	s.logger.Info("starting http server", logger.String("addr", s.server.Addr))
	s.server.Handler = s.gin
	// Open listener.
	trackedListener, err := conntrack.NewTrackedListener("tcp", s.addr, s.r)
	if err != nil {
		return err
	}
	return s.server.Serve(trackedListener)
}

// Close closes the server.
func (s *server) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

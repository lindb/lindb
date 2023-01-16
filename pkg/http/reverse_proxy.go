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
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

var (
	ProxyPath = "/proxy"
)

// ProxyParam represents proxy request params.
type ProxyParam struct {
	Target string `form:"target" binding:"required"`
	Path   string `form:"path" binding:"required"`
}

// ReverseProxy represents the http reverse proxy to target's api.
type ReverseProxy struct {
}

// NewReverseProxy creates a ReverseProxy instance.
func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{}
}

// Register adds proxy url route.
func (p *ReverseProxy) Register(route gin.IRoutes) {
	route.GET(ProxyPath, p.Proxy)
}

// Proxy forwards to target server api by given target ip and path.
//
// @Summary reverse proxy
// @Description Forward request to target server by given target ip and path.
// @Tags State
// @Accept json
// @Param param body models.ProxyParam ture "param data"
// @Produce json
// @Success 200 {object} object
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "validate failure"
// @Failure 500 {string} string "internal error"
// @Router /proxy [get]
func (p *ReverseProxy) Proxy(c *gin.Context) {
	var param ProxyParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		Error(c, err)
		return
	}
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = param.Target
		req.URL.Path = param.Path
		req.URL.RawQuery = c.Request.URL.RawQuery
	}
	proxy := &httputil.ReverseProxy{
		Director:  director,
		Transport: &transport{},
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

// transport implements http.RoundTripper.
type transport struct {
}

// RoundTrip removes cors http header.
func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	// if target server enable cors, maybe add duplicate Access-Control-Allow-Origin header.
	resp.Header.Del("Access-Control-Allow-Origin")
	return resp, err
}

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
	netHTTP "net/http"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

type parserFunc func(req *netHTTP.Request, enrichedTags tag.Tags, namespace string) (*metric.BrokerBatchRows, error)

type commonWriter struct {
	deps   *deps.HTTPDeps
	parser parserFunc
}

func (cw *commonWriter) Write(c *gin.Context) {
	if err := cw.deps.IngestLimiter.Do(func() error {
		return cw.realWrite(c)
	}); err != nil {
		http.Error(c, err)
	} else {
		http.NoContent(c)
	}
}

func (cw *commonWriter) realWrite(c *gin.Context) error {
	var param struct {
		Database  string `form:"db" binding:"required"`
		Namespace string `form:"ns"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		cw.deps.BrokerCfg.BrokerBase.Ingestion.IngestTimeout.Duration())
	defer cancel()

	if param.Namespace == "" {
		param.Namespace = constants.DefaultNamespace
	}
	enrichedTags, err := ingestCommon.ExtractEnrichTags(c.Request)
	if err != nil {
		return err
	}
	metrics, err := cw.parser(c.Request, enrichedTags, param.Namespace)
	if err != nil {
		return err
	}
	if err := cw.deps.CM.Write(ctx, param.Database, metrics); err != nil {
		return err
	}
	return nil
}

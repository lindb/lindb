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

package api

import (
	"testing"

	"github.com/gin-gonic/gin"

	prometheusIngest "github.com/lindb/lindb/app/broker/api/prometheus/ingest"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
)

func TestNewRouter(t *testing.T) {
	brokerCfg := &config.Broker{}
	httpDeps := &deps.HTTPDeps{BrokerCfg: brokerCfg}
	schema := prometheusIngest.DatabaseConfig{
		Namespace: brokerCfg.Prometheus.Namespace,
		Database:  brokerCfg.Prometheus.Database,
		Field:     brokerCfg.Prometheus.Field,
	}
	prometheusWriter := prometheusIngest.NewWriter(httpDeps, prometheusIngest.DefaultWriteOptions(schema))
	r := NewAPI(httpDeps, prometheusWriter)
	r.RegisterRouter(gin.New().Group(constants.APIRoot))
}

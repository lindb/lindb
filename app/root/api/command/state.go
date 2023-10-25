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

package command

import (
	"context"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var (
	metricCli = client.NewMetricCli()
)

// StateCommand executes the state query.
func StateCommand(_ context.Context, deps *depspkg.HTTPDeps,
	_ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	stateStmt := stmt.(*stmtpkg.State)
	switch stateStmt.Type {
	case stmtpkg.RootAlive:
		return deps.StateMgr.GetLiveNodes(), nil
	case stmtpkg.BrokerAlive:
		return deps.StateMgr.GetBrokerStates(), nil
	case stmtpkg.RootMetric:
		liveNodes := deps.StateMgr.GetLiveNodes()
		var nodes []models.Node
		for idx := range liveNodes {
			nodes = append(nodes, &liveNodes[idx])
		}
		return metricCli.FetchMetricData(nodes, stateStmt.MetricNames)
	default:
		return nil, nil
	}
}

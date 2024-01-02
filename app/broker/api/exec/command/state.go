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
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/lindb/common/pkg/logger"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
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
	case stmtpkg.Master:
		return deps.Master.GetMaster(), nil
	case stmtpkg.BrokerAlive:
		return deps.StateMgr.GetLiveNodes(), nil
	case stmtpkg.StorageAlive:
		return deps.StateMgr.GetStorage(), nil
	case stmtpkg.Replication:
		return getStateFromStorage(deps, stateStmt, "/state/replica", func() interface{} {
			var state []models.FamilyLogReplicaState
			return &state
		})
	case stmtpkg.MemoryDatabase:
		return getStateFromStorage(deps, stateStmt, "/state/tsdb/memory", func() interface{} {
			var state []models.DataFamilyState
			return &state
		})
	case stmtpkg.BrokerMetric:
		liveNodes := deps.StateMgr.GetLiveNodes()
		var nodes []models.Node
		for idx := range liveNodes {
			nodes = append(nodes, &liveNodes[idx])
		}
		return metricCli.FetchMetricData(nodes, stateStmt.MetricNames)
	case stmtpkg.StorageMetric:
		storage := deps.StateMgr.GetStorage()
		liveNodes := storage.LiveNodes
		var nodes []models.Node
		for id := range liveNodes {
			n := liveNodes[id]
			nodes = append(nodes, &n)
		}
		return metricCli.FetchMetricData(nodes, stateStmt.MetricNames)
	default:
		return nil, nil
	}
}

// getStateFromStorage returns the state from storage cluster.
func getStateFromStorage(deps *depspkg.HTTPDeps, stmt *stmtpkg.State, path string, newStateFn func() interface{}) (interface{}, error) {
	storage := deps.StateMgr.GetStorage()
	liveNodes := storage.LiveNodes
	var nodes []models.Node
	for id := range liveNodes {
		n := liveNodes[id]
		nodes = append(nodes, &n)
	}
	return fetchStateData(nodes, stmt, path, newStateFn)
}

// fetchStateData fetches the state metric from each live node.
func fetchStateData(nodes []models.Node, stmt *stmtpkg.State, path string, newStateFn func() interface{}) (interface{}, error) {
	size := len(nodes)
	if size == 0 {
		return nil, nil
	}
	result := make([]interface{}, size)
	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			address := node.HTTPAddress()
			state := newStateFn()
			_, err := resty.New().R().SetQueryParams(map[string]string{"db": stmt.Database}).
				SetHeader("Accept", "application/json").
				SetResult(&state).
				Get(address + constants.APIVersion1CliPath + path)
			if err != nil {
				log.Error("get state from storage node", logger.String("url", address), logger.Error(err))
				return
			}
			result[i] = state
		}()
	}
	wait.Wait()
	rs := make(map[string]interface{})
	for idx := range nodes {
		rs[nodes[idx].Indicator()] = result[idx]
	}
	return rs, nil
}

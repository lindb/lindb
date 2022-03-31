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
	"net/url"
	"strings"
	"sync"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
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
		return deps.StateMgr.GetStorageList(), nil
	case stmtpkg.Replication:
		return getReplicaState(deps, stateStmt)
	case stmtpkg.BrokerMetric:
		liveNodes := deps.StateMgr.GetLiveNodes()
		var nodes []models.Node
		for idx := range liveNodes {
			nodes = append(nodes, &liveNodes[idx])
		}
		return fetchMetricData(nodes, stateStmt.MetricNames)
	case stmtpkg.StorageMetric:
		storageName := strings.TrimSpace(stateStmt.StorageName)
		if storageName == "" {
			return nil, constants.ErrStorageNameRequired
		}
		storage, ok := deps.StateMgr.GetStorage(storageName)
		if !ok {
			return nil, nil
		}
		liveNodes := storage.LiveNodes
		var nodes []models.Node
		for id := range liveNodes {
			n := liveNodes[id]
			nodes = append(nodes, &n)
		}
		return fetchMetricData(nodes, stateStmt.MetricNames)
	default:
		return nil, nil
	}
}

// getReplicaState returns wal replica state.
func getReplicaState(deps *depspkg.HTTPDeps, stmt *stmtpkg.State) (interface{}, error) {
	storage, ok := deps.StateMgr.GetStorage(stmt.StorageName)
	if !ok {
		return nil, nil
	}
	liveNodes := storage.LiveNodes
	var nodes []models.Node
	for id := range liveNodes {
		n := liveNodes[id]
		nodes = append(nodes, &n)
	}
	return fetchStateData(nodes, stmt)
}

// fetchStateData fetches the state metric from each live nodes.
func fetchStateData(nodes []models.Node, stmt *stmtpkg.State) (interface{}, error) {
	size := len(nodes)
	if size == 0 {
		return nil, nil
	}
	result := make([][]models.FamilyLogReplicaState, size)
	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			address := node.HTTPAddress()
			var state []models.FamilyLogReplicaState
			_, err := NewRestyFn().R().SetQueryParams(map[string]string{"db": stmt.Database}).
				SetHeader("Accept", "application/json").
				SetResult(&state).
				Get(address + "/api/state/replica")
			if err != nil {
				log.Error("get replication state from storage node", logger.String("url", address), logger.Error(err))
				return
			}
			result[i] = state
		}()
	}
	wait.Wait()
	rs := make(map[string][]models.FamilyLogReplicaState)
	for idx := range nodes {
		rs[nodes[idx].Indicator()] = result[idx]
	}
	return rs, nil
}

// fetchMetricData fetches the state metric from each live nodes.
func fetchMetricData(nodes []models.Node, names []string) (interface{}, error) {
	size := len(nodes)
	if size == 0 {
		return nil, nil
	}
	result := make([]map[string][]*models.StateMetric, size)
	params := make(url.Values)
	for _, name := range names {
		params.Add("names", name)
	}

	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			address := node.HTTPAddress()
			metric := make(map[string][]*models.StateMetric)
			_, err := NewRestyFn().R().SetQueryParamsFromValues(params).
				SetHeader("Accept", "application/json").
				SetResult(&metric).
				Get(address + "/api/state/explore/current")
			if err != nil {
				log.Error("get current metric state from alive node", logger.String("url", address), logger.Error(err))
				return
			}
			result[i] = metric
		}()
	}
	wait.Wait()
	rs := make(map[string][]*models.StateMetric)
	for _, metricList := range result {
		if metricList == nil {
			continue
		}
		for name, list := range metricList {
			l, ok := rs[name]
			if ok {
				l = append(l, list...)
				rs[name] = l
			} else {
				rs[name] = list
			}
		}
	}
	return rs, nil
}

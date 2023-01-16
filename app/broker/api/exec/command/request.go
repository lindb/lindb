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
	"sort"
	"sync"

	"github.com/go-resty/resty/v2"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// RequestCommand executes requests/request related statement.
func RequestCommand(_ context.Context, deps *depspkg.HTTPDeps, _ *models.ExecuteParam, _ stmtpkg.Statement) (interface{}, error) {
	liveNodes := deps.StateMgr.GetLiveNodes()
	var nodes []models.Node
	for idx := range liveNodes {
		nodes = append(nodes, &liveNodes[idx])
	}
	size := len(nodes)
	if size == 0 {
		return nil, nil
	}
	result := make(map[string][]*models.Request)
	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			address := node.HTTPAddress()
			var stats []*models.Request
			_, err := resty.New().R().
				SetHeader("Accept", "application/json").
				SetResult(&stats).
				Get(address + constants.APIVersion1CliPath + "/state/requests")
			if err != nil {
				log.Error("get current alive reuqests from alive node", logger.String("url", address), logger.Error(err))
				return
			}
			result[node.Indicator()] = stats
		}()
	}
	wait.Wait()

	// build result set sort by request start time
	var rs []*models.Request
	for k, v := range result {
		for _, req := range v {
			req.Entry = k
			rs = append(rs, req)
		}
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Start < rs[j].Start
	})
	return rs, nil
}

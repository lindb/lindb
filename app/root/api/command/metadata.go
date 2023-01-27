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
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	// NewRestyFn represents new resty client.
	NewStateMachineCliFn = client.NewStateMachineCli
)

var brokerRole = map[string]models.StateMachineInfo{
	constants.BrokerState: {
		Path: constants.BrokerState,
		CreateState: func() interface{} {
			return &models.BrokerState{}
		},
	},
}

// MetadataCommand executes the metadata query.
func MetadataCommand(ctx context.Context, deps *depspkg.HTTPDeps,
	_ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	metadataStmt := stmt.(*stmtpkg.Metadata)
	if metadataStmt.MetadataType == stmtpkg.MetadataTypes {
		// returns metadata explore define info.
		return map[string]interface{}{
			constants.RootRole:   root.StateMachinePaths,
			constants.BrokerRole: brokerRole,
		}, nil
	}

	// explore metadata
	switch metadataStmt.Source {
	case stmtpkg.StateRepoSource:
		return exploreStateRepoData(ctx, deps, metadataStmt)
	case stmtpkg.StateMachineSource:
		return exploreStateMachineDate(metadataStmt, deps)
	}
	return nil, nil
}

// exploreStateRepoData explores state data from repo.
func exploreStateRepoData(ctx context.Context, deps *depspkg.HTTPDeps,
	metadataStmt *stmtpkg.Metadata) (interface{}, error) {
	stateMachineInfo, ok := root.StateMachinePaths[metadataStmt.Type]
	if !ok {
		return nil, nil
	}
	return discovery.ExploreData(ctx, deps.Repo, stateMachineInfo)
}

// exploreStateMachineDate explores the state from state machine of broker/master/storage.
func exploreStateMachineDate(metadataStmt *stmtpkg.Metadata, deps *depspkg.HTTPDeps) (interface{}, error) {
	param := map[string]string{
		"type":       metadataStmt.Type,
		"brokerName": metadataStmt.ClusterName,
	}
	var nodes []models.Node
	statelessNodes := deps.StateMgr.GetLiveNodes()
	for idx := range statelessNodes {
		nodes = append(nodes, &statelessNodes[idx])
	}
	// forward to live nodes
	cli := NewStateMachineCliFn()
	return cli.FetchStateByNodes(param, nodes), nil
}

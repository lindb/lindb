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
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	// NewRestyFn represents new resty client.
	NewRestyFn           = resty.New
	NewStateMachineCliFn = client.NewStateMachineCli
)

// MetadataCommand executes the metadata query.
func MetadataCommand(ctx context.Context, deps *depspkg.HTTPDeps,
	_ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	metadataStmt := stmt.(*stmtpkg.Metadata)
	if metadataStmt.MetadataType == stmtpkg.MetadataTypes {
		// returns metadata explore define info.
		return map[string]interface{}{
			constants.BrokerRole:  broker.StateMachinePaths,
			constants.MasterRole:  master.StateMachinePaths,
			constants.StorageRole: storage.StateMachinePaths,
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

func exploreStateRepoData(ctx context.Context, deps *depspkg.HTTPDeps,
	metadataStmt *stmtpkg.Metadata) (interface{}, error) {
	var stateMachineInfo models.StateMachineInfo
	var ok bool
	switch metadataStmt.MetadataType {
	case stmtpkg.BrokerMetadata:
		stateMachineInfo, ok = broker.StateMachinePaths[metadataStmt.Type]
	case stmtpkg.MasterMetadata:
		stateMachineInfo, ok = master.StateMachinePaths[metadataStmt.Type]
	case stmtpkg.StorageMetadata:
		storageName := strings.TrimSpace(metadataStmt.StorageName)
		if storageName == "" {
			return nil, constants.ErrStorageNameRequired
		}
		if deps.Master.IsMaster() {
			// if current node is master, explore storage data.
			stateMachineInfo, ok = storage.StateMachinePaths[metadataStmt.Type]
			if !ok {
				return nil, nil
			}
			stateMgr := deps.Master.GetStateManager()
			storageCluster := stateMgr.GetStorageCluster(storageName)
			if storageCluster == nil {
				return nil, nil
			}
			return exploreData(ctx, storageCluster.GetRepo(), stateMachineInfo)
		}
		// if current node is not master, reverse proxy to master
		masterNode := deps.Master.GetMaster()
		address := masterNode.Node.HTTPAddress()
		var meta []interface{}
		_, err := NewRestyFn().R().SetQueryParams(map[string]string{
			"sql": fmt.Sprintf("show storage metedata where path='%s' and storage='%s'",
				metadataStmt.Type, metadataStmt.StorageName)}).
			SetHeader("Accept", "application/json").
			SetResult(&meta).
			Get(address + "/api/exec")
		if err != nil {
			return nil, err
		}
		return meta, err
	}
	if !ok {
		return nil, nil
	}
	return exploreData(ctx, deps.Repo, stateMachineInfo)
}

// exploreData explores state repository data by given path.
func exploreData(ctx context.Context, repo state.Repository, stateMachineInfo models.StateMachineInfo) (interface{}, error) {
	var rs []interface{}
	err := repo.WalkEntry(ctx, stateMachineInfo.Path, func(key, value []byte) {
		r := stateMachineInfo.CreateState()
		err0 := encoding.JSONUnmarshal(value, r)
		if err0 != nil {
			log.Warn("unmarshal metadata info err, ignore it",
				logger.String("key", string(key)),
				logger.String("data", string(value)))
			return
		}
		rs = append(rs, r)
	})
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// exploreStateMachineDate explores the state from state machine of broker/master/storage.
func exploreStateMachineDate(metadataStmt *stmtpkg.Metadata, deps *depspkg.HTTPDeps) (interface{}, error) {
	param := map[string]string{
		"type":        metadataStmt.Type,
		"role":        strconv.Itoa(int(metadataStmt.MetadataType)),
		"storageName": metadataStmt.StorageName,
	}
	var nodes []models.Node
	switch metadataStmt.MetadataType {
	case stmtpkg.BrokerMetadata:
		statelessNodes := deps.StateMgr.GetLiveNodes()
		for idx := range statelessNodes {
			nodes = append(nodes, &statelessNodes[idx])
		}
	case stmtpkg.MasterMetadata:
		nodes = append(nodes, deps.Master.GetMaster().Node)
	case stmtpkg.StorageMetadata:
		// forward master
		cli := NewStateMachineCliFn()
		return cli.FetchStateByNode(param, deps.Master.GetMaster().Node)
	default:
		return nil, nil
	}
	// forward broker node
	cli := NewStateMachineCliFn()
	return cli.FetchStateByNodes(param, nodes), nil
}

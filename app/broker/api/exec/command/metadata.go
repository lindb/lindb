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
	"strings"

	"github.com/go-resty/resty/v2"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// MetadataCommand executes the metadata query.
func MetadataCommand(ctx context.Context, deps *depspkg.HTTPDeps,
	_ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	metadataStmt := stmt.(*stmtpkg.Metadata)
	var stateMachineInfo models.StateMachineInfo
	var ok bool
	switch metadataStmt.MetadataType {
	case stmtpkg.MetadataTypes:
		// returns metadata explore define info.
		return map[string]interface{}{
			constants.BrokerRole:  broker.StateMachinePaths,
			constants.MasterRole:  master.StateMachinePaths,
			constants.StorageRole: storage.StateMachinePaths,
		}, nil
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
		_, err := resty.New().R().SetQueryParams(map[string]string{
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

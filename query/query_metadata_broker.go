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

package query

import (
	"context"
	"errors"
	"sort"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/strutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/sql/stmt"
)

type metadataQuery struct {
	runtime *queryFactory
	ctx     context.Context

	database      string
	metaStmtQuery *stmt.Metadata

	results []string
}

// newMetadataQuery creates the execution which executes the job of parallel query
func newMetadataQuery(
	ctx context.Context,
	database string,
	stmt *stmt.Metadata,
	queryBuilder *queryFactory,
) MetaDataQuery {
	return &metadataQuery{
		metaStmtQuery: stmt,
		database:      database,
		ctx:           ctx,
		runtime:       queryBuilder,
	}
}

func (mq *metadataQuery) WaitResponse() ([]string, error) {
	physicalPlan, err := mq.makePlan()
	if err != nil {
		return nil, err
	}

	resultCh, err := mq.runtime.taskManager.SubmitMetaDataTask(physicalPlan, mq.metaStmtQuery)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case result, ok := <-resultCh:
			// received all data, break for loop
			if !ok {
				deduped := strutil.DeDupStringSlice(mq.results)
				sort.Strings(deduped)
				return deduped, nil
			}
			if result.ErrMsg != "" {
				return nil, errors.New(result.ErrMsg)
			}
			if err := mq.handleTaskResponse(result); err != nil {
				return nil, err
			}
		case <-mq.ctx.Done():
			return nil, ErrTimeout
		}
	}
}

// buildPhysicalPlan builds distribution physical execute plan
func (mq *metadataQuery) makePlan() (*models.PhysicalPlan, error) {
	//FIXME need using storage's replica state ???
	storageNodes := mq.runtime.replicaStateMachine.GetQueryableReplicas(mq.database)
	storageNodesLen := len(storageNodes)
	if storageNodesLen == 0 {
		return nil, errNoAvailableStorageNode
	}
	curBroker := mq.runtime.nodeStateMachine.GetCurrentNode()
	curBrokerIndicator := (&curBroker).Indicator()
	physicalPlan := &models.PhysicalPlan{
		Database: mq.database,
		Root: models.Root{
			Indicator: curBrokerIndicator,
			NumOfTask: int32(storageNodesLen),
		},
	}
	receivers := []models.Node{curBroker}
	for storageNode, shardIDs := range storageNodes {
		physicalPlan.AddLeaf(models.Leaf{
			BaseNode: models.BaseNode{
				Parent:    curBrokerIndicator,
				Indicator: storageNode,
			},
			ShardIDs:  shardIDs,
			Receivers: receivers,
		})
	}
	return physicalPlan, nil
}

func (mq *metadataQuery) handleTaskResponse(resp *protoCommonV1.TaskResponse) error {
	result := &models.SuggestResult{}
	if err := encoding.JSONUnmarshal(resp.Payload, result); err != nil {
		return err
	}
	mq.results = append(mq.results, result.Values...)
	return nil
}

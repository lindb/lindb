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

package parallel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// leafTask represents the leaf node's task, the leaf node is always storage node
// 1. receives the task request, and searches the data from time seres engine
// 2. sends the result to the parent node(root or intermediate)
type leafTask struct {
	currentNodeID     string
	storageService    service.StorageService
	executorFactory   ExecutorFactory
	taskServerFactory rpc.TaskServerFactory
}

// newLeafTask creates the leaf task
func newLeafTask(
	currentNode models.Node,
	storageService service.StorageService,
	executorFactory ExecutorFactory,
	taskServerFactory rpc.TaskServerFactory,
) TaskProcessor {
	return &leafTask{
		currentNodeID:     (&currentNode).Indicator(),
		storageService:    storageService,
		executorFactory:   executorFactory,
		taskServerFactory: taskServerFactory,
	}
}

// Process processes the task request, searches the metric's data from time series engine
func (p *leafTask) Process(ctx context.Context, req *pb.TaskRequest) error {
	physicalPlan := models.PhysicalPlan{}
	if err := json.Unmarshal(req.PhysicalPlan, &physicalPlan); err != nil {
		return errUnmarshalPlan
	}

	foundTask := false
	var curLeaf models.Leaf
	for _, leaf := range physicalPlan.Leafs {
		if leaf.Indicator == p.currentNodeID {
			foundTask = true
			curLeaf = leaf
			break
		}
	}
	if !foundTask {
		return errWrongRequest
	}
	db, ok := p.storageService.GetDatabase(physicalPlan.Database)
	if !ok {
		return errNoDatabase
	}
	stream := p.taskServerFactory.GetStream(curLeaf.Parent)
	if stream == nil {
		return errNoSendStream
	}

	switch req.RequestType {
	case pb.RequestType_Data:
		if err := p.processDataSearch(ctx, db, curLeaf.ShardIDs, req, stream); err != nil {
			return err
		}
	case pb.RequestType_Metadata:
		if err := p.processMetadataSuggest(db, curLeaf.ShardIDs, req, stream); err != nil {
			return err
		}
	}
	return nil
}

func (p *leafTask) processMetadataSuggest(db tsdb.Database, shardIDs []int32,
	req *pb.TaskRequest, stream pb.TaskService_HandleServer,
) error {
	payload := req.Payload
	query := &stmt.Metadata{}
	if err := encoding.JSONUnmarshal(payload, query); err != nil {
		return errUnmarshalSuggest
	}
	exec := p.executorFactory.NewMetadataStorageExecutor(db, shardIDs, query)
	result, err := exec.Execute()
	if err != nil && !errors.Is(err, constants.ErrNotFound) {
		return err
	}
	// send result to upstream
	if err := stream.Send(&pb.TaskResponse{
		JobID:     req.JobID,
		TaskID:    req.ParentTaskID,
		Completed: true,
		Payload:   encoding.JSONMarshal(&models.SuggestResult{Values: result}),
	}); err != nil {
		return err
	}
	return nil
}

func (p *leafTask) processDataSearch(ctx context.Context, db tsdb.Database, shardIDs []int32,
	req *pb.TaskRequest, stream pb.TaskService_HandleServer,
) error {
	payload := req.Payload
	query := stmt.Query{}
	if err := encoding.JSONUnmarshal(payload, &query); err != nil {
		return errUnmarshalQuery
	}

	// execute leaf task
	storageExecuteCtx := p.executorFactory.NewStorageExecuteContext(shardIDs, &query)
	queryFlow := NewStorageQueryFlow(ctx, storageExecuteCtx, &query, req, stream, db.ExecutorPool())
	exec := p.executorFactory.NewStorageExecutor(queryFlow, db, storageExecuteCtx)
	exec.Execute()
	return nil
}

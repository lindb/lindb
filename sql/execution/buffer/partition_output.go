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

package buffer

import (
	"context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PartitionOutputBuffer struct {
	fragment *plan.PlanFragment
	taskID   model.TaskID

	finished bool
}

func NewPartitionOutputBuffer(taskID model.TaskID, fragment *plan.PlanFragment) OutputBuffer {
	return &PartitionOutputBuffer{
		taskID:   taskID,
		fragment: fragment,
	}
}

// AddPage implements OutputBuffer
func (output *PartitionOutputBuffer) AddPage(page *types.Page) {
	output.finished = page.Error != ""
	output.sendResultSet(&model.TaskResultSet{
		TaskID: output.taskID,
		Node:   *output.fragment.ParentNode,
		Page:   page,
		NoMore: output.finished,
	})
}

func (output *PartitionOutputBuffer) Complete() {
	fmt.Println("partition complete complete")
	if !output.finished {
		output.finished = true
		// TODO: send complete sign
		output.sendResultSet(&model.TaskResultSet{
			TaskID: output.taskID,
			Node:   *output.fragment.ParentNode,
			NoMore: output.finished, // FIXME: set nomore
		})
	}
}

func (output *PartitionOutputBuffer) sendResultSet(rs *model.TaskResultSet) {
	// TODO: conn pool?
	receiver := output.fragment.Receivers[0]
	conn, err := grpc.NewClient(receiver.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protoCommandV1.NewResultSetServiceClient(conn)
	// TODO: handle resp/error?
	_, err = client.ResultSet(context.TODO(), &protoCommandV1.ResultSetRequest{
		Payload: encoding.JSONMarshal(rs),
	})
	if err != nil {
		panic(err)
	}
}

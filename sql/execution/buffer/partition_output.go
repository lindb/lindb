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

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PartitionOutputBuffer struct {
	fragment   *plan.PlanFragment
	taskID     model.TaskID
	nodeSource string // gRPC address (ip:port) of the node executing this task, e.g. "1.2.3.4:2891"

	serializer *larrow.Serializer

	finished bool
}

func NewPartitionOutputBuffer(taskID model.TaskID, fragment *plan.PlanFragment, nodeSource string) OutputBuffer {
	return &PartitionOutputBuffer{
		taskID:     taskID,
		fragment:   fragment,
		nodeSource: nodeSource,
	}
}

func (output *PartitionOutputBuffer) AddRecord(record arrow.RecordBatch) {
	if output.serializer == nil {
		output.serializer = larrow.NewSerializer(record.Schema())
	}
	data, err := output.serializer.Serialize(record)
	if err != nil {
		panic(err)
	}
	output.sendResultSet(&model.TaskResultSet{
		TaskID: output.taskID,
		Node:   *output.fragment.ParentNode,
		Page:   data,
		NoMore: output.finished,
	})
}

func (output *PartitionOutputBuffer) Complete() {
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

// Fail sends the execution error back to the broker via gRPC so the broker can
// surface it to the client. The error message is prefixed with the node source
// (e.g. "[storage@1.2.3.4:2891]") so the client can identify which node failed.
func (output *PartitionOutputBuffer) Fail(errMsg string) {
	if !output.finished {
		output.finished = true
		// Prepend the node source so the broker and client can tell which remote
		// node produced the error.
		taggedMsg := fmt.Sprintf("[%s] %s", output.nodeSource, errMsg)
		output.sendResultSet(&model.TaskResultSet{
			TaskID: output.taskID,
			Node:   *output.fragment.ParentNode,
			ErrMsg: taggedMsg,
			NoMore: true,
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

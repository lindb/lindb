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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestLeafTaskDispatcher_Dispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := pb.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	leafTaskDispatcher := NewLeafTaskDispatcher(models.Node{IP: "1.1.1.1", Port: 9000}, nil, nil, nil)
	leafTaskDispatcher.Dispatch(context.TODO(), server, &pb.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}

func TestIntermediateTaskDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewIntermediateTaskDispatcher()
	dispatcher.Dispatch(context.TODO(), nil, &pb.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}

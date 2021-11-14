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

package rpc

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/constants"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
	"github.com/lindb/lindb/replica"
)

func TestWriteHandler_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	replicaServer := protoWriteV1.NewMockWriteService_WriteServer(ctrl)
	replicaServer.EXPECT().Context().Return(context.TODO())
	r := NewWriteHandler(walMgr)

	// case 1: family state not exist
	err := r.Write(replicaServer)
	assert.Error(t, err)
	// case 3: shard decode err
	ctx := metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs(constants.RPCMetaKeyFamilyState, strconv.Itoa(1)))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Write(replicaServer)
	assert.Error(t, err)

	// case 3: create partition err
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs(constants.RPCMetaKeyFamilyState,
			`{  "database":"test-db",
				"shard":{
					"id":1,
					"leader":2,
					"replica":{"replicas":[1,2,3]}
				},
				"familyTime":12321
			}`))
	replicaServer.EXPECT().Context().Return(ctx).AnyTimes()
	wal := replica.NewMockWriteAheadLog(ctrl)
	walMgr.EXPECT().GetOrCreateLog(gomock.Any()).Return(wal).AnyTimes()
	wal.EXPECT().GetOrCreatePartition(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = r.Write(replicaServer)
	assert.Error(t, err)

	// case 6: build replica replica err
	p := replica.NewMockPartition(ctrl)
	wal.EXPECT().GetOrCreatePartition(gomock.Any(), gomock.Any(), gomock.Any()).Return(p, nil).AnyTimes()
	p.EXPECT().BuildReplicaForLeader(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = r.Write(replicaServer)
	assert.Error(t, err)

	// case 7: recv req err
	p.EXPECT().BuildReplicaForLeader(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	replicaServer.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 8: recv req EOF err
	replicaServer.EXPECT().Recv().Return(nil, io.EOF)
	err = r.Write(replicaServer)
	assert.NoError(t, err)
	// case 9: write wal err
	replicaServer.EXPECT().Recv().Return(&protoWriteV1.WriteRequest{}, nil)
	p.EXPECT().WriteLog(gomock.Any()).Return(fmt.Errorf("err"))
	replicaServer.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 10: write wal ok
	replicaServer.EXPECT().Recv().Return(&protoWriteV1.WriteRequest{}, nil)
	p.EXPECT().WriteLog(gomock.Any()).Return(nil)
	replicaServer.EXPECT().Send(gomock.Any()).Return(nil)
	replicaServer.EXPECT().Recv().Return(nil, io.EOF)
	err = r.Write(replicaServer)
	assert.NoError(t, err)
}

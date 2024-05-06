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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/constants"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/replica"
)

func TestReplicaHandler_Replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	replicaServer := protoReplicaV1.NewMockReplicaService_ReplicaServer(ctrl)
	r := NewReplicaHandler(walMgr)

	// case 5: create partition err
	ctx := metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs(
			constants.RPCMetaReplicaState, `{"database":"test-db","shardId":1,"leader":2,"follower":3}`,
		))
	replicaServer.EXPECT().Context().Return(ctx).AnyTimes()
	wal := replica.NewMockWriteAheadLog(ctrl)
	walMgr.EXPECT().GetOrCreateLog(gomock.Any()).Return(wal).AnyTimes()
	wal.EXPECT().GetOrCreatePartition(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := r.Replica(replicaServer)
	assert.Error(t, err)

	// case 6: build replica replica err
	p := replica.NewMockPartition(ctrl)
	wal.EXPECT().GetOrCreatePartition(gomock.Any(), gomock.Any(), gomock.Any()).Return(p, nil).AnyTimes()
	p.EXPECT().BuildReplicaForFollower(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = r.Replica(replicaServer)
	assert.Error(t, err)

	// case 7: recv req EOF
	p.EXPECT().BuildReplicaForFollower(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	replicaServer.EXPECT().Recv().Return(nil, io.EOF)
	err = r.Replica(replicaServer)
	assert.NoError(t, err)

	// case 8: recv req err
	replicaServer.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	err = r.Replica(replicaServer)
	assert.Error(t, err)

	// case 9: replica log err
	replicaServer.EXPECT().Recv().Return(&protoReplicaV1.ReplicaRequest{}, nil)
	p.EXPECT().ReplicaLog(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("err"))
	replicaServer.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	err = r.Replica(replicaServer)
	assert.Error(t, err)

	// case 9: replica log success
	replicaServer.EXPECT().Recv().Return(&protoReplicaV1.ReplicaRequest{}, nil)
	p.EXPECT().ReplicaLog(gomock.Any(), gomock.Any()).Return(int64(10), nil)
	replicaServer.EXPECT().Send(gomock.Any()).Return(nil)
	replicaServer.EXPECT().Recv().Return(nil, io.EOF)
	err = r.Replica(replicaServer)
	assert.NoError(t, err)
}

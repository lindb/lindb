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

package handler

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/replica"
	replicaRpc "github.com/lindb/lindb/rpc/proto/replica"
)

func TestReplicaHandler_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	replicaServer := replicaRpc.NewMockReplicaService_WriteServer(ctrl)
	replicaServer.EXPECT().Context().Return(context.TODO())
	r := NewReplicaHandler(walMgr, nil)

	// case 1: database not exist
	err := r.Write(replicaServer)
	assert.Error(t, err)
	// case 2: shard not exist
	ctx := metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db"))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 3: leader not exist
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1)))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 4: replicas not exist
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1),
			"metaKeyLeader", strconv.Itoa(2),
		))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 5: replicas is empty
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1),
			"metaKeyLeader", strconv.Itoa(2),
			"metaKeyReplicas", `[]`,
		))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 5: create partition err
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1),
			"metaKeyLeader", strconv.Itoa(2),
			"metaKeyReplicas", `[1,2,3]`,
		))
	replicaServer.EXPECT().Context().Return(ctx).AnyTimes()
	wal := replica.NewMockWriteAheadLog(ctrl)
	walMgr.EXPECT().GetOrCreateLog(gomock.Any()).Return(wal).AnyTimes()
	wal.EXPECT().GetOrCreatePartition(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = r.Write(replicaServer)
	assert.Error(t, err)

	// case 6: build replica replica err
	p := replica.NewMockPartition(ctrl)
	wal.EXPECT().GetOrCreatePartition(gomock.Any()).Return(p, nil).AnyTimes()
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
	replicaServer.EXPECT().Recv().Return(&replicaRpc.WriteRequest{}, nil)
	p.EXPECT().WriteLog(gomock.Any()).Return(fmt.Errorf("err"))
	replicaServer.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	err = r.Write(replicaServer)
	assert.Error(t, err)
	// case 10: write wal ok
	replicaServer.EXPECT().Recv().Return(&replicaRpc.WriteRequest{}, nil)
	p.EXPECT().WriteLog(gomock.Any()).Return(nil)
	replicaServer.EXPECT().Send(gomock.Any()).Return(nil)
	replicaServer.EXPECT().Recv().Return(nil, io.EOF)
	err = r.Write(replicaServer)
	assert.NoError(t, err)
}

func TestReplicaHandler_Replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	replicaServer := replicaRpc.NewMockReplicaService_ReplicaServer(ctrl)
	replicaServer.EXPECT().Context().Return(context.TODO())
	r := NewReplicaHandler(walMgr, nil)

	// case 1: database not exist
	err := r.Replica(replicaServer)
	assert.Error(t, err)
	// case 2: shard not exist
	ctx := metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db"))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Replica(replicaServer)
	assert.Error(t, err)
	// case 3: leader not exist
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1)))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Replica(replicaServer)
	assert.Error(t, err)
	// case 4: replicas not exist
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1),
			"metaKeyLeader", strconv.Itoa(2),
		))
	replicaServer.EXPECT().Context().Return(ctx)
	err = r.Replica(replicaServer)
	assert.Error(t, err)
	// case 5: create partition err
	ctx = metadata.NewIncomingContext(context.TODO(),
		metadata.Pairs("metaKeyDatabase", "test-db",
			"metaKeyShardID", strconv.Itoa(1),
			"metaKeyLeader", strconv.Itoa(2),
			"metaKeyReplica", `3`,
		))
	replicaServer.EXPECT().Context().Return(ctx).AnyTimes()
	wal := replica.NewMockWriteAheadLog(ctrl)
	walMgr.EXPECT().GetOrCreateLog(gomock.Any()).Return(wal).AnyTimes()
	wal.EXPECT().GetOrCreatePartition(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = r.Replica(replicaServer)
	assert.Error(t, err)

	// case 6: build replica replica err
	p := replica.NewMockPartition(ctrl)
	wal.EXPECT().GetOrCreatePartition(gomock.Any()).Return(p, nil).AnyTimes()
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
	replicaServer.EXPECT().Recv().Return(&replicaRpc.ReplicaRequest{}, nil)
	p.EXPECT().ReplicaLog(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("err"))
	replicaServer.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	err = r.Replica(replicaServer)
	assert.Error(t, err)

	// case 9: replica log success
	replicaServer.EXPECT().Recv().Return(&replicaRpc.ReplicaRequest{}, nil)
	p.EXPECT().ReplicaLog(gomock.Any(), gomock.Any()).Return(int64(10), nil)
	replicaServer.EXPECT().Send(gomock.Any()).Return(nil)
	replicaServer.EXPECT().Recv().Return(nil, io.EOF)
	err = r.Replica(replicaServer)
	assert.NoError(t, err)
}

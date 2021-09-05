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

package brokerquery

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

func Test_MetadataQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	stateMgr := broker.NewMockStateManager(ctrl)
	thisTaskManager := NewMockTaskManager(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	metaDataQuery := newMetadataQuery(
		ctx,
		"db",
		&stmt.Metadata{},
		&queryFactory{
			stateMgr:    stateMgr,
			taskManager: thisTaskManager,
		},
	)

	// GetQueryableReplicas return empty
	stateMgr.EXPECT().GetQueryableReplicas("db").
		Return(map[string][]models.ShardID{}, nil)
	results, err := metaDataQuery.WaitResponse()
	assert.Error(t, err)
	assert.Nil(t, results)

	stateMgr.EXPECT().GetQueryableReplicas("db").
		Return(map[string][]models.ShardID{
			"1.1.1.1:9000": {1, 2, 4},
			"1.1.1.2:9000": {3, 5, 6},
		}, nil).AnyTimes()
	stateMgr.EXPECT().GetCurrentNode().Return(models.StatelessNode{
		HostIP: "1.1.1.3", GRPCPort: 8000,
	}).AnyTimes()

	// wait error
	thisTaskManager.EXPECT().SubmitMetaDataTask(gomock.Any(), gomock.Any()).Return(nil, io.ErrClosedPipe)
	_, err = metaDataQuery.WaitResponse()
	assert.Error(t, err)

	// return error
	response1Ch := make(chan *protoCommonV1.TaskResponse)
	time.AfterFunc(time.Millisecond*200, func() {
		response1Ch <- &protoCommonV1.TaskResponse{ErrMsg: "error"}
	})
	thisTaskManager.EXPECT().SubmitMetaDataTask(gomock.Any(), gomock.Any()).Return(
		response1Ch,
		nil)
	_, err = metaDataQuery.WaitResponse()
	assert.Error(t, err)

	// bad data
	response2Ch := make(chan *protoCommonV1.TaskResponse)
	time.AfterFunc(time.Millisecond*200, func() {
		response2Ch <- &protoCommonV1.TaskResponse{Payload: nil}
	})
	thisTaskManager.EXPECT().SubmitMetaDataTask(gomock.Any(), gomock.Any()).Return(
		response2Ch, nil)
	_, err = metaDataQuery.WaitResponse()
	assert.Error(t, err)

	// ok data
	data := encoding.JSONMarshal(models.SuggestResult{Values: []string{"a"}})
	response3Ch := make(chan *protoCommonV1.TaskResponse)
	time.AfterFunc(time.Millisecond*200, func() {
		response3Ch <- &protoCommonV1.TaskResponse{Payload: data}
		response3Ch <- &protoCommonV1.TaskResponse{Payload: data}
		close(response3Ch)
	})
	thisTaskManager.EXPECT().SubmitMetaDataTask(gomock.Any(), gomock.Any()).Return(
		response3Ch, nil)
	results, err = metaDataQuery.WaitResponse()
	assert.Nil(t, err)
	assert.Len(t, results, 1)

	// timeout
	response4Ch := make(chan *protoCommonV1.TaskResponse)
	time.AfterFunc(time.Millisecond*200, cancel)
	thisTaskManager.EXPECT().SubmitMetaDataTask(gomock.Any(), gomock.Any()).Return(
		response4Ch,
		nil)
	_, err = metaDataQuery.WaitResponse()
	assert.Error(t, err)
}

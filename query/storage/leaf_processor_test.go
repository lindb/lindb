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

package storagequery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestLeafTaskProcessor_Process_sendStreamFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	leafTaskProcessor := NewLeafTaskProcessor(
		&models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
		nil,
		nil)
	leafTaskProcessor.Process(
		flow.NewTaskContextWithTimeout(context.Background(), time.Second),
		server,
		&protoCommonV1.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}

func TestLeafTask_Process_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	engine := tsdb.NewMockEngine(ctrl)
	serverStream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	mockDatabase := tsdb.NewMockDatabase(ctrl)

	currentNode := models.StatelessNode{HostIP: "1.1.1.3", GRPCPort: 8000}
	processorI := NewLeafTaskProcessor(&currentNode, engine, taskServerFactory)
	processor := processorI.(*leafTaskProcessor)

	cases := []struct {
		name    string
		req     *protoCommonV1.TaskRequest
		prepare func()
		assert  func(err error)
	}{
		{
			name: "unmarshal error",
			req:  &protoCommonV1.TaskRequest{PhysicalPlan: nil},
			assert: func(err error) {
				assert.True(t, errors.Is(err, query.ErrUnmarshalPlan))
			},
		},
		{
			name: "wrong request",
			req: &protoCommonV1.TaskRequest{PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
				Leaves: []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.4:8000"}}},
			})},
			assert: func(err error) {
				assert.True(t, errors.Is(err, query.ErrBadPhysicalPlan))
			},
		},
		{
			name: "db not exist",
			req: &protoCommonV1.TaskRequest{PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
				Database: "test_db",
				Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
			}), Payload: encoding.JSONMarshal(&stmt.Query{MetricName: "cpu"})},
			prepare: func() {
				engine.EXPECT().GetDatabase(gomock.Any()).Return(nil, false)
			},
			assert: func(err error) {
				assert.True(t, errors.Is(err, query.ErrNoDatabase))
			},
		},
		{
			name: "get upstream err",
			req: &protoCommonV1.TaskRequest{PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
				Database: "test_db",
				Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
			}), Payload: encoding.JSONMarshal(&stmt.Query{MetricName: "cpu"})},
			prepare: func() {
				engine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
				taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil)
			},
			assert: func(err error) {
				assert.True(t, errors.Is(err, query.ErrNoSendStream))
			},
		},
		{
			name: "unmarshal query err",
			req: &protoCommonV1.TaskRequest{PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
				Database: "test_db",
				Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
			}), Payload: []byte{1, 2, 3}},
			prepare: func() {
				engine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
				taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
			},
			assert: func(err error) {
				assert.True(t, errors.Is(err, query.ErrUnmarshalQuery))
			},
		},
		{
			name: "test executor fail",
			req: &protoCommonV1.TaskRequest{PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
				Database: "test_db",
				Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
			}), Payload: encoding.JSONMarshal(&stmt.Query{MetricName: "cpu"})},
			prepare: func() {
				mockDatabase.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{})
				taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
				engine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
			},
			assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "unknown request type",
			req: &protoCommonV1.TaskRequest{RequestType: protoCommonV1.RequestType(10),
				PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
					Database: "test_db",
					Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
				}), Payload: encoding.JSONMarshal(&stmt.Query{MetricName: "cpu"})},
			prepare: func() {
				taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
				engine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
			},
			assert: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			err := processor.process(
				flow.NewTaskContextWithTimeout(context.Background(), time.Second), tt.req)
			tt.assert(err)
		})
	}
}

func TestLeafProcessor_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	engine := tsdb.NewMockEngine(ctrl)

	currentNode := models.StatelessNode{HostIP: "1.1.1.3", GRPCPort: 8000}
	processorI := NewLeafTaskProcessor(&currentNode, engine, taskServerFactory)
	processor := processorI.(*leafTaskProcessor)
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	plan := encoding.JSONMarshal(&models.PhysicalPlan{
		Database: "test_db",
		Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	qry := stmt.Query{MetricName: "cpu"}
	data := encoding.JSONMarshal(&qry)

	mockDatabase.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{})
	engine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)

	serverStream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
	err := processor.process(flow.NewTaskContextWithTimeout(context.Background(), time.Second),
		&protoCommonV1.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.NoError(t, err)
}

func TestLeafTask_Suggest_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	engine := tsdb.NewMockEngine(ctrl)

	currentNode := models.StatelessNode{HostIP: "1.1.1.3", GRPCPort: 8000}
	processorI := NewLeafTaskProcessor(&currentNode, engine, taskServerFactory)
	processor := processorI.(*leafTaskProcessor)
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	plan := encoding.JSONMarshal(&models.PhysicalPlan{
		Database: "test_db",
		Leaves:   []*models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	engine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true).AnyTimes()
	serverStream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream).AnyTimes()

	cases := []struct {
		name    string
		payload []byte
		prepare func()
		assert  func(err error)
	}{
		{
			name:    "unmarshal err",
			payload: []byte{1, 2, 3},
			assert: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "stream err",
			payload: encoding.JSONMarshal(&stmt.MetricMetadata{}),
			prepare: func() {
				serverStream.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe)
			},
			assert: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "suggest successfully",
			payload: encoding.JSONMarshal(&stmt.MetricMetadata{}),
			prepare: func() {
				serverStream.EXPECT().Send(gomock.Any()).Return(nil)
			},
			assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "suggest not data",
			payload: encoding.JSONMarshal(&stmt.MetricMetadata{}),
			prepare: func() {
				serverStream.EXPECT().Send(gomock.Any()).Return(nil)
				q := NewMockstorageMetadataQuery(ctrl)
				q.EXPECT().Execute().Return(nil, constants.ErrNotFound)
				newStorageMetadataQueryFn = func(database tsdb.Database, shardIDs []models.ShardID,
					request *stmt.MetricMetadata) storageMetadataQuery {
					return q
				}
			},
			assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "suggest not data",
			payload: encoding.JSONMarshal(&stmt.MetricMetadata{}),
			prepare: func() {
				q := NewMockstorageMetadataQuery(ctrl)
				q.EXPECT().Execute().Return(nil, constants.ErrNoLiveNode)
				newStorageMetadataQueryFn = func(database tsdb.Database, shardIDs []models.ShardID,
					request *stmt.MetricMetadata) storageMetadataQuery {
					return q
				}
			},
			assert: func(err error) {
				assert.Equal(t, constants.ErrNoLiveNode, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newStorageMetadataQueryFn = newStorageMetadataQuery
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			err := processor.process(flow.NewTaskContextWithTimeout(context.Background(), time.Second),
				&protoCommonV1.TaskRequest{
					PhysicalPlan: plan,
					RequestType:  protoCommonV1.RequestType_Metadata,
					Payload:      tt.payload})

			if tt.assert != nil {
				tt.assert(err)
			}
		})
	}
}

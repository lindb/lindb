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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestLeafTaskProcessor_Process_sendStreamFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	leafTaskProcessor := NewLeafTaskProcessor(
		models.Node{IP: "1.1.1.1", Port: 9000},
		nil,
		nil)
	leafTaskProcessor.Process(
		context.Background(),
		server,
		&protoCommonV1.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}

func TestLeafTask_Process_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	storageService := service.NewMockStorageService(ctrl)
	serverStream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	mockDatabase := tsdb.NewMockDatabase(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processorI := NewLeafTaskProcessor(currentNode, storageService, taskServerFactory)
	processor := processorI.(*leafTaskProcessor)
	// unmarshal error
	err := processor.process(
		context.Background(),
		&protoCommonV1.TaskRequest{PhysicalPlan: nil})
	assert.True(t, errors.Is(err, errUnmarshalPlan))

	plan, _ := json.Marshal(&models.PhysicalPlan{
		Leafs: []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.4:8000"}}},
	})
	// wrong request
	err = processor.process(
		context.Background(),
		&protoCommonV1.TaskRequest{PhysicalPlan: plan})
	assert.True(t, errors.Is(err, errBadPhysicalPlan))

	plan, _ = json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	query := stmt.Query{MetricName: "cpu"}
	data := encoding.JSONMarshal(&query)

	// db not exist
	storageService.EXPECT().GetDatabase(gomock.Any()).Return(nil, false)
	err = processor.process(
		context.Background(),
		&protoCommonV1.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.True(t, errors.Is(err, errNoDatabase))

	// test get upstream err
	storageService.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil)
	err = processor.process(
		context.Background(),
		&protoCommonV1.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.True(t, errors.Is(err, errNoSendStream))

	// unmarshal query err
	storageService.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
	err = processor.process(
		context.Background(),
		&protoCommonV1.TaskRequest{PhysicalPlan: plan, Payload: []byte{1, 2, 3}})
	assert.Equal(t, errUnmarshalQuery, err)

	// test executor fail
	mockDatabase.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{})
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
	storageService.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true).AnyTimes()
	err = processor.process(
		context.Background(),
		&protoCommonV1.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.NoError(t, err)
}

func TestLeafProcessor_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	storageService := service.NewMockStorageService(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processorI := NewLeafTaskProcessor(currentNode, storageService, taskServerFactory)
	processor := processorI.(*leafTaskProcessor)
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	plan, _ := json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	query := stmt.Query{MetricName: "cpu"}
	data := encoding.JSONMarshal(&query)

	mockDatabase.EXPECT().ExecutorPool().Return(&tsdb.ExecutorPool{})
	storageService.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)

	serverStream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
	err := processor.process(context.Background(), &protoCommonV1.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.NoError(t, err)
}

func TestLeafTask_Suggest_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	storageService := service.NewMockStorageService(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processorI := NewLeafTaskProcessor(currentNode, storageService, taskServerFactory)
	processor := processorI.(*leafTaskProcessor)
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	plan, _ := json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	storageService.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true).AnyTimes()
	serverStream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream).AnyTimes()

	// test unmarshal err
	err := processor.process(context.Background(), &protoCommonV1.TaskRequest{
		PhysicalPlan: plan,
		RequestType:  protoCommonV1.RequestType_Metadata,
		Payload:      []byte{1, 2, 3}})
	assert.Error(t, err)

	// test stream err
	data := encoding.JSONMarshal(&stmt.Metadata{})
	serverStream.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe)
	err = processor.process(context.Background(), &protoCommonV1.TaskRequest{
		PhysicalPlan: plan,
		RequestType:  protoCommonV1.RequestType_Metadata,
		Payload:      data})
	assert.Error(t, err)
	// test send result ok
	serverStream.EXPECT().Send(gomock.Any()).Return(nil)
	err = processor.process(context.Background(), &protoCommonV1.TaskRequest{
		PhysicalPlan: plan,
		RequestType:  protoCommonV1.RequestType_Metadata,
		Payload:      data})
	assert.Nil(t, err)
}

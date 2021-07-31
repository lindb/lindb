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
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
)

func Test_Intermediate_decodePhysicalPlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskProcessor := intermediateTaskProcessor{}
	_, _, err := taskProcessor.decodePhysicalPlan(&protoCommonV1.TaskRequest{})
	assert.Error(t, err)

	plan := models.PhysicalPlan{Intermediates: []models.Intermediate{{
		BaseNode: models.BaseNode{Indicator: "1.1.1.1:80"},
	}}}
	data := encoding.JSONMarshal(plan)
	_, _, err = taskProcessor.decodePhysicalPlan(&protoCommonV1.TaskRequest{
		PhysicalPlan: data,
	})
	assert.Error(t, err)

	taskProcessor2 := intermediateTaskProcessor{currentNodeID: "1.1.1.1:80"}
	_, _, err = taskProcessor2.decodePhysicalPlan(&protoCommonV1.TaskRequest{
		PhysicalPlan: data,
	})
	assert.NoError(t, err)
}

func Test_Intermediate_process_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)

	taskProcessor := intermediateTaskProcessor{
		logger: logger.GetLogger("query", "Test"),
	}
	stream.EXPECT().Send(gomock.Any()).Return(nil)
	taskProcessor.Process(context.Background(), stream, &protoCommonV1.TaskRequest{})

	stream.EXPECT().Send(gomock.Any()).Return(nil)
	taskProcessor.Process(context.Background(), stream, &protoCommonV1.TaskRequest{
		RequestType: protoCommonV1.RequestType_Metadata,
	})

	stream.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe)
	taskProcessor.Process(context.Background(), stream, &protoCommonV1.TaskRequest{
		Type: protoCommonV1.TaskType_Leaf,
	})
}

func Test_Intermediate_processIntermediateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskManager := NewMockTaskManager(ctrl)
	taskProcessor := intermediateTaskProcessor{
		taskManager:   taskManager,
		currentNodeID: "1.1.1.1:80",
		logger:        logger.GetLogger("query", "Test"),
	}
	// decode stmt error
	err := taskProcessor.processIntermediateTask(context.Background(), &protoCommonV1.TaskRequest{})
	assert.Error(t, err)
	// decode plan error
	stmtData := []byte("{}")
	err = taskProcessor.processIntermediateTask(context.Background(), &protoCommonV1.TaskRequest{
		Payload: stmtData,
	})
	assert.Error(t, err)

	// task manager fail
	plan := models.PhysicalPlan{Intermediates: []models.Intermediate{{
		BaseNode: models.BaseNode{Indicator: "1.1.1.1:80"},
	}}}
	planData := encoding.JSONMarshal(plan)

	// closed error
	ch1 := make(chan *series.TimeSeriesEvent)
	taskManager.EXPECT().SubmitIntermediateMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(ch1)
	time.AfterFunc(time.Millisecond*200, func() {
		close(ch1)
	})
	assert.Error(t, taskProcessor.processIntermediateTask(context.Background(),
		&protoCommonV1.TaskRequest{
			Payload:      stmtData,
			PhysicalPlan: planData,
		}))

	// event error
	ch2 := make(chan *series.TimeSeriesEvent)
	taskManager.EXPECT().SubmitIntermediateMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(ch2)
	time.AfterFunc(time.Millisecond*200, func() {
		ch2 <- &series.TimeSeriesEvent{Err: io.ErrClosedPipe}
	})
	assert.Error(t, taskProcessor.processIntermediateTask(context.Background(),
		&protoCommonV1.TaskRequest{
			Payload:      stmtData,
			PhysicalPlan: planData,
		}))
	// context done
	ch3 := make(chan *series.TimeSeriesEvent)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Millisecond*200, cancel)
	taskManager.EXPECT().SubmitIntermediateMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(ch3)
	assert.Nil(t, taskProcessor.processIntermediateTask(ctx,
		&protoCommonV1.TaskRequest{
			Payload:      stmtData,
			PhysicalPlan: planData,
		}))
}

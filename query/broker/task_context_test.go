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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
)

func Test_TaskContext_metaDataTaskContext(t *testing.T) {
	ch := make(chan *protoCommonV1.TaskResponse)
	taskCtx1 := newMetaDataTaskContext(
		"1",
		RootTask,
		"",
		"",
		1,
		ch,
	)

	time.AfterFunc(time.Millisecond*10, func() {
		<-ch
		<-ch
	})
	// drop as there is  no reader is reading
	taskCtx1.WriteResponse(&protoCommonV1.TaskResponse{}, "")
	time.Sleep(time.Millisecond * 50)
	taskCtx1.WriteResponse(&protoCommonV1.TaskResponse{}, "")
	taskCtx1.WriteResponse(&protoCommonV1.TaskResponse{}, "")

	assert.Equal(t, "1", taskCtx1.TaskID())
	assert.Equal(t, RootTask, taskCtx1.TaskType())
	assert.Equal(t, "", taskCtx1.ParentNode())
	assert.Equal(t, "", taskCtx1.ParentTaskID())
	assert.True(t, taskCtx1.Expired(time.Nanosecond))
}

func Test_TaskContext_metricTaskContext(t *testing.T) {
	ch := make(chan *series.TimeSeriesEvent)
	taskCtx2 := newMetricTaskContext(
		"1",
		RootTask,
		"",
		"",
		nil,
		2,
		ch,
	)

	// sent omitted
	taskCtx2.WriteResponse(
		&protoCommonV1.TaskResponse{ErrMsg: "error"},
		"1.1.1.1",
	)
	// sent emitted
	time.AfterFunc(time.Millisecond, func() {
		<-ch
	})
	time.Sleep(time.Millisecond * 3)
	taskCtx2.WriteResponse(
		&protoCommonV1.TaskResponse{ErrMsg: "error"},
		"1.1.1.1",
	)
	// closed
	taskCtx2.WriteResponse(
		&protoCommonV1.TaskResponse{ErrMsg: "error2"},
		"1.1.1.2",
	)
}

func Test_TaskContext_handleStats(t *testing.T) {
	taskCtx3 := newMetricTaskContext(
		"1",
		RootTask,
		"",
		"",
		nil,
		2,
		nil,
	).(*metricTaskContext)
	//
	storageNodeStat1 := models.NewStorageStats()
	storageNodeStat1.NetPayload = 30000
	data1 := encoding.JSONMarshal(storageNodeStat1)

	storageNodeStat2 := models.NewStorageStats()
	storageNodeStat2.NetPayload = 40000
	data2 := encoding.JSONMarshal(storageNodeStat1)
	taskCtx3.handleStats(
		&protoCommonV1.TaskResponse{
			Stats: data1,
			Type:  protoCommonV1.TaskType_Leaf},
		"1.1.1.1")
	taskCtx3.handleStats(
		&protoCommonV1.TaskResponse{Stats: data2,
			Type: protoCommonV1.TaskType_Leaf},
		"1.1.1.2")

	queryStats := encoding.JSONMarshal(taskCtx3.stats)
	taskCtx3.stats = nil
	// query stats from intermediate
	taskCtx3.handleStats(
		&protoCommonV1.TaskResponse{Stats: queryStats,
			Type: protoCommonV1.TaskType_Intermediate},
		"1")
	taskCtx3.handleStats(
		&protoCommonV1.TaskResponse{Stats: queryStats,
			Type: protoCommonV1.TaskType_Intermediate},
		"2")
	assert.Len(t, taskCtx3.stats.BrokerNodes, 2)
}

func Test_TaskContext_metricTaskContext_notFound(t *testing.T) {
	ch := make(chan *series.TimeSeriesEvent)
	taskCtx3 := newMetricTaskContext(
		"1",
		RootTask,
		"",
		"",
		nil,
		2,
		ch,
	)
	// sent emitted
	time.AfterFunc(time.Millisecond*2, func() {
		<-ch
	})
	time.Sleep(time.Millisecond * 10)
	taskCtx3.WriteResponse(
		&protoCommonV1.TaskResponse{ErrMsg: "metricID not found"},
		"1.1.1.1",
	)
	taskCtx3.WriteResponse(
		&protoCommonV1.TaskResponse{ErrMsg: "metricID not found"},
		"1.1.1.1",
	)
}

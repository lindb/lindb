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

package replication

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
)

func TestChannel_New(t *testing.T) {
	ch, err := newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, "database", ch.Database())
	assert.Equal(t, models.ShardID(1), ch.ShardID())

	defer func() {
		newFanOutQueue = queue.NewFanOutQueue
	}()
	newFanOutQueue = func(dirPath string, dataSizeLimit int64, removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return nil, fmt.Errorf("err")
	}
	ch, err = newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func TestChannel_GetOrCreateReplicator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch.Startup()
	target := models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 12345}
	r, err := ch.GetOrCreateReplicator(&target)
	assert.NoError(t, err)
	assert.Equal(t, &target, r.Target())

	r2, err := ch.GetOrCreateReplicator(&target)
	assert.NoError(t, err)
	assert.Equal(t, r, r2)

	assert.Len(t, ch.Targets(), 1)
	assert.Equal(t, &target, ch.Targets()[0])

	ch1 := ch.(*channel)
	fanout := queue.NewMockFanOutQueue(ctrl)
	fanout.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, fmt.Errorf("err"))
	ch1.q = fanout
	r2, err = ch.GetOrCreateReplicator(&models.StatelessNode{HostIP: "err", GRPCPort: 12345})
	assert.Error(t, err)
	assert.Nil(t, r2)
	cancel()
	time.Sleep(300 * time.Millisecond)
}

func TestChannel_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch.Startup()

	ch1 := ch.(*channel)
	fanout := queue.NewMockFanOutQueue(ctrl)
	fanout.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	ch1.q = fanout

	metric := &protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}
	err = ch.Write(metric)
	assert.NoError(t, err)
	err = ch.Write(metric)
	assert.NoError(t, err)

	cancel()
	time.Sleep(time.Millisecond * 600)

	ch, err = newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch1 = ch.(*channel)
	// ignore data, after closed
	chunk := NewMockChunk(ctrl)
	ch1.chunk = chunk
	// make sure chan is full
	ch1.ch <- []byte{1, 2}
	ch1.ch <- []byte{1, 2}
	chunk.EXPECT().Append(gomock.Any())
	chunk.EXPECT().IsFull().Return(true)
	chunk.EXPECT().MarshalBinary().Return([]byte{1, 2, 3}, nil)
	err = ch.Write(metric)
	assert.Error(t, err)
	time.Sleep(time.Millisecond * 500)
}

func TestChannel_checkFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch.Startup()

	ch1 := ch.(*channel)
	fanout := queue.NewMockFanOutQueue(ctrl)
	fanout.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	ch1.q = fanout

	metric := &protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}
	err = ch.Write(metric)
	assert.NoError(t, err)

	time.Sleep(time.Second)
	cancel()
	time.Sleep(300 * time.Millisecond)
}

func TestChannel_write_pending_before_close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	metric := &protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}
	err = ch.Write(metric)
	assert.NoError(t, err)

	ch1 := ch.(*channel)
	ch1.ch <- []byte{1, 2, 3}
	fanOut := queue.NewMockFanOutQueue(ctrl)
	fanOut.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	ch1.q = fanOut
	ch1.writePendingBeforeClose()
}

func TestChannel_chunk_marshal_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	chunk := NewMockChunk(ctrl)
	ch1 := ch.(*channel)
	ch1.chunk = chunk

	metric := &protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}
	chunk.EXPECT().Append(gomock.Any())
	chunk.EXPECT().IsFull().Return(true)
	chunk.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
	err = ch.Write(metric)
	assert.Error(t, err)

	chunk.EXPECT().Append(gomock.Any())
	chunk.EXPECT().IsFull().Return(true)
	chunk.EXPECT().MarshalBinary().Return(nil, nil)
	err = ch.Write(metric)
	assert.NoError(t, err)

	chunk.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
	ch1.flushChunk()
	chunk.EXPECT().MarshalBinary().Return(nil, nil)
	ch1.flushChunk()
}

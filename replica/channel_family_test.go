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

package replica

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/metric"
)

func TestFamilyChannel_new(t *testing.T) {
	f := newFamilyChannel(context.TODO(), config.Write{}, "db", 1,
		1, nil, models.ShardState{}, nil)
	assert.NotNil(t, f)
	f.Stop(10)

	f = newFamilyChannel(context.TODO(), config.Write{}, "db", 1,
		1, nil, models.ShardState{}, nil)
	assert.NotNil(t, f)
	go func() {
		time.Sleep(100 * time.Millisecond)
		f1 := f.(*familyChannel)
		f1.stoppedSignal <- struct{}{}
	}()
	f.Stop(timeutil.OneSecond)
}

func TestFamilyChannel_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	chunk := NewMockChunk(ctrl)
	converter := metric.NewProtoConverter()
	var brokerRow metric.BrokerRow
	assert.NoError(t, converter.ConvertTo(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}, &brokerRow))

	cases := []struct {
		name    string
		rows    []metric.BrokerRow
		prepare func()
		wantErr bool
	}{
		{
			name:    "empty rows",
			wantErr: false,
		},
		{
			name: "batch failure",
			rows: []metric.BrokerRow{brokerRow},
			prepare: func() {
				chunk.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "batch successfully, but send failure",
			rows: []metric.BrokerRow{brokerRow},
			prepare: func() {
				chunk.EXPECT().Write(gomock.Any()).Return(0, nil)
				chunk.EXPECT().IsFull().Return(true)
				chunk.EXPECT().Compress().Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "batch successfully",
			rows: []metric.BrokerRow{brokerRow},
			prepare: func() {
				chunk.EXPECT().Write(gomock.Any()).Return(0, nil)
				chunk.EXPECT().IsFull().Return(false)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ch := &familyChannel{
				chunk:          chunk,
				stoppedSignal:  make(chan struct{}, 1),
				stoppingSignal: make(chan struct{}, 1),
				statistics:     metrics.NewBrokerFamilyWriteStatistics("db"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}

			err := ch.Write(context.TODO(), tt.rows)

			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFamilyChannel_leaderChanged(t *testing.T) {
	shard := models.ShardState{ID: 1}
	liveNodes := make(map[models.NodeID]models.StatefulNode)
	fc := &familyChannel{
		leaderChangedSignal: make(chan struct{}, 1),
		statistics:          metrics.NewBrokerFamilyWriteStatistics("db"),
	}
	fc.leaderChanged(shard, liveNodes)
	fc.lock4meta.Lock()
	assert.Equal(t, shard, fc.shardState)
	assert.Equal(t, liveNodes, fc.liveNodes)
	fc.lock4meta.Unlock()
}

func TestChannel_checkFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	chunk := NewMockChunk(ctrl)
	ctx, cancel := context.WithCancel(context.TODO())
	f := &familyChannel{
		cancel:         cancel,
		ctx:            ctx,
		chunk:          chunk,
		batchTimout:    5 * time.Second,
		lastFlushTime:  atomic.NewInt64(timeutil.Now()),
		stoppingSignal: make(chan struct{}, 1),
		stoppedSignal:  make(chan struct{}, 1),
		ch:             make(chan *compressedChunk),
		statistics:     metrics.NewBrokerFamilyWriteStatistics("test"),
		logger:         logger.GetLogger("test", "test"),
	}
	f.checkFlush()

	f.lastFlushTime.Store(timeutil.Now() - 6*timeutil.OneSecond)
	chunk.EXPECT().IsEmpty().Return(false)
	chunk.EXPECT().Compress().Return(nil, nil)
	f.checkFlush()

	f.Stop(10)
}

func TestFamilyChannel_flushChunkOnFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	chunk := NewMockChunk(ctrl)
	chunk.EXPECT().IsFull().Return(true).AnyTimes()
	chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil).AnyTimes()
	ctx, cancel := context.WithCancel(context.TODO())
	f := &familyChannel{
		cancel:        cancel,
		ctx:           ctx,
		chunk:         chunk,
		batchTimout:   5 * time.Second,
		lastFlushTime: atomic.NewInt64(timeutil.Now()),
		ch:            make(chan *compressedChunk, 1),
		statistics:    metrics.NewBrokerFamilyWriteStatistics("db"),
		logger:        logger.GetLogger("test", "test"),
	}
	assert.NoError(t, f.flushChunkOnFull(context.TODO()))
	ctx1, cancel1 := context.WithCancel(context.TODO())
	cancel1()
	assert.Equal(t, ErrIngestTimeout, f.flushChunkOnFull(ctx1))
	cancel()
	assert.Equal(t, ErrFamilyChannelCanceled, f.flushChunkOnFull(context.TODO()))
}

func TestFamilyChannel_isExpire(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	f := &familyChannel{
		ctx:            ctx,
		cancel:         cancel,
		familyTime:     1,
		ch:             make(chan *compressedChunk),
		stoppingSignal: make(chan struct{}, 1),
		statistics:     metrics.NewBrokerFamilyWriteStatistics("db"),
		lastFlushTime:  atomic.NewInt64(timeutil.Now()),
	}
	assert.Equal(t, int64(1), f.FamilyTime())

	assert.False(t, f.isExpire(timeutil.OneHour, 0))
	assert.False(t, f.isExpire(0, 0))
	f.lastFlushTime.Store(timeutil.Now() - timeutil.OneHour - 16*timeutil.OneMinute)
	assert.True(t, f.isExpire(timeutil.OneHour, 0))
	f.lastFlushTime.Store(timeutil.Now() - 16*timeutil.OneMinute)
	assert.True(t, f.isExpire(0, 0))

	f.Stop(10)
}

func TestFamilyChannel_flushChunk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	chunk := NewMockChunk(ctrl)
	ctx, cancel := context.WithCancel(context.TODO())
	f := &familyChannel{
		cancel:     cancel,
		ctx:        ctx,
		chunk:      chunk,
		ch:         make(chan *compressedChunk, 1),
		statistics: metrics.NewBrokerFamilyWriteStatistics("db"),
		logger:     logger.GetLogger("test", "test"),
	}
	// compress failure
	chunk.EXPECT().Compress().Return(nil, fmt.Errorf("err"))
	f.flushChunk()
	// compress data empty
	chunk.EXPECT().Compress().Return(nil, nil)
	f.flushChunk()
	// flush data
	chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
	f.flushChunk()

	cancel()
	chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
	// family is stopped
	f.flushChunk()
}

func TestFamilyChannel_writeTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	cases := []struct {
		name    string
		prepare func(f *familyChannel)
	}{
		{
			name: "stop family, no data need flush",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(true)
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, compress fail",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(nil, fmt.Errorf("err"))
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, compress nil result",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(nil, nil)
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, compress empty result",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(&compressedChunk{}, nil)
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, send fail",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return nil, fmt.Errorf("err")
				}
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, create stream ok, but send failure",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Close()
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, create stream ok, but send failure, EOF, close stream ok",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Close().Return(nil)
				stream.EXPECT().Send(gomock.Any()).Return(io.EOF)
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, create stream ok, but send failure, EOF, close stream failure",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Close().Return(fmt.Errorf("err"))
				stream.EXPECT().Send(gomock.Any()).Return(io.EOF)
				go func() {
					f.Stop(10)
				}()
			},
		},
		{
			name: "stop family, send successfully",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(false)
				chunk.EXPECT().Compress().Return(&compressedChunk{1, 2, 3}, nil)
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Close().Return(fmt.Errorf("err"))
				stream.EXPECT().Send(gomock.Any()).Return(nil)
				go func() {
					f.Stop(timeutil.OneSecond)
				}()
			},
		},
		{
			name: "stop family, send last message in ch fail",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(true).AnyTimes()
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return nil, fmt.Errorf("err")
				}
				go func() {
					f.cancel()
					f.ch <- &compressedChunk{1, 2, 3}
					time.Sleep(5 * time.Millisecond)
					close(f.ch)
				}()
			},
		},
		{
			name: "send msg",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				chunk.EXPECT().IsEmpty().Return(true).AnyTimes()
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Send(gomock.Any()).Return(nil)
				stream.EXPECT().Close().Return(fmt.Errorf("err"))
				f.ch <- &compressedChunk{1, 2, 3}

				go func() {
					time.Sleep(200 * time.Millisecond)
					f.leaderChangedSignal <- struct{}{} // mock leader change
					time.Sleep(20 * time.Millisecond)
					f.Stop(10)
				}()
			},
		},
		{
			name: "send msg failure, retry drop",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				f.maxRetryBuf = 0
				chunk.EXPECT().IsEmpty().Return(true).AnyTimes()
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
				stream.EXPECT().Close().Return(nil).AnyTimes()
				f.ch <- &compressedChunk{1, 2, 3}
				f.ch <- &compressedChunk{1, 2, 3}

				go func() {
					time.Sleep(200 * time.Millisecond)
					f.Stop(10)
				}()
			},
		},
		{
			name: "send msg failure, retry",
			prepare: func(f *familyChannel) {
				chunk := NewMockChunk(ctrl)
				f.chunk = chunk
				f.maxRetryBuf = 1
				chunk.EXPECT().IsEmpty().Return(true).AnyTimes()
				stream := rpc.NewMockWriteStream(ctrl)
				f.newWriteStreamFn = func(ctx context.Context, target models.Node,
					database string, shardState *models.ShardState, familyTime int64,
					fct rpc.ClientStreamFactory) (rpc.WriteStream, error) {
					return stream, nil
				}
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
				stream.EXPECT().Send(gomock.Any()).Return(nil)
				stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
				stream.EXPECT().Close().Return(nil).AnyTimes()
				f.ch <- &compressedChunk{1, 2, 3}
				f.ch <- &compressedChunk{1, 2, 3}

				go func() {
					time.Sleep(200 * time.Millisecond)
					f.Stop(10)
				}()
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			f := &familyChannel{
				cancel:              cancel,
				ctx:                 ctx,
				ch:                  make(chan *compressedChunk, 2),
				maxRetryBuf:         1,
				checkFlushInterval:  time.Millisecond * 100,
				lastFlushTime:       atomic.NewInt64(timeutil.Now()),
				shardState:          models.ShardState{ID: 0, Leader: 1},
				leaderChangedSignal: make(chan struct{}, 1),
				stoppedSignal:       make(chan struct{}, 1),
				stoppingSignal:      make(chan struct{}, 1),
				currentTarget:       &models.StatefulNode{},
				liveNodes: map[models.NodeID]models.StatefulNode{
					1: {},
				},
				statistics: metrics.NewBrokerFamilyWriteStatistics("db"),
				logger:     logger.GetLogger("test", "test"),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			f.writeTask(context.TODO())
		})
	}
}

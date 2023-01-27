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

package context

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/series/field"
)

func TestMetricContext_HandleResponse(t *testing.T) {
	payload, _ := (&protoCommonV1.TimeSeriesList{
		FieldAggSpecs: []*protoCommonV1.AggregatorSpec{
			{
				FieldName:    "test",
				FieldType:    uint32(field.Sum),
				FuncTypeList: []uint32{uint32(field.Sum)},
			},
		},
		TimeSeriesList: []*protoCommonV1.TimeSeries{{Fields: nil}},
	}).Marshal()
	emptyPayload, _ := (&protoCommonV1.TimeSeriesList{}).Marshal()
	payloadWithField, _ := (&protoCommonV1.TimeSeriesList{
		FieldAggSpecs: []*protoCommonV1.AggregatorSpec{
			{
				FieldName:    "test",
				FieldType:    uint32(field.Sum),
				FuncTypeList: []uint32{uint32(field.Sum)},
			},
		},
		TimeSeriesList: []*protoCommonV1.TimeSeries{{Fields: map[string][]byte{"test": nil}}},
	}).Marshal()
	stats := encoding.JSONMarshal(&models.NodeStats{})

	cases := []struct {
		name    string
		resp    *protoCommonV1.TaskResponse
		prepare func(metricCtx *MetricContext)
		wantErr bool
	}{
		{
			name:    "resp with err",
			resp:    &protoCommonV1.TaskResponse{ErrMsg: "err"},
			wantErr: true,
		},
		{
			name: "resp with not found, ignore it",
			prepare: func(metricCtx *MetricContext) {
				metricCtx.tolerantNotFounds = 2
			},
			resp: &protoCommonV1.TaskResponse{ErrMsg: "not found"},
		},
		{
			name:    "unmarshal payload failure",
			resp:    &protoCommonV1.TaskResponse{Payload: []byte("abc")},
			wantErr: true,
		},
		{
			name: "handle empty response",
			resp: &protoCommonV1.TaskResponse{Payload: emptyPayload},
		},
		{
			name: "handle task response without field data",
			resp: &protoCommonV1.TaskResponse{Payload: payload},
		},
		{
			name: "handle task response with field data",
			resp: &protoCommonV1.TaskResponse{Payload: payloadWithField, Stats: stats},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			metricCtx := newMetricContext(context.TODO(), nil)
			metricCtx.SetTracker(tracker.NewStageTracker(flow.NewTaskContextWithTimeout(context.TODO(), time.Minute)))
			if tt.prepare != nil {
				tt.prepare(&metricCtx)
			}
			metricCtx.HandleResponse(tt.resp, "leaf")
			if (metricCtx.err != nil) != tt.wantErr {
				t.Fatalf("fail test, case: %s", tt.name)
			}
		})
	}
}

func TestMetricContext_waitResponse(t *testing.T) {
	t.Run("time out", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		metricCtx := newMetricContext(ctx, nil)
		go func() {
			cancel()
		}()
		err := metricCtx.waitResponse()
		assert.Equal(t, constants.ErrTimeout, err)
	})
	t.Run("completed", func(t *testing.T) {
		metricCtx := newMetricContext(context.TODO(), nil)
		go func() {
			close(metricCtx.doneCh)
		}()
		err := metricCtx.waitResponse()
		assert.NoError(t, err)
	})
	t.Run("completed with err", func(t *testing.T) {
		metricCtx := newMetricContext(context.TODO(), nil)
		metricCtx.err = fmt.Errorf("err")
		go func() {
			close(metricCtx.doneCh)
		}()
		err := metricCtx.waitResponse()
		assert.Error(t, err)
	})
}

func TestMetricContext_checkErr(t *testing.T) {
	ctx := newMetricContext(context.TODO(), nil)

	cases := []struct {
		name     string
		errMsg   string
		prepare  func()
		assertFn func(ignore bool, err error)
	}{
		{
			name: "empty err msg",
			assertFn: func(ignore bool, err error) {
				assert.False(t, ignore)
				assert.NoError(t, err)
			},
		},
		{
			name:   "fail err msg",
			errMsg: "err",
			assertFn: func(ignore bool, err error) {
				assert.True(t, ignore)
				assert.Error(t, err)
			},
		},
		{
			name:   "ignore not found",
			errMsg: "not found",
			assertFn: func(ignore bool, err error) {
				assert.True(t, ignore)
				assert.NoError(t, err)
			},
		},
		{
			name:   "ignore not found",
			errMsg: "not found",
			prepare: func() {
				_, _ = ctx.checkError("not found")
			},
			assertFn: func(ignore bool, err error) {
				assert.True(t, ignore)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			ctx.tolerantNotFounds = 2
			if tt.prepare != nil {
				tt.prepare()
			}
			ignore, err := ctx.checkError(tt.errMsg)
			tt.assertFn(ignore, err)
		})
	}
}

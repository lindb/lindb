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

package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/series/tag"
)

func TestBaseRuntime_SystemCollector(t *testing.T) {
	r := NewBaseRuntime(context.TODO(), config.Monitor{}, linmetric.RootRegistry, tag.Tags{})
	r.SystemCollector()
}

func TestBaseRuntime_NativePusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newNativeProtoPusher = monitoring.NewNativeProtoPusher
		ctrl.Finish()
	}()

	r := NewBaseRuntime(context.TODO(), config.Monitor{}, linmetric.RootRegistry, tag.Tags{})
	r.NativePusher()
	assert.Nil(t, r.pusher)

	pusher := monitoring.NewMockNativePusher(ctrl)
	newNativeProtoPusher = func(_ context.Context, _ string, _, _ time.Duration,
		_ *linmetric.Registry, _ tag.Tags) monitoring.NativePusher {
		return pusher
	}
	r = NewBaseRuntime(context.TODO(), config.Monitor{ReportInterval: 1000}, linmetric.RootRegistry, tag.Tags{})
	ch := make(chan struct{})
	pusher.EXPECT().Start().Do(func() {
		close(ch)
	})
	pusher.EXPECT().Stop()
	r.NativePusher()
	assert.NotNil(t, r.pusher)
	<-ch
	r.Shutdown()
}

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

package kv

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestJobScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		InitStoreManager(nil)
		ctrl.Finish()
	}()

	storeMgr := NewMockStoreManager(ctrl)
	InitStoreManager(storeMgr)
	store := NewMockStore(ctrl)
	store.EXPECT().compact().AnyTimes()
	storeMgr.EXPECT().GetStores().Return([]Store{store}).AnyTimes()
	opt := StoreOptions{CompactCheckInterval: 1}

	t.Run("startup job", func(t *testing.T) {
		js := NewJobScheduler(context.TODO(), opt)
		assert.False(t, js.IsRunning())
		js.Startup()
		assert.True(t, js.IsRunning())
		js.Startup()
		assert.True(t, js.IsRunning())
		time.Sleep(1500 * time.Millisecond)
	})

	t.Run("shutdown job", func(t *testing.T) {
		js := NewJobScheduler(context.TODO(), opt)
		assert.False(t, js.IsRunning())
		js.Startup()
		assert.True(t, js.IsRunning())
		js.Shutdown()
		assert.False(t, js.IsRunning())
		js.Shutdown()
		assert.False(t, js.IsRunning())
		time.Sleep(100 * time.Millisecond)
	})
}

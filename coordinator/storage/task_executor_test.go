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

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/tsdb"
)

func TestTaskExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	engine := tsdb.NewMockEngine(ctrl)
	repo := state.NewMockRepository(ctrl)
	exec := NewTaskExecutor(context.TODO(), &models.Node{IP: "1.1.1.1", Port: 5000}, repo, engine)
	assert.NotNil(t, exec)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), true).Return(nil)
	exec.Run()
	time.Sleep(100 * time.Millisecond)
	err := exec.Close()
	if err != nil {
		t.Fatal(err)
	}
}

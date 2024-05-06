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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/kv/version"
)

func TestCompactionState_AddOutputFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	state := newCompactionState(100, snapshot, nil)
	file := version.NewFileMeta(1, 1, 199, 10)
	state.addOutputFile(file)
	assert.Equal(t, file, state.outputs[0])
	assert.Equal(t, 1, len(state.outputs))
}

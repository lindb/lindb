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

package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventType_String(t *testing.T) {
	assert.Equal(t, "unknown", EventType(-1).String())
	assert.Equal(t, "DatabaseConfigChanged", DatabaseConfigChanged.String())
	assert.Equal(t, "DatabaseConfigDeletion", DatabaseConfigDeletion.String())
	assert.Equal(t, "NodeStartup", NodeStartup.String())
	assert.Equal(t, "NodeFailure", NodeFailure.String())
	assert.Equal(t, "StorageStateChanged", StorageStateChanged.String())
	assert.Equal(t, "StorageStateDeletion", StorageStateDeletion.String())
	assert.Equal(t, "ShardAssignmentDeletion", ShardAssignmentDeletion.String())
	assert.Equal(t, "ShardAssignmentChanged", ShardAssignmentChanged.String())
	assert.Equal(t, "StorageConfigDeletion", StorageConfigDeletion.String())
	assert.Equal(t, "StorageConfigChanged", StorageConfigChanged.String())
}

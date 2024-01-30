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

package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pathName      = "name"
	slashPathName = "/name"
)

func TestGetDatabaseAssignPath(t *testing.T) {
	assert.Equal(t, ShardAssignmentPath+slashPathName, GetDatabaseAssignPath(pathName))
}

func TestGetDatabaseConfigPath(t *testing.T) {
	assert.Equal(t, DatabaseConfigPath+slashPathName, GetDatabaseConfigPath(pathName))
}

func TestGetDatabaseLimitPath(t *testing.T) {
	assert.Equal(t, DatabaseLimitPath+slashPathName, GetDatabaseLimitPath(pathName))
}

func TestGetNodePath(t *testing.T) {
	assert.Equal(t, StorageLiveNodesPath+slashPathName, GetStorageLiveNodePath(pathName))
}

func TestGetBrokerClusterConfigPath(t *testing.T) {
	assert.Equal(t, BrokerConfigPath+slashPathName, GetBrokerClusterConfigPath(pathName))
}

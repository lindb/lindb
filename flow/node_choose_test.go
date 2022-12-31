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

package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestBuildPhysicalPlan(t *testing.T) {
	plan := BuildPhysicalPlan("test", nil, 1)
	assert.Equal(t, "test", plan.Database)
	assert.Nil(t, plan.Targets)

	liveNodes := []models.StatelessNode{{HostIP: "1.1.2.3"}, {HostIP: "1.1.1.1"}, {HostIP: "2.2.2.2"}}
	plan = BuildPhysicalPlan("test", liveNodes, 1)
	assert.Len(t, plan.Targets, 1)

	plan = BuildPhysicalPlan("test", liveNodes, 2)
	assert.Len(t, plan.Targets, 2)

	plan = BuildPhysicalPlan("test", liveNodes, 5)
	assert.Len(t, plan.Targets, 3)
}

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

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
)

func TestNode_Indicator(t *testing.T) {
	node := &StatelessNode{HostIP: "1.1.1.1", GRPCPort: 19000}
	indicator := node.Indicator()
	assert.Equal(t, "1.1.1.1:19000", indicator)
	node2, err := ParseNode(indicator)
	assert.NoError(t, err)
	node3 := node2.(*StatelessNode)
	assert.Equal(t, node, node3)
}

func TestNode_ParseNode(t *testing.T) {
	_, err := ParseNode("xxx:123")
	assert.Error(t, err)

	_, err = ParseNode("1.1.1.1123")
	assert.Error(t, err)

	_, err = ParseNode("1.1.1.1:-1")
	assert.Error(t, err)

	_, err = ParseNode("1.1.1.1:65536")
	assert.Error(t, err)

	node, err := ParseNode("1.1.1.1:65535")
	assert.NoError(t, err)
	node1 := node.(*StatelessNode)

	assert.Equal(t, node1.HostIP, "1.1.1.1")
	assert.Equal(t, node1.GRPCPort, uint16(65535))

	_, err = ParseNode(":123")
	assert.Error(t, err)
}

func TestMaster_ToTable(t *testing.T) {
	assert.NotEmpty(t, (&Master{
		Node:      &StatelessNode{},
		ElectTime: timeutil.Now(),
	}).ToTable())
}

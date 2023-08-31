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

	"github.com/lindb/common/pkg/timeutil"
)

func TestNode_Indicator(t *testing.T) {
	node := &StatelessNode{HostIP: "1.1.1.1", HTTPPort: 19000}
	assert.Equal(t, "1.1.1.1:19000", node.Indicator())
	node = &StatelessNode{HostIP: "1.1.1.1", GRPCPort: 19000}
	indicator := node.Indicator()
	assert.Equal(t, "1.1.1.1:19000", indicator)
	assert.Equal(t, "http://1.1.1.1:8080", (&StatelessNode{HostIP: "1.1.1.1", HTTPPort: 8080}).HTTPAddress())
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
	rows, rs := (&Master{
		Node:      &StatelessNode{},
		ElectTime: timeutil.Now(),
	}).ToTable()
	assert.NotEmpty(t, rs)
	assert.Equal(t, rows, 1)
}

func TestStatelessNodes_ToTable(t *testing.T) {
	rows, rs := (StatelessNodes{}).ToTable()
	assert.Empty(t, rs)
	assert.Equal(t, rows, 0)

	rows, rs = (StatelessNodes{{OnlineTime: timeutil.Now()}}).ToTable()
	assert.NotEmpty(t, rs)
	assert.Equal(t, rows, 1)
}

func TestNodeID(t *testing.T) {
	assert.Equal(t, 1, NodeID(1).Int())
	assert.Equal(t, "1", NodeID(1).String())
	assert.Equal(t, NodeID(1), ParseNodeID("1"))
}

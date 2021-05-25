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
)

func TestNode_Indicator(t *testing.T) {
	node := &Node{IP: "1.1.1.1", Port: 19000}
	indicator := node.Indicator()
	assert.Equal(t, "1.1.1.1:19000", indicator)
	node2, err := ParseNode(indicator)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, *node, *node2)
}

func TestNode_ParseNode(t *testing.T) {
	if _, err := ParseNode("xxx:123"); err == nil {
		t.Fatal("should be error")
	}

	if _, err := ParseNode("1.1.1.1123"); err == nil {
		t.Fatal("should be error")
	}

	if _, err := ParseNode("1.1.1.1:-1"); err == nil {
		t.Fatal("should be error")
	}

	if _, err := ParseNode("1.1.1.1:65536"); err == nil {
		t.Fatal("should be error")
	}

	node, err := ParseNode("1.1.1.1:65535")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, node.IP, "1.1.1.1")
	assert.Equal(t, node.Port, uint16(65535))

	if _, err = ParseNode(":123"); err == nil {
		t.Fatal(err)
	}
}

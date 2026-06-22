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

package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestUseStatement(t *testing.T) {
	stmt, err := GetParser().CreateStatement("use test", NewNodeIDAllocator())
	assert.NoError(t, err)
	checkStatement(t, &Use{
		Database: &Identifier{Value: "test"},
	}, stmt)
}

// checkStatement compares two Statement trees for semantic equality, ignoring
// all BaseNode metadata fields (ID, Location, Text). These are assigned by the
// parser and must not be hard-coded in tests.
func checkStatement(t *testing.T, want, got Statement) {
	t.Helper()
	// IgnoreFields skips all parser-internal BaseNode fields in every node of the tree.
	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(BaseNode{}, "ID", "Location", "Text")); diff != "" {
		t.Errorf("statement mismatch (-want +got):\n%s", diff)
	}
}

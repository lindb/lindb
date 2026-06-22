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

package plan

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lindb/common/pkg/encoding"
	"github.com/stretchr/testify/require"
)

// checkPlanNode compares two plan trees for structural equality, ignoring
// BaseNode.ID values. IDs are assigned by the planner/optimizer in allocation
// order and must not be hard-coded in tests.
func checkPlanNode(t *testing.T, want, got PlanNode) {
	t.Helper()
	// IgnoreFields skips BaseNode.ID in every embedded BaseNode across the whole tree.
	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(BaseNode{}, "ID")); diff != "" {
		t.Errorf("plan node mismatch (-want +got):\n%s", diff)
	}
}

// TestEncoding_RoundTrip verifies that a plan tree survives a JSON marshal/unmarshal
// round-trip with the correct concrete type restored via the @type/@data wrapper.
func TestEncoding_RoundTrip(t *testing.T) {
	// Build a minimal plan tree: TableScanNode with one output symbol.
	want := &TableScanNode{
		BaseNode:      BaseNode{ID: 1},
		OutputSymbols: []*Symbol{{Name: "cpu"}},
	}

	// Wrap in PlanFragment so that the PlanNode interface encoder is exercised.
	fragment := &PlanFragment{Root: want}
	data := encoding.JSONMarshal(fragment)
	require.NotEmpty(t, data)

	// Unmarshal back and verify the Root node matches the original.
	var got PlanFragment
	require.NoError(t, encoding.JSONUnmarshal(data, &got))

	checkPlanNode(t, want, got.Root)
}

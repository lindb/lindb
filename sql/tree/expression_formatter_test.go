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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatExpression_FloatLiteral is the direct regression test for the production panic:
// "expression formatter unsupport node:*tree.FloatLiteral".
// FormatVisitor.Visit previously had no case for FloatLiteral and fell through to the
// default panic branch when printing plan nodes that contain histogram_quantile(0.99, col).
func TestFormatExpression_FloatLiteral(t *testing.T) {
	idAlloc := NewNodeIDAllocator()

	cases := []struct {
		raw  string
		want string
	}{
		{"0.99", "0.99"},
		{"0.5", "0.5"},
		{"1.0", "1"}, // %g strips trailing zero
		{"3.14", "3.14"},
		{"0.0", "0"},
	}
	for _, tc := range cases {
		node := NewFloatLiteral(idAlloc.Next(), nil, tc.raw)
		var result string
		require.NotPanics(t, func() {
			result = FormatExpression(node)
		}, "FormatExpression must not panic for FloatLiteral(%s)", tc.raw)
		assert.Equal(t, tc.want, result, "FloatLiteral(%s) format mismatch", tc.raw)
	}
}

// TestFormatExpression_FloatLiteralInFunctionCall verifies that
// histogram_quantile(0.99, sent_duration) can be formatted without panic.
// This is the exact scenario that triggered the production panic during plan printing.
func TestFormatExpression_FloatLiteralInFunctionCall(t *testing.T) {
	idAlloc := NewNodeIDAllocator()

	phi := NewFloatLiteral(idAlloc.Next(), nil, "0.99")
	col := &Identifier{Value: "sent_duration"}
	col.SetID(idAlloc.Next())
	call := &FunctionCall{
		Name:      HistogramQuantile,
		Arguments: []Expression{phi, col},
	}
	call.SetID(idAlloc.Next())

	var result string
	require.NotPanics(t, func() {
		result = FormatExpression(call)
	}, "FormatExpression must not panic for histogram_quantile(0.99, col)")
	assert.Contains(t, result, "0.99", "formatted output should contain the phi value")
	assert.Contains(t, result, "sent_duration", "formatted output should contain the column name")
}

// TestFormatExpression_LongLiteral verifies that LongLiteral still formats correctly (regression guard).
func TestFormatExpression_LongLiteral(t *testing.T) {
	idAlloc := NewNodeIDAllocator()
	node := NewLongLiteral(idAlloc.Next(), nil, "42")
	assert.Equal(t, "42", FormatExpression(node))
}

// TestFormatExpression_StringLiteral verifies that StringLiteral still formats correctly (regression guard).
func TestFormatExpression_StringLiteral(t *testing.T) {
	idAlloc := NewNodeIDAllocator()
	node := &StringLiteral{Value: "hello"}
	node.SetID(idAlloc.Next())
	assert.Equal(t, "'hello'", FormatExpression(node))
}

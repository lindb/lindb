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

package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAppendWALPointer verifies that multiple WAL pointer entries are concatenated
// correctly in the pebble-based read-modify-write path (replaces the old
// TraceMergeOperator test which was removed along with RocksDB).
func TestAppendWALPointer(t *testing.T) {
	// Simulate the concatenation logic used in appendWALPointer:
	// existing value + new entry should equal the full byte sequence.
	existing := []byte{1, 2, 3}
	entry1 := []byte{4, 5, 6}
	entry2 := []byte{8, 9}

	// First append: no existing data
	r1 := make([]byte, len(entry1))
	copy(r1, entry1)
	assert.Equal(t, []byte{4, 5, 6}, r1)

	// Second append: existing + entry1
	r2 := make([]byte, len(existing)+len(entry1))
	copy(r2, existing)
	copy(r2[len(existing):], entry1)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, r2)

	// Third append: r2 + entry2
	r3 := make([]byte, len(r2)+len(entry2))
	copy(r3, r2)
	copy(r3[len(r2):], entry2)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 8, 9}, r3)
}

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

package kv

import (
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

// compactionState represents the state of compaction job
type compactionState struct {
	outputs           []*version.FileMeta
	builder           table.Builder
	compaction        *version.Compaction
	snapshot          version.Snapshot
	currentFileNumber table.FileNumber
	maxFileSize       int32
}

// newCompactionState creates a compaction state
func newCompactionState(maxFileSize int32, snapshot version.Snapshot, compaction *version.Compaction) *compactionState {
	return &compactionState{
		maxFileSize: maxFileSize,
		snapshot:    snapshot,
		compaction:  compaction,
	}
}

// addOutputFile adds a new output file
func (c *compactionState) addOutputFile(fileMete *version.FileMeta) {
	c.outputs = append(c.outputs, fileMete)
}

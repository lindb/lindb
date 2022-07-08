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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
)

func TestCompaction(t *testing.T) {
	f1 := FileMeta{fileNumber: 1, minKey: 10, maxKey: 100}
	f2 := FileMeta{fileNumber: 2, minKey: 1000, maxKey: 1001}
	f4 := FileMeta{fileNumber: 4, minKey: 100, maxKey: 200}
	compaction := NewCompaction(1, 0,
		[]*FileMeta{&f1, &f2},
		[]*FileMeta{&f4},
	)
	assert.Equal(t, 0, compaction.GetLevel())
	assert.False(t, compaction.IsTrivialMove())
	assert.Equal(t, []*FileMeta{&f1, &f2}, compaction.GetLevelFiles())
	assert.Equal(t, [][]*FileMeta{{&f1, &f2}, {&f4}}, compaction.GetInputs())
	assert.True(t, compaction.GetEditLog().IsEmpty())
	compaction.MarkInputDeletes()
	compaction.AddFile(1, &FileMeta{fileNumber: 6, minKey: 10, maxKey: 1001})
	assert.False(t, compaction.GetEditLog().IsEmpty())

	compaction = NewCompaction(1, 0,
		[]*FileMeta{&f2},
		nil,
	)
	assert.True(t, compaction.IsTrivialMove())
	compaction.DeleteFile(0, 2)
	assert.False(t, compaction.GetEditLog().IsEmpty())
}

func TestCompaction_AddReferenceFiles(t *testing.T) {
	compaction := NewCompaction(1, 0,
		[]*FileMeta{},
		nil,
	)
	assert.True(t, compaction.GetEditLog().IsEmpty())
	compaction.AddReferenceFiles([]Log{CreateNewReferenceFile(FamilyID(10), table.FileNumber(10))})
	assert.False(t, compaction.GetEditLog().IsEmpty())
}

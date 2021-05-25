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
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestRollup_RollupFiles(t *testing.T) {
	rollup := newRollup()
	rollup.addRollupFile(10, 100)
	result := rollup.getRollupFiles()
	assert.Equal(t, map[table.FileNumber]timeutil.Interval{10: 100}, result)
	rollup.removeRollupFile(100)
	result = rollup.getRollupFiles()
	assert.Equal(t, map[table.FileNumber]timeutil.Interval{10: 100}, result)
	rollup.removeRollupFile(10)
	result = rollup.getRollupFiles()
	assert.Empty(t, result)
}

func TestRollup_Reference(t *testing.T) {
	rollup := newRollup()
	rollup.addReferenceFile(10, 100)
	rollup.addReferenceFile(10, 100)
	result := rollup.getReferenceFiles()
	assert.Equal(t, map[FamilyID][]table.FileNumber{10: {100}}, result)
	rollup.addReferenceFile(10, 200)
	result = rollup.getReferenceFiles()
	assert.Equal(t, map[FamilyID][]table.FileNumber{10: {100, 200}}, result)
	rollup.removeReferenceFile(100, 100)
	result = rollup.getReferenceFiles()
	assert.Equal(t, map[FamilyID][]table.FileNumber{10: {100, 200}}, result)
	rollup.removeReferenceFile(10, 200)
	result = rollup.getReferenceFiles()
	assert.Equal(t, map[FamilyID][]table.FileNumber{10: {100}}, result)
	rollup.removeReferenceFile(10, 100)
	result = rollup.getReferenceFiles()
	assert.Empty(t, result)
}

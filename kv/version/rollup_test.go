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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

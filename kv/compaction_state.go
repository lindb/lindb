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
	currentFileNumber int64
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

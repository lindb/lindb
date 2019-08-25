package kv

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/version"
)

func TestCompactionState_AddOutputFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	state := newCompactionState(100, snapshot, nil)
	file := version.NewFileMeta(1, 1, 199, 10)
	state.addOutputFile(file)
	assert.Equal(t, file, state.outputs[0])
	assert.Equal(t, 1, len(state.outputs))
}

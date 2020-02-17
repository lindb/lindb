package tsdb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/replication"
)

var _testSequencePath = filepath.Join(testPath, replicaDir)

func TestSequence_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
		listDirFunc = fileutil.ListDir
		newSequenceFunc = replication.NewSequence
		_ = fileutil.RemoveDir(testPath)
	}()

	// create dir err
	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	seq, err := newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	mkdirFunc = fileutil.MkDirIfNotExist

	// create seq success
	seq, err = newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)
	s, err := seq.getOrCreateSequence("remote-test")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// reopen list err
	listDirFunc = func(path string) (strings []string, e error) {
		return nil, fmt.Errorf("err")
	}
	seq, err = newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	// reopen new sequence err
	listDirFunc = fileutil.ListDir
	newSequenceFunc = func(dirPath string) (sequence replication.Sequence, e error) {
		return nil, fmt.Errorf("err")
	}
	seq, err = newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	// sync error
	s1 := replication.NewMockSequence(ctrl)
	newSequenceFunc = func(dirPath string) (sequence replication.Sequence, e error) {
		return s1, nil
	}
	s1.EXPECT().GetAckSeq().Return(int64(10))
	s1.EXPECT().SetHeadSeq(int64(10))
	s1.EXPECT().Sync().Return(fmt.Errorf("err"))
	seq, err = newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	// reopen success
	newSequenceFunc = replication.NewSequence
	seq, err = newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)
}

func TestSequence_getOrCreateSequence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		newSequenceFunc = replication.NewSequence
		_ = fileutil.RemoveDir(testPath)
	}()

	// create seq success
	seq, err := newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)
	s, err := seq.getOrCreateSequence("remote-test")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s2, err := seq.getOrCreateSequence("remote-test")
	assert.NoError(t, err)
	assert.Equal(t, s, s2)

	// create err
	newSequenceFunc = func(dirPath string) (sequence replication.Sequence, e error) {
		return nil, fmt.Errorf("err")
	}
	s2, err = seq.getOrCreateSequence("remote-test-2")
	assert.Error(t, err)
	assert.Nil(t, s2)
}

func TestSequence_ack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	// create seq success
	seq, err := newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)
	s, err := seq.getOrCreateSequence("remote-test")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	heads := seq.getAllHeads()
	err = seq.ack(heads)
	assert.NoError(t, err)

	// ack not match
	err = seq.ack(map[string]int64{"no": int64(10)})
	assert.NoError(t, err)

	seq1 := seq.(*replicaSequence)
	seq1.sequenceMap.Store("not-match", "test")
	err = seq.ack(map[string]int64{"not-match": int64(10)})
	assert.NoError(t, err)
}

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

package tsdb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue"
)

var _testSequencePath = filepath.Join(testPath, replicaDir)

func TestSequence_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		listDir = fileutil.ListDir
		newSequenceFunc = queue.NewSequence
		_ = fileutil.RemoveDir(testPath)
	}()

	// create dir err
	mkDirIfNotExist = func(path string) error {
		return fmt.Errorf("err")
	}
	seq, err := newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	mkDirIfNotExist = fileutil.MkDirIfNotExist

	// create seq success
	seq, err = newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)
	s, err := seq.getOrCreateSequence("remote-test")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// reopen list err
	listDir = func(path string) (strings []string, e error) {
		return nil, fmt.Errorf("err")
	}
	seq, err = newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	// reopen new sequence err
	listDir = fileutil.ListDir
	newSequenceFunc = func(dirPath string) (sequence queue.Sequence, e error) {
		return nil, fmt.Errorf("err")
	}
	seq, err = newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	// sync error
	s1 := queue.NewMockSequence(ctrl)
	newSequenceFunc = func(dirPath string) (sequence queue.Sequence, e error) {
		return s1, nil
	}
	s1.EXPECT().GetAckSeq().Return(int64(10))
	s1.EXPECT().SetHeadSeq(int64(10))
	s1.EXPECT().Sync().Return(fmt.Errorf("err"))
	seq, err = newReplicaSequence(_testSequencePath)
	assert.Error(t, err)
	assert.Nil(t, seq)

	// reopen success
	newSequenceFunc = queue.NewSequence
	seq, err = newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)
}

func TestSequence_getOrCreateSequence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		newSequenceFunc = queue.NewSequence
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
	newSequenceFunc = func(dirPath string) (sequence queue.Sequence, e error) {
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

	err = seq.Close()
	assert.NoError(t, err)
}

func TestReplicaSequence_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	// create seq success
	seq, err := newReplicaSequence(_testSequencePath)
	assert.NoError(t, err)
	assert.NotNil(t, seq)

	seq1 := seq.(*replicaSequence)
	mockSeq := queue.NewMockSequence(ctrl)
	mockSeq.EXPECT().Close().Return(fmt.Errorf("err"))
	seq1.sequenceMap.Store("test", mockSeq)
	err = seq.Close()
	assert.Error(t, err)
}

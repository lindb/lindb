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

package queue

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/queue/page"
)

func TestNewConsumerGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(t.TempDir(), t.Name())

	defer func() {
		newPageFactoryFunc = page.NewFactory
		ctrl.Finish()
	}()

	// case 1: new meta page factory
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return nil, fmt.Errorf("err")
	}
	fo, err := NewConsumerGroup(dir, "f1", nil)
	assert.Error(t, err)
	assert.Nil(t, fo)
	// case 2: acquire meta page err
	pageFct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return pageFct, nil
	}
	pageFct.EXPECT().Close().Return(fmt.Errorf("err"))
	pageFct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	fo, err = NewConsumerGroup(dir, "f1", nil)
	assert.Error(t, err)
	assert.Nil(t, fo)
}

func TestConsumerGroup_IsEmpty(t *testing.T) {
	dir := path.Join(t.TempDir(), t.Name())

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	f1, err := fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), f1.ConsumedSeq())
	assert.Equal(t, int64(-1), f1.AcknowledgedSeq())
	assert.True(t, f1.IsEmpty())
	msg := []byte("123")
	err = fq.Queue().Put(msg)
	assert.NoError(t, err)
	assert.False(t, f1.IsEmpty())
	idx := f1.Consume()
	assert.False(t, f1.IsEmpty())
	assert.Equal(t, idx, f1.ConsumedSeq())
	f1.Close()
	fq.Close()

	fq, err = NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	f1, err = fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), f1.ConsumedSeq())
	assert.Equal(t, int64(-1), f1.AcknowledgedSeq())
	assert.False(t, f1.IsEmpty())
	idx = f1.Consume()
	f1.Ack(idx)
	assert.True(t, f1.IsEmpty())
	f1.Close()
	fq.Close()

	fq, err = NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	f1, err = fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)
	assert.Equal(t, idx, f1.ConsumedSeq())
	assert.Equal(t, idx, f1.AcknowledgedSeq())
}

func TestConsumerGroup_one_consumer(t *testing.T) {
	dir := path.Join(t.TempDir(), t.Name())

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	assert.Empty(t, fq.ConsumerGroupNames())
	assert.Equal(t, int64(-1), fq.Queue().AppendedSeq())
	assert.Equal(t, int64(-1), fq.Queue().AcknowledgedSeq())

	f1, err := fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, consumerGroupDirName, "f1"), f1.Name())
	assert.Equal(t, int64(-1), f1.Queue().Queue().AppendedSeq())
	assert.Equal(t, int64(-1), f1.Queue().Queue().AcknowledgedSeq())
	assert.Equal(t, SeqNoNewMessageAvailable, f1.Consume())
	assert.Equal(t, int64(0), f1.Pending())

	// msg 0
	msg := []byte("123")
	err = fq.Queue().Put(msg)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), f1.Pending())

	fseq := f1.Consume()
	assert.Equal(t, int64(0), fseq)
	assert.Equal(t, int64(0), f1.ConsumedSeq())
	assert.Equal(t, int64(-1), f1.AcknowledgedSeq())
	assert.Equal(t, int64(0), f1.Pending())

	fmsg, err := f1.Queue().Queue().Get(0)
	assert.NoError(t, err)
	assert.Equal(t, msg, fmsg)
	assert.Equal(t, SeqNoNewMessageAvailable, f1.Consume())

	// msg 1
	msg1 := []byte("456")

	err = fq.Queue().Put(msg1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), f1.ConsumedSeq())
	assert.Equal(t, int64(1), f1.Pending())
	// msg 2
	msg2 := []byte("789")
	err = fq.Queue().Put(msg2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), fq.Queue().AppendedSeq())
	assert.Equal(t, int64(2), f1.Pending())

	fseq = f1.Consume()
	assert.Equal(t, int64(1), fseq)
	assert.Equal(t, int64(2), fq.Queue().AppendedSeq())
	assert.Equal(t, int64(1), f1.Pending())

	fmsg, err = f1.Queue().Queue().Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, msg1, fmsg)

	f1.Ack(fseq) // ack 1
	assert.Equal(t, fseq, f1.AcknowledgedSeq())

	fseq = f1.Consume()
	assert.Equal(t, int64(2), fseq)
	assert.Equal(t, int64(2), f1.ConsumedSeq())
	assert.Equal(t, int64(0), f1.Pending())

	fmsg, err = f1.Queue().Queue().Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, msg2, fmsg)
	f1.Ack(fseq) // akc 2
	assert.Equal(t, fseq, f1.AcknowledgedSeq())
	assert.Equal(t, int64(0), f1.Pending())

	fq.Close()
	// reopen
	fq, err = NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	f1, err = fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), f1.AcknowledgedSeq())
	assert.Equal(t, int64(2), f1.ConsumedSeq())
	assert.Equal(t, int64(0), f1.Pending())
	fq.Close()
}

func TestConsumerGroup_SetConsumedSeq(t *testing.T) {
	dir := path.Join(t.TempDir(), t.Name())

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)

	f1, err := fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)

	f1.SetConsumedSeq(1)
	assert.Equal(t, int64(1), f1.ConsumedSeq())

	err = fq.Queue().Put([]byte("123"))
	assert.NoError(t, err)

	err = fq.Queue().Put([]byte("456"))
	assert.NoError(t, err)

	seq := f1.Consume()
	assert.Equal(t, SeqNoNewMessageAvailable, seq)

	// reset head consume sequence
	f1.SetConsumedSeq(-1)
	seq = f1.Consume()
	assert.Equal(t, int64(0), seq)

	seq = f1.Consume()
	assert.Equal(t, int64(1), seq)

	f1.Ack(1)

	f1.SetConsumedSeq(0)
	assert.Equal(t, int64(0), f1.ConsumedSeq())
	fq.Close()
}
